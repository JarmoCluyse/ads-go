package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// GetDataType retrieves and builds a complete data type definition from the ADS server.
func (c *Client) GetDataType(name string) (types.AdsDataType, error) {
	c.logger.Debug("GetDataType: Requested data type", "name", name)

	// Check if the data type is already cached
	if dataType, ok := c.plcDataTypes[name]; ok {
		c.logger.Debug("GetDataType: Data type found in cache", "name", name)
		return dataType, nil
	}

	c.logger.Debug("GetDataType: Data type not cached, building from target", "name", name)

	dataType, err := c.BuildDataType(name)
	if err != nil {
		return types.AdsDataType{}, fmt.Errorf("GetDataType: failed to build data type: %w", err)
	}

	// Cache the data type
	c.plcDataTypes[name] = dataType

	c.logger.Debug("GetDataType: Data type built and cached", "name", name)
	return dataType, nil
}

// BuildDataType retrieves and builds a complete data type definition from the ADS server.
func (c *Client) BuildDataType(name string) (types.AdsDataType, error) {
	c.logger.Debug("BuildDataType: Building data type", "name", name)

	dataType, err := c.getDataTypeDeclaration(name)
	if err != nil {
		return types.AdsDataType{}, fmt.Errorf("BuildDataType: failed to get data type declaration: %w", err)
	}

	return c.buildDataTypeRecursive(dataType, true)
}

