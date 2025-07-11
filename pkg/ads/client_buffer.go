package ads

import (
	"encoding/binary"
	"io"

	"github.com/jarmoCluyse/ads-go/pkg/ads/constants"
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
		// Check if we have enough data for AMS/TCP header (6 bytes)
		if c.receiveBuffer.Len() < constants.AMSTCPHeaderLength {
			return // Not enough data for header
		}

		// Read packet length from AMS/TCP header (bytes 2-5)
		// We need to peek without advancing the buffer's read pointer
		headerBytes := c.receiveBuffer.Bytes()[:constants.AMSTCPHeaderLength]
		packetLength := binary.LittleEndian.Uint32(headerBytes[2:6])

		// Total length of the full packet (AMS/TCP header + AMS header + ADS data)
		totalPacketLength := constants.AMSTCPHeaderLength + packetLength

		if packetLength < constants.AMSHeaderLength {
			return // Not enough data for full packet
		}
		// Check if we have the full packet
		if c.receiveBuffer.Len() < int(totalPacketLength) {
			return // Not enough data for full packet
		}

		// Extract the full packet
		fullPacket := make([]byte, totalPacketLength)
		c.receiveBuffer.Read(fullPacket)

		// Now process the full packet
		amsPacket := fullPacket[constants.AMSTCPHeaderLength:]
		amsHeader := amsPacket[:constants.AMSHeaderLength]
		amsData := amsPacket[constants.AMSHeaderLength:]

		invokeID := binary.LittleEndian.Uint32(amsHeader[28:32])
		c.logger.Debug("receive: Received packet", "invokeID", invokeID, "length", totalPacketLength, "fullPacket", fullPacket)

		c.mutex.Lock()
		ch, ok := c.requests[invokeID]
		c.mutex.Unlock()

		if ok {
			c.logger.Debug("receive: Found channel for InvokeID, sending response.", "invokeID", invokeID)
			ch <- amsData
		} else {
			c.logger.Warn("receive: No channel found for InvokeID, discarding packet.", "invokeID", invokeID)
		}
	}
}
