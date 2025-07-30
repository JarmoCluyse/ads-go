package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// GetDataType retrieves and builds a complete data type definition from the ADS server.
func (c *Client) GetDataType(name string, port uint16) (types.AdsDataType, error) {
	c.logger.Debug("GetDataType: Requested data type", "name", name)

	dataType, err := c.BuildDataType(name, port)
	if err != nil {
		return types.AdsDataType{}, fmt.Errorf("GetDataType: failed to build data type: %w", err)
	}

	c.logger.Debug("GetDataType: Data type built", "name", name)
	return dataType, nil
}

// BuildDataType retrieves and builds a complete data type definition from the ADS server.
func (c *Client) BuildDataType(name string, port uint16) (types.AdsDataType, error) {
	c.logger.Debug("BuildDataType: Building data type", "name", name)

	dataType, err := c.getDataTypeDeclaration(name, port)
	if err != nil {
		return types.AdsDataType{}, fmt.Errorf("BuildDataType: failed to get data type declaration: %w", err)
	}

	return c.buildDataTypeRecursive(dataType, port, true)
}

func (c *Client) buildDataTypeRecursive(dataType types.AdsDataType, port uint16, isRootType bool) (types.AdsDataType, error) {
	c.logger.Debug("buildDataTypeRecursive: Building data type recursively", "Name", dataType.Name, "Type", dataType.Type, "isRootType", isRootType)
	// If the data type has sub-items, recursively build them
	if len(dataType.SubItems) > 0 {
		builtType := dataType
		builtType.SubItems = []types.AdsDataType{}

		for _, subItemDeclaration := range dataType.SubItems {
			c.logger.Debug("buildDataTypeRecursive: Building data type recursively", "subItemDeclaration", subItemDeclaration)
			// Recursively build the sub-item's data type using its Type field
			builtSubItemType, err := c.BuildDataType(subItemDeclaration.Type, port)
			if err != nil {
				return types.AdsDataType{}, err
			}

			// Copy relevant information from the sub-item declaration
			builtSubItemType.Name = subItemDeclaration.Name
			builtSubItemType.Offset = subItemDeclaration.Offset
			builtSubItemType.Comment = subItemDeclaration.Comment

			builtType.SubItems = append(builtType.SubItems, builtSubItemType)
		}
		return builtType, nil
	} else if dataType.ArrayDim > 0 {
		// Data type is an array - get array subtype
		builtType, err := c.BuildDataType(dataType.Type, port)
		if err != nil {
			return types.AdsDataType{}, err
		}

		// Append array info from the current data type to the built type
		builtType.ArrayInfo = append(dataType.ArrayInfo, builtType.ArrayInfo...)
		return builtType, nil
	}

	// For the root type, set the Type field from the Name and clear the Name
	if isRootType {
		dataType.Type = dataType.Name
		dataType.Name = ""
	}

	return dataType, nil
}

func (c *Client) getDataTypeDeclaration(name string, port uint16) (types.AdsDataType, error) {
	data, err := c.ReadWriteRaw(
		port,
		uint32(types.ADSReservedIndexGroupDataDataTypeInfoByNameEx),
		uint32(0),
		uint32(0xFFFFFFFF),
		[]byte(name),
	)
	if err != nil {
		return types.AdsDataType{}, fmt.Errorf("getDataTypeDeclaration: failed to send ADS command: %w", err)
	}
	dataType, err := c.ParseAdsDataTypeResponse(data)
	return dataType, err
}

