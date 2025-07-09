package ads

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// ADSReservedIndexGroups defines the reserved index groups.
type ADSReservedIndexGroups uint32

const (
	DeviceData    ADSReservedIndexGroups = 0xF100
	SymbolVersion ADSReservedIndexGroups = 0xF006
)

// AmsNetIDStrToByteArray converts an AmsNetId string to a byte array.
func AmsNetIDStrToByteArray(s string) ([]byte, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid AmsNetId: %s", s)
	}
	bytes := make([]byte, 6)
	for i, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid part in AmsNetId: %s", part)
		}
		bytes[i] = byte(val)
	}
	return bytes, nil
}

// ByteArrayToAmsNetIDStr converts a byte array to an AmsNetId string.
func ByteArrayToAmsNetIDStr(b []byte) string {
	parts := make([]string, len(b))
	for i, byte := range b {
		parts[i] = fmt.Sprintf("%d", byte)
	}
	return strings.Join(parts, ".")
}

// AmsTcpHeader represents the AMS/TCP header.
type AmsTcpHeader struct {
	Command types.AMSHeaderFlag
	Length  uint32
}

// AmsHeader represents the AMS header.
type AmsHeader struct {
	Target     types.AmsAddress
	Source     types.AmsAddress
	Command    types.ADSCommand
	StateFlags uint16
	DataLength uint32
	ErrorCode  uint32
	InvokeID   uint32
}

// AdsReadDeviceInfoResponse represents the response of a ReadDeviceInfo command.
type AdsReadDeviceInfoResponse struct {
	ErrorCode    uint32
	MajorVersion uint8
	MinorVersion uint8
	VersionBuild uint16
	DeviceName   string
}

// AdsReadStateResponse represents the response of a ReadState command.
type AdsReadStateResponse struct {
	ErrorCode   uint32
	AdsState    types.ADSState
	DeviceState uint16
}

// AdsWriteResponse represents the response of a Write command.
type AdsWriteResponse struct {
	ErrorCode uint32
}

// AdsReadResponse represents the response of a Read command.
type AdsReadResponse struct {
	ErrorCode uint32
	Length    uint32
	Data      []byte
}

// AdsReadWriteResponse represents the response of a ReadWrite command.
type AdsReadWriteResponse struct {
	ErrorCode uint32
	Length    uint32
	Data      []byte
}

// AdsWriteControlResponse represents the response of a WriteControl command.
type AdsWriteControlResponse struct {
	ErrorCode uint32
}

// createAmsTcpHeader creates the AMS/TCP header.
func createAmsTcpHeader(command types.AMSHeaderFlag, dataLength uint32) []byte {
	buf := make([]byte, AMSTCPHeaderLength)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(command))
	binary.LittleEndian.PutUint32(buf[2:6], dataLength)
	return buf
}

// createAmsHeader creates the AMS header.
func createAmsHeader(target types.AmsAddress, source types.AmsAddress, command types.ADSCommand, dataLength uint32, invokeID uint32) ([]byte, error) {
	buf := make([]byte, AMSHeaderLength)
	targetNetID, err := AmsNetIDStrToByteArray(target.NetID)
	if err != nil {
		return nil, err
	}
	copy(buf[0:6], targetNetID)
	binary.LittleEndian.PutUint16(buf[6:8], target.Port)

	sourceNetID, err := AmsNetIDStrToByteArray(source.NetID)
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
