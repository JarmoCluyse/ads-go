// Package adsserializer provides serialization and deserialization of ADS data types.
//
// This package handles the conversion between Go values and binary ADS data according
// to TwinCAT data type specifications. It supports primitive types, structs, arrays,
// and complex nested structures.
//
// # Supported Data Types
//
// Primitive types:
//   - ADST_BIT: bool
//   - ADST_INT8, ADST_INT16, ADST_INT32, ADST_INT64: signed integers
//   - ADST_UINT8, ADST_UINT16, ADST_UINT32, ADST_UINT64: unsigned integers
//   - ADST_REAL32, ADST_REAL64: floating-point numbers
//   - ADST_STRING, ADST_WSTRING: strings (WSTRING uses UTF-16LE encoding)
//   - ADST_VOID: no data
//
// Complex types:
//   - Structs: Represented as map[string]any in Go
//   - Arrays: Represented as []any in Go (supports multidimensional)
//   - Nested combinations of the above
//
// # Usage Examples
//
// ## Deserializing Primitives
//
// Reading a simple INT32 value:
//
//	dataType := types.AdsDataType{
//	    DataType: types.ADST_INT32,
//	}
//	data := []byte{0x64, 0x00, 0x00, 0x00} // 100 in little-endian
//
//	value, err := adsserializer.Deserialize(data, dataType)
//	if err != nil {
//	    return err
//	}
//	intValue := value.(int32) // 100
//
// ## Deserializing Structs
//
// Reading a struct with multiple fields:
//
//	// TwinCAT struct:
//	// TYPE Position :
//	// STRUCT
//	//     X : INT;
//	//     Y : INT;
//	// END_STRUCT
//	// END_TYPE
//
//	dataType := types.AdsDataType{
//	    SubItems: []types.AdsDataType{
//	        {Name: "X", DataType: types.ADST_INT16, Offset: 0, Size: 2},
//	        {Name: "Y", DataType: types.ADST_INT16, Offset: 2, Size: 2},
//	    },
//	}
//	data := []byte{0x0A, 0x00, 0x14, 0x00} // X=10, Y=20
//
//	value, err := adsserializer.Deserialize(data, dataType)
//	if err != nil {
//	    return err
//	}
//	structMap := value.(map[string]any)
//	x := structMap["X"].(int16) // 10
//	y := structMap["Y"].(int16) // 20
//
// ## Deserializing Arrays
//
// Reading an array of INT32:
//
//	dataType := types.AdsDataType{
//	    DataType:  types.ADST_INT32,
//	    Size:      4,
//	    ArrayInfo: []types.AdsArrayInfo{{Length: 3}},
//	}
//	data := []byte{
//	    0x01, 0x00, 0x00, 0x00, // 1
//	    0x02, 0x00, 0x00, 0x00, // 2
//	    0x03, 0x00, 0x00, 0x00, // 3
//	}
//
//	value, err := adsserializer.Deserialize(data, dataType)
//	if err != nil {
//	    return err
//	}
//	arr := value.([]any)
//	// arr[0].(int32) == 1
//	// arr[1].(int32) == 2
//	// arr[2].(int32) == 3
//
// ## Serializing Primitives
//
// Writing a simple UINT32 value:
//
//	dataType := types.AdsDataType{
//	    DataType: types.ADST_UINT32,
//	}
//	value := uint32(12345)
//
//	data, err := adsserializer.Serialize(value, dataType)
//	if err != nil {
//	    return err
//	}
//	// data == []byte{0x39, 0x30, 0x00, 0x00}
//
// ## Serializing Structs
//
// Writing a struct:
//
//	dataType := types.AdsDataType{
//	    SubItems: []types.AdsDataType{
//	        {Name: "Field1", DataType: types.ADST_INT32, Size: 4},
//	        {Name: "Field2", DataType: types.ADST_UINT16, Size: 2},
//	    },
//	}
//	value := map[string]any{
//	    "Field1": int32(100),
//	    "Field2": uint16(42),
//	}
//
//	data, err := adsserializer.Serialize(value, dataType)
//	if err != nil {
//	    return err
//	}
//	// data == []byte{0x64, 0x00, 0x00, 0x00, 0x2A, 0x00}
//
// ## Serializing Arrays
//
// Writing an array:
//
//	dataType := types.AdsDataType{
//	    DataType:  types.ADST_INT16,
//	    Size:      2,
//	    ArrayInfo: []types.AdsArrayInfo{{Length: 3}},
//	}
//	value := []any{int16(10), int16(20), int16(30)}
//
//	data, err := adsserializer.Serialize(value, dataType)
//	if err != nil {
//	    return err
//	}
//	// data == []byte{0x0A, 0x00, 0x14, 0x00, 0x1E, 0x00}
//
// ## Multidimensional Arrays
//
// Both serialization and deserialization support multidimensional arrays:
//
//	// 2x3 array of INT16
//	dataType := types.AdsDataType{
//	    DataType:  types.ADST_INT16,
//	    Size:      2,
//	    ArrayInfo: []types.AdsArrayInfo{
//	        {Length: 2}, // 2 rows
//	        {Length: 3}, // 3 columns
//	    },
//	}
//
//	// Serialize
//	value := []any{
//	    []any{int16(1), int16(2), int16(3)},
//	    []any{int16(4), int16(5), int16(6)},
//	}
//	data, err := adsserializer.Serialize(value, dataType)
//
//	// Deserialize
//	result, err := adsserializer.Deserialize(data, dataType)
//	arr := result.([]any)
//	row1 := arr[0].([]any)
//	// row1[0].(int16) == 1
//	// row1[1].(int16) == 2
//	// row1[2].(int16) == 3
//
// # Type Conversion
//
// The package includes flexible type conversion that allows passing different
// numeric types that can be safely converted. For example:
//
//	// These all work for ADST_INT32:
//	Serialize(int32(42), dataType)
//	Serialize(int(42), dataType)    // int converted to int32
//	Serialize(int64(42), dataType)  // int64 converted to int32 (if in range)
//
// The type converters perform range checking to ensure values fit in the
// target type. If a value is out of range, serialization will fail with an error.
//
// Supported conversions:
//   - int can convert to any integer type (with range checking)
//   - Larger integer types can convert to smaller ones (with range checking)
//   - int can convert to float32/float64
//   - float64 can convert to float32
//
// # String Handling
//
// ## ADST_STRING (Single-byte strings)
//
// Standard ADS strings are null-terminated and use single-byte encoding:
//
//	dataType := types.AdsDataType{
//	    DataType: types.ADST_STRING,
//	    Size:     80, // Buffer size in bytes
//	}
//
// When serializing, the string is:
//   - Truncated if longer than Size-1 (leaving room for null terminator)
//   - Null-terminated
//   - Padded with zeros to fill the buffer
//
// When deserializing, the string is read until the null terminator.
//
// ## ADST_WSTRING (Wide strings / UTF-16LE)
//
// Wide strings use UTF-16LE encoding (2 bytes per character):
//
//	dataType := types.AdsDataType{
//	    DataType: types.ADST_WSTRING,
//	    Size:     160, // Buffer size in bytes (80 characters * 2)
//	}
//
// The Size field specifies bytes, not characters. For 80 characters, you need
// Size = 160 bytes.
//
// # Error Handling
//
// Both Serialize and Deserialize return errors in these cases:
//
//   - Unsupported data type (e.g., ADST_BIGTYPE)
//   - Type mismatch (e.g., passing string when int expected)
//   - Missing struct fields during serialization
//   - Invalid NetID format
//   - Insufficient data during deserialization
//   - Value out of range during type conversion
//
// Example error handling:
//
//	value, err := adsserializer.Deserialize(data, dataType)
//	if err != nil {
//	    if errors.Is(err, adsserializer.ErrInsufficientData) {
//	        // Handle insufficient data
//	    }
//	    return fmt.Errorf("deserialization failed: %w", err)
//	}
//
// # Binary Format
//
// All data is encoded in little-endian format to match TwinCAT ADS protocol:
//
//   - Integers: Little-endian byte order
//   - Floats: IEEE 754 little-endian
//   - Strings: Null-terminated, zero-padded
//   - WStrings: UTF-16LE, null-terminated (2-byte null), zero-padded
//   - Structs: Fields laid out sequentially according to their Offset values
//   - Arrays: Elements laid out sequentially in row-major order
//
// # Performance Considerations
//
// Serialization and deserialization are designed to be efficient:
//
//   - Uses bytes.Buffer for efficient byte slice construction
//   - Minimizes allocations where possible
//   - Primitive operations delegate to ads-primitives module
//
// For high-performance scenarios, consider:
//   - Reusing dataType structures instead of recreating them
//   - Pre-allocating buffers if possible
//   - Caching symbol and data type information
//
// # Integration with Other Modules
//
// This package is designed to work with:
//   - ads-primitives: Low-level binary I/O for primitive types
//   - ads-datatype: Data type information parsing
//   - Client code: High-level Read/WriteValue operations
//
// # Thread Safety
//
// The Serialize and Deserialize functions are thread-safe as they don't
// maintain any internal state. However, the dataType parameter is read-only
// and should not be modified during serialization/deserialization.
package adsserializer
