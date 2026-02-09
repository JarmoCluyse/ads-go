package adssymbol

import (
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

	// Append comment (no null terminator)
	commentBytes := make([]byte, int(commentLen))
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
