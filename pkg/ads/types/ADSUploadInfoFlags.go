package types

// ADSUploadInfoFlags defines the flags for ADS symbol upload info.
// Source: TwinCAT.Ads.dll By Beckhoff
type ADSUploadInfoFlags uint32

const (
	ADSUploadInfoFlagNone                  ADSUploadInfoFlags = 0x0 // None / Uninitialized
	ADSUploadInfoFlagIs64BitPlatform       ADSUploadInfoFlags = 0x1 // Target is 64 Bit Platform
	ADSUploadInfoFlagIncludesBaseTypes     ADSUploadInfoFlags = 0x2 // Symbol Server includes Base types
	ADSUploadInfoFlagUtf8EncodedStringData ADSUploadInfoFlags = 0x4 // String data is encoded in UTF-8 (All PLC STRING data is encoded in UTF8. This flag is used from TwinCAT 4026 on.)
)

// ADSUploadInfoFlagsToStringArray converts the flags to a string array.
func ADSUploadInfoFlagsToStringArray(flags ADSUploadInfoFlags) []string {
	var result []string
	if flags&ADSUploadInfoFlagIs64BitPlatform != 0 {
		result = append(result, "Is64BitPlatform")
	}
	if flags&ADSUploadInfoFlagIncludesBaseTypes != 0 {
		result = append(result, "IncludesBaseTypes")
	}
	if flags&ADSUploadInfoFlagUtf8EncodedStringData != 0 {
		result = append(result, "Utf8EncodedStringData")
	}
	if len(result) == 0 {
		result = append(result, "None")
	}
	return result
}
