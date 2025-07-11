package ads

import (
	"encoding/binary"

	"github.com/jarmoCluyse/ads-go/pkg/ads/constants"
	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
	"github.com/jarmoCluyse/ads-go/pkg/ads/utils"
)

// ADSReservedIndexGroups defines the reserved index groups.
type ADSReservedIndexGroups uint32

const (
	DeviceData    ADSReservedIndexGroups = 0xF100
	SymbolVersion ADSReservedIndexGroups = 0xF006
)

type AmsAddress struct {
	NetID string
	Port  uint16
}

// AmsTcpHeader represents the AMS/TCP header.
type AmsTcpHeader struct {
	Command types.AMSHeaderFlag
	Length  uint32
}

// AmsHeader represents the AMS header.
type AmsHeader struct {
	Target     AmsAddress
	Source     AmsAddress
	Command    types.ADSCommand
	StateFlags uint16
	DataLength uint32
	ErrorCode  uint32
	InvokeID   uint32
}

// createAmsTcpHeader creates the AMS/TCP header.
func createAmsTcpHeader(command types.AMSHeaderFlag, dataLength uint32) []byte {
	buf := make([]byte, constants.AMSTCPHeaderLength)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(command))
	binary.LittleEndian.PutUint32(buf[2:6], dataLength)
	return buf
}

// createAmsHeader creates the AMS header.
func createAmsHeader(target AmsAddress, source AmsAddress, command types.ADSCommand, dataLength uint32, invokeID uint32) ([]byte, error) {
	buf := make([]byte, constants.AMSHeaderLength)
	targetNetID, err := utils.AmsNetIdStrToByteArray(target.NetID)
	if err != nil {
		return nil, err
	}
	copy(buf[0:6], targetNetID)
	binary.LittleEndian.PutUint16(buf[6:8], target.Port)

	sourceNetID, err := utils.AmsNetIdStrToByteArray(source.NetID)
	if err != nil {
		return nil, err
	}
	copy(buf[8:14], sourceNetID)
	binary.LittleEndian.PutUint16(buf[14:16], source.Port)
	binary.LittleEndian.PutUint16(buf[16:18], uint16(command))
	binary.LittleEndian.PutUint16(buf[18:20], 0x0004)
	binary.LittleEndian.PutUint32(buf[20:24], dataLength)
	binary.LittleEndian.PutUint32(buf[24:28], 0) // ErrorCode
	binary.LittleEndian.PutUint32(buf[28:32], invokeID)
	return buf, nil
}
