package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

func (c *Client) ReadValue(path string) (any, error) {
	c.logger.Debug("ReadValue: Reading value", "path", path)

	symbol, err := c.GetSymbol(path)
	if err != nil {
		return nil, fmt.Errorf("ReadValue: failed to get symbol: %w", err)
	}

	dataType, err := c.GetDataType(symbol.Type)
	if err != nil {
		return nil, fmt.Errorf("ReadValue: failed to get data type: %w", err)
	}

	data, err := c.ReadRaw(symbol.IndexGroup, symbol.IndexOffset, symbol.Size)
	if err != nil {
		return nil, fmt.Errorf("ReadValue: failed to read raw data: %w", err)
	}

	return c.convertBufferToValue(data, dataType)
}

func (c *Client) convertBufferToValue(data []byte, dataType types.AdsDataType) (any, error) {
	switch dataType.DataType {
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
		if len(dataType.SubItems) > 0 {
			// Handle structs
			result := make(map[string]any)
			for _, subItem := range dataType.SubItems {
				value, err := c.convertBufferToValue(data[subItem.Offset:subItem.Offset+subItem.Size], subItem)
				if err != nil {
					return nil, err
				}
				result[subItem.Name] = value
			}
			return result, nil
		} else if dataType.ArrayDim > 0 {
			// Handle arrays
			var result []any
			for i := 0; i < int(dataType.ArrayInfo[0].Elements); i++ {
				// TODO: This is not correct for multi-dimensional arrays
				elementSize := dataType.Size / uint32(dataType.ArrayInfo[0].Elements)
				value, err := c.convertBufferToValue(data[uint32(i)*elementSize:], dataType)
				if err != nil {
					return nil, err
				}
				result = append(result, value)
			}
			return result, nil
		}
	}
	return nil, fmt.Errorf("unsupported data type: %v", dataType.DataType)
}
