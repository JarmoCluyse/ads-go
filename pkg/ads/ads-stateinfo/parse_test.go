package adsstateinfo

import (
	"encoding/binary"
	"errors"
	"testing"

	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/stretchr/testify/assert"
)

func TestParseSystemState(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expected    SystemState
		expectError bool
		errorType   error
	}{
		{
			name: "Valid - Run state",
			data: buildSystemStateData(types.ADSStateRun, 0x1234),
			expected: SystemState{
				AdsState:    types.ADSStateRun,
				DeviceState: 0x1234,
			},
			expectError: false,
		},
		{
			name: "Valid - Stop state",
			data: buildSystemStateData(types.ADSStateStop, 0xABCD),
			expected: SystemState{
				AdsState:    types.ADSStateStop,
				DeviceState: 0xABCD,
			},
			expectError: false,
		},
		{
			name: "Valid - Config state with zero device state",
			data: buildSystemStateData(types.ADSStateConfig, 0x0000),
			expected: SystemState{
				AdsState:    types.ADSStateConfig,
				DeviceState: 0x0000,
			},
			expectError: false,
		},
		{
			name: "Valid - Invalid state with max device state",
			data: buildSystemStateData(types.ADSStateInvalid, 0xFFFF),
			expected: SystemState{
				AdsState:    types.ADSStateInvalid,
				DeviceState: 0xFFFF,
			},
			expectError: false,
		},
		{
			name:        "Invalid - less than 4 bytes",
			data:        []byte{0x01, 0x02, 0x03},
			expected:    SystemState{},
			expectError: true,
			errorType:   ErrInvalidStateLength,
		},
		{
			name:        "Invalid - empty data",
			data:        []byte{},
			expected:    SystemState{},
			expectError: true,
			errorType:   ErrInvalidStateLength,
		},
		{
			name:        "Invalid - exactly 2 bytes",
			data:        []byte{0x01, 0x02},
			expected:    SystemState{},
			expectError: true,
			errorType:   ErrInvalidStateLength,
		},
		{
			name: "Valid - more than 4 bytes (ignores extra)",
			data: append(buildSystemStateData(types.ADSStateRun, 0x5678), 0xFF, 0xFF),
			expected: SystemState{
				AdsState:    types.ADSStateRun,
				DeviceState: 0x5678,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, err := ParseSystemState(tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expected.AdsState, state.AdsState)
			assert.Equal(t, tt.expected.DeviceState, state.DeviceState)
		})
	}
}

