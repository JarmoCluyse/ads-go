package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// AdsReadStateResponse represents the response of a ReadState command.
type AdsReadStateResponse struct {
	ErrorCode   uint32
	AdsState    types.ADSState
	DeviceState uint16
}

// ReadState sends an AdsCommandReadState (command ID = 4) and returns the PLC’s ADS
// and device states.  Older/embedded runtimes may reply with **only the 4-byte
// result code** (0 = OK).  We accept that short frame and return zeros for the
// state fields so the behaviour matches the Bun implementation.
func (c *Client) ReadState() (*AdsReadStateResponse, error) {
	c.logger.Info("ReadState: Reading ADS state.")
	req := AdsCommandRequest{
		Command:     types.ADSCommandReadState,
		Data:        []byte{},
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}
	data, err := c.send(req)
	if err != nil {
		c.logger.Error("ReadState: Failed to send ReadState command", "error", err)
		return nil, err
	}

	c.logger.Debug("ReadState: Received raw response data", "length", len(data), "data", fmt.Sprintf("%x", data))

	if len(data) < 4 {
		c.logger.Error("ReadState: Invalid response length", "length", len(data), "expected", "at least 4")
		return nil, fmt.Errorf("invalid response length: %d", len(data))
	}

	/* NOTE: Legacy runtimes: reply is only the 4-byte result code (OK = 0) */
	if len(data) == 4 && binary.LittleEndian.Uint32(data[:4]) == 0 {
		return &AdsReadStateResponse{
				ErrorCode:   0,
				AdsState:    0, // unknown
				DeviceState: 0, // unknown
			},
			nil
	}
	/* NOTE: Normal case – need 8 bytes total */
	if len(data) < 8 {
		return nil, fmt.Errorf("state reply length %d (want 4 or ≥8)", len(data))
	}

	errorCode := binary.LittleEndian.Uint32(data[0:4])
	if errorCode != 0 {
		c.logger.Error("ReadState: ADS error received", "errorCode", fmt.Sprintf("0x%x", errorCode))
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	var adsState, deviceState uint16
	// The actual ADS state and device state are after the 4-byte error code
	if len(data) >= 8 { // ErrorCode (4 bytes) + ADS State (2 bytes) + Device State (2 bytes)
		adsState = binary.LittleEndian.Uint16(data[4:6])
		deviceState = binary.LittleEndian.Uint16(data[6:8])
		c.logger.Debug("ReadState: Parsed ADS State", "adsState", adsState, "deviceState", deviceState)
	} else {
		c.logger.Warn("ReadState: Not enough data in response to parse ADS and Device State", "length", len(data))
	}

	respState := &AdsReadStateResponse{
		ErrorCode:   errorCode,
		AdsState:    types.ADSState(adsState),
		DeviceState: deviceState,
	}
	c.logger.Info("ReadState: Successfully parsed ADS state response", "response", respState)
	return respState, nil
}

// AdsTcSystemStateResponse represents the TwinCAT system state response.
type AdsTcSystemStateResponse struct {
	ErrorCode   uint32
	AdsState    types.ADSState
	DeviceState uint16
}

// ReadTcSystemState reads the TwinCAT system state.
func (c *Client) ReadTcSystemState() (*AdsTcSystemStateResponse, error) {
	c.logger.Debug("ReadTcSystemState: Reading TwinCAT system state.")
	req := AdsCommandRequest{
		Command:     types.ADSCommandReadState,
		Data:        []byte{},
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  types.ADSReservedPortSystemService, // Explicitly target SystemService port
	}
	data, err := c.send(req)
	if err != nil {
		c.logger.Error("ReadTcSystemState: Failed to send ReadState command", "error", err)
		return nil, err
	}
	c.logger.Debug("ReadTcSystemState: Received raw response data", "length", len(data), "data", fmt.Sprintf("%x", data))

	if len(data) < 4 {
		c.logger.Error("ReadTcSystemState: Invalid response length", "length", len(data), "expected", "at least 4")
		return nil, fmt.Errorf("invalid response length: %d", len(data))
	}
	errorCode := binary.LittleEndian.Uint32(data[0:4])
	if errorCode != 0 {
		c.logger.Error("ReadTcSystemState: ADS error received", "errorCode", fmt.Sprintf("0x%x", errorCode))
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	var adsState, deviceState uint16
	// The actual ADS state and device state are after the 4-byte error code
	if len(data) >= 8 { // ErrorCode (4 bytes) + ADS State (2 bytes) + Device State (2 bytes)
		adsState = binary.LittleEndian.Uint16(data[4:6])
		deviceState = binary.LittleEndian.Uint16(data[6:8])
		c.logger.Debug("ReadTcSystemState: Parsed ADS State", "adsState", adsState, "deviceState", deviceState)
	} else {
		c.logger.Warn("ReadTcSystemState: Not enough data in response to parse ADS and Device State", "length", len(data))
	}

	respState := &AdsTcSystemStateResponse{
		ErrorCode:   errorCode,
		AdsState:    types.ADSState(adsState),
		DeviceState: deviceState,
	}
	c.logger.Info("ReadTcSystemState: Successfully parsed TwinCAT system state response", "response", respState)
	return respState, nil
}

// AdsStateResponse represents the ADS state response payload.
type AdsStateResponse struct {
	ADSState    types.ADSState
	DeviceState uint16
}

// ReadPlcRuntimeState reads the PLC runtime state.
func (c *Client) ReadPlcRuntimeState() (*AdsStateResponse, error) {
	c.logger.Debug("ReadPlcRuntimeState: Reading PLC runtime state.")
	resp, err := c.ReadTcSystemState()
	if err != nil {
		c.logger.Error("ReadPlcRuntimeState: Failed to read state", "error", err)
		return nil, err
	}
	c.logger.Info("ReadPlcRuntimeState: Successfully read PLC runtime state", "adsState", resp.AdsState)
	return &AdsStateResponse{
		ADSState:    types.ADSState(resp.AdsState),
		DeviceState: resp.DeviceState,
	}, nil
}

// AdsReadDeviceInfoResponse represents the response for an ADS ReadDeviceInfo command.
type AdsReadDeviceInfoResponse struct {
	ErrorCode    uint32
	MajorVersion uint8
	MinorVersion uint8
	VersionBuild uint16
	DeviceName   string
}

// ReadDeviceInfo reads the device information.
func (c *Client) ReadDeviceInfo() (*AdsReadDeviceInfoResponse, error) {
	c.logger.Info("ReadDeviceInfo: Sending ReadDeviceInfo command.")
	req := AdsCommandRequest{
		Command:     types.ADSCommandReadDeviceInfo,
		Data:        []byte{},
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}
	data, err := c.send(req)
	if err != nil {
		c.logger.Error("ReadDeviceInfo: Failed to send ReadDeviceInfo command", "error", err)
		return nil, err
	}
	c.logger.Debug("ReadDeviceInfo: Received raw response data", "length", len(data), "data", fmt.Sprintf("%x", data))

	// Handle cases where only the error code is returned (e.g., older/embedded runtimes)
	if len(data) == 4 {
		c.logger.Info("ReadDeviceInfo: Received only error code, returning default device info.")
		return &AdsReadDeviceInfoResponse{
				ErrorCode:    0,
				MajorVersion: 0,
				MinorVersion: 0,
				VersionBuild: 0,
				DeviceName:   "",
			},
			nil
	}

	if len(data) < 24 {
		c.logger.Error("ReadDeviceInfo: Invalid response length for device info", "length", len(data), "expected", "at least 24")
		return nil, fmt.Errorf("invalid response length for device info: %d", len(data))
	}

	resp := &AdsReadDeviceInfoResponse{
		MajorVersion: data[4],
		MinorVersion: data[5],
		VersionBuild: binary.LittleEndian.Uint16(data[6:8]),
		DeviceName:   string(bytes.Trim(data[8:24], "\x00")),
	}
	c.logger.Info("ReadDeviceInfo: Successfully parsed device info", "deviceInfo", resp)

	return resp, nil
}
