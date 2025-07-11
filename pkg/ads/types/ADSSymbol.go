package types

// AdsSymbol represents an ADS symbol.
type AdsSymbol struct {
	IndexGroup    uint32
	IndexOffset   uint32
	Size          uint32
	DataType      ADSDataType
	Flags         ADSSymbolFlags
	NameLength    uint16
	TypeLength    uint16
	CommentLength uint16
	Name          string
	Type          string
	Comment       string
}
