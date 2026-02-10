package types

// ADSRcpMethodParamFlags defines the flags for ADS RCP method parameters.
// Source: TwinCAT.Ads.dll By Beckhoff
type ADSRcpMethodParamFlags uint32

const (
	ADSRcpMethodParamFlagIn              ADSRcpMethodParamFlags = 0x1  // Input Parameter (ADSMETHODPARAFLAG_IN)
	ADSRcpMethodParamFlagOut             ADSRcpMethodParamFlags = 0x2  // Output Parameter (ADSMETHODPARAFLAG_OUT)
	ADSRcpMethodParamFlagByReference     ADSRcpMethodParamFlags = 0x4  // By reference Parameter (ADSMETHODPARAFLAG_BYREFERENCE / ADSMETHODPARAFLAG_RPC_OUTPTR)
	ADSRcpMethodParamFlagMaskRpcArrayDim ADSRcpMethodParamFlags = 0x30 // (ADSMETHODPARAFLAG_RPC_ARRAYDIM_MASK)
	ADSRcpMethodParamFlagAttributes      ADSRcpMethodParamFlags = 0x40 // Attributes (ADSMETHODPARAFLAG_ATTRIBUTES)
)

// ADSRcpMethodParamFlagsToStringArray converts the flags to a string array.
func ADSRcpMethodParamFlagsToStringArray(flags ADSRcpMethodParamFlags) []string {
	var result []string
	if flags&ADSRcpMethodParamFlagIn != 0 {
		result = append(result, "In")
	}
	if flags&ADSRcpMethodParamFlagOut != 0 {
		result = append(result, "Out")
	}
	if flags&ADSRcpMethodParamFlagByReference != 0 {
		result = append(result, "ByReference")
	}
	if flags&ADSRcpMethodParamFlagMaskRpcArrayDim != 0 {
		result = append(result, "MaskRpcArrayDim")
	}
	if flags&ADSRcpMethodParamFlagAttributes != 0 {
		result = append(result, "Attributes")
	}
	if len(result) == 0 {
		result = append(result, "None")
	}
	return result
}
