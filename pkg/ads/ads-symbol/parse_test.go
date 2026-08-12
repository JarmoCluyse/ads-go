package adssymbol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/stretchr/testify/assert"
)

func TestParseSymbol(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expected    AdsSymbol
		expectError bool
		errorType   error
	}{
		{
			name: "Valid - minimal symbol with empty strings",
			data: buildSymbolData(
				100,                    // dataLen
				0x1234,                 // indexGroup
				0x5678,                 // indexOffset
				4,                      // size
				uint32(types.ADST_BIT), // dataType
				0x0001,                 // flags
				0,                      // nameLength (empty name)
				0,                      // typeLength (empty type)
				0,                      // commentLength (empty comment)
				"",                     // name
				"",                     // type
				"",                     // comment
			),
			expected: AdsSymbol{
				IndexGroup:    0x1234,
				IndexOffset:   0x5678,
				Size:          4,
				DataType:      types.ADST_BIT,
				Flags:         0x0001,
				NameLength:    0,
				TypeLength:    0,
				CommentLength: 0,
				Name:          "",
				Type:          "",
				Comment:       "",
			},
			expectError: false,
		},
		{
			name: "Valid - full symbol with all fields populated",
			data: buildSymbolData(
				200,
				0xAABBCCDD,
				0x11223344,
				8,
				uint32(types.ADST_INT32),
				0x000F,
				8,  // nameLength for "MyVar123"
				6,  // typeLength for "INT32"
				12, // commentLength for "Test comment"
				"MyVar123",
				"INT32",
				"Test comment",
			),
			expected: AdsSymbol{
				IndexGroup:    0xAABBCCDD,
				IndexOffset:   0x11223344,
				Size:          8,
				DataType:      types.ADST_INT32,
				Flags:         0x000F,
				NameLength:    8,
				TypeLength:    6,
				CommentLength: 12,
				Name:          "MyVar123",
				Type:          "INT32",
				Comment:       "Test comment",
			},
			expectError: false,
		},
		{
			name: "Valid - symbol with long strings",
			data: buildSymbolData(
				500,
				0x1000,
				0x2000,
				256,
				uint32(types.ADST_STRING),
				0x0003,
				51,  // actual length of "This_Is_A_Very_Long_Variable_Name_With_Underscores"
				32,  // actual length of "ARRAY[0..10] OF STRUCT MyStruct"
				106, // actual length of the comment below
				"This_Is_A_Very_Long_Variable_Name_With_Underscores",
				"ARRAY[0..10] OF STRUCT MyStruct",
				"This is a very long comment that describes what this variable does in the PLC program in great detail",
			),
			expected: AdsSymbol{
				IndexGroup:    0x1000,
				IndexOffset:   0x2000,
				Size:          256,
				DataType:      types.ADST_STRING,
				Flags:         0x0003,
				NameLength:    51,
				TypeLength:    32,
				CommentLength: 106,
				Name:          "This_Is_A_Very_Long_Variable_Name_With_Underscores",
				Type:          "ARRAY[0..10] OF STRUCT MyStruct",
				Comment:       "This is a very long comment that describes what this variable does in the PLC program in great detail",
			},
			expectError: false,
		},
		{
			name: "Valid - symbol with comment but no name/type",
			data: buildSymbolData(
				100,
				0x3000,
				0x4000,
				16,
				uint32(types.ADST_REAL32),
				0x0001,
				0,  // no name
				0,  // no type
				14, // actual length of "Just a comment"
				"",
				"",
				"Just a comment",
			),
			expected: AdsSymbol{
				IndexGroup:    0x3000,
				IndexOffset:   0x4000,
				Size:          16,
				DataType:      types.ADST_REAL32,
				Flags:         0x0001,
				NameLength:    0,
				TypeLength:    0,
				CommentLength: 14,
				Name:          "",
				Type:          "",
				Comment:       "Just a comment",
			},
			expectError: false,
		},
		{
			name:        "Invalid - less than 30 bytes (incomplete header)",
			data:        make([]byte, 29),
			expected:    AdsSymbol{},
			expectError: true,
			errorType:   ErrInvalidSymbolLength,
		},
		{
			name:        "Invalid - empty data",
			data:        []byte{},
			expected:    AdsSymbol{},
			expectError: true,
			errorType:   ErrInvalidSymbolLength,
		},
		{
			name:        "Invalid - exactly 30 bytes but missing string data",
			data:        buildPartialSymbolData(30, 5, 5, 5), // declares strings but provides none
			expected:    AdsSymbol{},
			expectError: true,
			errorType:   ErrInsufficientData,
		},
		{
			name:        "Invalid - header present but name truncated",
			data:        buildPartialSymbolData(35, 10, 0, 0), // declares 10-char name but only 5 bytes after header
			expected:    AdsSymbol{},
			expectError: true,
			errorType:   ErrInsufficientData,
		},
		{
			name:        "Invalid - header and name present but type truncated",
			data:        buildPartialSymbolData(40, 5, 10, 0), // name fits but type doesn't
			expected:    AdsSymbol{},
			expectError: true,
			errorType:   ErrInsufficientData,
		},
		{
			name:        "Invalid - header, name, and type present but comment truncated",
			data:        buildPartialSymbolData(50, 5, 5, 20), // name and type fit but comment doesn't
			expected:    AdsSymbol{},
			expectError: true,
			errorType:   ErrInsufficientData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			symbol, err := ParseSymbol(tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expected.IndexGroup, symbol.IndexGroup)
			assert.Equal(t, tt.expected.IndexOffset, symbol.IndexOffset)
			assert.Equal(t, tt.expected.Size, symbol.Size)
			assert.Equal(t, tt.expected.DataType, symbol.DataType)
			assert.Equal(t, tt.expected.Flags, symbol.Flags)
			assert.Equal(t, tt.expected.NameLength, symbol.NameLength)
			assert.Equal(t, tt.expected.TypeLength, symbol.TypeLength)
			assert.Equal(t, tt.expected.CommentLength, symbol.CommentLength)
			assert.Equal(t, tt.expected.Name, symbol.Name)
			assert.Equal(t, tt.expected.Type, symbol.Type)
			assert.Equal(t, tt.expected.Comment, symbol.Comment)
		})
	}
}

