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
	err = c.WriteRaw(port, symbol.IndexGroup, symbol.IndexOffset, data)
	return err
}

func (c *Client) convertValueToBuffer(value any, dataType types.AdsDataType) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch dataType.DataType {
	case types.ADST_VOID:
		break
	case types.ADST_BIT:
		// Only accept bool for ADS_BIT
		b, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_BIT: %T (expected bool)", value)
		}
		var bitVal byte
		if b {
			bitVal = 1
		} else {
			bitVal = 0
		}
		buf.WriteByte(bitVal)
	case types.ADST_INT8:
		cast, ok := toInt8(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_INT8: %T", value)
		}
		buf.WriteByte(byte(cast))
	case types.ADST_UINT8:
		cast, ok := toUint8(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_UINT8: %T", value)
		}
		buf.WriteByte(cast)
	case types.ADST_INT16:
		cast, ok := toInt16(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_INT16: %T", value)
		}
		if err := binary.Write(buf, binary.LittleEndian, cast); err != nil {
			return nil, err
		}
	case types.ADST_UINT16:
		cast, ok := toUint16(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_UINT16: %T", value)
		}
		if err := binary.Write(buf, binary.LittleEndian, cast); err != nil {
			return nil, err
		}
	case types.ADST_INT32:
		cast, ok := toInt32(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_INT32: %T", value)
		}
		if err := binary.Write(buf, binary.LittleEndian, cast); err != nil {
			return nil, err
		}
	case types.ADST_UINT32:
		cast, ok := toUint32(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_UINT32: %T", value)
		}
		if err := binary.Write(buf, binary.LittleEndian, cast); err != nil {
			return nil, err
		}
	case types.ADST_INT64:
		cast, ok := toInt64(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_INT64: %T", value)
		}
		if err := binary.Write(buf, binary.LittleEndian, cast); err != nil {
			return nil, err
		}
	case types.ADST_UINT64:
		cast, ok := toUint64(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_UINT64: %T", value)
		}
		if err := binary.Write(buf, binary.LittleEndian, cast); err != nil {
			return nil, err
		}
	case types.ADST_REAL32:
		cast, ok := toFloat32(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_REAL32: %T", value)
		}
		if err := binary.Write(buf, binary.LittleEndian, cast); err != nil {
			return nil, err
		}
	case types.ADST_REAL64:
		cast, ok := toFloat64(value)
		if !ok {
			return nil, fmt.Errorf("invalid type for ADST_REAL64: %T", value)
		}
		if err := binary.Write(buf, binary.LittleEndian, cast); err != nil {
			return nil, err
		}
	case types.ADST_STRING:
		if val, ok := value.(string); ok {
			bufferSize := int(dataType.Size)
			if bufferSize <= 0 {
				bufferSize = 80 // Default ADS STRING length if not specified
			}
			strBuf := make([]byte, bufferSize)
			// Reserve last byte for null-terminator
			copyLen := min(len(val), bufferSize-1)
			copy(strBuf, val[:copyLen])
			// strBuf is already zero-padded (Go makes with zeroes)
			buf.Write(strBuf)
		} else {
			return nil, fmt.Errorf("invalid type for ADST_STRING: %T", value)
		}
	case types.ADST_WSTRING:
		if val, ok := value.(string); ok {
			bufferSize := int(dataType.Size)
			if bufferSize <= 0 {
				bufferSize = 160 // Default WSTRING size: 80 chars * 2 bytes (UTF-16LE)
			}
			wbuf := make([]byte, bufferSize)
			// TODO: Proper UTF-16LE encoding
			encoded := []byte(val)
			// Reserve last two bytes for null-terminator (UTF-16)
			copyLen := min(len(encoded), bufferSize-2)
			copy(wbuf, encoded[:copyLen])
			buf.Write(wbuf)
		} else {
			return nil, fmt.Errorf("invalid type for ADST_WSTRING: %T", value)
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
