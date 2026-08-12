package adssymbol

import "github.com/jarmocluyse/ads-go/pkg/ads/types"

// SymbolAttribute represents a single TwinCAT pragma attribute attached to a
// PLC variable, e.g. {attribute 'otelcol_role' := 'log_ring'}.
type SymbolAttribute struct {
	Name  string
	Value string
}

// ArrayInfoEntry describes one dimension of an array type.
type ArrayInfoEntry struct {
	StartIndex int32  // StartIndex is the lower bound of the array dimension
	Length     uint32 // Length is the number of elements in this dimension
}

// AdsSymbol represents an ADS symbol with metadata about a PLC variable.
//
// An ADS symbol contains information about a variable in the PLC including
// its memory location (IndexGroup/IndexOffset), size, data type, and descriptive
// information (name, type name, comment).
type AdsSymbol struct {
	IndexGroup    uint32               // IndexGroup is the ADS index group for accessing this symbol's data
	IndexOffset   uint32               // IndexOffset is the ADS index offset for accessing this symbol's data
	Size          uint32               // Size is the size of the variable in bytes
	DataType      types.ADSDataType    // DataType is the ADS data type identifier (see types.ADSDataType)
	Flags         types.ADSSymbolFlags // Flags contains symbol properties (persistent, static, read-only, etc.)
	NameLength    uint16               // NameLength is the length of the Name string (as reported by PLC)
	TypeLength    uint16               // TypeLength is the length of the Type string (as reported by PLC)
	CommentLength uint16               // CommentLength is the length of the Comment string (as reported by PLC)
	Name          string               // Name is the variable name in the PLC (e.g., "MAIN.MyVariable")
	Type          string               // Type is the variable type name (e.g., "INT", "ARRAY[0..10] OF REAL")
	Comment       string               // Comment is the descriptive comment from the PLC code
	ArrayInfo     []ArrayInfoEntry     // ArrayInfo holds dimension info for array types (present when arrayDimension > 0)
	TypeGUID      string               // TypeGUID is the 16-byte type GUID as a hex string (present when Flags & ADSSymbolFlagTypeGuid)
	Attributes    []SymbolAttribute    // Attributes holds TwinCAT pragma attributes (present when Flags & ADSSymbolFlagAttributes)
}
