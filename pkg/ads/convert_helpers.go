package ads

// NOTE: these are needed so a panic does not happen
// Type conversion helpers for ADS writing. These convert interface{} to the correct Go type for binary.Write, returning (value, ok).

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

func toInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	}
	return 0, false
}

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
