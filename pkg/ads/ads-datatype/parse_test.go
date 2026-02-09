package adsdatatype

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/stretchr/testify/assert"
)

// buildBasicDataType creates a minimal valid data type response for testing.
func buildBasicDataType(name, typeName, comment string) []byte {
	buf := new(bytes.Buffer)

	// Calculate data length (will be set at the end)
	dataLenPos := buf.Len()
	binary.Write(buf, binary.LittleEndian, uint32(0)) // Placeholder

	// Header fields (28 bytes)
	binary.Write(buf, binary.LittleEndian, uint32(1))                // Version
	binary.Write(buf, binary.LittleEndian, uint32(1234))             // HashValue
	binary.Write(buf, binary.LittleEndian, uint32(5678))             // TypeHash
	binary.Write(buf, binary.LittleEndian, uint32(4))                // Size
	binary.Write(buf, binary.LittleEndian, uint32(0))                // Offset
	binary.Write(buf, binary.LittleEndian, uint32(types.ADST_INT32)) // DataType
	binary.Write(buf, binary.LittleEndian, uint32(0))                // Flags

	// String length fields (10 bytes)
	binary.Write(buf, binary.LittleEndian, uint16(len(name)))     // Name length
	binary.Write(buf, binary.LittleEndian, uint16(len(typeName))) // Type length
	binary.Write(buf, binary.LittleEndian, uint16(len(comment)))  // Comment length
	binary.Write(buf, binary.LittleEndian, uint16(0))             // ArrayDim
	binary.Write(buf, binary.LittleEndian, uint16(0))             // NumSubItems

	// Strings (null-terminated)
	buf.WriteString(name)
	buf.WriteByte(0)
	buf.WriteString(typeName)
	buf.WriteByte(0)
	buf.WriteString(comment)
	buf.WriteByte(0)

	// Set the data length
	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[dataLenPos:], uint32(len(data)))

	return data
}

func TestParseDataType_Basic(t *testing.T) {
	data := buildBasicDataType("myVar", "INT32", "A test variable")

	result, err := ParseDataType(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, "myVar", result.Name)
	assert.Equal(t, "INT32", result.Type)
	assert.Equal(t, "A test variable", result.Comment)
	assert.Equal(t, uint32(1), result.Version)
	assert.Equal(t, uint32(1234), result.HashValue)
	assert.Equal(t, uint32(5678), result.TypeHash)
	assert.Equal(t, uint32(4), result.Size)
	assert.Equal(t, uint32(0), result.Offset)
	assert.Equal(t, types.ADST_INT32, result.DataType)
	assert.Equal(t, types.ADSDataTypeFlags(0), result.Flags)
}

func TestParseDataType_EmptyStrings(t *testing.T) {
	data := buildBasicDataType("", "", "")

	result, err := ParseDataType(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, "", result.Name)
	assert.Equal(t, "", result.Type)
	assert.Equal(t, "", result.Comment)
}

func TestParseDataType_InsufficientData(t *testing.T) {
	data := []byte{0x01, 0x02}

	_, err := ParseDataType(data)
	if err == nil {
		t.Fatal("Expected error for insufficient data")
	}

	assert.True(t, errors.Is(err, ErrInsufficientData))
}

