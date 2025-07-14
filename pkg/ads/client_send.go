package ads

import (
	"fmt"
	"time"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// AdsCommandRequest represents a request for an ADS command.
type AdsCommandRequest struct {
	Command     types.ADSCommand
	Data        []byte
	TargetNetID string
	TargetPort  uint16
}

// send sends a command to the ADS router.
func (c *Client) send(req AdsCommandRequest) ([]byte, error) {
	c.logger.Debug("send: Preparing to send command", "command", req.Command.String(), "dataLength", len(req.Data))
	c.mutex.Lock()
	c.invokeID++
	invokeID := c.invokeID
	ch := make(chan Response)
	c.requests[invokeID] = ch
	c.mutex.Unlock()
	c.logger.Debug("send: Assigned InvokeID", "invokeID", invokeID)

	defer func() {
		c.mutex.Lock()
		delete(c.requests, invokeID)
		c.mutex.Unlock()
		c.logger.Debug("send: Cleaned up request", "invokeID", invokeID)
	}()

	target := AmsAddress{NetID: req.TargetNetID, Port: req.TargetPort}
	c.logger.Debug("send: Target AMS Address", "netID", target.NetID, "port", target.Port)

	amsHeader, err := createAmsHeader(target, c.localAmsAddr, req.Command, uint32(len(req.Data)), invokeID)
	if err != nil {
		c.logger.Error("send: Failed to create AMS header", "error", err)
		return nil, err
	}
	amsTcpHeader := createAmsTcpHeader(types.AMSTCPPortAMSCommand, uint32(len(amsHeader)+len(req.Data)))

	packet := append(amsTcpHeader, amsHeader...)
	packet = append(packet, req.Data...)
	c.logger.Debug("send: Constructed packet", "totalLength", len(packet), "packet", fmt.Sprintf("%x", packet))

	_, err = c.conn.Write(packet)
	if err != nil {
		c.logger.Error("send: Failed to write packet to connection", "error", err)
		return nil, err
	}
	c.logger.Debug("send: Packet sent. Waiting for response or timeout.")

	select {
	case response := <-ch:
		c.logger.Debug("send: Received response", "invokeID", invokeID, "response", response)
		if response.Error != nil {
			return nil, response.Error
		}
		return response.Data, nil
	case <-time.After(c.settings.Timeout):
		c.logger.Warn("send: Timeout waiting for response", "invokeID", invokeID)
		return nil, fmt.Errorf("timeout waiting for response")
	}
}
