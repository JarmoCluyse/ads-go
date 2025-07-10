package types

// ADSSymbolFlags defines the flags for an ADS symbol.
type ADSSymbolFlags uint32

const (
	ADSSymbolFlagPersistent      ADSSymbolFlags = 0x1    // ADSSymbolFlagPersistent
	ADSSymbolFlagBitValue        ADSSymbolFlags = 0x2    // ADSSymbolFlagBitValue
	ADSSymbolFlagReferenceTo     ADSSymbolFlags = 0x4    // ADSSymbolFlagReferenceTo
	ADSSymbolFlagTypeGuid        ADSSymbolFlags = 0x8    // ADSSymbolFlagTypeGuid
	ADSSymbolFlagTComIfacePtr    ADSSymbolFlags = 0x10   // ADSSymbolFlagTComIfacePtr
	ADSSymbolFlagReadOnly        ADSSymbolFlags = 0x20   // ADSSymbolFlagReadOnly
	ADSSymbolFlagItfMethodAccess ADSSymbolFlags = 0x40   // ADSSymbolFlagItfMethodAccess
	ADSSymbolFlagMethodDeref     ADSSymbolFlags = 0x80   // ADSSymbolFlagMethodDeref
	ADSSymbolFlagContextMask     ADSSymbolFlags = 0xF00  // ADSSymbolFlagContextMask
	ADSSymbolFlagAttributes      ADSSymbolFlags = 0x1000 // ADSSymbolFlagAttributes
	ADSSymbolFlagStatic          ADSSymbolFlags = 0x2000 // ADSSymbolFlagStatic
	ADSSymbolFlagInitOnReset     ADSSymbolFlags = 0x4000 // ADSSymbolFlagInitOnReset
	ADSSymbolFlagExtendedFlags   ADSSymbolFlags = 0x8000 // ADSSymbolFlagExtendedFlags
)

// ADSSymbolFlagsToStringArray converts the flags to a string array.
func ADSSymbolFlagsToStringArray(flags ADSSymbolFlags) []string {
	var result []string
	if flags&ADSSymbolFlagPersistent != 0 {
		result = append(result, "Persistent")
	}
	if flags&ADSSymbolFlagBitValue != 0 {
		result = append(result, "BitValue")
	}
	if flags&ADSSymbolFlagReferenceTo != 0 {
		result = append(result, "ReferenceTo")
	}
	if flags&ADSSymbolFlagTypeGuid != 0 {
		result = append(result, "TypeGuid")
	}
	if flags&ADSSymbolFlagTComIfacePtr != 0 {
		result = append(result, "TComInterfacePtr")
	}
	if flags&ADSSymbolFlagReadOnly != 0 {
		result = append(result, "ReadOnly")
	}
	if flags&ADSSymbolFlagItfMethodAccess != 0 {
		result = append(result, "ItfMethodAccess")
	}
	if flags&ADSSymbolFlagMethodDeref != 0 {
		result = append(result, "MethodDeref")
	}
	if flags&ADSSymbolFlagContextMask != 0 {
		result = append(result, "ContextMask")
	}
	if flags&ADSSymbolFlagAttributes != 0 {
		result = append(result, "Attributes")
	}
	if flags&ADSSymbolFlagStatic != 0 {
		result = append(result, "Static")
	}
	if flags&ADSSymbolFlagInitOnReset != 0 {
		result = append(result, "InitOnReset")
	}
	if flags&ADSSymbolFlagExtendedFlags != 0 {
		result = append(result, "ExtendedFlags")
	}
	return result
}
