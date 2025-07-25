package ads

import (
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// SetTcSystemToConfig sets the TwinCAT system to config mode.
func (c *Client) SetTcSystemToConfig() error {
	c.logger.Debug("SetTcSystemToConfig: Setting TwinCAT system to config mode.")

	// Reading device state first as we don't want to change it (even though it's most probably 0)
	currentState, err := c.ReadTcSystemState()
	if err != nil {
		c.logger.Error("SetTcSystemToConfig: Failed to read current PLC runtime state", "error", err)
		return fmt.Errorf("failed to read current PLC runtime state: %w", err)
	}

	// Set ADS state to Config
	err = c.WriteControl(types.ADSStateReconfig, currentState.DeviceState, types.ADSReservedPortSystemService)
	if err != nil {
		c.logger.Error("SetTcSystemToConfig: Failed to send WriteControl command", "error", err)
		return err
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
	state, err := c.ReadTcSystemState()
	if err != nil {
		c.logger.Error("SetTcSystemToRun: Failed to read current PLC runtime state", "error", err)
		return fmt.Errorf("failed to read current PLC runtime state: %w", err)
	}

	err = c.WriteControl(types.ADSStateReset, state.DeviceState, types.ADSReservedPortSystemService)
	if err != nil {
		c.logger.Error("SetTcSystemToRun: Failed to send WriteControl command", "error", err)
		return err
	}

	c.logger.Info("SetTcSystemToRun: TwinCAT system successfully set to run mode.")
	return nil
}
