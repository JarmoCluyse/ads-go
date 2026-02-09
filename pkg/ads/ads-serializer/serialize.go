package adsserializer

import (
	"bytes"
	"fmt"

	adsprimitives "github.com/jarmocluyse/ads-go/pkg/ads/ads-primitives"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// Serialize converts a Go value to binary data according to the ADS data type.
//
// The function handles:
//   - Primitive types (bool, int8-64, uint8-64, float32/64, strings)
//   - Structs (expects map[string]any)
//   - Arrays (including multidimensional arrays)
//
// Parameters:
//   - value: Go value to serialize
//   - dataType: ADS data type information describing the target structure
//   - isArrayItem: Internal flag for recursion (should not be set by callers)
//
// Returns the serialized binary data and any error encountered.
//
// Example:
//
//	// Write a simple INT32
//	data, err := adsserializer.Serialize(int32(42), dataType)
//	if err != nil {
//	    return err
//	}
//
//	// Write a struct
//	structValue := map[string]any{
//	    "Field1": int32(100),
//	    "Field2": "Hello",
//	}
//	data, err := adsserializer.Serialize(structValue, structDataType)
func Serialize(value any, dataType types.AdsDataType, isArrayItem ...bool) ([]byte, error) {
	buf := new(bytes.Buffer)

	isArrItem := false
	if len(isArrayItem) > 0 {
		isArrItem = isArrayItem[0]
	}

	// First: handle structs/subitems
	if len(dataType.SubItems) > 0 && (len(dataType.ArrayInfo) == 0) {
		// Struct type - expects map[string]any
		valMap, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid type for struct: %T (expected map[string]any)", value)
		}
		for _, subItem := range dataType.SubItems {
			subItemValue, exists := valMap[subItem.Name]
			if !exists {
				return nil, fmt.Errorf("missing field %s for struct", subItem.Name)
			}
			subItemBuf, err := Serialize(subItemValue, subItem)
			if err != nil {
				return nil, err
			}
			buf.Write(subItemBuf)
		}
		return buf.Bytes(), nil
	}

	// Second: handle arrays (including multidimensional)
	if len(dataType.ArrayInfo) > 0 && !isArrItem {
		valSlice, ok := value.([]any)
		if !ok {
			// Try to convert strongly-typed slices to []any
			valSliceConv, isSlice := toAnySlice(value)
			if !isSlice {
				return nil, fmt.Errorf("invalid type for array: %T (expected []any or slice type)", value)
			}
			valSlice = valSliceConv
		}

		var writeArray func(dim int, dType types.AdsDataType, arr []any) error
		writeArray = func(dim int, dType types.AdsDataType, arr []any) error {
			for i := 0; i < int(dType.ArrayInfo[dim].Length); i++ {
				if dim+1 < len(dType.ArrayInfo) {
					// Nested dimension
					subArr, ok := arr[i].([]any)
					if !ok {
						// Attempt to convert using toAnySlice for nested slices
						if converted, isSubSlice := toAnySlice(arr[i]); isSubSlice {
							subArr = converted
						} else {
							return fmt.Errorf("invalid nested array type: %T (expected []any)", arr[i])
						}
					}
					if err := writeArray(dim+1, dType, subArr); err != nil {
						return err
					}
				} else {
					// Final dimension - serialize the element
					elementBuf, err := Serialize(arr[i], dType, true)
					if err != nil {
						return err
					}
					buf.Write(elementBuf)
				}
			}
			return nil
		}
		if err := writeArray(0, dataType, valSlice); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	// Handle primitive types last
	switch dataType.DataType {
	case types.ADST_VOID:
		// Void type, no value to write
		return buf.Bytes(), nil

	case types.ADST_BIT:
		b, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_BIT: %T (expected bool)", value)
		}
		buf.Write(adsprimitives.WriteBool(b))

	case types.ADST_INT8:
		cast, ok := toInt8(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_INT8: %T", value)
		}
		buf.Write(adsprimitives.WriteInt8(cast))

	case types.ADST_UINT8:
		cast, ok := toUint8(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_UINT8: %T", value)
		}
		buf.Write(adsprimitives.WriteUint8(cast))

	case types.ADST_INT16:
		cast, ok := toInt16(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_INT16: %T", value)
		}
		data, err := adsprimitives.WriteInt16(cast)
		if err != nil {
			return nil, err
		}
		buf.Write(data)

	case types.ADST_UINT16:
		cast, ok := toUint16(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_UINT16: %T", value)
		}
		data, err := adsprimitives.WriteUint16(cast)
		if err != nil {
			return nil, err
		}
		buf.Write(data)

	case types.ADST_INT32:
		cast, ok := toInt32(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_INT32: %T", value)
		}
		data, err := adsprimitives.WriteInt32(cast)
		if err != nil {
			return nil, err
		}
		buf.Write(data)

	case types.ADST_UINT32:
		cast, ok := toUint32(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_UINT32: %T", value)
		}
		data, err := adsprimitives.WriteUint32(cast)
		if err != nil {
			return nil, err
		}
		buf.Write(data)

	case types.ADST_INT64:
		cast, ok := toInt64(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_INT64: %T", value)
		}
		data, err := adsprimitives.WriteInt64(cast)
		if err != nil {
			return nil, err
		}
		buf.Write(data)

	case types.ADST_UINT64:
		cast, ok := toUint64(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_UINT64: %T", value)
		}
		data, err := adsprimitives.WriteUint64(cast)
		if err != nil {
			return nil, err
		}
		buf.Write(data)

	case types.ADST_REAL32:
		cast, ok := toFloat32(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_REAL32: %T", value)
		}
		data, err := adsprimitives.WriteFloat32(cast)
		if err != nil {
			return nil, err
		}
		buf.Write(data)

	case types.ADST_REAL64:
		cast, ok := toFloat64(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_REAL64: %T", value)
		}
		data, err := adsprimitives.WriteFloat64(cast)
		if err != nil {
			return nil, err
		}
		buf.Write(data)

	case types.ADST_STRING:
		val, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_STRING: %T (expected string)", value)
		}
		bufferSize := int(dataType.Size)
		if bufferSize <= 0 {
			bufferSize = 80 // Default ADS STRING length if not specified
		}
		data := adsprimitives.WriteString(val, bufferSize)
		buf.Write(data)

	case types.ADST_WSTRING:
		val, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_WSTRING: %T (expected string)", value)
		}
		bufferSize := int(dataType.Size)
		if bufferSize <= 0 {
			bufferSize = 160 // Default WSTRING size: 80 chars * 2 bytes (UTF-16LE)
		}
		wbuf := make([]byte, bufferSize)
		// Proper UTF-16LE encoding
		runes := []rune(val)
		utf16Units := make([]uint16, len(runes))
		for i, r := range runes {
			utf16Units[i] = uint16(r)
		}
		// Write encoded runes as little-endian bytes
		byteIdx := 0
		maxChars := (bufferSize / 2) - 1 // Last two bytes are for null-terminator
		for i := 0; i < len(utf16Units) && i < maxChars; i++ {
			b0 := byte(utf16Units[i] & 0xFF)
			b1 := byte(utf16Units[i] >> 8)
			wbuf[byteIdx] = b0
			wbuf[byteIdx+1] = b1
			byteIdx += 2
		}
		// Null-terminated, wbuf is already zero-padded
		buf.Write(wbuf)

	case types.ADST_BIGTYPE:
		return nil, fmt.Errorf("ADST_BIGTYPE is not yet supported")

	default:
		return nil, fmt.Errorf("unsupported data type: %v", dataType.DataType)
	}

	return buf.Bytes(), nil
}
