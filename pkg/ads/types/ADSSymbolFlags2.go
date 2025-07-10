package types

// ADSSymbolFlags2 defines the extended flags for an ADS symbol.
type ADSSymbolFlags2 uint32

const (
	ADSSymbolFlag2None                   ADSSymbolFlags2 = 0x0  // None
	ADSSymbolFlag2PlcPointerType         ADSSymbolFlags2 = 0x1  // PLC pointer type (ADSSYMBOLFLAG2_PLCPOINTERTYPE)
	ADSSymbolFlag2RedundancyIgnore       ADSSymbolFlags2 = 0x2  // Ignore symbol while equalizing redundancy projects (ADSSYMBOLFLAG2_RDIGNORE)
	ADSSymbolFlag2RefactorInfo           ADSSymbolFlags2 = 0x4  // Contains refactoring information (ADSSYMBOLFLAG2_REFACTORINFO)
	ADSSymbolFlag2OnlineChangePtrRefType ADSSymbolFlags2 = 0x8  // Online change PTR reference type (ADSSYMBOLFLAG2_OCPTRREFTYPE)
	ADSSymbolFlag2VariantType            ADSSymbolFlags2 = 0x10 // Symbol is a Variant Type (ADSSYMBOLFLAG2_VARIANT)
)

// ADSSymbolFlags2ToStringArray converts the extended flags to a string array.
func ADSSymbolFlags2ToStringArray(flags ADSSymbolFlags2) []string {
	var result []string
	if flags&ADSSymbolFlag2PlcPointerType != 0 {
		result = append(result, "PlcPointerType")
	}
	if flags&ADSSymbolFlag2RedundancyIgnore != 0 {
		result = append(result, "RedundancyIgnore")
	}
	if flags&ADSSymbolFlag2RefactorInfo != 0 {
		result = append(result, "RefactorInfo")
	}
	if flags&ADSSymbolFlag2OnlineChangePtrRefType != 0 {
		result = append(result, "OnlineChangePtrRefType")
	}
	if flags&ADSSymbolFlag2VariantType != 0 {
		result = append(result, "VariantType")
	}
	if len(result) == 0 {
		result = append(result, "None")
	}
	return result
}
