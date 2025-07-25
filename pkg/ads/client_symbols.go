package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
	"github.com/jarmoCluyse/ads-go/pkg/ads/utils"
)

// GetSymbol retrieves information about a symbol from the ADS server.
func (c *Client) GetSymbol(port uint16, path string) (*types.AdsSymbol, error) {
	c.logger.Debug("GetSymbol: Requested symbol", "path", path)

	// Create the request data
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, types.ADSReservedIndexGroupSymbolInfoByNameEx)
	binary.Write(data, binary.LittleEndian, uint32(0))
	binary.Write(data, binary.LittleEndian, uint32(0xFFFFFFFF))
	binary.Write(data, binary.LittleEndian, uint32(len(path)+1))
	data.Write(utils.EncodeStringToPlcStringBuffer(path, c.settings.AdsSymbolsUseUtf8))
	binary.Write(data, binary.LittleEndian, uint8(0))

	req := AdsCommandRequest{
		Command:    types.ADSCommandReadWrite,
		TargetPort: port,
		Data:       data.Bytes(),
	}
	response, err := c.send(req)
	if err != nil {
		c.logger.Error("GetSymbol: Failed to send ADS command", "error", err)
		return &types.AdsSymbol{}, fmt.Errorf("GetSymbol: failed to send ADS command: %w", err)
	}
	c.logger.Debug("GetSymbol: Full response received from ADS server", "response", fmt.Sprintf("%x", response), "length", len(response))
	errorCode := binary.LittleEndian.Uint32(response[0:4])
	if errorCode != 0 {
		errorString := types.ADSError[errorCode]
		c.logger.Error("ReadTcSystemState: ADS error received", "errorCode", fmt.Sprintf("0x%x", errorCode), "error", errorString)
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}
	symbolLen := binary.LittleEndian.Uint32(response[4:8])
	symbolData := response[8:]
	if len(symbolData) < int(symbolLen) {
		c.logger.Error("received to little data", "length", symbolLen, "receivedLen", len(symbolData))
		return nil, fmt.Errorf("received to little data")
	}
	symbol, err := c.parseAdsSymbol(symbolData)
	if err != nil {
		c.logger.Error("GetSymbol: Failed to parse symbol from response", "error", err)
		return &types.AdsSymbol{}, fmt.Errorf("GetSymbol: failed to parse symbol from response: %w", err)
	}

	c.logger.Debug("GetSymbol: Symbol read and parsed", "path", path)
	return &symbol, nil
}

func (c *Client) parseAdsSymbol(data []byte) (types.AdsSymbol, error) {
	var symbol types.AdsSymbol
	reader := bytes.NewReader(data)

	c.logger.Debug("parseAdsSymbol: Raw data for parsing", "data", fmt.Sprintf("%x", data), "length", len(data))

	// NOTE: I think beckhoff does something weird here,
	// Response is 0:4 -> ads error 4:8 -> data length 8: -> data
	// however the data 0:4 -> is again the length which makes no sense?
	var dataLen uint32
	if err := binary.Read(reader, binary.LittleEndian, &dataLen); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read dataLen: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read dataLen", "value", dataLen, "remainingBytes", reader.Len())

	// NOTE: 4:8 group
	if err := binary.Read(reader, binary.LittleEndian, &symbol.IndexGroup); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read IndexGroup: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read IndexGroup", "value", symbol.IndexGroup, "remainingBytes", reader.Len())

	// NOTE: 8:12 index
	if err := binary.Read(reader, binary.LittleEndian, &symbol.IndexOffset); err != nil {
		c.logger.Error("parseAdsSymbol: Failed to read IndexOffset", "error", err, "remainingBytes", reader.Len())
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read IndexOffset: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read IndexOffset", "value", symbol.IndexOffset, "remainingBytes", reader.Len())

	// NOTE: 12:16 size
	if err := binary.Read(reader, binary.LittleEndian, &symbol.Size); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Size: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read Size", "value", symbol.Size, "remainingBytes", reader.Len())

	// NOTE: 16:20 data type
	if err := binary.Read(reader, binary.LittleEndian, &symbol.DataType); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read DataType: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read DataType", "value", symbol.DataType, "remainingBytes", reader.Len())

	// NOTE: 20:24 flags
	if err := binary.Read(reader, binary.LittleEndian, &symbol.Flags); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Flags: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read Flags", "value", symbol.Flags, "remainingBytes", reader.Len())

	// NOTE: 24:26 name length
	if err := binary.Read(reader, binary.LittleEndian, &symbol.NameLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read NameLength: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read NameLength", "value", symbol.NameLength, "remainingBytes", reader.Len())

	// NOTE: 26:28 type length
	if err := binary.Read(reader, binary.LittleEndian, &symbol.TypeLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read TypeLength: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read TypeLength", "value", symbol.TypeLength, "remainingBytes", reader.Len())

	// NOTE: 28:30 comment length
	if err := binary.Read(reader, binary.LittleEndian, &symbol.CommentLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read CommentLength: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read CommentLength", "value", symbol.CommentLength, "remainingBytes", reader.Len())

	// NOTE: .. name
	name := make([]byte, symbol.NameLength+1)
	if err := binary.Read(reader, binary.LittleEndian, &name); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Name: %w (remaining bytes: %d)", err, reader.Len())
	}
	symbol.Name = utils.DecodePlcStringBuffer(name, c.settings.AdsSymbolsUseUtf8)
	c.logger.Debug("parseAdsSymbol: read Name", "value", symbol.Name, "raw", name, "remainingBytes", reader.Len())

	// NOTE: .. name
	typeName := make([]byte, symbol.TypeLength+1)
	if err := binary.Read(reader, binary.LittleEndian, &typeName); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read TypeName: %w (remaining bytes: %d)", err, reader.Len())
	}
	symbol.Type = utils.DecodePlcStringBuffer(typeName, c.settings.AdsSymbolsUseUtf8)
	c.logger.Debug("parseAdsSymbol: read Type", "value", symbol.Type, "raw", typeName, "remainingBytes", reader.Len())

	// NOTE: .. comment
	comment := make([]byte, symbol.CommentLength)
	if err := binary.Read(reader, binary.LittleEndian, &comment); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Comment: %w (remaining bytes: %d)", err, reader.Len())
	}
	symbol.Comment = utils.DecodePlcStringBuffer(comment, c.settings.AdsSymbolsUseUtf8)
	c.logger.Debug("parseAdsSymbol: read Comment", "value", symbol.Comment, "raw", comment, "remainingBytes", reader.Len())

	// rest
	rest := make([]byte, reader.Len())
	if err := binary.Read(reader, binary.LittleEndian, &rest); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Comment: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read rest", "value", rest)
	return symbol, nil
}
