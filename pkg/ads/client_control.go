package ads

import (
	"encoding/binary"
	"fmt"

	adserrors "github.com/jarmocluyse/ads-go/pkg/ads/ads-errors"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// WriteControl writes control data to an ADS device.
func (c *Client) WriteControl(adsState types.ADSState, deviceState uint16, targetPort uint16) error {
	c.logger.Debug("WriteControl: Setting ADS state", "adsState", adsState.String(), "deviceState", fmt.Sprintf("0x%x", deviceState))
	data := make([]byte, 8)
	binary.LittleEndian.PutUint16(data[0:2], uint16(adsState)) // stae to send
	binary.LittleEndian.PutUint16(data[2:4], deviceState)      // device state to send (normally 0)
	binary.LittleEndian.PutUint32(data[4:8], 0)                // DataLength, 0 for state change

	req := AdsCommandRequest{
		Command:    types.ADSCommandWriteControl,
		TargetPort: targetPort,
		Data:       data,
	}
	respData, err := c.send(req)
	if err != nil {
		c.logger.Error("WriteControl: Failed to send WriteControl command", "error", err)
		return err
	}
	c.logger.Debug("WriteControl: Received raw response data", "length", len(respData))

	if len(respData) < 4 {
		c.logger.Error("WriteControl: Invalid response length", "length", len(respData), "expected", "at least 4")
		return fmt.Errorf("invalid response length: %d", len(respData))
	}
	if err := adserrors.CheckAdsError(respData[0:4]); err != nil {
		c.logger.Error("WriteControl: ADS error received", "error", err)
		return err
	}
	c.logger.Info("WriteControl: Successfully wrote control")
	return nil
}
