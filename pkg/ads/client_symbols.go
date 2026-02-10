package ads

import (
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
		c.logger.Error("GetSymbol: Failed to send ADS command", "error", err)
		return &adssymbol.AdsSymbol{}, fmt.Errorf("GetSymbol: failed to send ADS command: %w", err)
	}
	symbol, err := adssymbol.ParseSymbol(data)
	if err != nil {
		c.logger.Error("GetSymbol: Failed to parse symbol from response", "error", err)
		return &adssymbol.AdsSymbol{}, fmt.Errorf("GetSymbol: failed to parse symbol from response: %w", err)
	}

	c.logger.Debug("GetSymbol: Symbol read and parsed", "path", path)
	return &symbol, nil
}
