package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"strings"
)

// PlcBaseDataType holds metadata and (de)serialization logic for a PLC base type.
type PlcBaseDataType struct {
	Names       []string
	AdsDataType uint32
	Size        int
	ToBuffer    func(value any) ([]byte, error)
	FromBuffer  func(buf []byte) (any, error)
}

var baseDataTypes = []PlcBaseDataType{
	{
		Names:       []string{"STRING"},
		AdsDataType: uint32(ADST_STRING),
		Size:        81,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 81)
			str, ok := value.(string)
			if !ok {
				return nil, errors.New("STRING expects string value")
			}
			copy(buf, str)
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return string(bytes.Trim(buf, "\x00")), nil
		},
	},
	{
		Names:       []string{"WSTRING"},
		AdsDataType: uint32(ADST_WSTRING),
		Size:        162,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 162)
			str, ok := value.(string)
			if !ok {
				return nil, errors.New("WSTRING expects string value")
			}
			copy(buf, []byte(str)) // TODO: Proper UTF-16LE encoding
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return string(bytes.Trim(buf, "\x00")), nil // TODO: Proper UTF-16LE decoding
		},
	},
	{
		Names:       []string{"BOOL", "BIT", "BIT8"},
		AdsDataType: uint32(ADST_BIT),
		Size:        1,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 1)
			b, ok := value.(bool)
			if !ok {
				return nil, errors.New("BOOL expects bool value")
			}
			if b {
				buf[0] = 1
			} else {
				buf[0] = 0
			}
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return buf[0] != 0, nil
		},
	},
	{
		Names:       []string{"BYTE", "USINT", "BITARR8", "UINT8"},
		AdsDataType: uint32(ADST_UINT8),
		Size:        1,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 1)
			num, ok := value.(uint8)
			if !ok {
				return nil, errors.New("BYTE expects uint8 value")
			}
			buf[0] = num
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return uint8(buf[0]), nil
		},
	},
	{
		Names:       []string{"SINT", "INT8"},
		AdsDataType: uint32(ADST_INT8),
		Size:        1,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 1)
			num, ok := value.(int8)
			if !ok {
				return nil, errors.New("SINT expects int8 value")
			}
			buf[0] = byte(num)
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return int8(buf[0]), nil
		},
	},
	{
		Names:       []string{"UINT", "WORD", "BITARR16", "UINT16"},
		AdsDataType: uint32(ADST_UINT16),
		Size:        2,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 2)
			num, ok := value.(uint16)
			if !ok {
				return nil, errors.New("UINT expects uint16 value")
			}
			binary.LittleEndian.PutUint16(buf, num)
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return binary.LittleEndian.Uint16(buf), nil
		},
	},
	{
		Names:       []string{"INT", "INT16"},
		AdsDataType: uint32(ADST_INT16),
		Size:        2,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 2)
			num, ok := value.(int16)
			if !ok {
				return nil, errors.New("INT expects int16 value")
			}
			binary.LittleEndian.PutUint16(buf, uint16(num))
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return int16(binary.LittleEndian.Uint16(buf)), nil
		},
	},
	{
		Names:       []string{"DINT", "INT32"},
		AdsDataType: uint32(ADST_INT32),
		Size:        4,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 4)
			num, ok := value.(int32)
			if !ok {
				return nil, errors.New("DINT expects int32 value")
			}
			binary.LittleEndian.PutUint32(buf, uint32(num))
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return int32(binary.LittleEndian.Uint32(buf)), nil
		},
	},
	{
		Names:       []string{"UDINT", "DWORD", "TIME", "TIME_OF_DAY", "TOD", "BITARR32", "UINT32"},
		AdsDataType: uint32(ADST_UINT32),
		Size:        4,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 4)
			num, ok := value.(uint32)
			if !ok {
				return nil, errors.New("UDINT expects uint32 value")
			}
			binary.LittleEndian.PutUint32(buf, num)
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return binary.LittleEndian.Uint32(buf), nil
		},
	},
	{
		Names:       []string{"REAL", "FLOAT"},
		AdsDataType: uint32(ADST_REAL32),
		Size:        4,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 4)
			f, ok := value.(float32)
			if !ok {
				return nil, errors.New("REAL expects float32 value")
			}
			binary.LittleEndian.PutUint32(buf, math.Float32bits(f))
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return math.Float32frombits(binary.LittleEndian.Uint32(buf)), nil
		},
	},
	{
		Names:       []string{"LREAL", "DOUBLE"},
		AdsDataType: uint32(ADST_REAL64),
		Size:        8,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 8)
			f, ok := value.(float64)
			if !ok {
				return nil, errors.New("LREAL expects float64 value")
			}
			binary.LittleEndian.PutUint64(buf, math.Float64bits(f))
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return math.Float64frombits(binary.LittleEndian.Uint64(buf)), nil
		},
	},
	{
		Names:       []string{"LWORD", "ULINT", "LTIME", "UINT64"},
		AdsDataType: uint32(ADST_UINT64),
		Size:        8,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 8)
			num, ok := value.(uint64)
			if !ok {
				return nil, errors.New("LWORD expects uint64 value")
			}
			binary.LittleEndian.PutUint64(buf, num)
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return binary.LittleEndian.Uint64(buf), nil
		},
	},
	{
		Names:       []string{"LINT", "INT64"},
		AdsDataType: uint32(ADST_INT64),
		Size:        8,
		ToBuffer: func(value any) ([]byte, error) {
			buf := make([]byte, 8)
			num, ok := value.(int64)
			if !ok {
				return nil, errors.New("LINT expects int64 value")
			}
			binary.LittleEndian.PutUint64(buf, uint64(num))
			return buf, nil
		},
		FromBuffer: func(buf []byte) (any, error) {
			return int64(binary.LittleEndian.Uint64(buf)), nil
		},
	},
}

