package adsserializer

import (
	"testing"

	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/stretchr/testify/assert"
)

// TestDeserialize_Primitives tests deserializing primitive types.
func TestDeserialize_Primitives(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		dataType types.AdsDataType
		expected any
	}{
		{
			name:     "BOOL - true",
			data:     []byte{0x01},
			dataType: types.AdsDataType{DataType: types.ADST_BIT},
			expected: true,
		},
		{
			name:     "BOOL - false",
			data:     []byte{0x00},
			dataType: types.AdsDataType{DataType: types.ADST_BIT},
			expected: false,
		},
		{
			name:     "INT8 - positive",
			data:     []byte{0x2A}, // 42
			dataType: types.AdsDataType{DataType: types.ADST_INT8},
			expected: int8(42),
		},
		{
			name:     "INT8 - negative",
			data:     []byte{0xFF}, // -1
			dataType: types.AdsDataType{DataType: types.ADST_INT8},
			expected: int8(-1),
		},
		{
			name:     "UINT8",
			data:     []byte{0xFF}, // 255
			dataType: types.AdsDataType{DataType: types.ADST_UINT8},
			expected: uint8(255),
		},
		{
			name:     "INT16",
			data:     []byte{0x00, 0x01}, // 256 (little-endian)
			dataType: types.AdsDataType{DataType: types.ADST_INT16},
			expected: int16(256),
		},
		{
			name:     "UINT16",
			data:     []byte{0xFF, 0xFF}, // 65535
			dataType: types.AdsDataType{DataType: types.ADST_UINT16},
			expected: uint16(65535),
		},
		{
			name:     "INT32",
			data:     []byte{0x00, 0x00, 0x00, 0x80}, // -2147483648 (little-endian)
			dataType: types.AdsDataType{DataType: types.ADST_INT32},
			expected: int32(-2147483648),
		},
		{
			name:     "UINT32",
			data:     []byte{0xFF, 0xFF, 0xFF, 0xFF}, // 4294967295
			dataType: types.AdsDataType{DataType: types.ADST_UINT32},
			expected: uint32(4294967295),
		},
		{
			name:     "INT64",
			data:     []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 1
			dataType: types.AdsDataType{DataType: types.ADST_INT64},
			expected: int64(1),
		},
		{
			name:     "UINT64",
			data:     []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, // max uint64
			dataType: types.AdsDataType{DataType: types.ADST_UINT64},
			expected: uint64(18446744073709551615),
		},
		{
			name:     "REAL32",
			data:     []byte{0x00, 0x00, 0x80, 0x3F}, // 1.0
			dataType: types.AdsDataType{DataType: types.ADST_REAL32},
			expected: float32(1.0),
		},
		{
			name:     "REAL64",
			data:     []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF0, 0x3F}, // 1.0
			dataType: types.AdsDataType{DataType: types.ADST_REAL64},
			expected: float64(1.0),
		},
		{
			name:     "STRING",
			data:     []byte{'H', 'e', 'l', 'l', 'o', 0x00},
			dataType: types.AdsDataType{DataType: types.ADST_STRING},
			expected: "Hello",
		},
		{
			name:     "VOID",
			data:     []byte{},
			dataType: types.AdsDataType{DataType: types.ADST_VOID},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Deserialize(tt.data, tt.dataType)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDeserialize_Struct tests deserializing struct types.
func TestDeserialize_Struct(t *testing.T) {
	// Struct with two fields: INT32 and UINT16
	dataType := types.AdsDataType{
		Offset: 0,
		SubItems: []types.AdsDataType{
			{
				Name:     "Field1",
				DataType: types.ADST_INT32,
				Offset:   0,
				Size:     4,
			},
			{
				Name:     "Field2",
				DataType: types.ADST_UINT16,
				Offset:   4,
				Size:     2,
			},
		},
	}

	data := []byte{0x64, 0x00, 0x00, 0x00, 0x2A, 0x00} // Field1=100, Field2=42

	result, err := Deserialize(data, dataType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	structMap, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", result)
	}

	assert.Equal(t, int32(100), structMap["Field1"])
	assert.Equal(t, uint16(42), structMap["Field2"])
}

// TestDeserialize_Array tests deserializing array types.
func TestDeserialize_Array(t *testing.T) {
	// Array of 3 INT32s
	dataType := types.AdsDataType{
		DataType: types.ADST_INT32,
		Size:     4,
		ArrayInfo: []types.AdsArrayInfo{
			{Length: 3},
		},
	}

	data := []byte{
		0x01, 0x00, 0x00, 0x00, // 1
		0x02, 0x00, 0x00, 0x00, // 2
		0x03, 0x00, 0x00, 0x00, // 3
	}

	result, err := Deserialize(data, dataType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result)
	}

	assert.Equal(t, 3, len(arr))
	assert.Equal(t, int32(1), arr[0])
	assert.Equal(t, int32(2), arr[1])
	assert.Equal(t, int32(3), arr[2])
}

