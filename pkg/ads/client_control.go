package ads

import (
	"fmt"

	adserrors "github.com/jarmocluyse/ads-go/pkg/ads/ads-errors"
	adsrequests "github.com/jarmocluyse/ads-go/pkg/ads/ads-requests"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// WriteControl writes control data to an ADS device.
func (c *Client) WriteControl(adsState types.ADSState, deviceState uint16, targetPort uint16) error {
	c.logger.Debug("WriteControl: Setting ADS state", "adsState", adsState.String(), "deviceState", fmt.Sprintf("0x%x", deviceState))

	payload := adsrequests.BuildWriteControlRequest(uint16(adsState), deviceState)

	req := AdsCommandRequest{
		Command:    types.ADSCommandWriteControl,
		TargetPort: targetPort,
		Data:       payload,
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
