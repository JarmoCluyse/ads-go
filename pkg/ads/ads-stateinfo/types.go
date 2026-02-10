package adsstateinfo

import "github.com/jarmocluyse/ads-go/pkg/ads/types"

// SystemState represents the TwinCAT system state response.
//
// The system state includes both the ADS state (indicating the operational
// state of the PLC) and a device-specific state value.
type SystemState struct {
	// AdsState indicates the current operational state of the ADS device
	// (e.g., Run, Stop, Config, etc.). See types.ADSState for valid values.
	AdsState types.ADSState

	// DeviceState is a device-specific state value. The interpretation
	// depends on the device type and manufacturer.
	DeviceState uint16
}

// ExtendedSystemState represents the TwinCAT extended system state response.
//
// This extended state includes additional information beyond the basic SystemState,
// most importantly the RestartIndex which increments when the TwinCAT system service
// restarts. This allows detection of system restarts even when the ADS state remains "Run".
//
// Extended state is read from ADS port 10000 (system service), IndexGroup 240, IndexOffset 0.
// Not all TwinCAT versions support extended state - use ReadTcSystemExtendedState() which
// gracefully falls back to basic state if extended state is not available.
type ExtendedSystemState struct {
	// AdsState indicates the current operational state of the ADS device
	// (e.g., Run, Stop, Config, etc.). See types.ADSState for valid values.
	AdsState types.ADSState

	// DeviceState is a device-specific state value. The interpretation
	// depends on the device type and manufacturer.
	DeviceState uint16

	// RestartIndex increments each time the TwinCAT system service restarts.
	// This is crucial for detecting system restarts even when AdsState stays "Run".
	RestartIndex uint16

	// Version is the major version number of the TwinCAT system
	Version uint8

	// Revision is the revision number of the TwinCAT system
	Revision uint8

	// Build is the build number of the TwinCAT system
	Build uint16

	// Platform is the platform ID of the TwinCAT system
	Platform uint8

	// OsType is the operating system type ID
	OsType uint8

	// Flags contains system service state flags (interpretation is system-specific)
	Flags uint16
}

// DeviceInfo represents the response for an ADS ReadDeviceInfo command.
//
// Device information includes version details and the device name,
// which helps identify the PLC or TwinCAT runtime being communicated with.
type DeviceInfo struct {
	// MajorVersion is the major version number of the device/runtime
	MajorVersion uint8

	// MinorVersion is the minor version number of the device/runtime
	MinorVersion uint8

	// VersionBuild is the build number of the device/runtime
	VersionBuild uint16

	// DeviceName is the name of the device (e.g., "PLC", "TwinCAT")
	// Maximum 16 characters, null-padded
	DeviceName string
}