func TestCheckSystemState(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
		errorType   error
	}{
		{
			name:        "Valid - exactly 4 bytes",
			data:        buildSystemStateData(types.ADSStateRun, 0x1234),
			expectError: false,
		},
		{
			name:        "Valid - more than 4 bytes",
			data:        append(buildSystemStateData(types.ADSStateStop, 0xABCD), 0xFF, 0xFF),
			expectError: false,
		},
		{
			name:        "Invalid - less than 4 bytes",
			data:        []byte{0x01, 0x02, 0x03},
			expectError: true,
			errorType:   ErrInvalidStateLength,
		},
		{
			name:        "Invalid - empty data",
			data:        []byte{},
			expectError: true,
			errorType:   ErrInvalidStateLength,
		},
		{
			name:        "Invalid - exactly 2 bytes",
			data:        []byte{0x01, 0x02},
			expectError: true,
			errorType:   ErrInvalidStateLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSystemState(tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
	}
}

func TestParseExtendedSystemState(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expected    ExtendedSystemState
		expectError bool
		errorType   error
	}{
		{
			name: "Valid - Run state with restart index",
			data: buildExtendedSystemStateData(types.ADSStateRun, 0x1234, 42, 3, 1, 4024, 1, 2, 0x0100),
			expected: ExtendedSystemState{
				AdsState:     types.ADSStateRun,
				DeviceState:  0x1234,
				RestartIndex: 42,
				Version:      3,
				Revision:     1,
				Build:        4024,
				Platform:     1,
				OsType:       2,
				Flags:        0x0100,
			},
			expectError: false,
		},
		{
			name: "Valid - Config state with zero restart index",
			data: buildExtendedSystemStateData(types.ADSStateConfig, 0xABCD, 0, 3, 1, 4026, 5, 10, 0x0000),
			expected: ExtendedSystemState{
				AdsState:     types.ADSStateConfig,
				DeviceState:  0xABCD,
				RestartIndex: 0,
				Version:      3,
				Revision:     1,
				Build:        4026,
				Platform:     5,
				OsType:       10,
				Flags:        0x0000,
			},
			expectError: false,
		},
		{
			name: "Valid - Stop state with high restart index",
			data: buildExtendedSystemStateData(types.ADSStateStop, 0xFFFF, 65535, 255, 255, 65535, 255, 255, 0xFFFF),
			expected: ExtendedSystemState{
				AdsState:     types.ADSStateStop,
				DeviceState:  0xFFFF,
				RestartIndex: 65535,
				Version:      255,
				Revision:     255,
				Build:        65535,
				Platform:     255,
				OsType:       255,
				Flags:        0xFFFF,
			},
			expectError: false,
		},
		{
			name: "Valid - Restart detected (index incremented)",
			data: buildExtendedSystemStateData(types.ADSStateRun, 0x0000, 100, 3, 1, 4024, 1, 2, 0x0100),
			expected: ExtendedSystemState{
				AdsState:     types.ADSStateRun,
				DeviceState:  0x0000,
				RestartIndex: 100,
				Version:      3,
				Revision:     1,
				Build:        4024,
				Platform:     1,
				OsType:       2,
				Flags:        0x0100,
			},
			expectError: false,
		},
		{
			name:        "Invalid - less than 16 bytes",
			data:        []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
			expected:    ExtendedSystemState{},
			expectError: true,
			errorType:   ErrInvalidExtendedStateLength,
		},
		{
			name:        "Invalid - empty data",
			data:        []byte{},
			expected:    ExtendedSystemState{},
			expectError: true,
			errorType:   ErrInvalidExtendedStateLength,
		},
		{
			name:        "Invalid - exactly 4 bytes (basic state only)",
			data:        []byte{0x05, 0x00, 0x00, 0x00},
			expected:    ExtendedSystemState{},
			expectError: true,
			errorType:   ErrInvalidExtendedStateLength,
		},
		{
			name: "Valid - more than 16 bytes (ignores extra)",
			data: append(buildExtendedSystemStateData(types.ADSStateRun, 0x5678, 10, 3, 1, 4024, 1, 2, 0x0100), 0xFF, 0xFF, 0xFF),
			expected: ExtendedSystemState{
				AdsState:     types.ADSStateRun,
				DeviceState:  0x5678,
				RestartIndex: 10,
				Version:      3,
				Revision:     1,
				Build:        4024,
				Platform:     1,
				OsType:       2,
				Flags:        0x0100,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, err := ParseExtendedSystemState(tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expected.AdsState, state.AdsState)
			assert.Equal(t, tt.expected.DeviceState, state.DeviceState)
			assert.Equal(t, tt.expected.RestartIndex, state.RestartIndex)
			assert.Equal(t, tt.expected.Version, state.Version)
			assert.Equal(t, tt.expected.Revision, state.Revision)
			assert.Equal(t, tt.expected.Build, state.Build)
			assert.Equal(t, tt.expected.Platform, state.Platform)
			assert.Equal(t, tt.expected.OsType, state.OsType)
			assert.Equal(t, tt.expected.Flags, state.Flags)
		})
	}
}

