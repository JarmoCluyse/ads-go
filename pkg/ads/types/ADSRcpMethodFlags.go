package types

// ADSRcpMethodFlags defines the flags for an ADS RCP method.
// Source: TwinCAT.Ads.dll By Beckhoff
type ADSRcpMethodFlags uint32

const (
	ADSRcpMethodFlagPlcCallingConvention ADSRcpMethodFlags = 0x1 // (ADSMETHODFLAG_PLC_CALLINGCONVENTION)
	ADSRcpMethodFlagCallUnlocked         ADSRcpMethodFlags = 0x2 // (ADSMETHODFLAG_CALL_UNLOCKED)
	ADSRcpMethodFlagNotCallable          ADSRcpMethodFlags = 0x4 // (ADSMETHODFLAG_NOTCALLABLE)
	ADSRcpMethodFlagAttributes           ADSRcpMethodFlags = 0x8 // (ADSMETHODFLAG_ATTRIBUTES)
)

// ADSRcpMethodFlagsToStringArray converts the flags to a string array.
func ADSRcpMethodFlagsToStringArray(flags ADSRcpMethodFlags) []string {
	var result []string
	if flags&ADSRcpMethodFlagPlcCallingConvention != 0 {
		result = append(result, "PlcCallingConvention")
	}
	if flags&ADSRcpMethodFlagCallUnlocked != 0 {
		result = append(result, "CallUnlocked")
	}
	if flags&ADSRcpMethodFlagNotCallable != 0 {
		result = append(result, "NotCallable")
	}
	if flags&ADSRcpMethodFlagAttributes != 0 {
		result = append(result, "Attributes")
	}
	if len(result) == 0 {
		result = append(result, "None")
	}
	return result
}
