package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// ReadRaw reads raw data from the ADS server.
func (c *Client) ReadRaw(port uint16, indexGroup uint32, indexOffset uint32, size uint32) ([]byte, error) {
	c.logger.Debug("ReadRaw: Reading raw data", "indexGroup", indexGroup, "indexOffset", indexOffset, "size", size)

	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, indexGroup)
	binary.Write(data, binary.LittleEndian, indexOffset)
	binary.Write(data, binary.LittleEndian, size)

	req := AdsCommandRequest{
		Command:     types.ADSCommandRead,
		Data:        data.Bytes(),
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  port,
	}

	response, err := c.send(req)
	if err != nil {
		return nil, fmt.Errorf("ReadRaw: failed to send ADS command: %w", err)
	}

	return response, nil
}

// WriteRaw writes raw data to the ADS server.
func (c *Client) WriteRaw(port uint16, indexGroup uint32, indexOffset uint32, data []byte) ([]byte, error) {
	c.logger.Debug("WriteRaw: Writing raw data", "indexGroup", indexGroup, "indexOffset", indexOffset, "size", len(data))

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, indexGroup)
	binary.Write(buf, binary.LittleEndian, indexOffset)
	binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	buf.Write(data)

	req := AdsCommandRequest{
		Command:     types.ADSCommandWrite,
		Data:        buf.Bytes(),
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  port,
	}

	res, err := c.send(req)
	if err != nil {
		return nil, fmt.Errorf("WriteRaw: failed to send ADS command: %w", err)
	}

	return res, nil
}

// ReadWriteRaw reads and writes raw data to the ADS server.
func (c *Client) ReadWriteRaw(port uint16, indexGroup uint32, indexOffset uint32, readLength uint32, writeData []byte) ([]byte, error) {
	c.logger.Debug("ReadWriteRaw: Reading and writing raw data", "indexGroup", indexGroup, "indexOffset", indexOffset, "readLength", readLength, "writeDataSize", len(writeData))

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, indexGroup)
	binary.Write(buf, binary.LittleEndian, indexOffset)
	binary.Write(buf, binary.LittleEndian, readLength)
	binary.Write(buf, binary.LittleEndian, uint32(len(writeData)))
	buf.Write(writeData)

	req := AdsCommandRequest{
		Command:     types.ADSCommandReadWrite,
		Data:        buf.Bytes(),
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  port,
	}

	response, err := c.send(req)
	if err != nil {
		return nil, fmt.Errorf("ReadWriteRaw: failed to send ADS command: %w", err)
	}

	return response, nil
}
