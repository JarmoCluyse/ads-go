package adsrequests

import (
	"bytes"
	"encoding/binary"
)

// BuildReadRequest builds the payload for an ADS Read command.
//
// Binary format (12 bytes):
//
//	Bytes 0-3:  Index Group (uint32, little-endian)
//	Bytes 4-7:  Index Offset (uint32, little-endian)
//	Bytes 8-11: Read Length (uint32, little-endian)
func BuildReadRequest(indexGroup, indexOffset, readLength uint32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, indexGroup)
	binary.Write(buf, binary.LittleEndian, indexOffset)
	binary.Write(buf, binary.LittleEndian, readLength)
	return buf.Bytes()
}

// BuildWriteRequest builds the payload for an ADS Write command.
//
// Binary format (12 + N bytes):
//
//	Bytes 0-3:   Index Group (uint32, little-endian)
//	Bytes 4-7:   Index Offset (uint32, little-endian)
//	Bytes 8-11:  Data Length (uint32, little-endian)
//	Bytes 12+:   Data to write
func BuildWriteRequest(indexGroup, indexOffset uint32, data []byte) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, indexGroup)
	binary.Write(buf, binary.LittleEndian, indexOffset)
	binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	buf.Write(data)
	return buf.Bytes()
}

// BuildReadWriteRequest builds the payload for an ADS ReadWrite command.
// This command writes data and reads a response in a single transaction.
//
// Binary format (16 + N bytes):
//
//	Bytes 0-3:   Index Group (uint32, little-endian)
//	Bytes 4-7:   Index Offset (uint32, little-endian)
//	Bytes 8-11:  Read Length (uint32, little-endian)
//	Bytes 12-15: Write Length (uint32, little-endian)
//	Bytes 16+:   Data to write
//
// Note: writeData should already include any null terminators if required.
func BuildReadWriteRequest(indexGroup, indexOffset, readLength uint32, writeData []byte) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, indexGroup)
	binary.Write(buf, binary.LittleEndian, indexOffset)
	binary.Write(buf, binary.LittleEndian, readLength)
	binary.Write(buf, binary.LittleEndian, uint32(len(writeData)))
	buf.Write(writeData)
	return buf.Bytes()
}

// BuildReadWriteRequestWithNullTerminator builds an ADS ReadWrite request
// and automatically appends a null terminator to the write data.
//
// This is commonly used for symbol and data type queries where the PLC
// expects null-terminated strings.
//
// Binary format (16 + N + 1 bytes):
//
//	Bytes 0-3:     Index Group (uint32, little-endian)
//	Bytes 4-7:     Index Offset (uint32, little-endian)
//	Bytes 8-11:    Read Length (uint32, little-endian)
//	Bytes 12-15:   Write Length including null terminator (uint32, little-endian)
//	Bytes 16+:     Data to write
//	Last byte:     Null terminator (0x00)
func BuildReadWriteRequestWithNullTerminator(indexGroup, indexOffset, readLength uint32, writeData []byte) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, indexGroup)
	binary.Write(buf, binary.LittleEndian, indexOffset)
	binary.Write(buf, binary.LittleEndian, readLength)
	binary.Write(buf, binary.LittleEndian, uint32(len(writeData)+1)) // +1 for null terminator
	buf.Write(writeData)
	binary.Write(buf, binary.LittleEndian, uint8(0)) // Null terminator
	return buf.Bytes()
}

// BuildWriteControlRequest builds the payload for an ADS WriteControl command.
// This command is used to change the ADS state of a device (e.g., Run, Stop, Reset).
//
// Binary format (8 bytes):
//
//	Bytes 0-1: ADS State (uint16, little-endian)
//	Bytes 2-3: Device State (uint16, little-endian)
//	Bytes 4-7: Data Length (uint32, little-endian) - typically 0 for state changes
//
// Common ADS states:
//   - 2: Reset
//   - 5: Run
//   - 15: Config
//   - 16: Reconfig
//
// deviceState is device-specific and typically set to 0.
func BuildWriteControlRequest(adsState, deviceState uint16) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint16(data[0:2], adsState)
	binary.LittleEndian.PutUint16(data[2:4], deviceState)
	binary.LittleEndian.PutUint32(data[4:8], 0) // Data length = 0 for state change
	return data
}

// BuildReadStateRequest builds the payload for an ADS ReadState command.
// This command queries the current ADS state and device state.
//
// The payload is empty (0 bytes) - the command itself is sufficient.
func BuildReadStateRequest() []byte {
	return []byte{}
}

// BuildReadDeviceInfoRequest builds the payload for an ADS ReadDeviceInfo command.
// This command queries device information (version, name).
//
// The payload is empty (0 bytes) - the command itself is sufficient.
func BuildReadDeviceInfoRequest() []byte {
	return []byte{}
}
