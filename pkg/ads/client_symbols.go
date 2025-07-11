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

	// The first 4 bytes of the response are the length of the data
	symbol, err := parseAdsSymbol(response[4:])
	if err != nil {
		c.logger.Error("GetSymbol: Failed to parse symbol from response", "error", err)
		return types.AdsSymbol{}, fmt.Errorf("GetSymbol: failed to parse symbol from response: %w", err)
	}

	// Cache the symbol
	c.plcSymbols[path] = symbol

	c.logger.Debug("GetSymbol: Symbol read and parsed", "path", path)
	return symbol, nil
}

func parseAdsSymbol(data []byte) (types.AdsSymbol, error) {
	var symbol types.AdsSymbol
	reader := bytes.NewReader(data)

	// Log the raw data being parsed
	// c.logger.Debug("parseAdsSymbol: Raw data for parsing", "data", fmt.Sprintf("%x", data))

	if err := binary.Read(reader, binary.LittleEndian, &symbol.IndexGroup); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read IndexGroup: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.IndexOffset); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read IndexOffset: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.Size); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Size: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.DataType); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read DataType: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.Flags); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Flags: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.NameLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read NameLength: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.TypeLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read TypeLength: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &symbol.CommentLength); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read CommentLength: %w", err)
	}

	name := make([]byte, symbol.NameLength)
	if err := binary.Read(reader, binary.LittleEndian, &name); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Name: %w", err)
	}
	symbol.Name = string(name)

	typeName := make([]byte, symbol.TypeLength)
	if err := binary.Read(reader, binary.LittleEndian, &typeName); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read TypeName: %w", err)
	}
	symbol.Type = string(typeName)

	comment := make([]byte, symbol.CommentLength)
	if err := binary.Read(reader, binary.LittleEndian, &comment); err != nil {
		return types.AdsSymbol{}, fmt.Errorf("parseAdsSymbol: failed to read Comment: %w", err)
	}
	symbol.Comment = string(comment)

	return symbol, nil
}