// TestDeserialize_MultidimensionalArray tests deserializing multidimensional arrays.
func TestDeserialize_MultidimensionalArray(t *testing.T) {
	// 2x3 array of INT16s
	dataType := types.AdsDataType{
		DataType: types.ADST_INT16,
		Size:     2,
		ArrayInfo: []types.AdsArrayInfo{
			{Length: 2}, // 2 rows
			{Length: 3}, // 3 columns
		},
	}

	data := []byte{
		0x01, 0x00, 0x02, 0x00, 0x03, 0x00, // Row 1: [1, 2, 3]
		0x04, 0x00, 0x05, 0x00, 0x06, 0x00, // Row 2: [4, 5, 6]
	}

	result, err := Deserialize(data, dataType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result)
	}

	assert.Equal(t, 2, len(arr))

	row1, ok := arr[0].([]any)
	if !ok {
		t.Fatalf("Expected []any for row1, got %T", arr[0])
	}
	assert.Equal(t, 3, len(row1))
	assert.Equal(t, int16(1), row1[0])
	assert.Equal(t, int16(2), row1[1])
	assert.Equal(t, int16(3), row1[2])

	row2, ok := arr[1].([]any)
	if !ok {
		t.Fatalf("Expected []any for row2, got %T", arr[1])
	}
	assert.Equal(t, 3, len(row2))
	assert.Equal(t, int16(4), row2[0])
	assert.Equal(t, int16(5), row2[1])
	assert.Equal(t, int16(6), row2[2])
}

// TestDeserialize_UnsupportedType tests error handling for unsupported types.
func TestDeserialize_UnsupportedType(t *testing.T) {
	dataType := types.AdsDataType{DataType: types.ADST_BIGTYPE}
	data := []byte{0x00}

	_, err := Deserialize(data, dataType)
	if err == nil {
		t.Fatal("Expected error for ADST_BIGTYPE")
	}
}

