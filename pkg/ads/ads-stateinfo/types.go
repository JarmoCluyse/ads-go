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