func (c *Client) buildDataTypeRecursive(dataType types.AdsDataType, isRootType bool) (types.AdsDataType, error) {
	// If the data type has sub-items, recursively build them
	if len(dataType.SubItems) > 0 {
		builtType := dataType
		builtType.SubItems = []types.AdsDataType{}

		for _, subItemDeclaration := range dataType.SubItems {
			// Recursively build the sub-item's data type using its Type field
			builtSubItemType, err := c.BuildDataType(subItemDeclaration.Type)
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
		builtType, err := c.BuildDataType(dataType.Type)
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

func (c *Client) getDataTypeDeclaration(name string) (types.AdsDataType, error) {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, types.ADSReservedIndexGroupSymbolDataTypeUpload)
	binary.Write(data, binary.LittleEndian, uint32(len(name)))
	binary.Write(data, binary.LittleEndian, uint32(0xFFFFFFFF))
	data.WriteString(name)

	req := AdsCommandRequest{
		Command:     types.ADSCommandReadWrite,
		Data:        data.Bytes(),
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}

	response, err := c.send(req)
	if err != nil {
		return types.AdsDataType{}, fmt.Errorf("getDataTypeDeclaration: failed to send ADS command: %w", err)
	}

	return parseAdsDataTypeResponse(response)
}

func parseAdsDataTypeResponse(data []byte) (types.AdsDataType, error) {
	reader := bytes.NewReader(data)

	// Skip the first 8 bytes (error code + length) which are part of the ADS response header
	_, err := reader.Seek(8, 0)
	if err != nil {
		return types.AdsDataType{}, err
	}

	return parseAdsDataTypeInternal(reader)
}

func parseAdsDataTypeInternal(reader *bytes.Reader) (types.AdsDataType, error) {
	var dataType types.AdsDataType

	if err := binary.Read(reader, binary.LittleEndian, &dataType.Version); err != nil {
		return types.AdsDataType{}, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.HashValue); err != nil {
		return types.AdsDataType{}, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.TypeHash); err != nil {
		return types.AdsDataType{}, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Size); err != nil {
		return types.AdsDataType{}, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Offset); err != nil {
		return types.AdsDataType{}, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.DataType); err != nil {
		return types.AdsDataType{}, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.Flags); err != nil {
		return types.AdsDataType{}, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &dataType.ArrayDim); err != nil {
		return types.AdsDataType{}, err
	}

	nameLen := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
		return types.AdsDataType{}, err
	}
	name := make([]byte, nameLen)
	if _, err := reader.Read(name); err != nil {
		return types.AdsDataType{}, err
	}
	dataType.Name = string(name)

	commentLen := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &commentLen); err != nil {
		return types.AdsDataType{}, err
	}
	comment := make([]byte, commentLen)
	if _, err := reader.Read(comment); err != nil {
		return types.AdsDataType{}, err
	}
	dataType.Comment = string(comment)

	// Parse Array Info
	for i := 0; i < int(dataType.ArrayDim); i++ {
		var arrayInfo types.AdsArrayInfo
		if err := binary.Read(reader, binary.LittleEndian, &arrayInfo.LowerBound); err != nil {
			return types.AdsDataType{}, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &arrayInfo.Elements); err != nil {
			return types.AdsDataType{}, err
		}
		dataType.ArrayInfo = append(dataType.ArrayInfo, arrayInfo)
	}

	// Parse Sub-Items (declarations only, not full types)
	numSubItems := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &numSubItems); err != nil {
		return types.AdsDataType{}, err
	}

	for i := 0; i < int(numSubItems); i++ {
		var subItem types.AdsDataTypeSubItem
		if err := binary.Read(reader, binary.LittleEndian, &subItem.EntryLength); err != nil {
			return types.AdsDataType{}, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &subItem.Version); err != nil {
			return types.AdsDataType{}, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &subItem.HashValue); err != nil {
			return types.AdsDataType{}, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &subItem.TypeHash); err != nil {
			return types.AdsDataType{}, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &subItem.Size); err != nil {
			return types.AdsDataType{}, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &subItem.Offset); err != nil {
			return types.AdsDataType{}, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &subItem.DataType); err != nil {
			return types.AdsDataType{}, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &subItem.Flags); err != nil {
			return types.AdsDataType{}, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &subItem.ArrayDim); err != nil {
			return types.AdsDataType{}, err
		}
		nameLen := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
			return types.AdsDataType{}, err
		}
		name := make([]byte, nameLen)
		if _, err := reader.Read(name); err != nil {
			return types.AdsDataType{}, err
		}
		subItem.Name = string(name)

		typeLen := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &typeLen); err != nil {
			return types.AdsDataType{}, err
		}
		typeName := make([]byte, typeLen)
		if _, err := reader.Read(typeName); err != nil {
			return types.AdsDataType{}, err
		}
		subItem.Type = string(typeName)

		commentLen := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &commentLen); err != nil {
			return types.AdsDataType{}, err
		}
		comment := make([]byte, commentLen)
		if _, err := reader.Read(comment); err != nil {
			return types.AdsDataType{}, err
		}
		subItem.Comment = string(comment)

		// Append the sub-item declaration. buildDataTypeRecursive will build the full type later.
		// For now, we need to create a new AdsDataType from AdsDataTypeSubItem.
		dataType.SubItems = append(dataType.SubItems, types.AdsDataType{
			Name: subItem.Name,
			Type: subItem.Type,
			Size: subItem.Size,
			Offset: subItem.Offset,
			DataType: subItem.DataType,
			Flags: types.ADSDataTypeFlags(subItem.Flags),
			ArrayDim: subItem.ArrayDim,
			Comment: subItem.Comment,
		})
	}

	// Parse Enum Info
	numEnumInfos := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &numEnumInfos); err != nil {
		return types.AdsDataType{}, err
	}

	for i := 0; i < int(numEnumInfos); i++ {
		var enumInfo types.AdsEnumInfo
		nameLen := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
			return types.AdsDataType{}, err
		}
		name := make([]byte, nameLen)
		if _, err := reader.Read(name); err != nil {
			return types.AdsDataType{}, err
		}
		enumInfo.Name = string(name)

		if err := binary.Read(reader, binary.LittleEndian, &enumInfo.Value); err != nil {
			return types.AdsDataType{}, err
		}
		dataType.EnumInfo = append(dataType.EnumInfo, enumInfo)
	}

	// Parse Attributes
	numAttributes := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &numAttributes); err != nil {
		return types.AdsDataType{}, err
	}

	for i := 0; i < int(numAttributes); i++ {
		var attribute types.AdsAttribute
		nameLen := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
			return types.AdsDataType{}, err
		}
		name := make([]byte, nameLen)
		if _, err := reader.Read(name); err != nil {
			return types.AdsDataType{}, err
		}
		attribute.Name = string(name)

		valueLen := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &valueLen); err != nil {
			return types.AdsDataType{}, err
		}
		value := make([]byte, valueLen)
		if _, err := reader.Read(value); err != nil {
			return types.AdsDataType{}, err
		}
		attribute.Value = string(value)
		dataType.Attributes = append(dataType.Attributes, attribute)
	}

	// Parse Methods
	numMethods := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &numMethods); err != nil {
		return types.AdsDataType{}, err
	}

	for i := 0; i < int(numMethods); i++ {
		var method types.AdsMethod
		// Parse method name
		nameLen := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
			return types.AdsDataType{}, err
		}
		name := make([]byte, nameLen)
		if _, err := reader.Read(name); err != nil {
			return types.AdsDataType{}, err
		}
		method.Name = string(name)

		// Parse return type
		returnTypeLen := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &returnTypeLen); err != nil {
			return types.AdsDataType{}, err
		}
		returnType := make([]byte, returnTypeLen)
		if _, err := reader.Read(returnType); err != nil {
			return types.AdsDataType{}, err
		}
		method.ReturnType = string(returnType)

		// Parse params
		numParams := uint16(0)
		if err := binary.Read(reader, binary.LittleEndian, &numParams); err != nil {
			return types.AdsDataType{}, err
		}

		for j := 0; j < int(numParams); j++ {
			var param types.AdsMethodParam
			// Parse param name
			paramNameLen := uint16(0)
			if err := binary.Read(reader, binary.LittleEndian, &paramNameLen); err != nil {
				return types.AdsDataType{}, err
			}
			paramName := make([]byte, paramNameLen)
			if _, err := reader.Read(paramName); err != nil {
				return types.AdsDataType{}, err
			}
			param.Name = string(paramName)

			// Parse param type
			paramTypeLen := uint16(0)
			if err := binary.Read(reader, binary.LittleEndian, &paramTypeLen); err != nil {
				return types.AdsDataType{}, err
			}
			paramType := make([]byte, paramTypeLen)
			if _, err := reader.Read(paramType); err != nil {
				return types.AdsDataType{}, err
			}
			param.Type = string(paramType)
			method.Params = append(method.Params, param)
		}
		dataType.Methods = append(dataType.Methods, method)
	}

	return dataType, nil
}