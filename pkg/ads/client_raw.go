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
	binary.Write(data, binary.LittleEndian, indexGroup)  // group
	binary.Write(data, binary.LittleEndian, indexOffset) // index
	binary.Write(data, binary.LittleEndian, size)        // size

	req := AdsCommandRequest{
		Command:    types.ADSCommandRead,
		TargetPort: port,
		Data:       data.Bytes(),
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
	binary.Write(buf, binary.LittleEndian, indexGroup)        // group
	binary.Write(buf, binary.LittleEndian, indexOffset)       // index
	binary.Write(buf, binary.LittleEndian, uint32(len(data))) // data length
	buf.Write(data)                                           // data

	req := AdsCommandRequest{
		Command:    types.ADSCommandWrite,
		TargetPort: port,
		Data:       buf.Bytes(),
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
	binary.Write(buf, binary.LittleEndian, indexGroup)               // group
	binary.Write(buf, binary.LittleEndian, indexOffset)              // index
	binary.Write(buf, binary.LittleEndian, readLength)               // read lenght
	binary.Write(buf, binary.LittleEndian, uint32(len(writeData)+1)) // write length
	buf.Write(writeData)                                             // write data
	binary.Write(buf, binary.LittleEndian, uint8(0))                 // write 0 after data

	req := AdsCommandRequest{
		Command:    types.ADSCommandReadWrite,
		TargetPort: port,
		Data:       buf.Bytes(),
	}
	response, err := c.send(req)
	if err != nil {
		return nil, fmt.Errorf("ReadWriteRaw: failed to send ADS command: %w", err)
	}

	errorCode := binary.LittleEndian.Uint32(response[0:4])
	if errorCode != 0 {
		errorString := types.ADSError[errorCode]
		c.logger.Error("ReadTcSystemState: ADS error received", "errorCode", fmt.Sprintf("0x%x", errorCode), "error", errorString)
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}
	Len := binary.LittleEndian.Uint32(response[4:8])
	Data := response[8:]
	if len(Data) < int(Len) {
		c.logger.Error("received to little data", "length", Len, "receivedLen", len(Data))
		return nil, fmt.Errorf("received to little data")
	}
	return Data, nil
}
