package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

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
	fmt.Println("dataType", dataType)

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
		return nil, nil // Void type, no value to return
	case types.ADST_BIT:
		return data[0] != 0, nil // Bit type, return boolean value
	case types.ADST_INT8, types.ADST_UINT8:
		return data[0], nil
	case types.ADST_INT16, types.ADST_UINT16:
		var value int16
		if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
			return nil, err
		}
		return value, nil
	case types.ADST_INT32, types.ADST_UINT32:
		var value int32
		if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
			return nil, err
		}
		return value, nil
	case types.ADST_INT64, types.ADST_UINT64:
		var value int64
		if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
			return nil, err
		}
		return value, nil
	case types.ADST_REAL32:
		var value float32
		if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
			return nil, err
		}
		return value, nil
	case types.ADST_REAL64:
		var value float64
		if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &value); err != nil {
			return nil, err
		}
		return value, nil
	case types.ADST_STRING, types.ADST_WSTRING:
		return string(data), nil
	case types.ADST_BIGTYPE:
		return nil, fmt.Errorf("todo: this data type")
	}
	return nil, fmt.Errorf("unsupported data type: %v", dataType.DataType)
}
