package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

	adserrors "github.com/jarmocluyse/ads-go/pkg/ads/ads-errors"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// AdsTcSystemStateResponse represents the TwinCAT system state response.
type AdsTcSystemStateResponse struct {
	AdsState    types.ADSState
	DeviceState uint16
}

// ReadTcSystemState reads the TwinCAT system state.
func (c *Client) ReadTcSystemState() (*AdsTcSystemStateResponse, error) {
	c.logger.Debug("ReadTcSystemState: Reading TwinCAT system state.")

	req := AdsCommandRequest{
		Command:    types.ADSCommandReadState,
		TargetPort: types.ADSReservedPortSystemService, // Explicitly target SystemService port
		Data:       []byte{},
	}
	data, err := c.send(req)
	if err != nil {
		c.logger.Error("ReadTcSystemState: Failed to send ReadState command", "error", err)
		return nil, err
	}

	c.logger.Debug("ReadTcSystemState: Received raw response data", "length", len(data), "data", fmt.Sprintf("%x", data))
	if len(data) < 8 {
		c.logger.Error("ReadTcSystemState: Invalid response length", "length", len(data), "expected", "at least 8")
		return nil, fmt.Errorf("invalid response length: %d", len(data))
	}

	payload, err := adserrors.StripAdsError(data)
	if err != nil {
		c.logger.Error("ReadTcSystemState: ADS error received", "error", err)
		return nil, err
	}

	adsState := binary.LittleEndian.Uint16(payload[0:2])
	deviceState := binary.LittleEndian.Uint16(payload[2:4])
	respState := &AdsTcSystemStateResponse{
		AdsState:    types.ADSState(adsState),
		DeviceState: deviceState,
	}
	c.logger.Info("ReadTcSystemState: Successfully parsed TwinCAT system state response", "response", respState)
	return respState, nil
}

// AdsReadDeviceInfoResponse represents the response for an ADS ReadDeviceInfo command.
type AdsReadDeviceInfoResponse struct {
	MajorVersion uint8
	MinorVersion uint8
	VersionBuild uint16
	DeviceName   string
}

// ReadDeviceInfo reads the device information.
func (c *Client) ReadDeviceInfo() (*AdsReadDeviceInfoResponse, error) {
	c.logger.Info("ReadDeviceInfo: Sending ReadDeviceInfo command.")
	req := AdsCommandRequest{
		Command:    types.ADSCommandReadDeviceInfo,
		TargetPort: types.ADSReservedPortSystemService, // Explicitly target SystemService port
		Data:       []byte{},
	}
	data, err := c.send(req)
	if err != nil {
		c.logger.Error("ReadDeviceInfo: Failed to send ReadDeviceInfo command", "error", err)
		return nil, err
	}
	c.logger.Debug("ReadDeviceInfo: Received raw response data", "length", len(data), "data", fmt.Sprintf("%x", data))

	if len(data) < 4 {
		c.logger.Error("ReadTcSystemState: Invalid response length", "length", len(data), "expected", "at least 4")
		return nil, fmt.Errorf("invalid response length: %d", len(data))
	}
	payload, err := adserrors.StripAdsError(data)
	if err != nil {
		c.logger.Error("ReadDeviceInfo: ADS error received", "error", err)
		return nil, err
	}

	// Handle cases where only the error code is returned (e.g., older/embedded runtimes)
	if len(payload) == 0 {
		return &AdsReadDeviceInfoResponse{
				MajorVersion: 0,
				MinorVersion: 0,
				VersionBuild: 0,
				DeviceName:   "",
			},
			nil
	}

	if len(payload) < 20 {
		c.logger.Error("ReadDeviceInfo: Invalid response length for device info", "length", len(payload), "expected", "at least 20")
		return nil, fmt.Errorf("invalid response length for device info: %d", len(payload))
	}
	resp := &AdsReadDeviceInfoResponse{
		MajorVersion: payload[0],
		MinorVersion: payload[1],
		VersionBuild: binary.LittleEndian.Uint16(payload[2:4]),
		DeviceName:   string(bytes.Trim(payload[4:20], "\x00")),
	}
	c.logger.Info("ReadDeviceInfo: Successfully parsed device info", "deviceInfo", resp)
	return resp, nil
}
