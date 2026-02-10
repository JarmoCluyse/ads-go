package adsstateinfo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// Sentinel errors for parsing failures
var (
	ErrInvalidStateLength      = errors.New("invalid system state data length")
	ErrInvalidDeviceInfoLength = errors.New("invalid device info data length")
)

// ParseSystemState parses TwinCAT system state from binary data.
//
// Binary format (all fields in little-endian):
//   - Bytes 0-1: AdsState (uint16)
//   - Bytes 2-3: DeviceState (uint16)
//
// Returns the parsed SystemState and any error encountered.
func ParseSystemState(data []byte) (SystemState, error) {
	if len(data) < 4 {
		return SystemState{}, fmt.Errorf("%w: expected 4 bytes, got %d", ErrInvalidStateLength, len(data))
	}

	adsState := binary.LittleEndian.Uint16(data[0:2])
	deviceState := binary.LittleEndian.Uint16(data[2:4])

	return SystemState{
		AdsState:    types.ADSState(adsState),
		DeviceState: deviceState,
	}, nil
}

// CheckSystemState validates system state data without parsing it.
// This is useful for validation before passing data to ParseSystemState.
//
// Returns nil if the data appears valid, or an error describing the issue.
func CheckSystemState(data []byte) error {
	if len(data) < 4 {
		return fmt.Errorf("%w: expected 4 bytes, got %d", ErrInvalidStateLength, len(data))
	}
	return nil
}

// ParseDeviceInfo parses ADS device information from binary data.
//
// Binary format:
//   - Byte 0:    MajorVersion (uint8)
//   - Byte 1:    MinorVersion (uint8)
//   - Bytes 2-3: VersionBuild (uint16) in little-endian
//   - Bytes 4-19: DeviceName (16-byte null-padded string)
//
// Special case: If data is empty (0 bytes), returns a zero-value DeviceInfo.
// This handles older/embedded runtimes that may return only an error code.
//
// Returns the parsed DeviceInfo and any error encountered.
func ParseDeviceInfo(data []byte) (DeviceInfo, error) {
	// Handle empty response (some older runtimes return only error code)
	if len(data) == 0 {
		return DeviceInfo{
			MajorVersion: 0,
			MinorVersion: 0,
			VersionBuild: 0,
			DeviceName:   "",
		}, nil
	}

	if len(data) < 20 {
		return DeviceInfo{}, fmt.Errorf("%w: expected 20 bytes, got %d", ErrInvalidDeviceInfoLength, len(data))
	}

	majorVersion := data[0]
	minorVersion := data[1]
	versionBuild := binary.LittleEndian.Uint16(data[2:4])
	deviceName := string(bytes.Trim(data[4:20], "\x00"))

	return DeviceInfo{
		MajorVersion: majorVersion,
		MinorVersion: minorVersion,
		VersionBuild: versionBuild,
		DeviceName:   deviceName,
	}, nil
}

// CheckDeviceInfo validates device info data without parsing it.
// This is useful for validation before passing data to ParseDeviceInfo.
//
// Returns nil if the data appears valid, or an error describing the issue.
// Note: Empty data (0 bytes) is considered valid for backward compatibility.
func CheckDeviceInfo(data []byte) error {
	// Empty response is valid (some older runtimes)
	if len(data) == 0 {
		return nil
	}

	if len(data) < 20 {
		return fmt.Errorf("%w: expected 20 bytes, got %d", ErrInvalidDeviceInfoLength, len(data))
	}
	return nil
}
