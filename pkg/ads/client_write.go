package ads

import (
	"fmt"

	adsserializer "github.com/jarmocluyse/ads-go/pkg/ads/ads-serializer"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

func (c *Client) WriteValue(port uint16, path string, value any) error {
	c.logger.Debug("WriteValue: Writing value", "path", path)

	// Check if system is in Run mode before writing
	if err := c.checkStateForOperation("WriteValue"); err != nil {
		return err
	}

	symbol, err := c.GetSymbol(port, path)
	if err != nil {
		return fmt.Errorf("WriteValue: failed to get symbol: %w", err)
	}

	dataType, err := c.GetDataType(symbol.Type, port)
	if err != nil {
		return fmt.Errorf("WriteValue: failed to get data type: %w", err)
	}

	data, err := c.convertValueToBuffer(value, dataType)
	if err != nil {
		return fmt.Errorf("WriteValue: failed to convert value to buffer: %w", err)
	}
	err = c.WriteRaw(port, symbol.IndexGroup, symbol.IndexOffset, data)
	return err
}

func (c *Client) convertValueToBuffer(value any, dataType types.AdsDataType, isArrayItem ...bool) ([]byte, error) {
	return adsserializer.Serialize(value, dataType, isArrayItem...)
}