func (c *Client) ParseAdsDataTypeResponse(data []byte) (types.AdsDataType, error) {
	reader := bytes.NewReader(data)
	var dataType types.AdsDataType

	c.logger.Debug("Parsing AdsDataTypeResponse: starting header fields")
	c.logger.Debug("parsing data type response", "data", data)

	// NOTE: I think beckhoff does something weird here,
	// Response is 0:4 -> ads error 4:8 -> data length 8: -> data
	// however the data 0:4 -> is again the length which makes no sense?
	var dataLen uint32
	if err := binary.Read(reader, binary.LittleEndian, &dataLen); err != nil {
		return types.AdsDataType{}, err
	}
	c.logger.Debug("ParseAdsDataTypeResponse: read dataLen", "value", dataLen, "remainingBytes", reader.Len())

	// NOTE: 0:4 Version
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Version); err != nil {
		c.logger.Debug("Failed parsing Version", "error", err)
		return types.AdsDataType{}, err
	}
	c.logger.Debug("ParseAdsDataTypeResponse", "Version", dataType.Version)
	// NOTE: 4:8 Hash value of datatype for comparison
	if err := binary.Read(reader, binary.LittleEndian, &dataType.HashValue); err != nil {
		c.logger.Debug("Failed parsing HashValue", "error", err)
		return types.AdsDataType{}, err
	}
	c.logger.Debug("ParseAdsDataTypeResponse", "HashValue", dataType.HashValue)
	// NOTE: 8:12 hashValue of base type / Code Offset to setter Method (typeHashValue or offsSetCode)
	if err := binary.Read(reader, binary.LittleEndian, &dataType.TypeHash); err != nil {
		c.logger.Debug("Failed parsing TypeHash", "error", err)
		return types.AdsDataType{}, err
	}
	c.logger.Debug("ParseAdsDataTypeResponse", "TypeHash", dataType.TypeHash)
	// NOTE: 12:16 Size
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Size); err != nil {
		c.logger.Debug("Failed parsing Size", "error", err)
		return types.AdsDataType{}, err
	}
	c.logger.Debug("ParseAdsDataTypeResponse", "Size", dataType.Size)
	// NOTE: 16:20 Offset
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Offset); err != nil {
		c.logger.Debug("Failed parsing Offset", "error", err)
		return types.AdsDataType{}, err
	}
	// NOTE: 20:24 ADS data type
	if err := binary.Read(reader, binary.LittleEndian, &dataType.DataType); err != nil {
		c.logger.Debug("Failed parsing DataType", "error", err)
		return types.AdsDataType{}, err
	}
	// NOTE: 24:28 Flags
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Flags); err != nil {
		c.logger.Debug("Failed parsing Flags", "error", err)
		return types.AdsDataType{}, err
	}
	c.logger.Debug("Parsed header fields", "Version", dataType.Version, "HashValue", dataType.HashValue, "TypeHash", dataType.TypeHash, "Size", dataType.Size, "Offset", dataType.Offset, "DataType", dataType.DataType, "Flags", dataType.Flags)

	// NOTE: 28:30 Name length
	nameLen := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
		c.logger.Debug("Failed parsing NameLen", "error", err)
		return types.AdsDataType{}, err
	}
	// NOTE: 30:32 Type length
	typeLen := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &typeLen); err != nil {
		c.logger.Debug("Failed parsing TypeLen", "error", err)
		return types.AdsDataType{}, err
	}
	// NOTE: 32:34 Comment length
	commentLen := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &commentLen); err != nil {
		c.logger.Debug("Failed parsing CommentLen", "error", err)
		return types.AdsDataType{}, err
	}
	// NOTE: 34:36 Array dimension
	if err := binary.Read(reader, binary.LittleEndian, &dataType.ArrayDim); err != nil {
		c.logger.Debug("Failed parsing ArrayDim", "error", err)
		return types.AdsDataType{}, err
	}
	// NOTE: 36:38 Subitem count
	numSubItems := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &numSubItems); err != nil {
		c.logger.Debug("Failed parsing NumSubItems", "error", err)
		return types.AdsDataType{}, err
	}
	c.logger.Debug("Parsed lengths and counts", "NameLen", nameLen, "TypeLen", typeLen, "CommentLen", commentLen, "ArrayDim", dataType.ArrayDim, "NumSubItems", numSubItems)

	// NOTE: 38.. Data type name
	name := make([]byte, nameLen+1)
	if err := binary.Read(reader, binary.LittleEndian, name); err != nil {
		c.logger.Debug("Failed parsing Name value", "error", err)
		return types.AdsDataType{}, err
	}
	dataType.Name = string(name[:nameLen])
	// NOTE: .. Data type type
	typeVal := make([]byte, typeLen+1)
	if err := binary.Read(reader, binary.LittleEndian, typeVal); err != nil {
		c.logger.Debug("Failed parsing Type value", "error", err)
		return types.AdsDataType{}, err
	}
	dataType.Type = string(typeVal[:typeLen])
	// NOTE: .. Data type comment
	comment := make([]byte, commentLen+1)
	if err := binary.Read(reader, binary.LittleEndian, comment); err != nil {
		c.logger.Debug("Failed parsing Comment value", "error", err)
		return types.AdsDataType{}, err
	}
	dataType.Comment = string(comment[:commentLen])
	c.logger.Debug("Parsed identifier strings", "Name", dataType.Name, "Type", dataType.Type, "Comment", dataType.Comment)

	// NOTE: Parse Array Info
	for i := 0; i < int(dataType.ArrayDim); i++ {
		var arrayInfo types.AdsArrayInfo
		if err := binary.Read(reader, binary.LittleEndian, &arrayInfo.StartIndex); err != nil {
			c.logger.Debug("Failed parsing ArrayInfo StartIndex", "i", i, "error", err)
			return types.AdsDataType{}, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &arrayInfo.Length); err != nil {
			c.logger.Debug("Failed parsing ArrayInfo Length", "i", i, "error", err)
			return types.AdsDataType{}, err
		}
		dataType.ArrayInfo = append(dataType.ArrayInfo, arrayInfo)
	}
	c.logger.Debug("Parsed ArrayInfo", "ArrayInfo", dataType.ArrayInfo)

	// NOTE: Parse Subitems (children data types)
	dataType.SubItems = make([]types.AdsDataType, 0, int(numSubItems))
	for i := 0; i < int(numSubItems); i++ {
		c.logger.Debug("parsing Subitem entryLenBuf", "i", i)
		// Each subitem starts with its entry length (uint32)
		entryLenBuf := make([]byte, 4)
		if _, err := reader.Read(entryLenBuf); err != nil {
			c.logger.Debug("Failed parsing Subitem entryLenBuf", "i", i, "error", err)
			return types.AdsDataType{}, err
		}
		entryLen := binary.LittleEndian.Uint32(entryLenBuf)
		if entryLen < 4 {
			c.logger.Debug("Invalid subitem entryLen", "i", i, "entryLen", entryLen)
			return types.AdsDataType{}, fmt.Errorf("invalid subitem entryLen: %d", entryLen)
		}
		subItemBuf := make([]byte, entryLen-4)
		if _, err := reader.Read(subItemBuf); err != nil {
			c.logger.Debug("Failed reading Subitem buffer", "i", i, "error", err)
			return types.AdsDataType{}, err
		}
		// Recursively parse the subitem (children)
		fullSubItem := append(entryLenBuf, subItemBuf...)
		subItem, err := c.ParseAdsDataTypeResponse(fullSubItem)
		if err != nil {
			c.logger.Debug("Failed parsing Subitem recursively", "i", i, "error", err)
			return types.AdsDataType{}, err
		}
		dataType.SubItems = append(dataType.SubItems, subItem)
	}
	c.logger.Debug("Parsed subitems", "SubItemsCount", len(dataType.SubItems))
	// All subitems (children) have been parsed recursively.

	// NOTE: If flag TypeGuid set
	if (dataType.Flags & types.ADSDataTypeFlagTypeGuid) != 0 {
		c.logger.Debug("Parsing TypeGuid")
		typeGuid := make([]byte, 16)
		if _, err := reader.Read(typeGuid); err != nil {
			c.logger.Debug("Failed reading TypeGuid", "error", err)
			return types.AdsDataType{}, err
		}
		dataType.GUID = fmt.Sprintf("%x", typeGuid)
		c.logger.Debug("Parsed TypeGuid", "GUID", dataType.GUID)
	}

	// NOTE: If flag CopyMask set
	if (dataType.Flags & types.ADSDataTypeFlagCopyMask) != 0 {
		c.logger.Debug("Parsing CopyMask", "Size", dataType.Size)
		// CopyMask is of size dataType.Size
		if dataType.Size > 0 {
			if _, err := reader.Seek(int64(dataType.Size), 1); err != nil {
				c.logger.Debug("Failed seeking CopyMask", "error", err)
				return types.AdsDataType{}, err
			}
			c.logger.Debug("Skipped CopyMask bytes", "count", dataType.Size)
		}
	}

	// NOTE: If flag MethodInfos set
	if (dataType.Flags & types.ADSDataTypeFlagMethodInfos) != 0 {
		c.logger.Debug("Parsing MethodInfos")
		methodCount := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &methodCount); err != nil {
			c.logger.Debug("Failed reading MethodInfo count", "error", err)
			return types.AdsDataType{}, err
		}
		for i := 0; i < int(methodCount); i++ {
			entryLen := uint32(0)
			if err := binary.Read(reader, binary.LittleEndian, &entryLen); err != nil {
				c.logger.Debug("Failed reading MethodInfo entryLen", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			methodBuf := make([]byte, entryLen-4)
			if _, err := reader.Read(methodBuf); err != nil {
				c.logger.Debug("Failed reading MethodInfo buffer", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			// Not parsed in detail
		}
		c.logger.Debug("Parsed MethodInfos", "Count", methodCount)
	}

	// NOTE: If flag Attributes set
	if (dataType.Flags & types.ADSDataTypeFlagAttributes) != 0 {
		c.logger.Debug("Parsing Attributes")
		attributeCount := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &attributeCount); err != nil {
			c.logger.Debug("Failed reading Attribute count", "error", err)
			return types.AdsDataType{}, err
		}
		dataType.Attributes = make([]types.AdsAttribute, 0, attributeCount)
		for i := 0; i < int(attributeCount); i++ {
			nameLen := uint8(0)
			if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
				c.logger.Debug("Failed reading Attribute nameLen", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			valLen := uint8(0)
			if err := binary.Read(reader, binary.LittleEndian, &valLen); err != nil {
				c.logger.Debug("Failed reading Attribute valLen", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			nameBuf := make([]byte, int(nameLen)+1)
			if _, err := reader.Read(nameBuf); err != nil {
				c.logger.Debug("Failed reading Attribute nameBuf", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			valBuf := make([]byte, int(valLen)+1)
			if _, err := reader.Read(valBuf); err != nil {
				c.logger.Debug("Failed reading Attribute valBuf", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			dataType.Attributes = append(dataType.Attributes, types.AdsAttribute{
				Name:  string(nameBuf[:len(nameBuf)-1]),
				Value: string(valBuf[:len(valBuf)-1]),
			})
		}
		c.logger.Debug("Parsed Attributes", "Attributes", dataType.Attributes)
	}

	// NOTE: If flag EnumInfos set
	if (dataType.Flags & types.ADSDataTypeFlagEnumInfos) != 0 {
		c.logger.Debug("Parsing EnumInfo")
		enumInfoCount := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &enumInfoCount); err != nil {
			c.logger.Debug("Failed reading EnumInfo count", "error", err)
			return types.AdsDataType{}, err
		}
		dataType.EnumInfo = make([]types.AdsEnumInfo, 0, enumInfoCount)
		for i := 0; i < int(enumInfoCount); i++ {
			nameLen := uint8(0)
			if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
				c.logger.Debug("Failed reading EnumInfo nameLen", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			nameBuf := make([]byte, int(nameLen)+1)
			if _, err := reader.Read(nameBuf); err != nil {
				c.logger.Debug("Failed reading EnumInfo nameBuf", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			valBuf := make([]byte, dataType.Size)
			if _, err := reader.Read(valBuf); err != nil {
				c.logger.Debug("Failed reading EnumInfo valBuf", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
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
			dataType.EnumInfo = append(dataType.EnumInfo, types.AdsEnumInfo{
				Name:  string(nameBuf[:len(nameBuf)-1]),
				Value: value,
			})
		}
		c.logger.Debug("Parsed EnumInfos", "EnumInfo", dataType.EnumInfo)
	}

	// NOTE: If flag RefactorInfo set
	if (dataType.Flags & types.ADSDataTypeFlagRefactorInfo) != 0 {
		c.logger.Warn("Flag RefactorInfo set, not implemented")
	}

	// NOTE: If flag ExtendedFlags set
	if (dataType.Flags & types.ADSDataTypeFlagExtendedFlags) != 0 {
		c.logger.Debug("Parsing ExtendedFlags")
		var extFlags uint32
		if err := binary.Read(reader, binary.LittleEndian, &extFlags); err != nil {
			c.logger.Debug("Failed reading ExtendedFlags", "error", err)
			return types.AdsDataType{}, err
		}
		dataType.ExtendedFlags = extFlags
		c.logger.Debug("Parsed ExtendedFlags", "ExtendedFlags", extFlags)
	}

	// NOTE: If flag DeRefTypeItem set
	if (dataType.Flags & types.ADSDataTypeFlagDeRefTypeItem) != 0 {
		c.logger.Debug("Parsing DeRefTypeItem")
		var count uint16
		if err := binary.Read(reader, binary.LittleEndian, &count); err != nil {
			c.logger.Debug("Failed reading DeRefTypeItem count", "error", err)
			return types.AdsDataType{}, err
		}
		for i := 0; i < int(count); i++ {
			guid := make([]byte, 16)
			if _, err := reader.Read(guid); err != nil {
				c.logger.Debug("Failed reading DeRefTypeItem guid", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			// Ignored for now
		}
		c.logger.Debug("Parsed DeRefTypeItem GUIDs", "Count", count)
	}

	// NOTE: If flag ExtendedEnumInfos set
	if (dataType.Flags&types.ADSDataTypeFlagExtendedEnumInfos) != 0 && dataType.EnumInfo != nil {
		c.logger.Debug("Parsing ExtendedEnumInfos")
		for i := 0; i < len(dataType.EnumInfo); i++ {
			entryLen := uint16(0)
			if err := binary.Read(reader, binary.LittleEndian, &entryLen); err != nil {
				c.logger.Debug("Failed reading ExtendedEnumInfos entryLen", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			commentLen := uint8(0)
			if err := binary.Read(reader, binary.LittleEndian, &commentLen); err != nil {
				c.logger.Debug("Failed reading ExtendedEnumInfos commentLen", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			attrCount := uint8(0)
			if err := binary.Read(reader, binary.LittleEndian, &attrCount); err != nil {
				c.logger.Debug("Failed reading ExtendedEnumInfos attrCount", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			commentBuf := make([]byte, int(commentLen)+1)
			if _, err := reader.Read(commentBuf); err != nil {
				c.logger.Debug("Failed reading ExtendedEnumInfos commentBuf", "i", i, "error", err)
				return types.AdsDataType{}, err
			}
			attributes := make([]types.AdsAttribute, 0, attrCount)
			for a := 0; a < int(attrCount); a++ {
				nameLen := uint8(0)
				if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
					c.logger.Debug("Failed reading ExtendedEnumInfos attr nameLen", "i", i, "a", a, "error", err)
					return types.AdsDataType{}, err
				}
				valLen := uint8(0)
				if err := binary.Read(reader, binary.LittleEndian, &valLen); err != nil {
					c.logger.Debug("Failed reading ExtendedEnumInfos attr valLen", "i", i, "a", a, "error", err)
					return types.AdsDataType{}, err
				}
				nameBuf := make([]byte, int(nameLen)+1)
				if _, err := reader.Read(nameBuf); err != nil {
					c.logger.Debug("Failed reading ExtendedEnumInfos attr nameBuf", "i", i, "a", a, "error", err)
					return types.AdsDataType{}, err
				}
				valBuf := make([]byte, int(valLen)+1)
				if _, err := reader.Read(valBuf); err != nil {
					c.logger.Debug("Failed reading ExtendedEnumInfos attr valBuf", "i", i, "a", a, "error", err)
					return types.AdsDataType{}, err
				}
				attributes = append(attributes, types.AdsAttribute{
					Name:  string(nameBuf[:len(nameBuf)-1]),
					Value: string(valBuf[:len(valBuf)-1]),
				})
			}
			dataType.EnumInfo[i].Comment = string(commentBuf[:len(commentBuf)-1])
			dataType.EnumInfo[i].Attributes = attributes
		}
		c.logger.Debug("Parsed ExtendedEnumInfos", "EnumInfo", dataType.EnumInfo)
	}

	// NOTE: If flag SoftwareProtectionLevels set
	if (dataType.Flags & types.ADSDataTypeFlagSoftwareProtectionLvls) != 0 {
		c.logger.Warn("Flag SoftwareProtectionLevels set, not implemented")
	}

	c.logger.Debug("ParseAdsDataTypeResponse complete", "DataType.Name", dataType.Name, "Type", dataType.Type, "Subitems", len(dataType.SubItems), "Flags", dataType.Flags, "ArrayInfo", dataType.ArrayInfo)
	return dataType, nil
}
