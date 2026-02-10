package ads

import (
	"fmt"
	"io"

	adserrors "github.com/jarmocluyse/ads-go/pkg/ads/ads-errors"
	amsheader "github.com/jarmocluyse/ads-go/pkg/ads/ams-header"
)

// receive handles incoming data from the ADS router.
func (c *Client) receive() {
	c.logger.Info("receive: Starting receive goroutine.")
	defer func() {
		c.conn.Close()
		c.logger.Info("receive: Receive goroutine terminated.")
	}()

	// Temporary buffer for reading from connection
	tempBuf := make([]byte, 4096) // Read in chunks

	for {
		n, err := c.conn.Read(tempBuf)
		if err != nil {
			if err == io.EOF {
				c.logger.Info("receive: Connection closed by remote.")
			} else {
				c.logger.Error("receive: Error reading from connection", "error", err)
			}
			return // Exit goroutine on error or EOF
		}

		// Write read data to the receive buffer
		c.receiveBuffer.Write(tempBuf[:n])

		// Process packets from the receive buffer
		c.processReceiveBuffer()
	}
}

// process the received data
func (c *Client) processReceiveBuffer() {
	for {
		totalPacketLength, err := c.checkTcpPacketLength()
		if err != nil {
			return // Not enough data for full packet
		}
		// Extract the full packet
		fullPacket := make([]byte, totalPacketLength)
		c.receiveBuffer.Read(fullPacket)

		packet := c.parseAmsPacket(fullPacket)
		c.logger.Debug("receive: Parsed AMS packet", "invokeID", packet.InvokeId, "data", packet)

		c.mutex.Lock()
		ch, ok := c.requests[packet.InvokeId]
		c.mutex.Unlock()

		if ok {
			c.logger.Debug("receive: Found channel for InvokeID, sending response.", "invokeID", packet.InvokeId)
			if packet.ErrorCode != 0 {
				errorString := adserrors.ErrorCodeToString(packet.ErrorCode)
				c.logger.Error("receive: ADS error received", "invokeID", packet.InvokeId, "errorCode", packet.ErrorCode, "errorDesc", errorString)
				ch <- Response{Error: fmt.Errorf("ADS error: %s", errorString)}
			} else {
				ch <- Response{Data: packet.Data}
			}
		} else {
			c.logger.Warn("receive: No channel found for InvokeID, discarding packet.", "invokeID", packet.InvokeId)
		}
	}
}

// Read packet length from AMS/TCP header (bytes 2-5)
// We need to peek without advancing the buffer's read pointer
// to check if we received the full packet
func (c *Client) checkTcpPacketLength() (packetLenght uint32, error error) {
	// Use the ams-header module to check packet length
	totalPacketLength, err := amsheader.CheckTCPPacketLength(c.receiveBuffer.Bytes())
	if err != nil {
		return 0, err
	}
	c.logger.Debug("Full package in buffer", "totalPacketLength", totalPacketLength)
	return totalPacketLength, nil
}

type AmsPacket struct {
	TargetAmsAddress AmsAddress // target address
	SourceAmsAddress AmsAddress // source address
	AdsCommand       uint16     // ADS command to be send
	StateFlags       uint16     // state flags
	DataLength       uint32     // length og the data
	ErrorCode        uint32     // AMS header error code (not ADS payload error)
	InvokeId         uint32     // invoke Id
	Data             []byte     //received data
}

// parse the ams header
// NOTE: we now at this point the length is correct
func (c *Client) parseAmsPacket(data []byte) AmsPacket {
	// Use the ams-header module to parse
	packet, err := amsheader.ParsePacket(data)
	if err != nil {
		c.logger.Error("parseAmsPacket: Failed to parse packet", "error", err)
		return AmsPacket{}
	}

	return AmsPacket{
		TargetAmsAddress: AmsAddress{
			NetID: packet.TargetNetID,
			Port:  packet.TargetPort,
		},
		SourceAmsAddress: AmsAddress{
			NetID: packet.SourceNetID,
			Port:  packet.SourcePort,
		},
		AdsCommand: uint16(packet.Command),
		StateFlags: uint16(packet.StateFlags),
		DataLength: packet.DataLength,
		ErrorCode:  packet.ErrorCode,
		InvokeId:   packet.InvokeID,
		Data:       packet.Data,
	}
}
