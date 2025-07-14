package types

// AdsDataType represents an ADS data type.
type AdsDataType struct {
	Name          string
	Type          string // Added Type field
	Version       uint32
	HashValue     uint32
	TypeHash      uint32
	Size          uint32
	Offset        uint32
	DataType      ADSDataType
	Flags         ADSDataTypeFlags
	ArrayDim      uint16
	SubItems      []AdsDataType // Changed to AdsDataType
	Comment       string
	ArrayInfo     []AdsArrayInfo
	EnumInfo      []AdsEnumInfo
	Attributes    []AdsAttribute
	Methods       []AdsMethod
	GUID          string
	ExtendedFlags uint32
	CopyMask      uint32
}

// AdsDataTypeSubItem represents a sub-item within an ADS data type declaration.
type AdsDataTypeSubItem struct {
	EntryLength   uint32
	Version       uint32
	HashValue     uint32
	TypeHash      uint32
	Size          uint32
	Offset        uint32
	DataType      ADSDataType
	Flags         ADSSymbolFlags
	ArrayDim      uint16
	NameLength    uint16
	TypeLength    uint16
	CommentLength uint16
	Name          string
	Type          string
	Comment       string
}

// AdsArrayInfo represents information about an array.
type AdsArrayInfo struct {
	StartIndex int32
	Length     uint32
}

// AdsEnumInfo represents information about an enumeration.
type AdsEnumInfo struct {
	Name       string
	Value      int64
	Comment    string
	Attributes []AdsAttribute
}

// AdsAttribute represents an attribute of a data type.
type AdsAttribute struct {
	Name  string
	Value string
}

// AdsMethod represents a method of a data type.
type AdsMethod struct {
	Name       string
	ReturnType string
	Params     []AdsMethodParam
}

// AdsMethodParam represents a parameter of a method.
type AdsMethodParam struct {
	Name string
	Type string
}
