package ads

import (
	"fmt"

	adserrors "github.com/jarmocluyse/ads-go/pkg/ads/ads-errors"
	adsheader "github.com/jarmocluyse/ads-go/pkg/ads/ads-header"
	adsrequests "github.com/jarmocluyse/ads-go/pkg/ads/ads-requests"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// ReadRaw reads raw data from the ADS server.
func (c *Client) ReadRaw(port uint16, indexGroup uint32, indexOffset uint32, size uint32) ([]byte, error) {
	c.logger.Debug("ReadRaw: Reading raw data", "port", port, "indexGroup", indexGroup, "indexOffset", indexOffset, "size", size)

	payload := adsrequests.BuildReadRequest(indexGroup, indexOffset, size)

	req := AdsCommandRequest{
		Command:    types.ADSCommandRead,
		TargetPort: port,
		Data:       payload,
	}
	response, err := c.send(req)
	if err != nil {
		return nil, fmt.Errorf("ReadRaw: failed to send ADS command: %w", err)
	}
	payload, err = adsheader.StripAdsHeader(response)
	if err != nil {
		c.logger.Error("ReadRaw: ADS header error", "port", port, "indexGroup", indexGroup, "indexOffset", indexOffset, "error", err)
		return nil, err
	}
	return payload, nil
}

// WriteRaw writes raw data to the ADS server.
func (c *Client) WriteRaw(port uint16, indexGroup uint32, indexOffset uint32, data []byte) error {
	c.logger.Debug("WriteRaw: Writing raw data", "port", port, "indexGroup", indexGroup, "indexOffset", indexOffset, "size", len(data))

	payload := adsrequests.BuildWriteRequest(indexGroup, indexOffset, data)

	req := AdsCommandRequest{
		Command:    types.ADSCommandWrite,
		TargetPort: port,
		Data:       payload,
	}
	res, err := c.send(req)
	if err != nil {
		return fmt.Errorf("WriteRaw: failed to send ADS command: %w", err)
	}
	_, err = adserrors.StripAdsError(res)
	if err != nil {
		c.logger.Error("WriteRaw: ADS error received", "port", port, "indexGroup", indexGroup, "indexOffset", indexOffset, "error", err)
		return err
	}

	return nil
}

// ReadWriteRawBinary sends an ADS ReadWrite with exact binary write data (no
// null terminator appended). Use this for binary-protocol index groups such as
// the ADS logger consumer registration (IG 0x0000F090).
func (c *Client) ReadWriteRawBinary(port uint16, indexGroup uint32, indexOffset uint32, readLength uint32, writeData []byte) ([]byte, error) {
	c.logger.Debug("ReadWriteRawBinary: Reading and writing binary data", "port", port, "indexGroup", indexGroup, "indexOffset", indexOffset, "readLength", readLength, "writeLength", len(writeData))

	payload := adsrequests.BuildReadWriteRequest(indexGroup, indexOffset, readLength, writeData)

	req := AdsCommandRequest{
		Command:    types.ADSCommandReadWrite,
		TargetPort: port,
		Data:       payload,
	}
	response, err := c.send(req)
	if err != nil {
		return nil, fmt.Errorf("ReadWriteRawBinary: failed to send ADS command: %w", err)
	}

	data, err := adsheader.StripAdsHeader(response)
	if err != nil {
		c.logger.Error("ReadWriteRawBinary: ADS header error", "port", port, "indexGroup", indexGroup, "indexOffset", indexOffset, "error", err)
		return nil, err
	}
	return data, nil
}

// ReadWriteRaw reads and writes raw data to the ADS server.
func (c *Client) ReadWriteRaw(port uint16, indexGroup uint32, indexOffset uint32, readLength uint32, writeData []byte) ([]byte, error) {
	c.logger.Debug("ReadWriteRaw: Reading and writing raw data", "port", port, "indexGroup", indexGroup, "indexOffset", indexOffset, "readLength", readLength, "writeLength", len(writeData))

	payload := adsrequests.BuildReadWriteRequestWithNullTerminator(indexGroup, indexOffset, readLength, writeData)

	req := AdsCommandRequest{
		Command:    types.ADSCommandReadWrite,
		TargetPort: port,
		Data:       payload,
	}
	response, err := c.send(req)
	if err != nil {
		return nil, fmt.Errorf("ReadWriteRaw: failed to send ADS command: %w", err)
	}

	data, err := adsheader.StripAdsHeader(response)
	if err != nil {
		c.logger.Error("ReadWriteRaw: ADS header error", "port", port, "indexGroup", indexGroup, "indexOffset", indexOffset, "error", err)
		return nil, err
	}
	return data, nil
}
