package types

// ADSDataTypeFlags defines the flags for an ADS data type.
// Source: TwinCAT.Ads.dll By Beckhoff
type ADSDataTypeFlags uint32

const (
	ADSDataTypeFlagNone                       ADSDataTypeFlags = 0x0        // None / No Flag set
	ADSDataTypeFlagDataType                   ADSDataTypeFlags = 0x1        // ADSDATATYPEFLAG_DATATYPE
	ADSDataTypeFlagDataItem                   ADSDataTypeFlags = 0x2        // ADSDATATYPEFLAG_DATAITEM
	ADSDataTypeFlagReferenceTo                ADSDataTypeFlags = 0x4        // ADSDATATYPEFLAG_REFERENCETO
	ADSDataTypeFlagMethodDeref                ADSDataTypeFlags = 0x8        // ADSDATATYPEFLAG_METHODDEREF
	ADSDataTypeFlagOversample                 ADSDataTypeFlags = 0x10       // ADSDATATYPEFLAG_OVERSAMPLE
	ADSDataTypeFlagBitValues                  ADSDataTypeFlags = 0x20       // ADSDATATYPEFLAG_BITVALUES
	ADSDataTypeFlagPropItem                   ADSDataTypeFlags = 0x40       // ADSDATATYPEFLAG_PROPITEM
	ADSDataTypeFlagTypeGuid                   ADSDataTypeFlags = 0x80       // ADSDATATYPEFLAG_TYPEGUID
	ADSDataTypeFlagPersistent                 ADSDataTypeFlags = 0x100      // ADSDATATYPEFLAG_PERSISTENT
	ADSDataTypeFlagCopyMask                   ADSDataTypeFlags = 0x200      // ADSDATATYPEFLAG_COPYMASK
	ADSDataTypeFlagTComInterfacePtr           ADSDataTypeFlags = 0x400      // ADSDATATYPEFLAG_TCCOMIFACEPTR
	ADSDataTypeFlagMethodInfos                ADSDataTypeFlags = 0x800      // ADSDATATYPEFLAG_METHODINFOS
	ADSDataTypeFlagAttributes                 ADSDataTypeFlags = 0x1000     // ADSDATATYPEFLAG_ATTRIBUTES
	ADSDataTypeFlagEnumInfos                  ADSDataTypeFlags = 0x2000     // ADSDATATYPEFLAG_ENUMINFOS
	ADSDataTypeFlagAligned                    ADSDataTypeFlags = 0x10000    // ADSDATATYPEFLAG_ALIGNED
	ADSDataTypeFlagStatic                     ADSDataTypeFlags = 0x20000    // ADSDATATYPEFLAG_STATIC
	ADSDataTypeFlagSoftwareProtectionLvls     ADSDataTypeFlags = 0x40000    // Has Software Protection Levels for DataTypes
	ADSDataTypeFlagIgnorePersist              ADSDataTypeFlags = 0x80000    // ADSDATATYPEFLAG_IGNOREPERSIST
	ADSDataTypeFlagAnySizeArray               ADSDataTypeFlags = 0x100000   // ADSDATATYPEFLAG_ANYSIZEARRAY
	ADSDataTypeFlagPersistantDatatype         ADSDataTypeFlags = 0x200000   // ADSDATATYPEFLAG_PERSIST_DT
	ADSDataTypeFlagInitOnReset                ADSDataTypeFlags = 0x400000   // ADSDATATYPEFLAG_INITONRESET
	ADSDataTypeFlagPlcPointerType             ADSDataTypeFlags = 0x800000   // ADSDATATYPEFLAG_PLCPOINTERTYPE
	ADSDataTypeFlagRefactorInfo               ADSDataTypeFlags = 0x1000000  // ADSDATATYPEFLAG_REFACTORINFO
	ADSDataTypeFlagHideSubItems               ADSDataTypeFlags = 0x2000000  // ADSDATATYPEFLAG_HIDESUBITEMS
	ADSDataTypeFlagIncomplete                 ADSDataTypeFlags = 0x4000000  // ADSDATATYPEFLAG_INCOMPLETE
	ADSDataTypeFlagContainsOnlineChangePtrRef ADSDataTypeFlags = 0x8000000  // ADSDATATYPEFLAG_OCPTRREFTYPE
	ADSDataTypeFlagDeRefTypeItem              ADSDataTypeFlags = 0x10000000 // ADSDATATYPEFLAG_DEREFTYPEITEM
	ADSDataTypeFlagExtendedEnumInfos          ADSDataTypeFlags = 0x20000000 // ADSDATATYPEFLAG_EXTENUMINFOS
	ADSDataTypeFlagExtendedFlags              ADSDataTypeFlags = 0x80000000 // ADSDATATYPEFLAG_EXTENDEDFLAGS
)

