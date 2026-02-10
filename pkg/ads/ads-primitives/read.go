package adsprimitives

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// Sentinel errors for type checking with errors.Is()
var (
	ErrInsufficientData = errors.New("insufficient data for primitive type")
	ErrInvalidType      = errors.New("invalid data type")
)

// ReadBool reads a boolean value from binary data.
func ReadBool(data []byte) (bool, error) {
	if len(data) < 1 {
		return false, fmt.Errorf("%w: need 1 byte for bool", ErrInsufficientData)
	}
	return data[0] != 0, nil
}

// ReadInt8 reads an int8 value from binary data.
func ReadInt8(data []byte) (int8, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("%w: need 1 byte for int8", ErrInsufficientData)
	}
	return int8(data[0]), nil
}

// ReadUint8 reads a uint8 value from binary data.
func ReadUint8(data []byte) (uint8, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("%w: need 1 byte for uint8", ErrInsufficientData)
	}
	return data[0], nil
}

// ReadInt16 reads an int16 value from binary data (little-endian).
func ReadInt16(data []byte) (int16, error) {
	if len(data) < 2 {
		return 0, fmt.Errorf("%w: need 2 bytes for int16", ErrInsufficientData)
	}
	var value int16
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadUint16 reads a uint16 value from binary data (little-endian).
func ReadUint16(data []byte) (uint16, error) {
	if len(data) < 2 {
		return 0, fmt.Errorf("%w: need 2 bytes for uint16", ErrInsufficientData)
	}
	var value uint16
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadInt32 reads an int32 value from binary data (little-endian).
func ReadInt32(data []byte) (int32, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("%w: need 4 bytes for int32", ErrInsufficientData)
	}
	var value int32
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadUint32 reads a uint32 value from binary data (little-endian).
func ReadUint32(data []byte) (uint32, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("%w: need 4 bytes for uint32", ErrInsufficientData)
	}
	var value uint32
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadInt64 reads an int64 value from binary data (little-endian).
func ReadInt64(data []byte) (int64, error) {
	if len(data) < 8 {
		return 0, fmt.Errorf("%w: need 8 bytes for int64", ErrInsufficientData)
	}
	var value int64
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadUint64 reads a uint64 value from binary data (little-endian).
func ReadUint64(data []byte) (uint64, error) {
	if len(data) < 8 {
		return 0, fmt.Errorf("%w: need 8 bytes for uint64", ErrInsufficientData)
	}
	var value uint64
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadFloat32 reads a float32 value from binary data (little-endian).
func ReadFloat32(data []byte) (float32, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("%w: need 4 bytes for float32", ErrInsufficientData)
	}
	var value float32
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadFloat64 reads a float64 value from binary data (little-endian).
func ReadFloat64(data []byte) (float64, error) {
	if len(data) < 8 {
		return 0, fmt.Errorf("%w: need 8 bytes for float64", ErrInsufficientData)
	}
	var value float64
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// ReadString reads a null-terminated string from binary data.
// Returns the string without the null terminator.
func ReadString(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	// Find null terminator
	for i, b := range data {
		if b == 0 {
			return string(data[:i]), nil
		}
	}
	// No null terminator found, return entire buffer as string
	return string(data), nil
}