// IsKnownType returns true if given data type name is known
func IsKnownType(name string) bool {
	return FindType(name) != nil
}

// FindType returns pointer to PlcBaseDataType if found by name/alias (case insensitive, trims input)
func FindType(name string) *PlcBaseDataType {
	name = strings.ToUpper(strings.TrimSpace(name))
	for i := range baseDataTypes {
		typeDef := &baseDataTypes[i]
		for _, n := range typeDef.Names {
			if strings.Contains(n, name) || strings.Contains(name, n) {
				return typeDef
			}
		}
	}
	return nil
}

// ToBuffer serializes a value to []byte according to the named base data type
func ToBuffer(name string, value any) ([]byte, error) {
	typeDef := FindType(name)
	if typeDef == nil {
		return nil, errors.New("base type not found: " + name)
	}
	return typeDef.ToBuffer(value)
}

// FromBuffer deserializes a buffer to value according to the named base data type
func FromBuffer(name string, buf []byte) (any, error) {
	typeDef := FindType(name)
	if typeDef == nil {
		return nil, errors.New("base type not found: " + name)
	}
	return typeDef.FromBuffer(buf)
}

// PseudoTypeDef represents a pseudo data type mapping by size.
type PseudoTypeDef struct {
	Names            []string
	ActualTypesBySize map[int]string
}

// List of pseudo types for size-dependent mapping.
var pseudoTypes = []PseudoTypeDef{
	{
		Names: []string{"XINT", "__XINT"},
		ActualTypesBySize: map[int]string{
			4: "DINT",
			8: "LINT",
		},
	},
	{
		Names: []string{"UXINT", "__UXINT", "POINTER TO", "REFERENCE TO", "PVOID"},
		ActualTypesBySize: map[int]string{
			4: "UDINT",
			8: "ULINT",
		},
	},
	{
		Names: []string{"XWORD", "__XWORD"},
		ActualTypesBySize: map[int]string{
			4: "DWORD",
			8: "LWORD",
		},
	},
}

// FindPseudoType returns the pseudo type definition matching the name (case-insensitive).
func FindPseudoType(name string) *PseudoTypeDef {
	name = strings.ToUpper(strings.TrimSpace(name))
	for i := range pseudoTypes {
		typeDef := &pseudoTypes[i]
		for _, n := range typeDef.Names {
			if strings.HasPrefix(name, n) {
				return typeDef
			}
		}
	}
	return nil
}

// GetTypeByPseudoType returns the mapped actual type name for a pseudo type and size.
func GetTypeByPseudoType(name string, byteSize int) (string, bool) {
	if pseudo := FindPseudoType(name); pseudo != nil {
		if t, ok := pseudo.ActualTypesBySize[byteSize]; ok {
			return t, true
		}
	}
	return "", false
}
