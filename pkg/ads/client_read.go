package ads

import (
	"fmt"

	adsprimitives "github.com/jarmocluyse/ads-go/pkg/ads/ads-primitives"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

func (c *Client) ReadValue(port uint16, path string) (any, error) {
	c.logger.Debug("ReadValue: Reading value", "path", path)

	symbol, err := c.GetSymbol(port, path)
	if err != nil {
		return nil, fmt.Errorf("ReadValue: failed to get symbol: %w", err)
	}
	c.logger.Debug("symbol received", "symbol", symbol)

	dataType, err := c.GetDataType(symbol.Type, port)
	if err != nil {
		return nil, fmt.Errorf("ReadValue: failed to get data type: %w", err)
	}

	data, err := c.ReadRaw(port, symbol.IndexGroup, symbol.IndexOffset, symbol.Size)
	if err != nil {
		return nil, fmt.Errorf("ReadValue: failed to read raw data: %w", err)
	}
	return c.convertBufferToValue(data, dataType)
}

func (c *Client) convertBufferToValue(data []byte, dataType types.AdsDataType, isArrayItem ...bool) (any, error) {
	c.logger.Info("convertBufferToValue", "dataType", dataType, "data", data)

	isArrItem := false
	if len(isArrayItem) > 0 {
		isArrItem = isArrayItem[0]
	}

	// If struct or array item
	if (len(dataType.ArrayInfo) == 0 || isArrItem) && len(dataType.SubItems) > 0 {
		result := make(map[string]any)
		data = data[dataType.Offset:]
		for _, subItem := range dataType.SubItems {
			value, err := c.convertBufferToValue(data, subItem)
			if err != nil {
				return nil, err
			}
			result[subItem.Name] = value
		}
		return result, nil
	} else if len(dataType.ArrayInfo) > 0 && !isArrItem {
		// Handle arrays (including multidimensional)
		var convertArrayDimension func(dim int, d []byte) []any
		convertArrayDimension = func(dim int, d []byte) []any {
			temp := []any{}
			for i := 0; i < int(dataType.ArrayInfo[dim].Length); i++ {
				if dim+1 < len(dataType.ArrayInfo) {
					temp = append(temp, convertArrayDimension(dim+1, d))
				} else {
					value, err := c.convertBufferToValue(d, dataType, true)
					if err != nil {
						temp = append(temp, nil)
					} else {
						temp = append(temp, value)
					}
					d = d[dataType.Size:]
				}
			}
			return temp
		}
		return convertArrayDimension(0, data), nil
	}
	switch dataType.DataType {
	case types.ADST_VOID:
		return nil, nil
	case types.ADST_BIT:
		return adsprimitives.ReadBool(data)
	case types.ADST_INT8:
		return adsprimitives.ReadInt8(data)
	case types.ADST_UINT8:
		return adsprimitives.ReadUint8(data)
	case types.ADST_INT16:
		return adsprimitives.ReadInt16(data)
	case types.ADST_UINT16:
		return adsprimitives.ReadUint16(data)
	case types.ADST_INT32:
		return adsprimitives.ReadInt32(data)
	case types.ADST_UINT32:
		return adsprimitives.ReadUint32(data)
	case types.ADST_INT64:
		return adsprimitives.ReadInt64(data)
	case types.ADST_UINT64:
		return adsprimitives.ReadUint64(data)
	case types.ADST_REAL32:
		return adsprimitives.ReadFloat32(data)
	case types.ADST_REAL64:
		return adsprimitives.ReadFloat64(data)
	case types.ADST_STRING, types.ADST_WSTRING:
		return adsprimitives.ReadString(data)
	case types.ADST_BIGTYPE:
		return nil, fmt.Errorf("todo: this data type")
	}
	return nil, fmt.Errorf("unsupported data type: %v", dataType.DataType)
}
