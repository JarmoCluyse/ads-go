// Package adsprimitives provides functions for reading and writing primitive
// data types from/to binary ADS protocol data.
//
// This package handles the conversion between Go primitive types (bool, int8-64,
// uint8-64, float32, float64, string) and their binary representations in
// little-endian byte order as used by the ADS protocol.
//
// # Reading Primitives
//
// Read functions parse binary data into Go types:
//
//	data := []byte{0x39, 0x30} // 12345 in little-endian
//	value, err := ReadInt16(data)
//	// value = 12345
//
// All read functions validate that sufficient data is available and return
// ErrInsufficientData if the buffer is too small.
//
// # Writing Primitives
//
// Write functions convert Go types into binary data:
//
//	data, err := WriteInt16(12345)
//	// data = []byte{0x39, 0x30}
//
// Most write functions return an error for consistency, though errors are rare
// in practice (only if binary.Write fails).
//
// # Strings
//
// ReadString reads null-terminated strings. If no null terminator is found,
// the entire buffer is returned as a string.
//
// WriteString supports two modes:
//   - maxLength = 0: Appends a null terminator to the string
//   - maxLength > 0: Creates a fixed-size buffer (padded with zeros or truncated)
//
// Example:
//
//	// Null-terminated string
//	data := WriteString("Hello", 0)
//	// data = []byte{'H', 'e', 'l', 'l', 'o', 0x00}
//
//	// Fixed-length buffer (ADS STRING)
//	data := WriteString("Hi", 80)
//	// data = []byte{'H', 'i', 0x00, 0x00, ..., 0x00} (80 bytes total)
//
// # Binary Format
//
// All multi-byte values use little-endian byte order:
//   - int16/uint16: 2 bytes
//   - int32/uint32: 4 bytes
//   - int64/uint64: 8 bytes
//   - float32: 4 bytes (IEEE 754)
//   - float64: 8 bytes (IEEE 754)
//
// Single-byte values (bool, int8, uint8) are stored directly.
// Boolean values: 0x01 for true, 0x00 for false (any non-zero is read as true).
//
// # Error Handling
//
// The package defines sentinel errors that can be checked with errors.Is():
//   - ErrInsufficientData: Buffer too small for the requested type
//   - ErrInvalidType: Invalid data type (reserved for future use)
//
// Example:
//
//	value, err := ReadInt32(data)
//	if errors.Is(err, ErrInsufficientData) {
//	    // Handle buffer too small
//	}
package adsprimitives
