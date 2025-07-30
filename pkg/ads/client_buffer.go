package ads

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/jarmocluyse/ads-go/pkg/ads/constants"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/jarmocluyse/ads-go/pkg/ads/utils"
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
				c.logger.Error("receive: ADS error received", "invokeID", packet.InvokeId, "errorCode", packet.ErrorCode, "errorDesc", types.ADSError[packet.ErrorCode])
				ch <- Response{Error: fmt.Errorf("ADS error: %s", types.ADSError[packet.ErrorCode])}
			} else {
				ch <- Response{Data: packet.Data}
			}
		} else {
			c.logger.Warn("receive: No channel found for InvokeID, discarding packet.", "invokeID", packet.InvokeId)
		}
	}
}

// AmsTcpHeader represents the AMS/TCP header.
type AmsTcpHeader struct {
	Command types.AMSHeaderFlag // Ams commands to be send
	Length  uint32              // ams of the packet
}

// Read packet length from AMS/TCP header (bytes 2-5)
// We need to peek without advancing the buffer's read pointer
// to check if we received the full packet
func (c *Client) checkTcpPacketLength() (packetLenght uint32, error error) {
	// Check if we have enough data for AMS/TCP header (6 bytes)
	if c.receiveBuffer.Len() < constants.AMSTCPHeaderLength {
		return 0, fmt.Errorf("parseAmsTcpPacket: not enough data for AMSTCPHeaderLength")
	}
	// Read packet length from AMS/TCP header (bytes 2-5)
	// We need to peek without advancing the buffer's read pointer
	headerBytes := c.receiveBuffer.Bytes()[:constants.AMSTCPHeaderLength]
	packetLength := binary.LittleEndian.Uint32(headerBytes[2:6])
	// Total length of the full packet (AMS/TCP header + AMS header + ADS data)
	totalPacketLength := constants.AMSTCPHeaderLength + packetLength
	if packetLength < constants.AMSHeaderLength {
		return 0, fmt.Errorf("parseAmsTcpPacket: not enough data for AMSHeaderLength")
	}
	// Check if we have the full packet
	if c.receiveBuffer.Len() < int(totalPacketLength) {
		return 0, fmt.Errorf("parseAmsTcpPacket: not enough data for full packet")
	}
	c.logger.Debug("Full package in buffer", "totalPacketLength", totalPacketLength)
	return totalPacketLength, nil
}

type AmsPacket struct {
	TargetAmsAddress AmsAddress          // target address
	SourceAmsAddress AmsAddress          // source address
	AdsCommand       types.ADSCommand    // ADS command to be send
	StateFlags       types.ADSStateFlags // state flags
	DataLength       uint32              // length og the data
	ErrorCode        uint32              // AMS header error code (not ADS payload error)
	InvokeId         uint32              // invoke Id
	Data             []byte              //received data
}

// parse the ams header
// NOTE: we now at this point the length is correct
func (c *Client) parseAmsPacket(data []byte) AmsPacket {
	amsPacket := data[constants.AMSTCPHeaderLength:]
	amsHeader := amsPacket[:constants.AMSHeaderLength]
	amsData := amsPacket[constants.AMSHeaderLength:]

	targetNetID := utils.ByteArrayToAmsNetIdStr(amsHeader[0:6])
	targetPort := binary.LittleEndian.Uint16(amsHeader[6:8])
	sourceNetID := utils.ByteArrayToAmsNetIdStr(amsHeader[8:14])
	sourcePort := binary.LittleEndian.Uint16(amsHeader[14:16])
	adsCommand := types.ADSCommand(binary.LittleEndian.Uint16(amsHeader[16:18]))
	stateFlags := types.ADSStateFlags(binary.LittleEndian.Uint16(amsHeader[18:20]))
	dataLength := binary.LittleEndian.Uint32(amsHeader[20:24])
	errorCode := binary.LittleEndian.Uint32(amsHeader[24:28])
	invokeID := binary.LittleEndian.Uint32(amsHeader[28:32])

	return AmsPacket{
		TargetAmsAddress: AmsAddress{
			NetID: targetNetID,
			Port:  targetPort,
		},
		SourceAmsAddress: AmsAddress{
			NetID: sourceNetID,
			Port:  sourcePort,
		},
		AdsCommand: adsCommand,
		StateFlags: stateFlags,
		DataLength: dataLength,
		ErrorCode:  errorCode,
		InvokeId:   invokeID,
		Data:       amsData,
	}

}
