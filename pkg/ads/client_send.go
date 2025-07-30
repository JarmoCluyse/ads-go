package ads

import (
	"fmt"
	"time"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// AdsCommandRequest represents a request for an ADS command.
type AdsCommandRequest struct {
	Command    types.ADSCommand // Ads Command to send
	TargetPort uint16           // port to send to
	Data       []byte           // data to send
}

// send sends a command to the ADS router.
func (c *Client) send(req AdsCommandRequest) ([]byte, error) {
	if c.conn == nil {
		c.logger.Error("send: Connection is nil, cannot send command")
		return nil, fmt.Errorf("connection is not established")
	}
	c.logger.Debug("send: Preparing to send command", "command", req.Command.String(), "dataLength", len(req.Data))

	invokeID, channel := c.getInvokeID()
	defer c.removeInvokeId(invokeID)

	target := AmsAddress{NetID: c.settings.TargetNetID, Port: req.TargetPort}
	c.logger.Debug("send: Target AMS Address", "netID", target.NetID, "port", target.Port)

	amsHeader, err := createAmsHeader(target, c.localAmsAddr, req.Command, uint32(len(req.Data)), invokeID)
	if err != nil {
		c.logger.Error("send: Failed to create AMS header", "error", err)
		return nil, err
	}
	dataLen := uint32(len(amsHeader) + len(req.Data))
	amsTcpHeader := createAmsTcpHeader(types.AMSTCPPortAMSCommand, dataLen)

	packet := append(amsTcpHeader, amsHeader...)
	packet = append(packet, req.Data...)
	c.logger.Debug("send: Constructed packet", "totalLength", len(packet), "packet", fmt.Sprintf("%x", packet))

	_, err = c.conn.Write(packet)
	if err != nil {
		c.logger.Error("send: Failed to write packet to connection", "error", err)
		return nil, err
	}
	c.logger.Debug("send: Packet sent. Waiting for response or timeout.", "packet", packet)

	select {
	case response := <-channel:
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
