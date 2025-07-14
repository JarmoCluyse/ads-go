package types

// Reserved ADS index groups
// Source: TwinCAT.Ads.dll By Beckhoff
type ADSReservedIndexGroup uint32

const (
	ADSReservedIndexGroupPlcRWIB                  ADSReservedIndexGroup = 16384  // 0x00004000
	ADSReservedIndexGroupPlcRWOB                  ADSReservedIndexGroup = 16400  // 0x00004010
	ADSReservedIndexGroupPlcRWMB                  ADSReservedIndexGroup = 16416  // 0x00004020
	ADSReservedIndexGroupPlcRWRB                  ADSReservedIndexGroup = 16432  // 0x00004030
	ADSReservedIndexGroupPlcRWDB                  ADSReservedIndexGroup = 16448  // 0x00004040
	ADSReservedIndexGroupSymbolTable              ADSReservedIndexGroup = 61440  // 0x0000F000
	ADSReservedIndexGroupSymbolName               ADSReservedIndexGroup = 61441  // 0x0000F001
	ADSReservedIndexGroupSymbolValue              ADSReservedIndexGroup = 61442  // 0x0000F002
	ADSReservedIndexGroupSymbolHandleByName       ADSReservedIndexGroup = 61443  // 0x0000F003
	ADSReservedIndexGroupSymbolValueByName        ADSReservedIndexGroup = 61444  // 0x0000F004
	ADSReservedIndexGroupSymbolValueByHandle      ADSReservedIndexGroup = 61445  // 0x0000F005
	ADSReservedIndexGroupSymbolReleaseHandle      ADSReservedIndexGroup = 61446  // 0x0000F006
	ADSReservedIndexGroupSymbolInfoByName         ADSReservedIndexGroup = 61447  // 0x0000F007
	ADSReservedIndexGroupSymbolVersion            ADSReservedIndexGroup = 61448  // 0x0000F008
	ADSReservedIndexGroupSymbolInfoByNameEx       ADSReservedIndexGroup = 61449  // 0x0000F009
	ADSReservedIndexGroupSymbolDownload           ADSReservedIndexGroup = 61450  // 0x0000F00A
	ADSReservedIndexGroupSymbolUpload             ADSReservedIndexGroup = 61451  // 0x0000F00B
	ADSReservedIndexGroupSymbolUploadInfo         ADSReservedIndexGroup = 61452  // 0x0000F00C
	ADSReservedIndexGroupSymbolDownload2          ADSReservedIndexGroup = 0xF00D // Added, not from .dll
	ADSReservedIndexGroupSymbolDataTypeUpload     ADSReservedIndexGroup = 0xF00E // Added, not from .dll
	ADSReservedIndexGroupSymbolUploadInfo2        ADSReservedIndexGroup = 0xF00F // Added, not from .dll - 24 bytes of info, uploadinfo3 would contain 64 bytes
	ADSReservedIndexGroupSymbolNote               ADSReservedIndexGroup = 61456  // 0x0000F010
	ADSReservedIndexGroupDataDataTypeInfoByNameEx ADSReservedIndexGroup = 0xF011 // Added, not from .dll
	ADSReservedIndexGroupIOImageRWIB              ADSReservedIndexGroup = 61472  // 0x0000F020
	ADSReservedIndexGroupIOImageRWIX              ADSReservedIndexGroup = 61473  // 0x0000F021
	ADSReservedIndexGroupIOImageRWOB              ADSReservedIndexGroup = 61488  // 0x0000F030
	ADSReservedIndexGroupIOImageRWOX              ADSReservedIndexGroup = 61489  // 0x0000F031
	ADSReservedIndexGroupIOImageClearI            ADSReservedIndexGroup = 61504  // 0x0000F040
	ADSReservedIndexGroupIOImageClearO            ADSReservedIndexGroup = 61520  // 0x0000F050
	ADSReservedIndexGroupSumCommandRead           ADSReservedIndexGroup = 61568  // 0x0000F080
	ADSReservedIndexGroupSumCommandWrite          ADSReservedIndexGroup = 61569  // 0x0000F081
	ADSReservedIndexGroupSumCommandReadWrite      ADSReservedIndexGroup = 61570  // 0x0000F082
	ADSReservedIndexGroupSumCommandReadEx         ADSReservedIndexGroup = 61571  // 0x0000F083
	ADSReservedIndexGroupSumCommandReadEx2        ADSReservedIndexGroup = 61572  // 0x0000F084
	ADSReservedIndexGroupSumCommandAddDevNote     ADSReservedIndexGroup = 61573  // 0x0000F085
	ADSReservedIndexGroupSumCommandDelDevNote     ADSReservedIndexGroup = 61574  // 0x0000F086
	ADSReservedIndexGroupDeviceData               ADSReservedIndexGroup = 61696  // 0x0000F100
)

