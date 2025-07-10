package ads

import (
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// AdsReadResponse represents the response for an ADS Read command.
type AdsReadResponse struct {
	ErrorCode uint32
	Length    uint32
	Data      []byte
}

// Read reads data from an ADS device.
func (c *Client) Read(indexGroup, indexOffset, length uint32) (*AdsReadResponse, error) {
	c.logger.Debug("Read: Reading data", "indexGroup", fmt.Sprintf("0x%x", indexGroup), "indexOffset", fmt.Sprintf("0x%x", indexOffset), "length", length)
	data := make([]byte, 12)
	binary.LittleEndian.PutUint32(data[0:4], indexGroup)
	binary.LittleEndian.PutUint32(data[4:8], indexOffset)
	binary.LittleEndian.PutUint32(data[8:12], length)

	req := AdsCommandRequest{
		Command:     types.ADSCommandRead,
		Data:        data,
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}
	respData, err := c.send(req)
	if err != nil {
		c.logger.Error("Read: Failed to send command", "error", err)
		return nil, err
	}
	c.logger.Debug("Read: Received raw response data", "length", len(respData))

	if len(respData) < 8 {
		c.logger.Error("Read: Invalid response length", "length", len(respData), "expected", "at least 8")
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	errorCode := binary.LittleEndian.Uint32(respData[0:4])
	if errorCode != 0 {
		c.logger.Error("Read: ADS error received", "errorCode", fmt.Sprintf("0x%x", errorCode))
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	readLength := binary.LittleEndian.Uint32(respData[4:8])
	c.logger.Debug("Read: Reported read length", "length", readLength)

	resp := &AdsReadResponse{
		ErrorCode: errorCode,
		Length:    readLength,
		Data:      respData[8 : 8+readLength],
	}
	c.logger.Info("Read: Successfully parsed read response", "dataLength", len(resp.Data))

	return resp, nil
}

// AdsWriteResponse represents the response of a Write command.
type AdsWriteResponse struct {
	ErrorCode uint32
}

// Write writes data to an ADS device.
func (c *Client) Write(indexGroup, indexOffset uint32, data []byte) (*AdsWriteResponse, error) {
	c.logger.Info("Write: Writing data", "indexGroup", fmt.Sprintf("0x%x", indexGroup), "indexOffset", fmt.Sprintf("0x%x", indexOffset), "dataLength", len(data))
	reqData := make([]byte, 12+len(data))
	binary.LittleEndian.PutUint32(reqData[0:4], indexGroup)
	binary.LittleEndian.PutUint32(reqData[4:8], indexOffset)
	binary.LittleEndian.PutUint32(reqData[8:12], uint32(len(data)))
	copy(reqData[12:], data)

	req := AdsCommandRequest{
		Command:     types.ADSCommandWrite,
		Data:        reqData,
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}
	respData, err := c.send(req)
	if err != nil {
		c.logger.Error("Write: Failed to send command", "error", err)
		return nil, err
	}
	c.logger.Debug("Write: Received raw response data", "length", len(respData))

	if len(respData) < 4 {
		c.logger.Error("Write: Invalid response length", "length", len(respData), "expected", "at least 4")
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	resp := &AdsWriteResponse{
		ErrorCode: binary.LittleEndian.Uint32(respData[0:4]),
	}
	c.logger.Info("Write: Successfully parsed write response", "errorCode", fmt.Sprintf("0x%x", resp.ErrorCode))

	return resp, nil
}

// AdsReadWriteResponse represents the response for an ADS ReadWrite command.
type AdsReadWriteResponse struct {
	ErrorCode uint32
	Length    uint32
	Data      []byte
}

// ReadWrite reads and writes data to an ADS device.
func (c *Client) ReadWrite(indexGroup, indexOffset, readLength uint32, dataToWrite []byte) (*AdsReadWriteResponse, error) {
	c.logger.Info("ReadWrite: Reading and writing data", "indexGroup", fmt.Sprintf("0x%x", indexGroup), "indexOffset", fmt.Sprintf("0x%x", indexOffset), "readLength", readLength, "dataLength", len(dataToWrite))
	reqData := make([]byte, 16+len(dataToWrite))
	binary.LittleEndian.PutUint32(reqData[0:4], indexGroup)
	binary.LittleEndian.PutUint32(reqData[4:8], indexOffset)
	binary.LittleEndian.PutUint32(reqData[8:12], readLength)
	binary.LittleEndian.PutUint32(reqData[12:16], uint32(len(dataToWrite)))
	copy(reqData[16:], dataToWrite)

	req := AdsCommandRequest{
		Command:     types.ADSCommandReadWrite,
		Data:        reqData,
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}
	respData, err := c.send(req)
	if err != nil {
		c.logger.Error("ReadWrite: Failed to send command", "error", err)
		return nil, err
	}
	c.logger.Debug("ReadWrite: Received raw response data", "length", len(respData))

	if len(respData) < 8 {
		c.logger.Error("ReadWrite: Invalid response length", "length", len(respData), "expected", "at least 8")
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	errorCode := binary.LittleEndian.Uint32(respData[0:4])
	if errorCode != 0 {
		c.logger.Error("ReadWrite: ADS error received", "errorCode", fmt.Sprintf("0x%x", errorCode))
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	length := binary.LittleEndian.Uint32(respData[4:8])
	c.logger.Debug("ReadWrite: Reported response length", "length", length)

	resp := &AdsReadWriteResponse{
		ErrorCode: errorCode,
		Length:    length,
		Data:      respData[8 : 8+length],
	}
	c.logger.Info("ReadWrite: Successfully parsed read/write response", "dataLength", len(resp.Data))

	return resp, nil
}
