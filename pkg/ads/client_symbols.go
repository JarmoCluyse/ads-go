package ads

import (
	"encoding/binary"
	"fmt"

	adssymbol "github.com/jarmocluyse/ads-go/pkg/ads/ads-symbol"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/jarmocluyse/ads-go/pkg/ads/utils"
)

// GetSymbol retrieves information about a symbol from the ADS server.
func (c *Client) GetSymbol(port uint16, path string) (*adssymbol.AdsSymbol, error) {
	c.logger.Debug("GetSymbol: Requested symbol", "path", path)
	// Create the request data
	data, err := c.ReadWriteRaw(
		port,
		uint32(types.ADSReservedIndexGroupSymbolInfoByNameEx),
		uint32(0),
		uint32(0xFFFFFFFF),
		utils.EncodeStringToPlcStringBuffer(path),
	)
	if err != nil {
		c.logger.Error("GetSymbol: Failed to send ADS command", "path", path, "error", err)
		return &adssymbol.AdsSymbol{}, fmt.Errorf("GetSymbol(%q): failed to send ADS command: %w", path, err)
	}
	symbol, err := adssymbol.ParseSymbol(data)
	if err != nil {
		c.logger.Error("GetSymbol: Failed to parse symbol from response", "path", path, "error", err)
		return &adssymbol.AdsSymbol{}, fmt.Errorf("GetSymbol(%q): failed to parse symbol from response: %w", path, err)
	}

	c.logger.Debug("GetSymbol: Symbol read and parsed", "path", path)
	return &symbol, nil
}

// UploadInfo returns the number of symbols and the total byte size of the
// symbol table blob, read from ADS index group 0xF00F (SymbolUploadInfo2).
// The response is 24 bytes: symCount(4), symSize(4), dataTypeCount(4),
// dataTypeSize(4), extendedDataSize(4), padding(4).
func (c *Client) UploadInfo(port uint16) (symCount uint32, symSize uint32, err error) {
	data, err := c.ReadRaw(port, uint32(types.ADSReservedIndexGroupSymbolUploadInfo2), 0, 24)
	if err != nil {
		return 0, 0, fmt.Errorf("UploadInfo: %w", err)
	}
	if len(data) < 8 {
		return 0, 0, fmt.Errorf("UploadInfo: short response (%d bytes)", len(data))
	}
	symCount = binary.LittleEndian.Uint32(data[0:4])
	symSize = binary.LittleEndian.Uint32(data[4:8])
	return symCount, symSize, nil
}

// UploadSymbols fetches the full symbol table from the PLC and returns all
// symbols. Reads ADS index group 0xF00B (SymbolUpload). Each entry in the
// blob is prefixed with a uint32 entry length (the same value already stored
// in AdsSymbol.DataLen / the first 4 bytes of each entry), which is used to
// step through the flat blob without requiring the full entry to be valid.
func (c *Client) UploadSymbols(port uint16) ([]adssymbol.AdsSymbol, error) {
	_, symSize, err := c.UploadInfo(port)
	if err != nil {
		return nil, err
	}
	if symSize == 0 {
		return nil, nil
	}

	blob, err := c.ReadRaw(port, uint32(types.ADSReservedIndexGroupSymbolUpload), 0, symSize)
	if err != nil {
		return nil, fmt.Errorf("UploadSymbols: %w", err)
	}

	var symbols []adssymbol.AdsSymbol
	pos := 0
	for pos+4 <= len(blob) {
		entryLen := int(binary.LittleEndian.Uint32(blob[pos : pos+4]))
		if entryLen <= 0 || pos+entryLen > len(blob) {
			break
		}
		sym, err := adssymbol.ParseSymbol(blob[pos : pos+entryLen])
		if err == nil {
			symbols = append(symbols, sym)
		}
		pos += entryLen
	}
	return symbols, nil
}
