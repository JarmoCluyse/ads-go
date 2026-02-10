// Package adsstateinfo provides parsing and validation for ADS system state and device information.
//
// This package handles parsing of two common ADS response types:
//   - System State: The current operational state of the PLC/device
//   - Device Info: Version and identification information about the device
//
// These are fundamental operations used to query the status and identity of
// Beckhoff TwinCAT devices.
//
// # System State Format
//
// System state responses follow this binary format (all fields in little-endian):
//   - Bytes 0-1: AdsState (uint16) - Operational state (Run, Stop, Config, etc.)
//   - Bytes 2-3: DeviceState (uint16) - Device-specific state value
//
// The AdsState indicates the current mode of the PLC. Common values include:
//   - ADSStateRun: PLC is running
//   - ADSStateStop: PLC is stopped
//   - ADSStateConfig: PLC is in configuration mode
//   - ADSStateInvalid: Invalid or error state
//
// See types.ADSState for all possible state values.
//
// # Device Info Format
//
// Device info responses follow this binary format:
//   - Byte 0:      MajorVersion (uint8)
//   - Byte 1:      MinorVersion (uint8)
//   - Bytes 2-3:   VersionBuild (uint16) in little-endian
//   - Bytes 4-19:  DeviceName (16-byte null-padded string)
//
// Special case: Some older or embedded TwinCAT runtimes may return an empty
// response (0 bytes) after the error code. ParseDeviceInfo handles this by
// returning a zero-value DeviceInfo.
//
// # Core Functions
//
// ParseSystemState parses the system state from binary data:
//
//	// Query system state
//	response, err := client.send(readStateRequest)
//	if err != nil {
//	    return err
//	}
//
//	// Strip error code and parse state
//	payload, err := adserrors.StripAdsError(response)
//	if err != nil {
//	    return err
//	}
//
//	state, err := adsstateinfo.ParseSystemState(payload)
//	if err != nil {
//	    return err
//	}
//
//	fmt.Printf("PLC State: %v, Device State: %d\n",
//	    state.AdsState, state.DeviceState)
//
// ParseDeviceInfo parses device information from binary data:
//
//	// Query device info
//	response, err := client.send(readDeviceInfoRequest)
//	if err != nil {
//	    return err
//	}
//
//	// Strip error code and parse info
//	payload, err := adserrors.StripAdsError(response)
//	if err != nil {
//	    return err
//	}
//
//	info, err := adsstateinfo.ParseDeviceInfo(payload)
//	if err != nil {
//	    return err
//	}
//
//	fmt.Printf("Device: %s v%d.%d.%d\n",
//	    info.DeviceName, info.MajorVersion, info.MinorVersion, info.VersionBuild)
//
// CheckSystemState and CheckDeviceInfo provide validation without parsing:
//
//	if err := adsstateinfo.CheckSystemState(data); err != nil {
//	    log.Printf("Invalid state data: %v", err)
//	    return err
//	}
//
// # Error Types
//
// The package exports sentinel errors for error inspection:
//   - ErrInvalidStateLength: System state data is not 4 bytes
//   - ErrInvalidDeviceInfoLength: Device info data is less than 20 bytes (and not empty)
//
// Use errors.Is() for error type checking:
//
//	if errors.Is(err, adsstateinfo.ErrInvalidStateLength) {
//	    // Handle invalid length
//	}
//
// # Common Use Cases
//
// Checking if PLC is running:
//
//	state, err := client.ReadTcSystemState()
//	if err != nil {
//	    return err
//	}
//
//	if state.AdsState == types.ADSStateRun {
//	    fmt.Println("PLC is running")
//	} else {
//	    fmt.Printf("PLC is not running (state: %v)\n", state.AdsState)
//	}
//
// Identifying the device:
//
//	info, err := client.ReadDeviceInfo()
//	if err != nil {
//	    return err
//	}
//
//	fmt.Printf("Connected to: %s\n", info.DeviceName)
//	fmt.Printf("Version: %d.%d.%d\n",
//	    info.MajorVersion, info.MinorVersion, info.VersionBuild)
//
// # Backward Compatibility
//
// The ParseDeviceInfo function handles empty responses (0 bytes) for backward
// compatibility with older TwinCAT runtimes that may not return device info.
// In this case, a zero-value DeviceInfo is returned without error.
package adsstateinfo
