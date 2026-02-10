package adsprimitives

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test ReadBool
func TestReadBool(t *testing.T) {
	t.Run("true value", func(t *testing.T) {
		data := []byte{0x01}
		result, err := ReadBool(data)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("false value", func(t *testing.T) {
		data := []byte{0x00}
		result, err := ReadBool(data)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("non-zero is true", func(t *testing.T) {
		data := []byte{0xFF}
		result, err := ReadBool(data)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{}
		_, err := ReadBool(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadInt8
func TestReadInt8(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		data := []byte{0x2A}
		result, err := ReadInt8(data)
		assert.NoError(t, err)
		assert.Equal(t, int8(42), result)
	})

	t.Run("negative value", func(t *testing.T) {
		data := []byte{0xFF}
		result, err := ReadInt8(data)
		assert.NoError(t, err)
		assert.Equal(t, int8(-1), result)
	})

	t.Run("zero", func(t *testing.T) {
		data := []byte{0x00}
		result, err := ReadInt8(data)
		assert.NoError(t, err)
		assert.Equal(t, int8(0), result)
	})

	t.Run("max value", func(t *testing.T) {
		data := []byte{0x7F}
		result, err := ReadInt8(data)
		assert.NoError(t, err)
		assert.Equal(t, int8(127), result)
	})

	t.Run("min value", func(t *testing.T) {
		data := []byte{0x80}
		result, err := ReadInt8(data)
		assert.NoError(t, err)
		assert.Equal(t, int8(-128), result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{}
		_, err := ReadInt8(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadUint8
func TestReadUint8(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		data := []byte{0x2A}
		result, err := ReadUint8(data)
		assert.NoError(t, err)
		assert.Equal(t, uint8(42), result)
	})

	t.Run("max value", func(t *testing.T) {
		data := []byte{0xFF}
		result, err := ReadUint8(data)
		assert.NoError(t, err)
		assert.Equal(t, uint8(255), result)
	})

	t.Run("zero", func(t *testing.T) {
		data := []byte{0x00}
		result, err := ReadUint8(data)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0), result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{}
		_, err := ReadUint8(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadInt16
func TestReadInt16(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		data := []byte{0x39, 0x30} // 12345 in little-endian
		result, err := ReadInt16(data)
		assert.NoError(t, err)
		assert.Equal(t, int16(12345), result)
	})

	t.Run("negative value", func(t *testing.T) {
		data := []byte{0xFF, 0xFF} // -1 in little-endian
		result, err := ReadInt16(data)
		assert.NoError(t, err)
		assert.Equal(t, int16(-1), result)
	})

	t.Run("zero", func(t *testing.T) {
		data := []byte{0x00, 0x00}
		result, err := ReadInt16(data)
		assert.NoError(t, err)
		assert.Equal(t, int16(0), result)
	})

	t.Run("max value", func(t *testing.T) {
		data := []byte{0xFF, 0x7F} // 32767
		result, err := ReadInt16(data)
		assert.NoError(t, err)
		assert.Equal(t, int16(32767), result)
	})

	t.Run("min value", func(t *testing.T) {
		data := []byte{0x00, 0x80} // -32768
		result, err := ReadInt16(data)
		assert.NoError(t, err)
		assert.Equal(t, int16(-32768), result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{0x01}
		_, err := ReadInt16(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadUint16
func TestReadUint16(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		data := []byte{0x39, 0x30} // 12345 in little-endian
		result, err := ReadUint16(data)
		assert.NoError(t, err)
		assert.Equal(t, uint16(12345), result)
	})

	t.Run("max value", func(t *testing.T) {
		data := []byte{0xFF, 0xFF} // 65535
		result, err := ReadUint16(data)
		assert.NoError(t, err)
		assert.Equal(t, uint16(65535), result)
	})

	t.Run("zero", func(t *testing.T) {
		data := []byte{0x00, 0x00}
		result, err := ReadUint16(data)
		assert.NoError(t, err)
		assert.Equal(t, uint16(0), result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{0x01}
		_, err := ReadUint16(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadInt32
func TestReadInt32(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		data := []byte{0x15, 0xCD, 0x5B, 0x07} // 123456789 in little-endian
		result, err := ReadInt32(data)
		assert.NoError(t, err)
		assert.Equal(t, int32(123456789), result)
	})

	t.Run("negative value", func(t *testing.T) {
		data := []byte{0xFF, 0xFF, 0xFF, 0xFF} // -1
		result, err := ReadInt32(data)
		assert.NoError(t, err)
		assert.Equal(t, int32(-1), result)
	})

	t.Run("zero", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00}
		result, err := ReadInt32(data)
		assert.NoError(t, err)
		assert.Equal(t, int32(0), result)
	})

	t.Run("max value", func(t *testing.T) {
		data := []byte{0xFF, 0xFF, 0xFF, 0x7F} // 2147483647
		result, err := ReadInt32(data)
		assert.NoError(t, err)
		assert.Equal(t, int32(2147483647), result)
	})

	t.Run("min value", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x80} // -2147483648
		result, err := ReadInt32(data)
		assert.NoError(t, err)
		assert.Equal(t, int32(-2147483648), result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03}
		_, err := ReadInt32(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadUint32
func TestReadUint32(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		data := []byte{0x15, 0xCD, 0x5B, 0x07} // 123456789 in little-endian
		result, err := ReadUint32(data)
		assert.NoError(t, err)
		assert.Equal(t, uint32(123456789), result)
	})

	t.Run("max value", func(t *testing.T) {
		data := []byte{0xFF, 0xFF, 0xFF, 0xFF} // 4294967295
		result, err := ReadUint32(data)
		assert.NoError(t, err)
		assert.Equal(t, uint32(4294967295), result)
	})

	t.Run("zero", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00}
		result, err := ReadUint32(data)
		assert.NoError(t, err)
		assert.Equal(t, uint32(0), result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03}
		_, err := ReadUint32(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadInt64
func TestReadInt64(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		data := []byte{0x15, 0x81, 0xE9, 0x7D, 0xF4, 0x10, 0x22, 0x11} // 1234567890123456789
		result, err := ReadInt64(data)
		assert.NoError(t, err)
		assert.Equal(t, int64(1234567890123456789), result)
	})

	t.Run("negative value", func(t *testing.T) {
		data := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} // -1
		result, err := ReadInt64(data)
		assert.NoError(t, err)
		assert.Equal(t, int64(-1), result)
	})

	t.Run("zero", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		result, err := ReadInt64(data)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), result)
	})

	t.Run("max value", func(t *testing.T) {
		data := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F}
		result, err := ReadInt64(data)
		assert.NoError(t, err)
		assert.Equal(t, int64(9223372036854775807), result)
	})

	t.Run("min value", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}
		result, err := ReadInt64(data)
		assert.NoError(t, err)
		assert.Equal(t, int64(-9223372036854775808), result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
		_, err := ReadInt64(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadUint64
func TestReadUint64(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		data := []byte{0x15, 0x81, 0xE9, 0x7D, 0xF4, 0x10, 0x22, 0x11} // 1234567890123456789
		result, err := ReadUint64(data)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1234567890123456789), result)
	})

	t.Run("max value", func(t *testing.T) {
		data := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
		result, err := ReadUint64(data)
		assert.NoError(t, err)
		assert.Equal(t, uint64(18446744073709551615), result)
	})

	t.Run("zero", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		result, err := ReadUint64(data)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
		_, err := ReadUint64(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadFloat32
func TestReadFloat32(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		data := []byte{0xD0, 0x0F, 0x49, 0x40} // 3.14159 in IEEE 754
		result, err := ReadFloat32(data)
		assert.NoError(t, err)
		assert.InDelta(t, float32(3.14159), result, 0.00001)
	})

	t.Run("negative value", func(t *testing.T) {
		data := []byte{0xD0, 0x0F, 0x49, 0xC0} // -3.14159 in IEEE 754
		result, err := ReadFloat32(data)
		assert.NoError(t, err)
		assert.InDelta(t, float32(-3.14159), result, 0.00001)
	})

	t.Run("zero", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00}
		result, err := ReadFloat32(data)
		assert.NoError(t, err)
		assert.Equal(t, float32(0), result)
	})

	t.Run("max value", func(t *testing.T) {
		data := []byte{0xFF, 0xFF, 0x7F, 0x7F}
		result, err := ReadFloat32(data)
		assert.NoError(t, err)
		assert.Equal(t, float32(math.MaxFloat32), result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03}
		_, err := ReadFloat32(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadFloat64
func TestReadFloat64(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		data := []byte{0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0x40} // 3.141592653589793
		result, err := ReadFloat64(data)
		assert.NoError(t, err)
		assert.InDelta(t, 3.141592653589793, result, 0.00000000000001)
	})

	t.Run("negative value", func(t *testing.T) {
		data := []byte{0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0xC0} // -3.141592653589793
		result, err := ReadFloat64(data)
		assert.NoError(t, err)
		assert.InDelta(t, -3.141592653589793, result, 0.00000000000001)
	})

	t.Run("zero", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		result, err := ReadFloat64(data)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), result)
	})

	t.Run("max value", func(t *testing.T) {
		data := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xEF, 0x7F}
		result, err := ReadFloat64(data)
		assert.NoError(t, err)
		assert.Equal(t, math.MaxFloat64, result)
	})

	t.Run("insufficient data", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
		_, err := ReadFloat64(data)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInsufficientData))
	})
}

// Test ReadString
func TestReadString(t *testing.T) {
	t.Run("null-terminated string", func(t *testing.T) {
		data := []byte{'H', 'e', 'l', 'l', 'o', 0x00}
		result, err := ReadString(data)
		assert.NoError(t, err)
		assert.Equal(t, "Hello", result)
	})

	t.Run("string with trailing data", func(t *testing.T) {
		data := []byte{'H', 'i', 0x00, 'X', 'X', 'X'}
		result, err := ReadString(data)
		assert.NoError(t, err)
		assert.Equal(t, "Hi", result)
	})

	t.Run("string without null terminator", func(t *testing.T) {
		data := []byte{'T', 'e', 's', 't'}
		result, err := ReadString(data)
		assert.NoError(t, err)
		assert.Equal(t, "Test", result)
	})

	t.Run("empty string with null", func(t *testing.T) {
		data := []byte{0x00}
		result, err := ReadString(data)
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("empty data", func(t *testing.T) {
		data := []byte{}
		result, err := ReadString(data)
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})
}

// Test WriteBool
func TestWriteBool(t *testing.T) {
	t.Run("true value", func(t *testing.T) {
		result := WriteBool(true)
		assert.Equal(t, []byte{0x01}, result)
	})

	t.Run("false value", func(t *testing.T) {
		result := WriteBool(false)
		assert.Equal(t, []byte{0x00}, result)
	})
}

// Test WriteInt8
func TestWriteInt8(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		result := WriteInt8(42)
		assert.Equal(t, []byte{0x2A}, result)
	})

	t.Run("negative value", func(t *testing.T) {
		result := WriteInt8(-1)
		assert.Equal(t, []byte{0xFF}, result)
	})

	t.Run("zero", func(t *testing.T) {
		result := WriteInt8(0)
		assert.Equal(t, []byte{0x00}, result)
	})

	t.Run("max value", func(t *testing.T) {
		result := WriteInt8(127)
		assert.Equal(t, []byte{0x7F}, result)
	})

	t.Run("min value", func(t *testing.T) {
		result := WriteInt8(-128)
		assert.Equal(t, []byte{0x80}, result)
	})
}

// Test WriteUint8
func TestWriteUint8(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		result := WriteUint8(42)
		assert.Equal(t, []byte{0x2A}, result)
	})

	t.Run("max value", func(t *testing.T) {
		result := WriteUint8(255)
		assert.Equal(t, []byte{0xFF}, result)
	})

	t.Run("zero", func(t *testing.T) {
		result := WriteUint8(0)
		assert.Equal(t, []byte{0x00}, result)
	})
}

