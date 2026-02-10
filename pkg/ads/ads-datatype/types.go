package adsdatatype

import "github.com/jarmocluyse/ads-go/pkg/ads/types"

// DataType represents a complete ADS data type definition.
type DataType struct {
	Name          string
	Type          string
	Version       uint32
	HashValue     uint32
	TypeHash      uint32
	Size          uint32
	Offset        uint32
	DataType      types.ADSDataType
	Flags         types.ADSDataTypeFlags
	ArrayDim      uint16
	SubItems      []DataType
	Comment       string
	ArrayInfo     []ArrayInfo
	EnumInfo      []EnumInfo
	Attributes    []Attribute
	Methods       []Method
	GUID          string
	ExtendedFlags uint32
	CopyMask      uint32
}

// ArrayInfo represents information about an array dimension.
type ArrayInfo struct {
	StartIndex int32
	Length     uint32
}

// EnumInfo represents information about an enumeration value.
type EnumInfo struct {
	Name       string
	Value      int64
	Comment    string
	Attributes []Attribute
}

// Attribute represents an attribute of a data type.
type Attribute struct {
	Name  string
	Value string
}

// Method represents a method of a data type.
type Method struct {
	Name       string
	ReturnType string
	Params     []MethodParam
}

// MethodParam represents a parameter of a method.
type MethodParam struct {
	Name string
	Type string
}