// ADSReservedIndexGroupToString returns the name of the reserved ADS index group by value, or "UNKNOWN"
func (group ADSReservedIndexGroup) ADSReservedIndexGroupToString() string {
	switch group {
	case ADSReservedIndexGroupPlcRWIB:
		return "PlcRWIB"
	case ADSReservedIndexGroupPlcRWOB:
		return "PlcRWOB"
	case ADSReservedIndexGroupPlcRWMB:
		return "PlcRWMB"
	case ADSReservedIndexGroupPlcRWRB:
		return "PlcRWRB"
	case ADSReservedIndexGroupPlcRWDB:
		return "PlcRWDB"
	case ADSReservedIndexGroupSymbolTable:
		return "SymbolTable"
	case ADSReservedIndexGroupSymbolName:
		return "SymbolName"
	case ADSReservedIndexGroupSymbolValue:
		return "SymbolValue"
	case ADSReservedIndexGroupSymbolHandleByName:
		return "SymbolHandleByName"
	case ADSReservedIndexGroupSymbolValueByName:
		return "SymbolValueByName"
	case ADSReservedIndexGroupSymbolValueByHandle:
		return "SymbolValueByHandle"
	case ADSReservedIndexGroupSymbolReleaseHandle:
		return "SymbolReleaseHandle"
	case ADSReservedIndexGroupSymbolInfoByName:
		return "SymbolInfoByName"
	case ADSReservedIndexGroupSymbolVersion:
		return "SymbolVersion"
	case ADSReservedIndexGroupSymbolInfoByNameEx:
		return "SymbolInfoByNameEx"
	case ADSReservedIndexGroupSymbolDownload:
		return "SymbolDownload"
	case ADSReservedIndexGroupSymbolUpload:
		return "SymbolUpload"
	case ADSReservedIndexGroupSymbolUploadInfo:
		return "SymbolUploadInfo"
	case ADSReservedIndexGroupSymbolDownload2:
		return "SymbolDownload2"
	case ADSReservedIndexGroupSymbolDataTypeUpload:
		return "SymbolDataTypeUpload"
	case ADSReservedIndexGroupSymbolUploadInfo2:
		return "SymbolUploadInfo2"
	case ADSReservedIndexGroupSymbolNote:
		return "SymbolNote"
	case ADSReservedIndexGroupDataDataTypeInfoByNameEx:
		return "DataDataTypeInfoByNameEx"
	case ADSReservedIndexGroupIOImageRWIB:
		return "IOImageRWIB"
	case ADSReservedIndexGroupIOImageRWIX:
		return "IOImageRWIX"
	case ADSReservedIndexGroupIOImageRWOB:
		return "IOImageRWOB"
	case ADSReservedIndexGroupIOImageRWOX:
		return "IOImageRWOX"
	case ADSReservedIndexGroupIOImageClearI:
		return "IOImageClearI"
	case ADSReservedIndexGroupIOImageClearO:
		return "IOImageClearO"
	case ADSReservedIndexGroupSumCommandRead:
		return "SumCommandRead"
	case ADSReservedIndexGroupSumCommandWrite:
		return "SumCommandWrite"
	case ADSReservedIndexGroupSumCommandReadWrite:
		return "SumCommandReadWrite"
	case ADSReservedIndexGroupSumCommandReadEx:
		return "SumCommandReadEx"
	case ADSReservedIndexGroupSumCommandReadEx2:
		return "SumCommandReadEx2"
	case ADSReservedIndexGroupSumCommandAddDevNote:
		return "SumCommandAddDevNote"
	case ADSReservedIndexGroupSumCommandDelDevNote:
		return "SumCommandDelDevNote"
	case ADSReservedIndexGroupDeviceData:
		return "DeviceData"
	default:
		return "UNKNOWN"
	}
}
