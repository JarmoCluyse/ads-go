package adsdatatype

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// Sentinel errors for type checking with errors.Is()
var (
	ErrInvalidData        = errors.New("invalid data type response data")
	ErrInsufficientData   = errors.New("insufficient data for data type response")
	ErrInvalidEntryLength = errors.New("invalid subitem entry length")
)

// ParseDataType parses an ADS data type response from binary data.
// The data should contain the complete response including the initial length field.
//
// Binary format:
//
//	0:4   -> Data length (uint32)
//	4:8   -> Version (uint32)
//	8:12  -> Hash value (uint32)
//	12:16 -> Type hash (uint32)
//	16:20 -> Size (uint32)
//	20:24 -> Offset (uint32)
//	24:28 -> ADS data type enum (uint32)
//	28:32 -> Flags (uint32)
//	32:34 -> Name length (uint16)
//	34:36 -> Type length (uint16)
//	36:38 -> Comment length (uint16)
//	38:40 -> Array dimension (uint16)
//	40:42 -> Subitem count (uint16)
//	42:.. -> Name string (null-terminated)
//	..    -> Type string (null-terminated)
//	..    -> Comment string (null-terminated)
//	..    -> Array info entries (8 bytes each: StartIndex int32, Length uint32)
//	..    -> Subitems (recursive, each prefixed with uint32 entry length)
//	..    -> Optional fields based on flags (GUID, CopyMask, MethodInfos, Attributes, EnumInfos, etc.)
//
// Returns the parsed DataType or an error if parsing fails.
func ParseDataType(data []byte) (DataType, error) {
	if len(data) < 4 {
		return DataType{}, ErrInsufficientData
	}

	reader := bytes.NewReader(data)
	var dataType DataType

	// Parse initial data length field
	var dataLen uint32
	if err := binary.Read(reader, binary.LittleEndian, &dataLen); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read data length: %v", ErrInvalidData, err)
	}

	// Parse header fields (28 bytes total)
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Version); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read version: %v", ErrInvalidData, err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.HashValue); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read hash value: %v", ErrInvalidData, err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.TypeHash); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read type hash: %v", ErrInvalidData, err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Size); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read size: %v", ErrInvalidData, err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Offset); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read offset: %v", ErrInvalidData, err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.DataType); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read data type enum: %v", ErrInvalidData, err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Flags); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read flags: %v", ErrInvalidData, err)
	}

	// Parse string length fields (10 bytes total)
	var nameLen, typeLen, commentLen, numSubItems uint16
	if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read name length: %v", ErrInvalidData, err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &typeLen); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read type length: %v", ErrInvalidData, err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &commentLen); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read comment length: %v", ErrInvalidData, err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.ArrayDim); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read array dimension: %v", ErrInvalidData, err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &numSubItems); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read subitem count: %v", ErrInvalidData, err)
	}

	// Parse strings (null-terminated)
	if err := parseString(reader, &dataType.Name, nameLen); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read name: %v", ErrInvalidData, err)
	}
	if err := parseString(reader, &dataType.Type, typeLen); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read type: %v", ErrInvalidData, err)
	}
	if err := parseString(reader, &dataType.Comment, commentLen); err != nil {
		return DataType{}, fmt.Errorf("%w: failed to read comment: %v", ErrInvalidData, err)
	}

	// Parse array info
	if err := parseArrayInfo(reader, &dataType, int(dataType.ArrayDim)); err != nil {
		return DataType{}, err
	}

	// Parse subitems recursively
	if err := parseSubItems(reader, &dataType, int(numSubItems)); err != nil {
		return DataType{}, err
	}

	// Parse optional fields based on flags
	if err := parseOptionalFields(reader, &dataType); err != nil {
		return DataType{}, err
	}

	return dataType, nil
}

// CheckDataType validates that the data can be parsed as a DataType without
// actually parsing the entire structure. This is a lightweight validation.
func CheckDataType(data []byte) error {
	if len(data) < 42 {
		return ErrInsufficientData
	}
	return nil
}

// parseString reads a null-terminated string from the reader.
func parseString(reader *bytes.Reader, target *string, length uint16) error {
	buf := make([]byte, length+1)
	if err := binary.Read(reader, binary.LittleEndian, buf); err != nil {
		return err
	}
	*target = string(buf[:length])
	return nil
}