func TestCheckSymbol(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
		errorType   error
	}{
		{
			name:        "Valid - minimal symbol",
			data:        buildSymbolData(100, 0x1234, 0x5678, 4, uint32(types.ADST_BIT), 0x0001, 0, 0, 0, "", "", ""),
			expectError: false,
		},
		{
			name:        "Valid - full symbol",
			data:        buildSymbolData(200, 0xAABB, 0xCCDD, 8, uint32(types.ADST_INT32), 0x000F, 5, 5, 10, "MyVar", "INT32", "A comment!"),
			expectError: false,
		},
		{
			name: "Valid - long strings",
			data: buildSymbolData(500, 0x1000, 0x2000, 256, uint32(types.ADST_STRING), 0x0003, 51, 32, 106,
				"This_Is_A_Very_Long_Variable_Name_With_Underscores",
				"ARRAY[0..10] OF STRUCT MyStruct",
				"This is a very long comment that describes what this variable does in the PLC program in great detail"),
			expectError: false,
		},
		{
			name:        "Invalid - less than 30 bytes",
			data:        make([]byte, 29),
			expectError: true,
			errorType:   ErrInvalidSymbolLength,
		},
		{
			name:        "Invalid - empty data",
			data:        []byte{},
			expectError: true,
			errorType:   ErrInvalidSymbolLength,
		},
		{
			name:        "Invalid - 30 bytes but strings don't fit",
			data:        buildPartialSymbolData(30, 10, 10, 10),
			expectError: true,
			errorType:   ErrInsufficientData,
		},
		{
			name:        "Invalid - name length exceeds available data",
			data:        buildPartialSymbolData(35, 20, 0, 0),
			expectError: true,
			errorType:   ErrInsufficientData,
		},
		{
			name:        "Invalid - type length exceeds available data",
			data:        buildPartialSymbolData(40, 5, 20, 0),
			expectError: true,
			errorType:   ErrInsufficientData,
		},
		{
			name:        "Invalid - comment length exceeds available data",
			data:        buildPartialSymbolData(45, 5, 5, 50),
			expectError: true,
			errorType:   ErrInsufficientData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSymbol(tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
	}
}

// Helper function to build valid symbol data for testing
func buildSymbolData(dataLen uint32, indexGroup uint32, indexOffset uint32, size uint32, dataType uint32, flags uint32, nameLen uint16, typeLen uint16, commentLen uint16, name string, typeName string, comment string) []byte {
	data := make([]byte, 30)

	// Write header (30 bytes)
	binary.LittleEndian.PutUint32(data[0:4], dataLen)
	binary.LittleEndian.PutUint32(data[4:8], indexGroup)
	binary.LittleEndian.PutUint32(data[8:12], indexOffset)
	binary.LittleEndian.PutUint32(data[12:16], size)
	binary.LittleEndian.PutUint32(data[16:20], dataType)
	binary.LittleEndian.PutUint32(data[20:24], flags)
	binary.LittleEndian.PutUint16(data[24:26], nameLen)
	binary.LittleEndian.PutUint16(data[26:28], typeLen)
	binary.LittleEndian.PutUint16(data[28:30], commentLen)

	// Append name (with null terminator)
	nameBytes := make([]byte, int(nameLen)+1)
	copy(nameBytes, []byte(name))
	data = append(data, nameBytes...)

	// Append type (with null terminator)
	typeBytes := make([]byte, int(typeLen)+1)
	copy(typeBytes, []byte(typeName))
	data = append(data, typeBytes...)

	// Append comment (with null terminator)
	commentBytes := make([]byte, int(commentLen)+1)
	copy(commentBytes, []byte(comment))
	data = append(data, commentBytes...)

	return data
}

// Helper function to build partial (invalid) symbol data for error testing
func buildPartialSymbolData(totalLen int, nameLen uint16, typeLen uint16, commentLen uint16) []byte {
	data := make([]byte, 30)

	// Write minimal header
	binary.LittleEndian.PutUint32(data[0:4], 100)                      // dataLen
	binary.LittleEndian.PutUint32(data[4:8], 0x1234)                   // indexGroup
	binary.LittleEndian.PutUint32(data[8:12], 0x5678)                  // indexOffset
	binary.LittleEndian.PutUint32(data[12:16], 4)                      // size
	binary.LittleEndian.PutUint32(data[16:20], uint32(types.ADST_BIT)) // dataType
	binary.LittleEndian.PutUint32(data[20:24], 0x0001)                 // flags
	binary.LittleEndian.PutUint16(data[24:26], nameLen)
	binary.LittleEndian.PutUint16(data[26:28], typeLen)
	binary.LittleEndian.PutUint16(data[28:30], commentLen)

	// Extend to totalLen with zeros (but don't provide enough for the declared string lengths)
	if totalLen > 30 {
		padding := make([]byte, totalLen-30)
		data = append(data, padding...)
	}

	return data
}

// buildSymbolDataWithAttrs builds symbol binary data including optional TypeGuid and attributes.
// flags should include the appropriate flag bits (TypeGuid=0x8, Attributes=0x1000).
func buildSymbolDataWithAttrs(name, typeName, comment string, flags uint32, typeGuid []byte, attrs []struct{ name, value string }) []byte {
	var buf bytes.Buffer

	nameLen := uint16(len(name))
	typeLen := uint16(len(typeName))
	commentLen := uint16(len(comment))

	// Header (30 bytes): dataLen, indexGroup, indexOffset, size, dataType, flags, nameLen, typeLen, commentLen
	binary.Write(&buf, binary.LittleEndian, uint32(0))                 // dataLen placeholder
	binary.Write(&buf, binary.LittleEndian, uint32(0x1000))            // indexGroup
	binary.Write(&buf, binary.LittleEndian, uint32(0x2000))            // indexOffset
	binary.Write(&buf, binary.LittleEndian, uint32(256))               // size
	binary.Write(&buf, binary.LittleEndian, uint32(types.ADST_STRING)) // dataType
	binary.Write(&buf, binary.LittleEndian, flags)                     // flags (low16=symbol flags, high16=arrayDimension)
	binary.Write(&buf, binary.LittleEndian, nameLen)
	binary.Write(&buf, binary.LittleEndian, typeLen)
	binary.Write(&buf, binary.LittleEndian, commentLen)

	// Name (null-terminated)
	buf.WriteString(name)
	buf.WriteByte(0)

	// Type (null-terminated)
	buf.WriteString(typeName)
	buf.WriteByte(0)

	// Comment (null-terminated)
	buf.WriteString(comment)
	buf.WriteByte(0)

	// TypeGuid (16 bytes) if TypeGuid flag set
	if flags&uint32(types.ADSSymbolFlagTypeGuid) != 0 {
		guid := typeGuid
		if len(guid) == 0 {
			guid = make([]byte, 16)
		}
		buf.Write(guid[:16])
	}

	// Attributes if Attributes flag set
	if flags&uint32(types.ADSSymbolFlagAttributes) != 0 {
		binary.Write(&buf, binary.LittleEndian, uint16(len(attrs)))
		for _, attr := range attrs {
			buf.WriteByte(byte(len(attr.name)))
			buf.WriteByte(byte(len(attr.value)))
			buf.WriteString(attr.name)
			buf.WriteByte(0) // name null terminator
			buf.WriteString(attr.value)
			buf.WriteByte(0) // value null terminator
		}
	}

	return buf.Bytes()
}

func TestParseSymbol_WithAttributes(t *testing.T) {
	t.Run("attribute with TypeGuid and no comment", func(t *testing.T) {
		data := buildSymbolDataWithAttrs(
			"MAIN.var_with_custom_attribute",
			"STRING(255)",
			"",
			uint32(types.ADSSymbolFlagTypeGuid)|uint32(types.ADSSymbolFlagAttributes), // 0x1008
			nil, // zeroed TypeGuid
			[]struct{ name, value string }{
				{"my_custom_attribute", ""},
			},
		)

		symbol, err := ParseSymbol(data)
		assert.NoError(t, err)
		assert.Equal(t, "MAIN.var_with_custom_attribute", symbol.Name)
		assert.Equal(t, "STRING(255)", symbol.Type)
		assert.Equal(t, "00000000000000000000000000000000", symbol.TypeGUID)
		assert.Len(t, symbol.Attributes, 1)
		assert.Equal(t, "my_custom_attribute", symbol.Attributes[0].Name)
		assert.Equal(t, "", symbol.Attributes[0].Value)
	})

	t.Run("attribute with known TypeGuid", func(t *testing.T) {
		guid := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
		data := buildSymbolDataWithAttrs(
			"MAIN.test_var",
			"INT",
			"",
			uint32(types.ADSSymbolFlagTypeGuid)|uint32(types.ADSSymbolFlagAttributes),
			guid,
			[]struct{ name, value string }{
				{"otelcol_role", "log_ring"},
			},
		)

		symbol, err := ParseSymbol(data)
		assert.NoError(t, err)
		assert.Equal(t, "0102030405060708090a0b0c0d0e0f10", symbol.TypeGUID)
		assert.Len(t, symbol.Attributes, 1)
		assert.Equal(t, "otelcol_role", symbol.Attributes[0].Name)
		assert.Equal(t, "log_ring", symbol.Attributes[0].Value)
	})

	t.Run("multiple attributes", func(t *testing.T) {
		data := buildSymbolDataWithAttrs(
			"GVL.multi_attr_var",
			"BOOL",
			"",
			uint32(types.ADSSymbolFlagTypeGuid)|uint32(types.ADSSymbolFlagAttributes),
			nil,
			[]struct{ name, value string }{
				{"attr_one", ""},
				{"attr_two", "value_two"},
			},
		)

		symbol, err := ParseSymbol(data)
		assert.NoError(t, err)
		assert.Len(t, symbol.Attributes, 2)
		assert.Equal(t, "attr_one", symbol.Attributes[0].Name)
		assert.Equal(t, "", symbol.Attributes[0].Value)
		assert.Equal(t, "attr_two", symbol.Attributes[1].Name)
		assert.Equal(t, "value_two", symbol.Attributes[1].Value)
	})

	t.Run("attributes without TypeGuid", func(t *testing.T) {
		data := buildSymbolDataWithAttrs(
			"MAIN.simple_attr_var",
			"INT",
			"some comment",
			uint32(types.ADSSymbolFlagAttributes), // no TypeGuid
			nil,
			[]struct{ name, value string }{
				{"my_attr", "my_val"},
			},
		)

		symbol, err := ParseSymbol(data)
		assert.NoError(t, err)
		assert.Equal(t, "some comment", symbol.Comment)
		assert.Equal(t, "", symbol.TypeGUID)
		assert.Len(t, symbol.Attributes, 1)
		assert.Equal(t, "my_attr", symbol.Attributes[0].Name)
		assert.Equal(t, "my_val", symbol.Attributes[0].Value)
	})

	t.Run("no attributes flag - attributes slice is nil", func(t *testing.T) {
		data := buildSymbolData(100, 0x1234, 0x5678, 4, uint32(types.ADST_BIT), 0x0001, 3, 3, 0, "foo", "INT", "")

		symbol, err := ParseSymbol(data)
		assert.NoError(t, err)
		assert.Nil(t, symbol.Attributes)
		assert.Equal(t, "", symbol.TypeGUID)
		assert.Nil(t, symbol.ArrayInfo)
	})
}

func TestParseSymbol_WithArrayInfo(t *testing.T) {
	t.Run("1D array symbol", func(t *testing.T) {
		// flags: high16 = arrayDimension=1, low16 = 0x0001
		flags := uint32(1<<16) | 0x0001
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, uint32(0))                // dataLen
		binary.Write(&buf, binary.LittleEndian, uint32(0x1000))           // indexGroup
		binary.Write(&buf, binary.LittleEndian, uint32(0x2000))           // indexOffset
		binary.Write(&buf, binary.LittleEndian, uint32(40))               // size
		binary.Write(&buf, binary.LittleEndian, uint32(types.ADST_INT16)) // dataType
		binary.Write(&buf, binary.LittleEndian, flags)
		binary.Write(&buf, binary.LittleEndian, uint16(7)) // nameLen "arr_var" (7 chars)
		binary.Write(&buf, binary.LittleEndian, uint16(3)) // typeLen "INT" (3 chars)
		binary.Write(&buf, binary.LittleEndian, uint16(0)) // commentLen
		buf.WriteString("arr_var\x00")                     // name (7+1 bytes)
		buf.WriteString("INT\x00")                         // type (3+1 bytes)
		buf.WriteByte(0)                                   // comment null terminator
		// Array info: startIndex=-5, length=20
		binary.Write(&buf, binary.LittleEndian, int32(-5))
		binary.Write(&buf, binary.LittleEndian, uint32(20))

		symbol, err := ParseSymbol(buf.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, "arr_var", symbol.Name)
		assert.Len(t, symbol.ArrayInfo, 1)
		assert.Equal(t, int32(-5), symbol.ArrayInfo[0].StartIndex)
		assert.Equal(t, uint32(20), symbol.ArrayInfo[0].Length)
	})
}
