package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
	"github.com/jarmoCluyse/ads-go/pkg/ads/utils"
)

// GetSymbol retrieves information about a symbol from the ADS server.
func (c *Client) GetSymbol(path string) (types.AdsSymbol, error) {
	c.logger.Debug("GetSymbol: Requested symbol", "path", path)

	// Create the request data
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, types.ADSReservedIndexGroupSymbolInfoByNameEx)
	binary.Write(data, binary.LittleEndian, uint32(0))
	binary.Write(data, binary.LittleEndian, uint32(0xFFFFFFFF))
	binary.Write(data, binary.LittleEndian, uint32(len(path)+1))
	data.WriteString(path)

	req := AdsCommandRequest{
		Command:     types.ADSCommandReadWrite,
		Data:        data.Bytes(),
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}

	response, err := c.send(req)
	if err != nil {
		c.logger.Error("GetSymbol: Failed to send ADS command", "error", err)
		return types.AdsSymbol{}, fmt.Errorf("GetSymbol: failed to send ADS command: %w", err)
	}
	c.logger.Debug("GetSymbol: Full response received from ADS server", "response", fmt.Sprintf("%x", response), "length", len(response))

	symbol, err := c.parseAdsSymbol(response[4:])
	if err != nil {
		c.logger.Error("GetSymbol: Failed to parse symbol from response", "error", err)
		return types.AdsSymbol{}, fmt.Errorf("GetSymbol: failed to parse symbol from response: %w", err)
	}

	c.logger.Debug("GetSymbol: Symbol read and parsed", "path", path)
	return symbol, nil
}

func (c *Client) parseAdsSymbol(data []byte) (types.AdsSymbol, error) {
	var symbol types.AdsSymbol
	reader := bytes.NewReader(data)

	// Log the raw data being parsed and its length
	c.logger.Debug("parseAdsSymbol: Raw data for parsing", "data", fmt.Sprintf("%x", data), "length", len(data))

	if err := binary.Read(reader, binary.LittleEndian, &symbol.IndexGroup); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read IndexGroup: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read IndexGroup", "value", symbol.IndexGroup, "remainingBytes", reader.Len())

	if err := binary.Read(reader, binary.LittleEndian, &symbol.IndexOffset); err != nil {
		c.logger.Error("parseAdsSymbol: Failed to read IndexOffset", "error", err, "remainingBytes", reader.Len())
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read IndexOffset: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read IndexOffset", "value", symbol.IndexOffset, "remainingBytes", reader.Len())

	if err := binary.Read(reader, binary.LittleEndian, &symbol.Size); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Size: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read Size", "value", symbol.Size, "remainingBytes", reader.Len())

	if err := binary.Read(reader, binary.LittleEndian, &symbol.DataType); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read DataType: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read DataType", "value", symbol.DataType, "remainingBytes", reader.Len())

	if err := binary.Read(reader, binary.LittleEndian, &symbol.Flags); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Flags: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read Flags", "value", symbol.Flags, "remainingBytes", reader.Len())

	if err := binary.Read(reader, binary.LittleEndian, &symbol.NameLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read NameLength: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read NameLength", "value", symbol.NameLength, "remainingBytes", reader.Len())

	if err := binary.Read(reader, binary.LittleEndian, &symbol.TypeLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read TypeLength: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read TypeLength", "value", symbol.TypeLength, "remainingBytes", reader.Len())

	if err := binary.Read(reader, binary.LittleEndian, &symbol.CommentLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read CommentLength: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: read CommentLength", "value", symbol.CommentLength, "remainingBytes", reader.Len())

	name := make([]byte, symbol.NameLength)
	if err := binary.Read(reader, binary.LittleEndian, &name); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Name: %w (remaining bytes: %d)", err, reader.Len())
	}
	symbol.Name = utils.DecodePlcStringBuffer(name, c.settings.AdsSymbolsUseUtf8)

	typeName := make([]byte, symbol.TypeLength)
	if err := binary.Read(reader, binary.LittleEndian, &typeName); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read TypeName: %w (remaining bytes: %d)", err, reader.Len())
	}
	symbol.Type = utils.DecodePlcStringBuffer(typeName, c.settings.AdsSymbolsUseUtf8)

	comment := make([]byte, symbol.CommentLength)
	if err := binary.Read(reader, binary.LittleEndian, &comment); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Comment: %w (remaining bytes: %d)", err, reader.Len())
	}
	symbol.Comment = utils.DecodePlcStringBuffer(comment, c.settings.AdsSymbolsUseUtf8)

	return symbol, nil
}