// Test WriteInt16
func TestWriteInt16(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		result, err := WriteInt16(12345)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x39, 0x30}, result)
	})

	t.Run("negative value", func(t *testing.T) {
		result, err := WriteInt16(-1)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0xFF}, result)
	})

	t.Run("zero", func(t *testing.T) {
		result, err := WriteInt16(0)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00}, result)
	})

	t.Run("max value", func(t *testing.T) {
		result, err := WriteInt16(32767)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0x7F}, result)
	})

	t.Run("min value", func(t *testing.T) {
		result, err := WriteInt16(-32768)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x80}, result)
	})
}

// Test WriteUint16
func TestWriteUint16(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		result, err := WriteUint16(12345)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x39, 0x30}, result)
	})

	t.Run("max value", func(t *testing.T) {
		result, err := WriteUint16(65535)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0xFF}, result)
	})

	t.Run("zero", func(t *testing.T) {
		result, err := WriteUint16(0)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00}, result)
	})
}

// Test WriteInt32
func TestWriteInt32(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		result, err := WriteInt32(123456789)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x15, 0xCD, 0x5B, 0x07}, result)
	})

	t.Run("negative value", func(t *testing.T) {
		result, err := WriteInt32(-1)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF}, result)
	})

	t.Run("zero", func(t *testing.T) {
		result, err := WriteInt32(0)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, result)
	})

	t.Run("max value", func(t *testing.T) {
		result, err := WriteInt32(2147483647)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0x7F}, result)
	})

	t.Run("min value", func(t *testing.T) {
		result, err := WriteInt32(-2147483648)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x80}, result)
	})
}