func TestCheckExtendedSystemState(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
		errorType   error
	}{
		{
			name:        "Valid - exactly 16 bytes",
			data:        buildExtendedSystemStateData(types.ADSStateRun, 0x1234, 42, 3, 1, 4024, 1, 2, 0x0100),
			expectError: false,
		},
		{
			name:        "Valid - more than 16 bytes",
			data:        append(buildExtendedSystemStateData(types.ADSStateStop, 0xABCD, 10, 3, 1, 4026, 5, 10, 0x0000), 0xFF, 0xFF),
			expectError: false,
		},
		{
			name:        "Invalid - less than 16 bytes",
			data:        []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
			expectError: true,
			errorType:   ErrInvalidExtendedStateLength,
		},
		{
			name:        "Invalid - empty data",
			data:        []byte{},
			expectError: true,
			errorType:   ErrInvalidExtendedStateLength,
		},
		{
			name:        "Invalid - exactly 4 bytes (basic state only)",
			data:        []byte{0x05, 0x00, 0x00, 0x00},
			expectError: true,
			errorType:   ErrInvalidExtendedStateLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckExtendedSystemState(tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
	}
}

func TestParseDeviceInfo(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expected    DeviceInfo
		expectError bool
		errorType   error
	}{
		{
			name: "Valid - full device info",
			data: buildDeviceInfoData(3, 1, 4024, "PLC Runtime"),
			expected: DeviceInfo{
				MajorVersion: 3,
				MinorVersion: 1,
				VersionBuild: 4024,
				DeviceName:   "PLC Runtime",
			},
			expectError: false,
		},
		{
			name: "Valid - TwinCAT device",
			data: buildDeviceInfoData(3, 1, 4024, "TwinCAT"),
			expected: DeviceInfo{
				MajorVersion: 3,
				MinorVersion: 1,
				VersionBuild: 4024,
				DeviceName:   "TwinCAT",
			},
			expectError: false,
		},
		{
			name: "Valid - zero version",
			data: buildDeviceInfoData(0, 0, 0, "Device"),
			expected: DeviceInfo{
				MajorVersion: 0,
				MinorVersion: 0,
				VersionBuild: 0,
				DeviceName:   "Device",
			},
			expectError: false,
		},
		{
			name: "Valid - max version values",
			data: buildDeviceInfoData(255, 255, 65535, "MaxVersionDevice"),
			expected: DeviceInfo{
				MajorVersion: 255,
				MinorVersion: 255,
				VersionBuild: 65535,
				DeviceName:   "MaxVersionDevice",
			},
			expectError: false,
		},
		{
			name: "Valid - empty device name (null-padded)",
			data: buildDeviceInfoData(1, 0, 100, ""),
			expected: DeviceInfo{
				MajorVersion: 1,
				MinorVersion: 0,
				VersionBuild: 100,
				DeviceName:   "",
			},
			expectError: false,
		},
		{
			name: "Valid - device name with trailing nulls (trimmed)",
			data: buildDeviceInfoDataRaw(2, 5, 3000, []byte("TestDevice\x00\x00\x00\x00\x00\x00")),
			expected: DeviceInfo{
				MajorVersion: 2,
				MinorVersion: 5,
				VersionBuild: 3000,
				DeviceName:   "TestDevice",
			},
			expectError: false,
		},
		{
			name: "Valid - empty response (older runtimes)",
			data: []byte{},
			expected: DeviceInfo{
				MajorVersion: 0,
				MinorVersion: 0,
				VersionBuild: 0,
				DeviceName:   "",
			},
			expectError: false,
		},
		{
			name:        "Invalid - less than 20 bytes",
			data:        []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			expected:    DeviceInfo{},
			expectError: true,
			errorType:   ErrInvalidDeviceInfoLength,
		},
		{
			name:        "Invalid - exactly 19 bytes",
			data:        make([]byte, 19),
			expected:    DeviceInfo{},
			expectError: true,
			errorType:   ErrInvalidDeviceInfoLength,
		},
		{
			name: "Valid - exactly 20 bytes",
			data: buildDeviceInfoData(1, 2, 3, ""),
			expected: DeviceInfo{
				MajorVersion: 1,
				MinorVersion: 2,
				VersionBuild: 3,
				DeviceName:   "",
			},
			expectError: false,
		},
		{
			name: "Valid - more than 20 bytes (ignores extra)",
			data: append(buildDeviceInfoData(4, 5, 6000, "PLC"), 0xFF, 0xFF, 0xFF),
			expected: DeviceInfo{
				MajorVersion: 4,
				MinorVersion: 5,
				VersionBuild: 6000,
				DeviceName:   "PLC",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseDeviceInfo(tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			assert.Equal(t, tt.expected.MajorVersion, info.MajorVersion)
			assert.Equal(t, tt.expected.MinorVersion, info.MinorVersion)
			assert.Equal(t, tt.expected.VersionBuild, info.VersionBuild)
			assert.Equal(t, tt.expected.DeviceName, info.DeviceName)
		})
	}
}

func TestCheckDeviceInfo(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
		errorType   error
	}{
		{
			name:        "Valid - exactly 20 bytes",
			data:        buildDeviceInfoData(1, 2, 3, "Test"),
			expectError: false,
		},
		{
			name:        "Valid - more than 20 bytes",
			data:        append(buildDeviceInfoData(1, 2, 3, "Test"), 0xFF, 0xFF),
			expectError: false,
		},
		{
			name:        "Valid - empty data (older runtimes)",
			data:        []byte{},
			expectError: false,
		},
		{
			name:        "Invalid - less than 20 bytes",
			data:        []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			expectError: true,
			errorType:   ErrInvalidDeviceInfoLength,
		},
		{
			name:        "Invalid - exactly 19 bytes",
			data:        make([]byte, 19),
			expectError: true,
			errorType:   ErrInvalidDeviceInfoLength,
		},
		{
			name:        "Invalid - exactly 10 bytes",
			data:        make([]byte, 10),
			expectError: true,
			errorType:   ErrInvalidDeviceInfoLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckDeviceInfo(tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
	}
}

// Helper function to build system state data for testing
func buildSystemStateData(adsState types.ADSState, deviceState uint16) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint16(data[0:2], uint16(adsState))
	binary.LittleEndian.PutUint16(data[2:4], deviceState)
	return data
}

// Helper function to build device info data for testing
func buildDeviceInfoData(major uint8, minor uint8, build uint16, name string) []byte {
	data := make([]byte, 20)
	data[0] = major
	data[1] = minor
	binary.LittleEndian.PutUint16(data[2:4], build)
	// Copy name (up to 16 bytes), rest stays null-padded
	copy(data[4:20], []byte(name))
	return data
}

// Helper function to build device info data with raw name bytes for testing
func buildDeviceInfoDataRaw(major uint8, minor uint8, build uint16, nameBytes []byte) []byte {
	data := make([]byte, 20)
	data[0] = major
	data[1] = minor
	binary.LittleEndian.PutUint16(data[2:4], build)
	// Copy exactly 16 bytes for name field
	copy(data[4:20], nameBytes)
	return data
}

// Helper function to build extended system state data for testing
func buildExtendedSystemStateData(adsState types.ADSState, deviceState, restartIndex uint16, version, revision uint8, build uint16, platform, osType uint8, flags uint16) []byte {
	data := make([]byte, 16)
	binary.LittleEndian.PutUint16(data[0:2], uint16(adsState))
	binary.LittleEndian.PutUint16(data[2:4], deviceState)
	binary.LittleEndian.PutUint16(data[4:6], restartIndex)
	data[6] = version
	data[7] = revision
	binary.LittleEndian.PutUint16(data[8:10], build)
	data[10] = platform
	data[11] = osType
	binary.LittleEndian.PutUint16(data[12:14], flags)
	// Bytes 14-15 are reserved (zeros)
	return data
}
