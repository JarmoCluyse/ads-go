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
		c.logger.Error("ReadTcSystemState: Invalid response length", "length", len(data), "expected", "at least 4")
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
