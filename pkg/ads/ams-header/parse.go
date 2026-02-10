package amsheader

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/jarmocluyse/ads-go/pkg/ads/utils"
)

const (
	// AMSTCPHeaderLength is the length of the AMS/TCP header (6 bytes)
	AMSTCPHeaderLength = 6
	// AMSHeaderLength is the length of the AMS header (32 bytes)
	AMSHeaderLength = 32
	// MinPacketLength is the minimum length of a complete AMS packet
	MinPacketLength = AMSTCPHeaderLength + AMSHeaderLength
)

// Sentinel errors for type checking with errors.Is()
var (
	ErrInsufficientData = errors.New("insufficient data for AMS packet")
	ErrInvalidData      = errors.New("invalid AMS packet data")
	ErrInvalidLength    = errors.New("invalid packet length in AMS/TCP header")
)

// ParsePacket parses a complete AMS packet including AMS/TCP header, AMS header, and data.
//
// Binary format:
//
//	AMS/TCP Header (6 bytes):
//	  0:2   -> Reserved (uint16)
//	  2:6   -> Packet length (uint32, excludes AMS/TCP header itself)
//
//	AMS Header (32 bytes):
//	  0:6   -> Target AMS Net ID (6 bytes)
//	  6:8   -> Target Port (uint16)
//	  8:14  -> Source AMS Net ID (6 bytes)
//	  14:16 -> Source Port (uint16)
//	  16:18 -> ADS Command (uint16)
//	  18:20 -> State Flags (uint16)
//	  20:24 -> Data Length (uint32)
//	  24:28 -> Error Code (uint32)
//	  28:32 -> Invoke ID (uint32)
//
//	Data (variable length):
//	  32:.. -> ADS data payload
//
// The input data must include the complete packet starting from the AMS/TCP header.
func ParsePacket(data []byte) (Packet, error) {
	if len(data) < MinPacketLength {
		return Packet{}, fmt.Errorf("%w: need at least %d bytes, got %d", ErrInsufficientData, MinPacketLength, len(data))
	}

	// Parse AMS/TCP header
	packetLength := binary.LittleEndian.Uint32(data[2:6])

	// Validate packet length
	if packetLength < AMSHeaderLength {
		return Packet{}, fmt.Errorf("%w: packet length %d is less than AMS header size %d", ErrInvalidLength, packetLength, AMSHeaderLength)
	}

	// Check if we have the complete packet
	totalExpectedLength := AMSTCPHeaderLength + packetLength
	if uint32(len(data)) < totalExpectedLength {
		return Packet{}, fmt.Errorf("%w: expected %d bytes, got %d", ErrInsufficientData, totalExpectedLength, len(data))
	}

	// Extract AMS packet (skip AMS/TCP header)
	amsPacket := data[AMSTCPHeaderLength:]
	amsHeader := amsPacket[:AMSHeaderLength]
	amsData := amsPacket[AMSHeaderLength:packetLength]

	// Parse AMS header fields
	targetNetID := utils.ByteArrayToAmsNetIdStr(amsHeader[0:6])
	targetPort := binary.LittleEndian.Uint16(amsHeader[6:8])
	sourceNetID := utils.ByteArrayToAmsNetIdStr(amsHeader[8:14])
	sourcePort := binary.LittleEndian.Uint16(amsHeader[14:16])
	adsCommand := types.ADSCommand(binary.LittleEndian.Uint16(amsHeader[16:18]))
	stateFlags := types.ADSStateFlags(binary.LittleEndian.Uint16(amsHeader[18:20]))
	dataLength := binary.LittleEndian.Uint32(amsHeader[20:24])
	errorCode := binary.LittleEndian.Uint32(amsHeader[24:28])
	invokeID := binary.LittleEndian.Uint32(amsHeader[28:32])

	return Packet{
		TargetNetID: targetNetID,
		TargetPort:  targetPort,
		SourceNetID: sourceNetID,
		SourcePort:  sourcePort,
		Command:     adsCommand,
		StateFlags:  stateFlags,
		DataLength:  dataLength,
		ErrorCode:   errorCode,
		InvokeID:    invokeID,
		Data:        amsData,
	}, nil
}

// CheckPacket validates that the data contains a valid AMS packet header
// without fully parsing the entire packet.
func CheckPacket(data []byte) error {
	if len(data) < MinPacketLength {
		return fmt.Errorf("%w: need at least %d bytes, got %d", ErrInsufficientData, MinPacketLength, len(data))
	}

	// Check packet length
	packetLength := binary.LittleEndian.Uint32(data[2:6])
	if packetLength < AMSHeaderLength {
		return fmt.Errorf("%w: packet length %d is less than AMS header size %d", ErrInvalidLength, packetLength, AMSHeaderLength)
	}

	totalExpectedLength := AMSTCPHeaderLength + packetLength
	if uint32(len(data)) < totalExpectedLength {
		return fmt.Errorf("%w: expected %d bytes, got %d", ErrInsufficientData, totalExpectedLength, len(data))
	}

	return nil
}

// CheckTCPPacketLength checks if enough data is available for a complete AMS packet
// and returns the total packet length. This is useful for buffering scenarios where
// data arrives in chunks.
func CheckTCPPacketLength(data []byte) (uint32, error) {
	if len(data) < AMSTCPHeaderLength {
		return 0, fmt.Errorf("%w: need at least %d bytes for AMS/TCP header, got %d", ErrInsufficientData, AMSTCPHeaderLength, len(data))
	}

	packetLength := binary.LittleEndian.Uint32(data[2:6])
	if packetLength < AMSHeaderLength {
		return 0, fmt.Errorf("%w: packet length %d is less than AMS header size %d", ErrInvalidLength, packetLength, AMSHeaderLength)
	}

	totalPacketLength := AMSTCPHeaderLength + packetLength

	if uint32(len(data)) < totalPacketLength {
		return 0, fmt.Errorf("%w: expected %d bytes, got %d", ErrInsufficientData, totalPacketLength, len(data))
	}

	return totalPacketLength, nil
}
