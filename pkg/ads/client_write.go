package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

func (c *Client) WriteValue(port uint16, path string, value any) error {
	c.logger.Debug("WriteValue: Writing value", "path", path)

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
	_, err = c.WriteRaw(port, symbol.IndexGroup, symbol.IndexOffset, data)
	return err
}

func (c *Client) convertValueToBuffer(value any, dataType types.AdsDataType) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch dataType.DataType {
	case types.ADST_INT8, types.ADST_UINT8:
		if val, ok := value.(byte); ok {
			buf.WriteByte(val)
		} else {
			return nil, fmt.Errorf("invalid type for ADST_INT8/UINT8: %T", value)
		}
	case types.ADST_INT16, types.ADST_UINT16:
		if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
			return nil, err
		}
	case types.ADST_INT32, types.ADST_UINT32:
		if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
			return nil, err
		}
	case types.ADST_INT64, types.ADST_UINT64:
		if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
			return nil, err
		}
	case types.ADST_REAL32:
		if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
			return nil, err
		}
	case types.ADST_REAL64:
		if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
			return nil, err
		}
	case types.ADST_STRING, types.ADST_WSTRING:
		if val, ok := value.(string); ok {
			buf.WriteString(val)
		} else {
			return nil, fmt.Errorf("invalid type for ADST_STRING/WSTRING: %T", value)
		}
	case types.ADST_BIGTYPE:
		if len(dataType.SubItems) > 0 {
			// Handle structs
			if valMap, ok := value.(map[string]any); ok {
				for _, subItem := range dataType.SubItems {
					subItemValue, exists := valMap[subItem.Name]
					if !exists {
						return nil, fmt.Errorf("missing field %s for struct", subItem.Name)
					}
					subItemBuf, err := c.convertValueToBuffer(subItemValue, subItem)
					if err != nil {
						return nil, err
					}
					buf.Write(subItemBuf)
				}
			} else {
				return nil, fmt.Errorf("invalid type for struct: %T", value)
			}
		} else if dataType.ArrayDim > 0 {
			// Handle arrays
			if valSlice, ok := value.([]any); ok {
				for _, element := range valSlice {
					elementBuf, err := c.convertValueToBuffer(element, dataType)
					if err != nil {
						return nil, err
					}
					buf.Write(elementBuf)
				}
			} else {
				return nil, fmt.Errorf("invalid type for array: %T", value)
			}
		}
	default:
		return nil, fmt.Errorf("unsupported data type for conversion to buffer: %v", dataType.DataType)
	}

	return buf.Bytes(), nil
}