func TestParseDataType_WithArrayInfo(t *testing.T) {
	buf := new(bytes.Buffer)

	// Data length placeholder
	binary.Write(buf, binary.LittleEndian, uint32(0))

	// Header fields
	binary.Write(buf, binary.LittleEndian, uint32(1))                // Version
	binary.Write(buf, binary.LittleEndian, uint32(1234))             // HashValue
	binary.Write(buf, binary.LittleEndian, uint32(5678))             // TypeHash
	binary.Write(buf, binary.LittleEndian, uint32(40))               // Size (10 elements * 4 bytes)
	binary.Write(buf, binary.LittleEndian, uint32(0))                // Offset
	binary.Write(buf, binary.LittleEndian, uint32(types.ADST_INT32)) // DataType
	binary.Write(buf, binary.LittleEndian, uint32(0))                // Flags

	// String lengths
	binary.Write(buf, binary.LittleEndian, uint16(5)) // Name length
	binary.Write(buf, binary.LittleEndian, uint16(5)) // Type length (just "INT32")
	binary.Write(buf, binary.LittleEndian, uint16(0)) // Comment length
	binary.Write(buf, binary.LittleEndian, uint16(1)) // ArrayDim = 1
	binary.Write(buf, binary.LittleEndian, uint16(0)) // NumSubItems

	// Strings
	buf.WriteString("array")
	buf.WriteByte(0)
	buf.WriteString("INT32")
	buf.WriteByte(0)
	buf.WriteByte(0)

	// Array info: StartIndex=0, Length=10
	binary.Write(buf, binary.LittleEndian, int32(0))   // StartIndex
	binary.Write(buf, binary.LittleEndian, uint32(10)) // Length

	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[0:], uint32(len(data)))

	result, err := ParseDataType(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, uint16(1), result.ArrayDim)
	assert.Len(t, result.ArrayInfo, 1)
	assert.Equal(t, int32(0), result.ArrayInfo[0].StartIndex)
	assert.Equal(t, uint32(10), result.ArrayInfo[0].Length)
}

func TestParseDataType_WithMultiDimensionalArray(t *testing.T) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(1))
	binary.Write(buf, binary.LittleEndian, uint32(1234))
	binary.Write(buf, binary.LittleEndian, uint32(5678))
	binary.Write(buf, binary.LittleEndian, uint32(80))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADST_INT32))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint16(6))
	binary.Write(buf, binary.LittleEndian, uint16(5))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(2)) // 2D array
	binary.Write(buf, binary.LittleEndian, uint16(0))

	buf.WriteString("matrix")
	buf.WriteByte(0)
	buf.WriteString("INT32")
	buf.WriteByte(0)
	buf.WriteByte(0)

	// First dimension: [0..4]
	binary.Write(buf, binary.LittleEndian, int32(0))
	binary.Write(buf, binary.LittleEndian, uint32(5))
	// Second dimension: [0..3]
	binary.Write(buf, binary.LittleEndian, int32(0))
	binary.Write(buf, binary.LittleEndian, uint32(4))

	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[0:], uint32(len(data)))

	result, err := ParseDataType(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, uint16(2), result.ArrayDim)
	assert.Len(t, result.ArrayInfo, 2)
	assert.Equal(t, int32(0), result.ArrayInfo[0].StartIndex)
	assert.Equal(t, uint32(5), result.ArrayInfo[0].Length)
	assert.Equal(t, int32(0), result.ArrayInfo[1].StartIndex)
	assert.Equal(t, uint32(4), result.ArrayInfo[1].Length)
}

func TestParseDataType_WithSubItems(t *testing.T) {
	// Build a subitem first - buildBasicDataType already includes the length field
	subItemData := buildBasicDataType("field1", "INT32", "")

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(1))
	binary.Write(buf, binary.LittleEndian, uint32(1234))
	binary.Write(buf, binary.LittleEndian, uint32(5678))
	binary.Write(buf, binary.LittleEndian, uint32(4))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADST_INT32))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint16(6))
	binary.Write(buf, binary.LittleEndian, uint16(6))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(1)) // 1 subitem

	buf.WriteString("struct")
	buf.WriteByte(0)
	buf.WriteString("STRUCT")
	buf.WriteByte(0)
	buf.WriteByte(0)

	// The subitem data already includes the length at the start,
	// so we just write it directly
	buf.Write(subItemData)

	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[0:], uint32(len(data)))

	result, err := ParseDataType(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Len(t, result.SubItems, 1)
	assert.Equal(t, "field1", result.SubItems[0].Name)
	assert.Equal(t, "INT32", result.SubItems[0].Type)
}

