package adssymbol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/jarmocluyse/ads-go/pkg/ads/utils"
)

// Sentinel errors for parsing failures
var (
	ErrInvalidSymbolLength = errors.New("invalid symbol data length")
	ErrInsufficientData    = errors.New("insufficient data for symbol parsing")
)

// ParseSymbol parses an ADS symbol from binary data.
// The data should contain the complete symbol structure starting from the redundant length field.
//
// Binary format (all fields in little-endian):
//   - Bytes 0-3:   Data Length (uint32) - redundant length field
//   - Bytes 4-7:   IndexGroup (uint32)
//   - Bytes 8-11:  IndexOffset (uint32)
//   - Bytes 12-15: Size (uint32)
//   - Bytes 16-19: DataType (uint32)
//   - Bytes 20-23: Flags (uint32)
//   - Bytes 24-25: NameLength (uint16)
//   - Bytes 26-27: TypeLength (uint16)
//   - Bytes 28-29: CommentLength (uint16)
//   - Bytes 30+:   Name (null-terminated string, length = NameLength + 1)
//   - Bytes ...:   TypeName (null-terminated string, length = TypeLength + 1)
//   - Bytes ...:   Comment (string, length = CommentLength)
//
// Returns the parsed AdsSymbol and any error encountered.
func ParseSymbol(data []byte) (AdsSymbol, error) {
	if len(data) < 30 {
		return AdsSymbol{}, fmt.Errorf("%w: expected at least 30 bytes, got %d", ErrInvalidSymbolLength, len(data))
	}

	var symbol AdsSymbol
	reader := bytes.NewReader(data)

	// Read redundant data length field (bytes 0-3)
	// Note: Beckhoff includes a redundant length field at the start of the symbol data.
	// This duplicates the length already provided in the ADS response header.
	var dataLen uint32
	if err := binary.Read(reader, binary.LittleEndian, &dataLen); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read data length: %w", err)
	}

	// Read fixed-size fields (bytes 4-29)
	if err := binary.Read(reader, binary.LittleEndian, &symbol.IndexGroup); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read IndexGroup: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.IndexOffset); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read IndexOffset: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.Size); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read Size: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.DataType); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read DataType: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.Flags); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read Flags: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.NameLength); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read NameLength: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.TypeLength); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read TypeLength: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.CommentLength); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read CommentLength: %w", err)
	}

	// Validate remaining data length for variable-length fields
	nameLen := int(symbol.NameLength) + 1   // +1 for null terminator
	typeLen := int(symbol.TypeLength) + 1   // +1 for null terminator
	commentLen := int(symbol.CommentLength) // Comment is not null-terminated
	requiredLen := nameLen + typeLen + commentLen

	if reader.Len() < requiredLen {
		return AdsSymbol{}, fmt.Errorf("%w: need %d bytes for strings, have %d", ErrInsufficientData, requiredLen, reader.Len())
	}

	// Read Name (null-terminated string)
	name := make([]byte, nameLen)
	if err := binary.Read(reader, binary.LittleEndian, &name); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read Name: %w", err)
	}
	symbol.Name = utils.DecodePlcStringBuffer(name)

	// Read TypeName (null-terminated string)
	typeName := make([]byte, typeLen)
	if err := binary.Read(reader, binary.LittleEndian, &typeName); err != nil {
		return AdsSymbol{}, fmt.Errorf("failed to read TypeName: %w", err)
	}
	symbol.Type = utils.DecodePlcStringBuffer(typeName)

	// Read Comment (not null-terminated)
	comment := make([]byte, commentLen)
	if commentLen > 0 {
		if err := binary.Read(reader, binary.LittleEndian, &comment); err != nil {
			return AdsSymbol{}, fmt.Errorf("failed to read Comment: %w", err)
		}
		symbol.Comment = utils.DecodePlcStringBuffer(comment)
	}

	return symbol, nil
}

// CheckSymbol validates an ADS symbol response without parsing it.
// This is useful for validation before passing data to ParseSymbol.
//
// Returns nil if the data appears valid, or an error describing the issue.
func CheckSymbol(data []byte) error {
	if len(data) < 30 {
		return fmt.Errorf("%w: expected at least 30 bytes, got %d", ErrInvalidSymbolLength, len(data))
	}

	// Read the length fields to validate total data length
	nameLen := binary.LittleEndian.Uint16(data[24:26])
	typeLen := binary.LittleEndian.Uint16(data[26:28])
	commentLen := binary.LittleEndian.Uint16(data[28:30])

	// Calculate required total length
	// 30 bytes (header) + nameLen + 1 (null term) + typeLen + 1 (null term) + commentLen
	requiredLen := 30 + int(nameLen) + 1 + int(typeLen) + 1 + int(commentLen)
	if len(data) < requiredLen {
		return fmt.Errorf("%w: expected at least %d bytes (30 header + %d name + %d type + %d comment), got %d",
			ErrInsufficientData, requiredLen, nameLen+1, typeLen+1, commentLen, len(data))
	}

	return nil
}