// parseArrayInfo parses array dimension information.
func parseArrayInfo(reader *bytes.Reader, dataType *DataType, arrayDim int) error {
	dataType.ArrayInfo = make([]ArrayInfo, 0, arrayDim)
	for i := 0; i < arrayDim; i++ {
		var arrayInfo ArrayInfo
		if err := binary.Read(reader, binary.LittleEndian, &arrayInfo.StartIndex); err != nil {
			return fmt.Errorf("%w: failed to read array info start index at %d: %v", ErrInvalidData, i, err)
		}
		if err := binary.Read(reader, binary.LittleEndian, &arrayInfo.Length); err != nil {
			return fmt.Errorf("%w: failed to read array info length at %d: %v", ErrInvalidData, i, err)
		}
		dataType.ArrayInfo = append(dataType.ArrayInfo, arrayInfo)
	}
	return nil
}

// parseSubItems recursively parses sub-items (nested data types).
func parseSubItems(reader *bytes.Reader, dataType *DataType, numSubItems int) error {
	dataType.SubItems = make([]DataType, 0, numSubItems)
	for i := 0; i < numSubItems; i++ {
		// Each subitem starts with its entry length (uint32)
		entryLenBuf := make([]byte, 4)
		if _, err := reader.Read(entryLenBuf); err != nil {
			return fmt.Errorf("%w: failed to read subitem entry length at %d: %v", ErrInvalidData, i, err)
		}
		entryLen := binary.LittleEndian.Uint32(entryLenBuf)
		if entryLen < 4 {
			return fmt.Errorf("%w: %d at subitem %d", ErrInvalidEntryLength, entryLen, i)
		}

		// Read the subitem data
		subItemBuf := make([]byte, entryLen-4)
		if _, err := reader.Read(subItemBuf); err != nil {
			return fmt.Errorf("%w: failed to read subitem buffer at %d: %v", ErrInvalidData, i, err)
		}

		// Recursively parse the subitem
		fullSubItem := append(entryLenBuf, subItemBuf...)
		subItem, err := ParseDataType(fullSubItem)
		if err != nil {
			return fmt.Errorf("failed to parse subitem %d: %w", i, err)
		}
		dataType.SubItems = append(dataType.SubItems, subItem)
	}
	return nil
}

// parseOptionalFields parses optional fields based on flags.
func parseOptionalFields(reader *bytes.Reader, dataType *DataType) error {
	// TypeGuid flag
	if (dataType.Flags & types.ADSDataTypeFlagTypeGuid) != 0 {
		if err := parseTypeGuid(reader, dataType); err != nil {
			return err
		}
	}

	// CopyMask flag
	if (dataType.Flags & types.ADSDataTypeFlagCopyMask) != 0 {
		if err := parseCopyMask(reader, dataType); err != nil {
			return err
		}
	}

	// MethodInfos flag
	if (dataType.Flags & types.ADSDataTypeFlagMethodInfos) != 0 {
		if err := parseMethodInfos(reader); err != nil {
			return err
		}
	}

	// Attributes flag
	if (dataType.Flags & types.ADSDataTypeFlagAttributes) != 0 {
		if err := parseAttributes(reader, dataType); err != nil {
			return err
		}
	}

	// EnumInfos flag
	if (dataType.Flags & types.ADSDataTypeFlagEnumInfos) != 0 {
		if err := parseEnumInfos(reader, dataType); err != nil {
			return err
		}
	}

	// ExtendedFlags flag
	if (dataType.Flags & types.ADSDataTypeFlagExtendedFlags) != 0 {
		if err := parseExtendedFlags(reader, dataType); err != nil {
			return err
		}
	}

	// DeRefTypeItem flag
	if (dataType.Flags & types.ADSDataTypeFlagDeRefTypeItem) != 0 {
		if err := parseDeRefTypeItem(reader); err != nil {
			return err
		}
	}

	// ExtendedEnumInfos flag
	if (dataType.Flags&types.ADSDataTypeFlagExtendedEnumInfos) != 0 && dataType.EnumInfo != nil {
		if err := parseExtendedEnumInfos(reader, dataType); err != nil {
			return err
		}
	}

	return nil
}