func TestParseDataType_WithTypeGuid(t *testing.T) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(1))
	binary.Write(buf, binary.LittleEndian, uint32(1234))
	binary.Write(buf, binary.LittleEndian, uint32(5678))
	binary.Write(buf, binary.LittleEndian, uint32(4))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADST_INT32))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADSDataTypeFlagTypeGuid)) // TypeGuid flag
	binary.Write(buf, binary.LittleEndian, uint16(4))
	binary.Write(buf, binary.LittleEndian, uint16(5))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))

	buf.WriteString("test")
	buf.WriteByte(0)
	buf.WriteString("INT32")
	buf.WriteByte(0)
	buf.WriteByte(0)

	// Write 16-byte GUID
	guid := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
	buf.Write(guid)

	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[0:], uint32(len(data)))

	result, err := ParseDataType(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, "0102030405060708090a0b0c0d0e0f10", result.GUID)
}

func TestParseDataType_WithAttributes(t *testing.T) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(1))
	binary.Write(buf, binary.LittleEndian, uint32(1234))
	binary.Write(buf, binary.LittleEndian, uint32(5678))
	binary.Write(buf, binary.LittleEndian, uint32(4))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADST_INT32))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADSDataTypeFlagAttributes)) // Attributes flag
	binary.Write(buf, binary.LittleEndian, uint16(4))
	binary.Write(buf, binary.LittleEndian, uint16(5))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))

	buf.WriteString("test")
	buf.WriteByte(0)
	buf.WriteString("INT32")
	buf.WriteByte(0)
	buf.WriteByte(0)

	// Write 2 attributes
	binary.Write(buf, binary.LittleEndian, uint16(2)) // attribute count

	// Attribute 1: key="units", value="meters"
	binary.Write(buf, binary.LittleEndian, uint8(5)) // name length
	binary.Write(buf, binary.LittleEndian, uint8(6)) // value length
	buf.WriteString("units")
	buf.WriteByte(0)
	buf.WriteString("meters")
	buf.WriteByte(0)

	// Attribute 2: key="range", value="0-100"
	binary.Write(buf, binary.LittleEndian, uint8(5)) // name length
	binary.Write(buf, binary.LittleEndian, uint8(5)) // value length
	buf.WriteString("range")
	buf.WriteByte(0)
	buf.WriteString("0-100")
	buf.WriteByte(0)

	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[0:], uint32(len(data)))

	result, err := ParseDataType(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Len(t, result.Attributes, 2)
	assert.Equal(t, "units", result.Attributes[0].Name)
	assert.Equal(t, "meters", result.Attributes[0].Value)
	assert.Equal(t, "range", result.Attributes[1].Name)
	assert.Equal(t, "0-100", result.Attributes[1].Value)
}

func TestParseDataType_WithEnumInfos(t *testing.T) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(1))
	binary.Write(buf, binary.LittleEndian, uint32(1234))
	binary.Write(buf, binary.LittleEndian, uint32(5678))
	binary.Write(buf, binary.LittleEndian, uint32(4)) // Size = 4 bytes (INT32 enum)
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADST_INT32))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADSDataTypeFlagEnumInfos)) // EnumInfos flag
	binary.Write(buf, binary.LittleEndian, uint16(5))
	binary.Write(buf, binary.LittleEndian, uint16(4))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))

	buf.WriteString("State")
	buf.WriteByte(0)
	buf.WriteString("ENUM")
	buf.WriteByte(0)
	buf.WriteByte(0)

	// Write 3 enum values
	binary.Write(buf, binary.LittleEndian, uint16(3)) // enum count

	// IDLE = 0
	binary.Write(buf, binary.LittleEndian, uint8(4))
	buf.WriteString("IDLE")
	buf.WriteByte(0)
	binary.Write(buf, binary.LittleEndian, uint32(0))

	// RUNNING = 1
	binary.Write(buf, binary.LittleEndian, uint8(7))
	buf.WriteString("RUNNING")
	buf.WriteByte(0)
	binary.Write(buf, binary.LittleEndian, uint32(1))

	// STOPPED = 2
	binary.Write(buf, binary.LittleEndian, uint8(7))
	buf.WriteString("STOPPED")
	buf.WriteByte(0)
	binary.Write(buf, binary.LittleEndian, uint32(2))

	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[0:], uint32(len(data)))

	result, err := ParseDataType(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Len(t, result.EnumInfo, 3)
	assert.Equal(t, "IDLE", result.EnumInfo[0].Name)
	assert.Equal(t, int64(0), result.EnumInfo[0].Value)
	assert.Equal(t, "RUNNING", result.EnumInfo[1].Name)
	assert.Equal(t, int64(1), result.EnumInfo[1].Value)
	assert.Equal(t, "STOPPED", result.EnumInfo[2].Name)
	assert.Equal(t, int64(2), result.EnumInfo[2].Value)
}