// ADSDataTypeFlagsToStringArray converts the flags to a string array.
func ADSDataTypeFlagsToStringArray(flags ADSDataTypeFlags) []string {
	var result []string
	if flags&ADSDataTypeFlagDataType != 0 {
		result = append(result, "DataType")
	}
	if flags&ADSDataTypeFlagDataItem != 0 {
		result = append(result, "DataItem")
	}
	if flags&ADSDataTypeFlagReferenceTo != 0 {
		result = append(result, "ReferenceTo")
	}
	if flags&ADSDataTypeFlagMethodDeref != 0 {
		result = append(result, "MethodDeref")
	}
	if flags&ADSDataTypeFlagOversample != 0 {
		result = append(result, "Oversample")
	}
	if flags&ADSDataTypeFlagBitValues != 0 {
		result = append(result, "BitValues")
	}
	if flags&ADSDataTypeFlagPropItem != 0 {
		result = append(result, "PropItem")
	}
	if flags&ADSDataTypeFlagTypeGuid != 0 {
		result = append(result, "TypeGuid")
	}
	if flags&ADSDataTypeFlagPersistent != 0 {
		result = append(result, "Persistent")
	}
	if flags&ADSDataTypeFlagCopyMask != 0 {
		result = append(result, "CopyMask")
	}
	if flags&ADSDataTypeFlagTComInterfacePtr != 0 {
		result = append(result, "TComInterfacePtr")
	}
	if flags&ADSDataTypeFlagMethodInfos != 0 {
		result = append(result, "MethodInfos")
	}
	if flags&ADSDataTypeFlagAttributes != 0 {
		result = append(result, "Attributes")
	}
	if flags&ADSDataTypeFlagEnumInfos != 0 {
		result = append(result, "EnumInfos")
	}
	if flags&ADSDataTypeFlagAligned != 0 {
		result = append(result, "Aligned")
	}
	if flags&ADSDataTypeFlagStatic != 0 {
		result = append(result, "Static")
	}
	if flags&ADSDataTypeFlagSoftwareProtectionLvls != 0 {
		result = append(result, "SoftwareProtectionLevels")
	}
	if flags&ADSDataTypeFlagIgnorePersist != 0 {
		result = append(result, "IgnorePersist")
	}
	if flags&ADSDataTypeFlagAnySizeArray != 0 {
		result = append(result, "AnySizeArray")
	}
	if flags&ADSDataTypeFlagPersistantDatatype != 0 {
		result = append(result, "PersistantDatatype")
	}
	if flags&ADSDataTypeFlagInitOnReset != 0 {
		result = append(result, "InitOnReset")
	}
	if flags&ADSDataTypeFlagPlcPointerType != 0 {
		result = append(result, "PlcPointerType")
	}
	if flags&ADSDataTypeFlagRefactorInfo != 0 {
		result = append(result, "RefactorInfo")
	}
	if flags&ADSDataTypeFlagHideSubItems != 0 {
		result = append(result, "HideSubItems")
	}
	if flags&ADSDataTypeFlagIncomplete != 0 {
		result = append(result, "Incomplete")
	}
	if flags&ADSDataTypeFlagContainsOnlineChangePtrRef != 0 {
		result = append(result, "ContainsOnlineChangePtrRef")
	}
	if flags&ADSDataTypeFlagDeRefTypeItem != 0 {
		result = append(result, "DeRefTypeItem")
	}
	if flags&ADSDataTypeFlagExtendedEnumInfos != 0 {
		result = append(result, "ExtendedEnumInfos")
	}
	if flags&ADSDataTypeFlagExtendedFlags != 0 {
		result = append(result, "ExtendedFlags")
	}
	if len(result) == 0 {
		result = append(result, "None")
	}
	return result
}