// parseTypeGuid parses the 16-byte type GUID.
func parseTypeGuid(reader *bytes.Reader, dataType *DataType) error {
	typeGuid := make([]byte, 16)
	if _, err := reader.Read(typeGuid); err != nil {
		return fmt.Errorf("%w: failed to read type GUID: %v", ErrInvalidData, err)
	}
	dataType.GUID = fmt.Sprintf("%x", typeGuid)
	return nil
}

// parseCopyMask skips the copy mask data.
func parseCopyMask(reader *bytes.Reader, dataType *DataType) error {
	if dataType.Size > 0 {
		if _, err := reader.Seek(int64(dataType.Size), 1); err != nil {
			return fmt.Errorf("%w: failed to skip copy mask: %v", ErrInvalidData, err)
		}
	}
	return nil
}

// parseMethodInfos skips method info data (not parsed in detail).
func parseMethodInfos(reader *bytes.Reader) error {
	var methodCount uint16
	if err := binary.Read(reader, binary.LittleEndian, &methodCount); err != nil {
		return fmt.Errorf("%w: failed to read method count: %v", ErrInvalidData, err)
	}
	for i := 0; i < int(methodCount); i++ {
		var entryLen uint32
		if err := binary.Read(reader, binary.LittleEndian, &entryLen); err != nil {
			return fmt.Errorf("%w: failed to read method entry length at %d: %v", ErrInvalidData, i, err)
		}
		methodBuf := make([]byte, entryLen-4)
		if _, err := reader.Read(methodBuf); err != nil {
			return fmt.Errorf("%w: failed to skip method buffer at %d: %v", ErrInvalidData, i, err)
		}
	}
	return nil
}

// parseAttributes parses attribute name-value pairs.
func parseAttributes(reader *bytes.Reader, dataType *DataType) error {
	var attributeCount uint16
	if err := binary.Read(reader, binary.LittleEndian, &attributeCount); err != nil {
		return fmt.Errorf("%w: failed to read attribute count: %v", ErrInvalidData, err)
	}
	dataType.Attributes = make([]Attribute, 0, attributeCount)
	for i := 0; i < int(attributeCount); i++ {
		var nameLen, valLen uint8
		if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
			return fmt.Errorf("%w: failed to read attribute name length at %d: %v", ErrInvalidData, i, err)
		}
		if err := binary.Read(reader, binary.LittleEndian, &valLen); err != nil {
			return fmt.Errorf("%w: failed to read attribute value length at %d: %v", ErrInvalidData, i, err)
		}
		nameBuf := make([]byte, int(nameLen)+1)
		if _, err := reader.Read(nameBuf); err != nil {
			return fmt.Errorf("%w: failed to read attribute name at %d: %v", ErrInvalidData, i, err)
		}
		valBuf := make([]byte, int(valLen)+1)
		if _, err := reader.Read(valBuf); err != nil {
			return fmt.Errorf("%w: failed to read attribute value at %d: %v", ErrInvalidData, i, err)
		}
		dataType.Attributes = append(dataType.Attributes, Attribute{
			Name:  string(nameBuf[:len(nameBuf)-1]),
			Value: string(valBuf[:len(valBuf)-1]),
		})
	}
	return nil
}

// parseEnumInfos parses enumeration information.
func parseEnumInfos(reader *bytes.Reader, dataType *DataType) error {
	var enumInfoCount uint16
	if err := binary.Read(reader, binary.LittleEndian, &enumInfoCount); err != nil {
		return fmt.Errorf("%w: failed to read enum info count: %v", ErrInvalidData, err)
	}
	dataType.EnumInfo = make([]EnumInfo, 0, enumInfoCount)
	for i := 0; i < int(enumInfoCount); i++ {
		var nameLen uint8
		if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
			return fmt.Errorf("%w: failed to read enum name length at %d: %v", ErrInvalidData, i, err)
		}
		nameBuf := make([]byte, int(nameLen)+1)
		if _, err := reader.Read(nameBuf); err != nil {
			return fmt.Errorf("%w: failed to read enum name at %d: %v", ErrInvalidData, i, err)
		}
		valBuf := make([]byte, dataType.Size)
		if _, err := reader.Read(valBuf); err != nil {
			return fmt.Errorf("%w: failed to read enum value at %d: %v", ErrInvalidData, i, err)
		}

		// Convert value based on size
		var value int64
		bufLen := len(valBuf)
		if bufLen > 0 {
			switch bufLen {
			case 1:
				value = int64(valBuf[0])
			case 2:
				value = int64(binary.LittleEndian.Uint16(valBuf))
			case 4:
				value = int64(binary.LittleEndian.Uint32(valBuf))
			case 8:
				value = int64(binary.LittleEndian.Uint64(valBuf))
			default:
				padded := make([]byte, 8)
				copy(padded, valBuf)
				value = int64(binary.LittleEndian.Uint64(padded))
			}
		}
		dataType.EnumInfo = append(dataType.EnumInfo, EnumInfo{
			Name:  string(nameBuf[:len(nameBuf)-1]),
			Value: value,
		})
	}
	return nil
}

