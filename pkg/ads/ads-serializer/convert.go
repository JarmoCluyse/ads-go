package adsserializer

import "reflect"

// Type conversion helpers for ADS value serialization.
// These convert interface{} to the correct Go type with range checking,
// preventing panics and ensuring data integrity.

// toInt8 converts any to int8 with range checking.
// Supports: int8, int, int16, int32, int64 (with bounds checking).
func toInt8(value any) (int8, bool) {
	switch v := value.(type) {
	case int8:
		return v, true
	case int:
		if v >= -128 && v <= 127 {
			return int8(v), true
		}
	case int16, int32, int64:
		vi := int64(0)
		switch t := v.(type) {
		case int16:
			vi = int64(t)
		case int32:
			vi = int64(t)
		case int64:
			vi = t
		}
		if vi >= -128 && vi <= 127 {
			return int8(vi), true
		}
	}
	return 0, false
}

// toUint8 converts any to uint8 with range checking.
// Supports: uint8, int, uint16, uint32, uint64 (with bounds checking).
func toUint8(value any) (uint8, bool) {
	switch v := value.(type) {
	case uint8:
		return v, true
	case int:
		if v >= 0 && v <= 255 {
			return uint8(v), true
		}
	case uint16, uint32, uint64:
		vu := uint64(0)
		switch t := v.(type) {
		case uint16:
			vu = uint64(t)
		case uint32:
			vu = uint64(t)
		case uint64:
			vu = t
		}
		if vu <= 255 {
			return uint8(vu), true
		}
	}
	return 0, false
}

// toInt16 converts any to int16 with range checking.
// Supports: int16, int, int32, int64 (with bounds checking).
func toInt16(value any) (int16, bool) {
	switch v := value.(type) {
	case int16:
		return v, true
	case int:
		if v >= -32768 && v <= 32767 {
			return int16(v), true
		}
	case int32, int64:
		vi := int64(0)
		switch t := v.(type) {
		case int32:
			vi = int64(t)
		case int64:
			vi = t
		}
		if vi >= -32768 && vi <= 32767 {
			return int16(vi), true
		}
	}
	return 0, false
}

// toUint16 converts any to uint16 with range checking.
// Supports: uint16, int, uint32, uint64 (with bounds checking).
func toUint16(value any) (uint16, bool) {
	switch v := value.(type) {
	case uint16:
		return v, true
	case int:
		if v >= 0 && v <= 65535 {
			return uint16(v), true
		}
	case uint32, uint64:
		vu := uint64(0)
		switch t := v.(type) {
		case uint32:
			vu = uint64(t)
		case uint64:
			vu = t
		}
		if vu <= 65535 {
			return uint16(vu), true
		}
	}
	return 0, false
}

// toInt32 converts any to int32 with range checking.
// Supports: int32, int, int64 (with bounds checking).
func toInt32(value any) (int32, bool) {
	switch v := value.(type) {
	case int32:
		return v, true
	case int:
		if v >= -2147483648 && v <= 2147483647 {
			return int32(v), true
		}
	case int64:
		if v >= -2147483648 && v <= 2147483647 {
			return int32(v), true
		}
	}
	return 0, false
}

// toUint32 converts any to uint32 with range checking.
// Supports: uint32, int, uint64 (with bounds checking).
func toUint32(value any) (uint32, bool) {
	switch v := value.(type) {
	case uint32:
		return v, true
	case int:
		if v >= 0 && v <= 4294967295 {
			return uint32(v), true
		}
	case uint64:
		if v <= 4294967295 {
			return uint32(v), true
		}
	}
	return 0, false
}

// toInt64 converts any to int64.
// Supports: int64, int.
func toInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	}
	return 0, false
}

// toUint64 converts any to uint64.
// Supports: uint64, int (non-negative).
func toUint64(value any) (uint64, bool) {
	switch v := value.(type) {
	case uint64:
		return v, true
	case int:
		if v >= 0 {
			return uint64(v), true
		}
	}
	return 0, false
}

// toFloat32 converts any to float32.
// Supports: float32, float64, int.
func toFloat32(value any) (float32, bool) {
	switch v := value.(type) {
	case float32:
		return v, true
	case float64:
		return float32(v), true
	case int:
		return float32(v), true
	}
	return 0, false
}

// toFloat64 converts any to float64.
// Supports: float64, float32, int.
func toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	}
	return 0, false
}

// toAnySlice converts a slice of any element type (including nested slices) to []any recursively.
// This is used for handling arrays when the user provides strongly-typed slices.
//
// Example:
//
//	[]int{1, 2, 3} -> []any{1, 2, 3}
//	[][]int{{1, 2}, {3, 4}} -> []any{[]any{1, 2}, []any{3, 4}}
func toAnySlice(value any) ([]any, bool) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, false
	}
	length := v.Len()
	result := make([]any, length)
	for i := range length {
		elem := v.Index(i).Interface()
		// Recursively convert nested slices
		if reflect.ValueOf(elem).Kind() == reflect.Slice || reflect.ValueOf(elem).Kind() == reflect.Array {
			if subSlice, ok := toAnySlice(elem); ok {
				result[i] = subSlice
			} else {
				result[i] = elem
			}
		} else {
			result[i] = elem
		}
	}
	return result, true
}