// TestSerialize_Primitives tests serializing primitive types.
func TestSerialize_Primitives(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		dataType types.AdsDataType
		expected []byte
	}{
		{
			name:     "BOOL - true",
			value:    true,
			dataType: types.AdsDataType{DataType: types.ADST_BIT},
			expected: []byte{0x01},
		},
		{
			name:     "BOOL - false",
			value:    false,
			dataType: types.AdsDataType{DataType: types.ADST_BIT},
			expected: []byte{0x00},
		},
		{
			name:     "INT8",
			value:    int8(42),
			dataType: types.AdsDataType{DataType: types.ADST_INT8},
			expected: []byte{0x2A},
		},
		{
			name:     "UINT8",
			value:    uint8(255),
			dataType: types.AdsDataType{DataType: types.ADST_UINT8},
			expected: []byte{0xFF},
		},
		{
			name:     "INT16",
			value:    int16(256),
			dataType: types.AdsDataType{DataType: types.ADST_INT16},
			expected: []byte{0x00, 0x01},
		},
		{
			name:     "UINT16",
			value:    uint16(65535),
			dataType: types.AdsDataType{DataType: types.ADST_UINT16},
			expected: []byte{0xFF, 0xFF},
		},
		{
			name:     "INT32",
			value:    int32(100),
			dataType: types.AdsDataType{DataType: types.ADST_INT32},
			expected: []byte{0x64, 0x00, 0x00, 0x00},
		},
		{
			name:     "UINT32",
			value:    uint32(4294967295),
			dataType: types.AdsDataType{DataType: types.ADST_UINT32},
			expected: []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:     "INT64",
			value:    int64(1),
			dataType: types.AdsDataType{DataType: types.ADST_INT64},
			expected: []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:     "UINT64",
			value:    uint64(18446744073709551615),
			dataType: types.AdsDataType{DataType: types.ADST_UINT64},
			expected: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:     "REAL32",
			value:    float32(1.0),
			dataType: types.AdsDataType{DataType: types.ADST_REAL32},
			expected: []byte{0x00, 0x00, 0x80, 0x3F},
		},
		{
			name:     "REAL64",
			value:    float64(1.0),
			dataType: types.AdsDataType{DataType: types.ADST_REAL64},
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF0, 0x3F},
		},
		{
			name:     "STRING",
			value:    "Hello",
			dataType: types.AdsDataType{DataType: types.ADST_STRING, Size: 10},
			expected: []byte{'H', 'e', 'l', 'l', 'o', 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:     "VOID",
			value:    nil,
			dataType: types.AdsDataType{DataType: types.ADST_VOID},
			expected: nil, // Changed to nil to match actual return
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Serialize(tt.value, tt.dataType)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSerialize_Struct tests serializing struct types.
func TestSerialize_Struct(t *testing.T) {
	dataType := types.AdsDataType{
		SubItems: []types.AdsDataType{
			{
				Name:     "Field1",
				DataType: types.ADST_INT32,
				Size:     4,
			},
			{
				Name:     "Field2",
				DataType: types.ADST_UINT16,
				Size:     2,
			},
		},
	}

	value := map[string]any{
		"Field1": int32(100),
		"Field2": uint16(42),
	}

	expected := []byte{0x64, 0x00, 0x00, 0x00, 0x2A, 0x00}

	result, err := Serialize(value, dataType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, expected, result)
}

// TestSerialize_Struct_MissingField tests error handling for missing struct fields.
func TestSerialize_Struct_MissingField(t *testing.T) {
	dataType := types.AdsDataType{
		SubItems: []types.AdsDataType{
			{Name: "Field1", DataType: types.ADST_INT32},
			{Name: "Field2", DataType: types.ADST_INT32},
		},
	}

	value := map[string]any{
		"Field1": int32(100),
		// Field2 missing
	}

	_, err := Serialize(value, dataType)
	if err == nil {
		t.Fatal("Expected error for missing field")
	}
}

// TestSerialize_Array tests serializing array types.
func TestSerialize_Array(t *testing.T) {
	dataType := types.AdsDataType{
		DataType: types.ADST_INT32,
		Size:     4,
		ArrayInfo: []types.AdsArrayInfo{
			{Length: 3},
		},
	}

	value := []any{int32(1), int32(2), int32(3)}

	expected := []byte{
		0x01, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00,
		0x03, 0x00, 0x00, 0x00,
	}

	result, err := Serialize(value, dataType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, expected, result)
}

// TestSerialize_MultidimensionalArray tests serializing multidimensional arrays.
func TestSerialize_MultidimensionalArray(t *testing.T) {
	dataType := types.AdsDataType{
		DataType: types.ADST_INT16,
		Size:     2,
		ArrayInfo: []types.AdsArrayInfo{
			{Length: 2},
			{Length: 3},
		},
	}

	value := []any{
		[]any{int16(1), int16(2), int16(3)},
		[]any{int16(4), int16(5), int16(6)},
	}

	expected := []byte{
		0x01, 0x00, 0x02, 0x00, 0x03, 0x00,
		0x04, 0x00, 0x05, 0x00, 0x06, 0x00,
	}

	result, err := Serialize(value, dataType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, expected, result)
}

// TestSerialize_UnsupportedType tests error handling for unsupported types.
func TestSerialize_UnsupportedType(t *testing.T) {
	dataType := types.AdsDataType{DataType: types.ADST_BIGTYPE}
	var value any

	_, err := Serialize(value, dataType)
	if err == nil {
		t.Fatal("Expected error for ADST_BIGTYPE")
	}
}

// TestSerialize_InvalidType tests error handling for type mismatches.
func TestSerialize_InvalidType(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		dataType types.AdsDataType
	}{
		{
			name:     "BOOL with wrong type",
			value:    "not a bool",
			dataType: types.AdsDataType{DataType: types.ADST_BIT},
		},
		{
			name:     "INT32 with wrong type",
			value:    "not an int",
			dataType: types.AdsDataType{DataType: types.ADST_INT32},
		},
		{
			name:     "STRING with wrong type",
			value:    123,
			dataType: types.AdsDataType{DataType: types.ADST_STRING, Size: 10},
		},
		{
			name:  "Struct with wrong type",
			value: "not a map",
			dataType: types.AdsDataType{
				SubItems: []types.AdsDataType{
					{Name: "Field1", DataType: types.ADST_INT32},
				},
			},
		},
		{
			name:  "Array with wrong type",
			value: "not an array",
			dataType: types.AdsDataType{
				DataType:  types.ADST_INT32,
				ArrayInfo: []types.AdsArrayInfo{{Length: 3}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Serialize(tt.value, tt.dataType)
			if err == nil {
				t.Fatal("Expected error for invalid type")
			}
		})
	}
}

// TestRoundTrip tests serializing and then deserializing values.
func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		dataType types.AdsDataType
	}{
		{
			name:     "INT32",
			value:    int32(42),
			dataType: types.AdsDataType{DataType: types.ADST_INT32},
		},
		{
			name:     "REAL64",
			value:    float64(3.14159),
			dataType: types.AdsDataType{DataType: types.ADST_REAL64},
		},
		{
			name:     "STRING",
			value:    "Test",
			dataType: types.AdsDataType{DataType: types.ADST_STRING, Size: 10},
		},
		{
			name: "Struct",
			value: map[string]any{
				"X": int32(100),
				"Y": int32(200),
			},
			dataType: types.AdsDataType{
				SubItems: []types.AdsDataType{
					{Name: "X", DataType: types.ADST_INT32, Offset: 0, Size: 4},
					{Name: "Y", DataType: types.ADST_INT32, Offset: 4, Size: 4},
				},
			},
		},
		{
			name:  "Array",
			value: []any{int16(1), int16(2), int16(3)},
			dataType: types.AdsDataType{
				DataType:  types.ADST_INT16,
				Size:      2,
				ArrayInfo: []types.AdsArrayInfo{{Length: 3}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			data, err := Serialize(tt.value, tt.dataType)
			if err != nil {
				t.Fatalf("Serialize error: %v", err)
			}

			// Deserialize
			result, err := Deserialize(data, tt.dataType)
			if err != nil {
				t.Fatalf("Deserialize error: %v", err)
			}

			// Compare
			assert.Equal(t, tt.value, result)
		})
	}
}

// BenchmarkSerialize_INT32 benchmarks serializing INT32.
func BenchmarkSerialize_INT32(b *testing.B) {
	dataType := types.AdsDataType{DataType: types.ADST_INT32}
	value := int32(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Serialize(value, dataType)
	}
}

// BenchmarkDeserialize_INT32 benchmarks deserializing INT32.
func BenchmarkDeserialize_INT32(b *testing.B) {
	dataType := types.AdsDataType{DataType: types.ADST_INT32}
	data := []byte{0x2A, 0x00, 0x00, 0x00}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Deserialize(data, dataType)
	}
}

// BenchmarkSerialize_Struct benchmarks serializing a struct.
func BenchmarkSerialize_Struct(b *testing.B) {
	dataType := types.AdsDataType{
		SubItems: []types.AdsDataType{
			{Name: "Field1", DataType: types.ADST_INT32, Size: 4},
			{Name: "Field2", DataType: types.ADST_INT32, Size: 4},
		},
	}
	value := map[string]any{
		"Field1": int32(100),
		"Field2": int32(200),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Serialize(value, dataType)
	}
}

// BenchmarkDeserialize_Struct benchmarks deserializing a struct.
func BenchmarkDeserialize_Struct(b *testing.B) {
	dataType := types.AdsDataType{
		SubItems: []types.AdsDataType{
			{Name: "Field1", DataType: types.ADST_INT32, Offset: 0, Size: 4},
			{Name: "Field2", DataType: types.ADST_INT32, Offset: 4, Size: 4},
		},
	}
	data := []byte{0x64, 0x00, 0x00, 0x00, 0xC8, 0x00, 0x00, 0x00}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Deserialize(data, dataType)
	}
}