// parseExtendedFlags parses extended flags.
func parseExtendedFlags(reader *bytes.Reader, dataType *DataType) error {
	var extFlags uint32
	if err := binary.Read(reader, binary.LittleEndian, &extFlags); err != nil {
		return fmt.Errorf("%w: failed to read extended flags: %v", ErrInvalidData, err)
	}
	dataType.ExtendedFlags = extFlags
	return nil
}

// parseDeRefTypeItem skips DeRefTypeItem GUIDs.
func parseDeRefTypeItem(reader *bytes.Reader) error {
	var count uint16
	if err := binary.Read(reader, binary.LittleEndian, &count); err != nil {
		return fmt.Errorf("%w: failed to read deref type item count: %v", ErrInvalidData, err)
	}
	for i := 0; i < int(count); i++ {
		guid := make([]byte, 16)
		if _, err := reader.Read(guid); err != nil {
			return fmt.Errorf("%w: failed to skip deref type item GUID at %d: %v", ErrInvalidData, i, err)
		}
	}
	return nil
}

// parseExtendedEnumInfos parses extended enum information (comments and attributes).
func parseExtendedEnumInfos(reader *bytes.Reader, dataType *DataType) error {
	for i := 0; i < len(dataType.EnumInfo); i++ {
		var entryLen uint16
		if err := binary.Read(reader, binary.LittleEndian, &entryLen); err != nil {
			return fmt.Errorf("%w: failed to read extended enum entry length at %d: %v", ErrInvalidData, i, err)
		}
		var commentLen, attrCount uint8
		if err := binary.Read(reader, binary.LittleEndian, &commentLen); err != nil {
			return fmt.Errorf("%w: failed to read extended enum comment length at %d: %v", ErrInvalidData, i, err)
		}
		if err := binary.Read(reader, binary.LittleEndian, &attrCount); err != nil {
			return fmt.Errorf("%w: failed to read extended enum attr count at %d: %v", ErrInvalidData, i, err)
		}
		commentBuf := make([]byte, int(commentLen)+1)
		if _, err := reader.Read(commentBuf); err != nil {
			return fmt.Errorf("%w: failed to read extended enum comment at %d: %v", ErrInvalidData, i, err)
		}

		// Parse attributes
		attributes := make([]Attribute, 0, attrCount)
		for a := 0; a < int(attrCount); a++ {
			var nameLen, valLen uint8
			if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
				return fmt.Errorf("%w: failed to read extended enum attr name length at %d,%d: %v", ErrInvalidData, i, a, err)
			}
			if err := binary.Read(reader, binary.LittleEndian, &valLen); err != nil {
				return fmt.Errorf("%w: failed to read extended enum attr value length at %d,%d: %v", ErrInvalidData, i, a, err)
			}
			nameBuf := make([]byte, int(nameLen)+1)
			if _, err := reader.Read(nameBuf); err != nil {
				return fmt.Errorf("%w: failed to read extended enum attr name at %d,%d: %v", ErrInvalidData, i, a, err)
			}
			valBuf := make([]byte, int(valLen)+1)
			if _, err := reader.Read(valBuf); err != nil {
				return fmt.Errorf("%w: failed to read extended enum attr value at %d,%d: %v", ErrInvalidData, i, a, err)
			}
			attributes = append(attributes, Attribute{
				Name:  string(nameBuf[:len(nameBuf)-1]),
				Value: string(valBuf[:len(valBuf)-1]),
			})
		}

		dataType.EnumInfo[i].Comment = string(commentBuf[:len(commentBuf)-1])
		dataType.EnumInfo[i].Attributes = attributes
	}
	return nil
}
