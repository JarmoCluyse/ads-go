package adsprimitives

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// WriteBool writes a boolean value to binary format.
// Returns 1 byte: 0x01 for true, 0x00 for false.
func WriteBool(value bool) []byte {
	if value {
		return []byte{0x01}
	}
	return []byte{0x00}
}

// WriteInt8 writes an int8 value to binary format.
// Returns 1 byte in two's complement format.
func WriteInt8(value int8) []byte {
	return []byte{byte(value)}
}

// WriteUint8 writes a uint8 value to binary format.
// Returns 1 byte.
func WriteUint8(value uint8) []byte {
	return []byte{value}
}

// WriteInt16 writes an int16 value to binary format (little-endian).
// Returns 2 bytes.
func WriteInt16(value int16) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, fmt.Errorf("failed to write int16: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteUint16 writes a uint16 value to binary format (little-endian).
// Returns 2 bytes.
func WriteUint16(value uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, fmt.Errorf("failed to write uint16: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteInt32 writes an int32 value to binary format (little-endian).
// Returns 4 bytes.
func WriteInt32(value int32) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, fmt.Errorf("failed to write int32: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteUint32 writes a uint32 value to binary format (little-endian).
// Returns 4 bytes.
func WriteUint32(value uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, fmt.Errorf("failed to write uint32: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteInt64 writes an int64 value to binary format (little-endian).
// Returns 8 bytes.
func WriteInt64(value int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, fmt.Errorf("failed to write int64: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteUint64 writes a uint64 value to binary format (little-endian).
// Returns 8 bytes.
func WriteUint64(value uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, fmt.Errorf("failed to write uint64: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteFloat32 writes a float32 value to binary format (little-endian).
// Returns 4 bytes in IEEE 754 format.
func WriteFloat32(value float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, fmt.Errorf("failed to write float32: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteFloat64 writes a float64 value to binary format (little-endian).
// Returns 8 bytes in IEEE 754 format.
func WriteFloat64(value float64) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, fmt.Errorf("failed to write float64: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteString writes a string as null-terminated binary data.
// If maxLength is 0, writes the string plus a null terminator.
// If maxLength > 0, writes exactly maxLength bytes (padded with zeros or truncated).
func WriteString(value string, maxLength int) []byte {
	if maxLength <= 0 {
		// No length limit, add null terminator
		return append([]byte(value), 0x00)
	}

	// Fixed length buffer
	buf := make([]byte, maxLength)
	copy(buf, []byte(value))
	return buf
}
