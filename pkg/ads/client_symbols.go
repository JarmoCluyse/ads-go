package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// GetSymbol retrieves information about a symbol from the ADS server.
func (c *Client) GetSymbol(path string) (types.AdsSymbol, error) {
	c.logger.Debug("GetSymbol: Requested symbol", "path", path)

	// Check if the symbol is already cached
	if symbol, ok := c.plcSymbols[path]; ok {
		c.logger.Debug("GetSymbol: Symbol found in cache", "path", path)
		return symbol, nil
	}

	c.logger.Debug("GetSymbol: Symbol not cached, reading from target", "path", path)

	// Create the request data
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, types.ADSReservedIndexGroupSymbolInfoByNameEx)
	binary.Write(data, binary.LittleEndian, uint32(0))
	binary.Write(data, binary.LittleEndian, uint32(0xFFFFFFFF))
	binary.Write(data, binary.LittleEndian, uint32(len(path)+1))
	data.WriteString(path)
	data.WriteByte(0)

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

	// Check for minimum response length (4 bytes for error code + 4 bytes for data length)
	if len(response) < 8 {
		return types.AdsSymbol{}, fmt.Errorf("GetSymbol: invalid response length: %d, expected at least 8 bytes", len(response))
	}

	errorCode := binary.LittleEndian.Uint32(response[0:4])
	if errorCode != 0 {
		return types.AdsSymbol{}, fmt.Errorf("GetSymbol: ADS error received: 0x%x", errorCode)
	}

	// Read the length of the actual symbol data
	symbolDataLength := binary.LittleEndian.Uint32(response[4:8])
	c.logger.Debug("GetSymbol: Symbol data length from response", "length", symbolDataLength)

	// Ensure we have enough data for the symbol
	if len(response) < int(8+symbolDataLength) {
		return types.AdsSymbol{}, fmt.Errorf("GetSymbol: incomplete symbol data received. Expected %d bytes, got %d", 8+symbolDataLength, len(response))
	}

	// Pass only the actual symbol data to parseAdsSymbol
	symbol, err := c.parseAdsSymbol(response[8 : 8+symbolDataLength])
	if err != nil {
		c.logger.Error("GetSymbol: Failed to parse symbol from response", "error", err)
		return types.AdsSymbol{}, fmt.Errorf("GetSymbol: failed to parse symbol from response: %w", err)
	}

	// Cache the symbol
	c.plcSymbols[path] = symbol

	c.logger.Debug("GetSymbol: Symbol read and parsed", "path", path)
	return symbol, nil
}

func (c *Client) parseAdsSymbol(data []byte) (types.AdsSymbol, error) {
	var symbol types.AdsSymbol
	reader := bytes.NewReader(data)

	// Log the raw data being parsed and its length
	// c.logger.Debug("parseAdsSymbol: Raw data for parsing", "data", fmt.Sprintf("%x", data), "length", len(data))

	if err := binary.Read(reader, binary.LittleEndian, &symbol.IndexGroup); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read IndexGroup: %w (remaining bytes: %d)", err, reader.Len())
	}
	// c.logger.Debug("parseAdsSymbol: Read IndexGroup", "value", symbol.IndexGroup, "remainingBytes", reader.Len())

	c.logger.Debug("parseAdsSymbol: Attempting to read IndexOffset", "remainingBytes", reader.Len())
	if err := binary.Read(reader, binary.LittleEndian, &symbol.IndexOffset); err != nil {
		c.logger.Error("parseAdsSymbol: Failed to read IndexOffset", "error", err, "remainingBytes", reader.Len())
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read IndexOffset: %w (remaining bytes: %d)", err, reader.Len())
	}
	c.logger.Debug("parseAdsSymbol: Successfully read IndexOffset", "value", symbol.IndexOffset, "remainingBytes", reader.Len())
	// c.logger.Debug("parseAdsSymbol: Read IndexOffset", "value", symbol.IndexOffset, "remainingBytes", reader.Len())

	if err := binary.Read(reader, binary.LittleEndian, &symbol.Size); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Size: %w (remaining bytes: %d)", err, reader.Len())
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.DataType); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read DataType: %w (remaining bytes: %d)", err, reader.Len())
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.Flags); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Flags: %w (remaining bytes: %d)", err, reader.Len())
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.NameLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read NameLength: %w (remaining bytes: %d)", err, reader.Len())
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.TypeLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read TypeLength: %w (remaining bytes: %d)", err, reader.Len())
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.CommentLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read CommentLength: %w (remaining bytes: %d)", err, reader.Len())
	}

	name := make([]byte, symbol.NameLength)
	if err := binary.Read(reader, binary.LittleEndian, &name); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Name: %w (remaining bytes: %d)", err, reader.Len())
	}
	symbol.Name = string(name)

	typeName := make([]byte, symbol.TypeLength)
	if err := binary.Read(reader, binary.LittleEndian, &typeName); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read TypeName: %w (remaining bytes: %d)", err, reader.Len())
	}
	symbol.Type = string(typeName)

	comment := make([]byte, symbol.CommentLength)
	if err := binary.Read(reader, binary.LittleEndian, &comment); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Comment: %w (remaining bytes: %d)", err, reader.Len())
	}
	symbol.Comment = string(comment)

	return symbol, nil
}
