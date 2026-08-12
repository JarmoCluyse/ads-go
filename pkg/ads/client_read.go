package ads

import (
	"fmt"

	adsserializer "github.com/jarmocluyse/ads-go/pkg/ads/ads-serializer"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

func (c *Client) ReadValue(port uint16, path string) (any, error) {
	c.logger.Debug("ReadValue: Reading value", "path", path)

	// Check if system is in Run mode before reading
	if err := c.checkStateForOperation("ReadValue"); err != nil {
		return nil, err
	}

	symbol, err := c.GetSymbol(port, path)
	if err != nil {
		return nil, fmt.Errorf("ReadValue(%q): failed to get symbol: %w", path, err)
	}
	c.logger.Debug("ReadValue: Symbol received", "path", path, "symbol", symbol)

	dataType, err := c.GetDataType(symbol.Type, port)
	if err != nil {
		return nil, fmt.Errorf("ReadValue(%q): failed to get data type: %w", path, err)
	}

	data, err := c.ReadRaw(port, symbol.IndexGroup, symbol.IndexOffset, symbol.Size)
	if err != nil {
		return nil, fmt.Errorf("ReadValue(%q): failed to read raw data: %w", path, err)
	}
	return c.convertBufferToValue(data, dataType)
}

func (c *Client) convertBufferToValue(data []byte, dataType types.AdsDataType, isArrayItem ...bool) (any, error) {
	return adsserializer.Deserialize(data, dataType, isArrayItem...)
}