// Test WriteUint32
func TestWriteUint32(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		result, err := WriteUint32(123456789)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x15, 0xCD, 0x5B, 0x07}, result)
	})

	t.Run("max value", func(t *testing.T) {
		result, err := WriteUint32(4294967295)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF}, result)
	})

	t.Run("zero", func(t *testing.T) {
		result, err := WriteUint32(0)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, result)
	})
}

// Test WriteInt64
func TestWriteInt64(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		result, err := WriteInt64(1234567890123456789)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x15, 0x81, 0xE9, 0x7D, 0xF4, 0x10, 0x22, 0x11}, result)
	})

	t.Run("negative value", func(t *testing.T) {
		result, err := WriteInt64(-1)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, result)
	})

	t.Run("zero", func(t *testing.T) {
		result, err := WriteInt64(0)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, result)
	})

	t.Run("max value", func(t *testing.T) {
		result, err := WriteInt64(9223372036854775807)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F}, result)
	})

	t.Run("min value", func(t *testing.T) {
		result, err := WriteInt64(-9223372036854775808)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, result)
	})
}

// Test WriteUint64
func TestWriteUint64(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		result, err := WriteUint64(1234567890123456789)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x15, 0x81, 0xE9, 0x7D, 0xF4, 0x10, 0x22, 0x11}, result)
	})

	t.Run("max value", func(t *testing.T) {
		result, err := WriteUint64(18446744073709551615)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, result)
	})

	t.Run("zero", func(t *testing.T) {
		result, err := WriteUint64(0)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, result)
	})
}

