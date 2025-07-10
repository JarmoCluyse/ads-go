package ads

import (
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// SetTcSystemToConfig sets the TwinCAT system to config mode.
func (c *Client) SetTcSystemToConfig() error {
	c.logger.Debug("SetTcSystemToConfig: Setting TwinCAT system to config mode.")

	// Read current state to ensure it's not already in config mode or other unexpected state
	currentState, err := c.ReadPlcRuntimeState()
	if err != nil {
		c.logger.Error("SetTcSystemToConfig: Failed to read current PLC runtime state", "error", err)
		return fmt.Errorf("failed to read current PLC runtime state: %w", err)
	}
	if currentState.ADSState == types.ADSStateConfig {
		c.logger.Info("SetTcSystemToConfig: TwinCAT system is already in config mode.")
		return nil
	}

	// Set ADS state to Config
	resp, err := c.WriteControl(types.ADSStateReconfig, currentState.DeviceState, types.ADSReservedPortSystemService)
	if err != nil {
		c.logger.Error("SetTcSystemToConfig: Failed to send WriteControl command", "error", err)
		return err
	}
	if resp.ErrorCode != 0 {
		c.logger.Error("SetTcSystemToConfig: WriteControl command returned ADS error", "errorCode", fmt.Sprintf("0x%x", resp.ErrorCode))
		return fmt.Errorf("WriteControl command returned ADS error: 0x%x", resp.ErrorCode)
	}

	c.logger.Info("SetTcSystemToConfig: TwinCAT system successfully set to config mode.")
	return nil
}

// SetTcSystemToRun sets the TwinCAT system to run mode.
func (c *Client) SetTcSystemToRun() error {
	c.logger.Info("SetTcSystemToRun: Setting TwinCAT system to run mode")
	if c.conn == nil {
		return fmt.Errorf("SetTcSystemToRun: Client is not connected. Use Connect() to connect to the target first.")
	}

	// Reading device state first as we don't want to change it (even though it's most probably 0)
	state, err := c.ReadPlcRuntimeState()
	if err != nil {
		c.logger.Error("SetTcSystemToRun: Failed to read current PLC runtime state", "error", err)
		return fmt.Errorf("failed to read current PLC runtime state: %w", err)
	}

	resp, err := c.WriteControl(types.ADSStateReset, state.DeviceState, types.ADSReservedPortSystemService)
	if err != nil {
		c.logger.Error("SetTcSystemToRun: Failed to send WriteControl command", "error", err)
		return err
	}
	if resp.ErrorCode != 0 {
		c.logger.Error("SetTcSystemToRun: WriteControl command returned ADS error", "errorCode", fmt.Sprintf("0x%x", resp.ErrorCode))
		return fmt.Errorf("WriteControl command returned ADS error: 0x%x", resp.ErrorCode)
	}

	c.logger.Info("SetTcSystemToRun: TwinCAT system successfully set to run mode.")
	return nil
}
