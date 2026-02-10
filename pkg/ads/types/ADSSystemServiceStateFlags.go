package types

// ADSSystemServiceStateFlags defines the flags for system service state.
// Source: TwinCAT.Ads.dll By Beckhoff
type ADSSystemServiceStateFlags uint32

const (
	ADSSystemServiceStateFlagRouterModeOnly      ADSSystemServiceStateFlags = 0x1   // The router mode only
	ADSSystemServiceStateFlagRedundancySystem    ADSSystemServiceStateFlags = 0x2   // System is part of a controller redundancy
	ADSSystemServiceStateFlagRedundancyPrimary   ADSSystemServiceStateFlags = 0x4   // System is the primary controller
	ADSSystemServiceStateFlagRedundancyActive    ADSSystemServiceStateFlags = 0x10  // System is currently active -> controlling the machine
	ADSSystemServiceStateFlagDataFolderSupport   ADSSystemServiceStateFlags = 0x20  // The data folder support
	ADSSystemServiceStateFlagRedundancyInOp      ADSSystemServiceStateFlags = 0x40  // Redundancy is currently down -> not synchronized
	ADSSystemServiceStateFlagRedundancySuspended ADSSystemServiceStateFlags = 0x80  // Standby system is currently suspended - e.g. while online change
	ADSSystemServiceStateFlagNewCurrentConfig    ADSSystemServiceStateFlags = 0x100 // Creates new currentconfig
)

// ADSSystemServiceStateFlagsToStringArray converts the flags to a string array.
func ADSSystemServiceStateFlagsToStringArray(flags ADSSystemServiceStateFlags) []string {
	var result []string
	if flags&ADSSystemServiceStateFlagRouterModeOnly != 0 {
		result = append(result, "RouterModeOnly")
	}
	if flags&ADSSystemServiceStateFlagRedundancySystem != 0 {
		result = append(result, "RedundancySystem")
	}
	if flags&ADSSystemServiceStateFlagRedundancyPrimary != 0 {
		result = append(result, "RedundancyPrimary")
	}
	if flags&ADSSystemServiceStateFlagRedundancyActive != 0 {
		result = append(result, "RedundancyActive")
	}
	if flags&ADSSystemServiceStateFlagDataFolderSupport != 0 {
		result = append(result, "DataFolderSupport")
	}
	if flags&ADSSystemServiceStateFlagRedundancyInOp != 0 {
		result = append(result, "RedundancyInOp")
	}
	if flags&ADSSystemServiceStateFlagRedundancySuspended != 0 {
		result = append(result, "RedundancySuspended")
	}
	if flags&ADSSystemServiceStateFlagNewCurrentConfig != 0 {
		result = append(result, "NewCurrentConfig")
	}
	if len(result) == 0 {
		result = append(result, "None")
	}
	return result
}