// Test WriteFloat32
func TestWriteFloat32(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		result, err := WriteFloat32(3.14159)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xD0, 0x0F, 0x49, 0x40}, result)
	})

	t.Run("negative value", func(t *testing.T) {
		result, err := WriteFloat32(-3.14159)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xD0, 0x0F, 0x49, 0xC0}, result)
	})

	t.Run("zero", func(t *testing.T) {
		result, err := WriteFloat32(0)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, result)
	})

	t.Run("max value", func(t *testing.T) {
		result, err := WriteFloat32(math.MaxFloat32)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0xFF, 0x7F, 0x7F}, result)
	})
}

// Test WriteFloat64
func TestWriteFloat64(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		result, err := WriteFloat64(3.141592653589793)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0x40}, result)
	})

	t.Run("negative value", func(t *testing.T) {
		result, err := WriteFloat64(-3.141592653589793)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0xC0}, result)
	})

	t.Run("zero", func(t *testing.T) {
		result, err := WriteFloat64(0)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, result)
	})

	t.Run("max value", func(t *testing.T) {
		result, err := WriteFloat64(math.MaxFloat64)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xEF, 0x7F}, result)
	})
}

// Test WriteString
func TestWriteString(t *testing.T) {
	t.Run("simple string with null terminator", func(t *testing.T) {
		result := WriteString("Hello", 0)
		assert.Equal(t, []byte{'H', 'e', 'l', 'l', 'o', 0x00}, result)
	})

	t.Run("empty string with null terminator", func(t *testing.T) {
		result := WriteString("", 0)
		assert.Equal(t, []byte{0x00}, result)
	})

	t.Run("fixed length buffer exact fit", func(t *testing.T) {
		result := WriteString("Test", 4)
		assert.Equal(t, []byte{'T', 'e', 's', 't'}, result)
	})

	t.Run("fixed length buffer with padding", func(t *testing.T) {
		result := WriteString("Hi", 5)
		assert.Equal(t, []byte{'H', 'i', 0x00, 0x00, 0x00}, result)
	})

	t.Run("fixed length buffer truncates", func(t *testing.T) {
		result := WriteString("TooLong", 4)
		assert.Equal(t, []byte{'T', 'o', 'o', 'L'}, result)
	})

	t.Run("fixed length buffer zero", func(t *testing.T) {
		result := WriteString("Test", 0)
		assert.Equal(t, []byte{'T', 'e', 's', 't', 0x00}, result)
	})
}