func TestParseDataType_WithExtendedFlags(t *testing.T) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(1))
	binary.Write(buf, binary.LittleEndian, uint32(1234))
	binary.Write(buf, binary.LittleEndian, uint32(5678))
	binary.Write(buf, binary.LittleEndian, uint32(4))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADST_INT32))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADSDataTypeFlagExtendedFlags))
	binary.Write(buf, binary.LittleEndian, uint16(4))
	binary.Write(buf, binary.LittleEndian, uint16(5))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))

	buf.WriteString("test")
	buf.WriteByte(0)
	buf.WriteString("INT32")
	buf.WriteByte(0)
	buf.WriteByte(0)

	// Extended flags value
	binary.Write(buf, binary.LittleEndian, uint32(0x12345678))

	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[0:], uint32(len(data)))

	result, err := ParseDataType(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, uint32(0x12345678), result.ExtendedFlags)
}

func TestParseDataType_InvalidSubitemEntryLength(t *testing.T) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(1))
	binary.Write(buf, binary.LittleEndian, uint32(1234))
	binary.Write(buf, binary.LittleEndian, uint32(5678))
	binary.Write(buf, binary.LittleEndian, uint32(4))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADST_INT32))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint16(4))
	binary.Write(buf, binary.LittleEndian, uint16(5))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(1)) // 1 subitem

	buf.WriteString("test")
	buf.WriteByte(0)
	buf.WriteString("INT32")
	buf.WriteByte(0)
	buf.WriteByte(0)

	// Invalid entry length < 4
	binary.Write(buf, binary.LittleEndian, uint32(2))

	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[0:], uint32(len(data)))

	_, err := ParseDataType(data)
	if err == nil {
		t.Fatal("Expected error for invalid entry length")
	}

	assert.True(t, errors.Is(err, ErrInvalidEntryLength))
}

func TestCheckDataType_Valid(t *testing.T) {
	data := buildBasicDataType("test", "INT32", "")

	err := CheckDataType(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCheckDataType_InsufficientData(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}

	err := CheckDataType(data)
	if err == nil {
		t.Fatal("Expected error for insufficient data")
	}

	assert.True(t, errors.Is(err, ErrInsufficientData))
}

func TestParseDataType_TruncatedInHeader(t *testing.T) {
	data := buildBasicDataType("test", "INT32", "")
	truncated := data[:20] // Cut off in the middle of header

	_, err := ParseDataType(truncated)
	if err == nil {
		t.Fatal("Expected error for truncated header")
	}

	assert.True(t, errors.Is(err, ErrInvalidData))
}

func TestParseDataType_TruncatedInStrings(t *testing.T) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(1))
	binary.Write(buf, binary.LittleEndian, uint32(1234))
	binary.Write(buf, binary.LittleEndian, uint32(5678))
	binary.Write(buf, binary.LittleEndian, uint32(4))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(types.ADST_INT32))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint16(10)) // Say name is 10 bytes
	binary.Write(buf, binary.LittleEndian, uint16(5))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))

	// Only write 5 bytes of name instead of 11 (10 + null terminator)
	buf.WriteString("short")

	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[0:], uint32(len(data)))

	_, err := ParseDataType(data)
	if err == nil {
		t.Fatal("Expected error for truncated string")
	}

	assert.True(t, errors.Is(err, ErrInvalidData))
}
