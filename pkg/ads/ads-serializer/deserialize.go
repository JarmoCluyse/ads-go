package adsserializer

import (
	"fmt"

	adsprimitives "github.com/jarmocluyse/ads-go/pkg/ads/ads-primitives"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// Deserialize converts binary data to a Go value according to the ADS data type.
//
// The function handles:
//   - Primitive types (bool, int8-64, uint8-64, float32/64, strings)
//   - Structs (returned as map[string]any)
//   - Arrays (including multidimensional arrays)
//
// Parameters:
//   - data: Binary data to deserialize
//   - dataType: ADS data type information describing the structure
//   - isArrayItem: Internal flag for recursion (should not be set by callers)
//
// Returns the deserialized value and any error encountered.
//
// Example:
//
//	// Read a simple INT32
//	value, err := adsserializer.Deserialize(data, dataType)
//	if err != nil {
//	    return err
//	}
//	intValue := value.(int32)
//
//	// Read a struct
//	value, err := adsserializer.Deserialize(data, structDataType)
//	structMap := value.(map[string]any)
//	field1 := structMap["Field1"].(int32)
func Deserialize(data []byte, dataType types.AdsDataType, isArrayItem ...bool) (any, error) {
	isArrItem := false
	if len(isArrayItem) > 0 {
		isArrItem = isArrayItem[0]
	}

	// Handle structs: if not an array item and has subitems, treat as struct
	if (len(dataType.ArrayInfo) == 0 || isArrItem) && len(dataType.SubItems) > 0 {
		result := make(map[string]any)
		data = data[dataType.Offset:]
		for _, subItem := range dataType.SubItems {
			value, err := Deserialize(data, subItem)
			if err != nil {
				return nil, err
			}
			result[subItem.Name] = value
		}
		return result, nil
	}

	// Handle arrays (including multidimensional)
	if len(dataType.ArrayInfo) > 0 && !isArrItem {
		// Track data position globally for multidimensional arrays
		dataPos := 0
		var convertArrayDimension func(dim int) []any
		convertArrayDimension = func(dim int) []any {
			temp := []any{}
			for i := 0; i < int(dataType.ArrayInfo[dim].Length); i++ {
				if dim+1 < len(dataType.ArrayInfo) {
					// Nested dimension - recurse without advancing data
					temp = append(temp, convertArrayDimension(dim+1))
				} else {
					// Final dimension - deserialize the actual element
					value, err := Deserialize(data[dataPos:], dataType, true)
					if err != nil {
						temp = append(temp, nil)
					} else {
						temp = append(temp, value)
					}
					dataPos += int(dataType.Size)
				}
			}
			return temp
		}
		return convertArrayDimension(0), nil
	}

	// Handle primitive types
	// Apply offset for primitives (used in structs where each field has an offset)
	data = data[dataType.Offset:]

	switch dataType.DataType {
	case types.ADST_VOID:
		return nil, nil
	case types.ADST_BIT:
		return adsprimitives.ReadBool(data)
	case types.ADST_INT8:
		return adsprimitives.ReadInt8(data)
	case types.ADST_UINT8:
		return adsprimitives.ReadUint8(data)
	case types.ADST_INT16:
		return adsprimitives.ReadInt16(data)
	case types.ADST_UINT16:
		return adsprimitives.ReadUint16(data)
	case types.ADST_INT32:
		return adsprimitives.ReadInt32(data)
	case types.ADST_UINT32:
		return adsprimitives.ReadUint32(data)
	case types.ADST_INT64:
		return adsprimitives.ReadInt64(data)
	case types.ADST_UINT64:
		return adsprimitives.ReadUint64(data)
	case types.ADST_REAL32:
		return adsprimitives.ReadFloat32(data)
	case types.ADST_REAL64:
		return adsprimitives.ReadFloat64(data)
	case types.ADST_STRING, types.ADST_WSTRING:
		return adsprimitives.ReadString(data)
	case types.ADST_BIGTYPE:
		return nil, fmt.Errorf("ADST_BIGTYPE is not yet supported")
	default:
		return nil, fmt.Errorf("unsupported data type: %v", dataType.DataType)
	}
}
