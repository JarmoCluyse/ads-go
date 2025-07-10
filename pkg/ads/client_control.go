package ads

import (
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// AdsWriteControlResponse represents the response for an ADS WriteControl command.
type AdsWriteControlResponse struct {
	ErrorCode uint32
}

// WriteControl writes control data to an ADS device.
// It can optionally target a specific ADS port. If targetPort is 0, c.settings.TargetPort is used.
func (c *Client) WriteControl(adsState types.ADSState, deviceState uint16, targetPort ...uint16) (*AdsWriteControlResponse, error) {
	port := c.settings.TargetPort
	if len(targetPort) > 0 && targetPort[0] != 0 {
		port = targetPort[0]
	}

	c.logger.Debug("WriteControl: Setting ADS state", "adsState", adsState.String(), "deviceState", fmt.Sprintf("0x%x", deviceState), "targetPort", port)
	data := make([]byte, 8)
	binary.LittleEndian.PutUint16(data[0:2], uint16(adsState))
	binary.LittleEndian.PutUint16(data[2:4], deviceState)
	binary.LittleEndian.PutUint32(data[4:8], 0) // DataLength, 0 for state change

	req := AdsCommandRequest{
		Command:     types.ADSCommandWriteControl,
		Data:        data,
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  port,
	}
	respData, err := c.send(req)
	if err != nil {
		c.logger.Error("WriteControl: Failed to send WriteControl command", "error", err)
		return nil, err
	}
	c.logger.Debug("WriteControl: Received raw response data", "length", len(respData))

	if len(respData) < 4 {
		c.logger.Error("WriteControl: Invalid response length", "length", len(respData), "expected", "at least 4")
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	resp := &AdsWriteControlResponse{
		ErrorCode: binary.LittleEndian.Uint32(respData[0:4]),
	}
	c.logger.Info("WriteControl: Successfully parsed write control response", "errorCode", fmt.Sprintf("0x%x", resp.ErrorCode))
	return resp, nil
}