// Test round-trip consistency
func TestRoundTrip(t *testing.T) {
	t.Run("int16 round trip", func(t *testing.T) {
		original := int16(-12345)
		written, err := WriteInt16(original)
		assert.NoError(t, err)

		read, err := ReadInt16(written)
		assert.NoError(t, err)
		assert.Equal(t, original, read)
	})

	t.Run("uint32 round trip", func(t *testing.T) {
		original := uint32(987654321)
		written, err := WriteUint32(original)
		assert.NoError(t, err)

		read, err := ReadUint32(written)
		assert.NoError(t, err)
		assert.Equal(t, original, read)
	})

	t.Run("float32 round trip", func(t *testing.T) {
		original := float32(3.14159)
		written, err := WriteFloat32(original)
		assert.NoError(t, err)

		read, err := ReadFloat32(written)
		assert.NoError(t, err)
		assert.Equal(t, original, read)
	})

	t.Run("float64 round trip", func(t *testing.T) {
		original := float64(3.141592653589793)
		written, err := WriteFloat64(original)
		assert.NoError(t, err)

		read, err := ReadFloat64(written)
		assert.NoError(t, err)
		assert.Equal(t, original, read)
	})

	t.Run("string round trip", func(t *testing.T) {
		original := "Test String"
		written := WriteString(original, 0)

		read, err := ReadString(written)
		assert.NoError(t, err)
		assert.Equal(t, original, read)
	})

	t.Run("bool round trip", func(t *testing.T) {
		original := true
		written := WriteBool(original)

		read, err := ReadBool(written)
		assert.NoError(t, err)
		assert.Equal(t, original, read)
	})
}
