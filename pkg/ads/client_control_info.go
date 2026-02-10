package ads

import (
	"fmt"

	adserrors "github.com/jarmocluyse/ads-go/pkg/ads/ads-errors"
	adsstateinfo "github.com/jarmocluyse/ads-go/pkg/ads/ads-stateinfo"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// ReadTcSystemState reads the TwinCAT system state.
func (c *Client) ReadTcSystemState() (*adsstateinfo.SystemState, error) {
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

	state, err := adsstateinfo.ParseSystemState(payload)
	if err != nil {
		c.logger.Error("ReadTcSystemState: Failed to parse system state", "error", err)
		return nil, err
	}

	c.logger.Info("ReadTcSystemState: Successfully parsed TwinCAT system state response", "response", state)
	return &state, nil
}

// ReadTcSystemExtendedState reads the TwinCAT extended system state.
//
// Extended state includes additional information beyond the basic state, most importantly
// the RestartIndex which increments each time the TwinCAT system service restarts.
// This allows detection of system restarts even when AdsState remains "Run".
//
// Extended state is read from ADS port 10000 (system service), IndexGroup 240, IndexOffset 0.
//
// Not all TwinCAT versions support extended state. If the read fails, use ReadTcSystemState()
// as a fallback to get basic state information.
//
// Example:
//
//	extState, err := client.ReadTcSystemExtendedState()
//	if err != nil {
//	    // Fall back to basic state
//	    basicState, err := client.ReadTcSystemState()
//	    // ...
//	}
func (c *Client) ReadTcSystemExtendedState() (*adsstateinfo.ExtendedSystemState, error) {
	c.logger.Debug("ReadTcSystemExtendedState: Reading TwinCAT extended system state.")

	// Extended state is read from system service port (10000)
	// IndexGroup 240, IndexOffset 0, Size 16 bytes
	// We can use ReadRaw which handles the request building for us
	data, err := c.ReadRaw(10000, 240, 0, 16)
	if err != nil {
		c.logger.Error("ReadTcSystemExtendedState: Failed to read extended state", "error", err)
		return nil, err
	}

	c.logger.Debug("ReadTcSystemExtendedState: Received raw response data", "length", len(data), "data", fmt.Sprintf("%x", data))

	// Parse extended state
	state, err := adsstateinfo.ParseExtendedSystemState(data)
	if err != nil {
		c.logger.Error("ReadTcSystemExtendedState: Failed to parse extended system state", "error", err)
		return nil, err
	}

	c.logger.Info("ReadTcSystemExtendedState: Successfully parsed TwinCAT extended system state",
		"adsState", state.AdsState.String(),
		"deviceState", state.DeviceState,
		"restartIndex", state.RestartIndex,
		"version", fmt.Sprintf("%d.%d.%d", state.Version, state.Revision, state.Build))

	return &state, nil
}

// ReadDeviceInfo reads the device information.
func (c *Client) ReadDeviceInfo() (*adsstateinfo.DeviceInfo, error) {
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
		c.logger.Error("ReadDeviceInfo: Invalid response length", "length", len(data), "expected", "at least 4")
		return nil, fmt.Errorf("invalid response length: %d", len(data))
	}
	payload, err := adserrors.StripAdsError(data)
	if err != nil {
		c.logger.Error("ReadDeviceInfo: ADS error received", "error", err)
		return nil, err
	}

	info, err := adsstateinfo.ParseDeviceInfo(payload)
	if err != nil {
		c.logger.Error("ReadDeviceInfo: Failed to parse device info", "error", err)
		return nil, err
	}

	c.logger.Info("ReadDeviceInfo: Successfully parsed device info", "deviceInfo", info)
	return &info, nil
}
