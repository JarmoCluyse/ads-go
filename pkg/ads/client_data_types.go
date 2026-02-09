package ads

import (
	"fmt"

	adsdatatype "github.com/jarmocluyse/ads-go/pkg/ads/ads-datatype"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
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
	// Use the ads-datatype module to parse
	parsedType, err := adsdatatype.ParseDataType(data)
	if err != nil {
		return types.AdsDataType{}, err
	}

	// Convert from adsdatatype.DataType to types.AdsDataType
	return convertDataType(parsedType), nil
}

// convertDataType converts from adsdatatype.DataType to types.AdsDataType
func convertDataType(dt adsdatatype.DataType) types.AdsDataType {
	result := types.AdsDataType{
		Name:          dt.Name,
		Type:          dt.Type,
		Version:       dt.Version,
		HashValue:     dt.HashValue,
		TypeHash:      dt.TypeHash,
		Size:          dt.Size,
		Offset:        dt.Offset,
		DataType:      dt.DataType,
		Flags:         dt.Flags,
		ArrayDim:      dt.ArrayDim,
		Comment:       dt.Comment,
		GUID:          dt.GUID,
		ExtendedFlags: dt.ExtendedFlags,
		CopyMask:      dt.CopyMask,
	}

	// Convert SubItems recursively
	if len(dt.SubItems) > 0 {
		result.SubItems = make([]types.AdsDataType, len(dt.SubItems))
		for i, subItem := range dt.SubItems {
			result.SubItems[i] = convertDataType(subItem)
		}
	}

	// Convert ArrayInfo
	if len(dt.ArrayInfo) > 0 {
		result.ArrayInfo = make([]types.AdsArrayInfo, len(dt.ArrayInfo))
		for i, ai := range dt.ArrayInfo {
			result.ArrayInfo[i] = types.AdsArrayInfo{
				StartIndex: ai.StartIndex,
				Length:     ai.Length,
			}
		}
	}

	// Convert EnumInfo
	if len(dt.EnumInfo) > 0 {
		result.EnumInfo = make([]types.AdsEnumInfo, len(dt.EnumInfo))
		for i, ei := range dt.EnumInfo {
			var attrs []types.AdsAttribute
			if len(ei.Attributes) > 0 {
				attrs = make([]types.AdsAttribute, len(ei.Attributes))
				for j, attr := range ei.Attributes {
					attrs[j] = types.AdsAttribute{
						Name:  attr.Name,
						Value: attr.Value,
					}
				}
			}
			result.EnumInfo[i] = types.AdsEnumInfo{
				Name:       ei.Name,
				Value:      ei.Value,
				Comment:    ei.Comment,
				Attributes: attrs,
			}
		}
	}

	// Convert Attributes
	if len(dt.Attributes) > 0 {
		result.Attributes = make([]types.AdsAttribute, len(dt.Attributes))
		for i, attr := range dt.Attributes {
			result.Attributes[i] = types.AdsAttribute{
				Name:  attr.Name,
				Value: attr.Value,
			}
		}
	}

	// Convert Methods
	if len(dt.Methods) > 0 {
		result.Methods = make([]types.AdsMethod, len(dt.Methods))
		for i, method := range dt.Methods {
			var params []types.AdsMethodParam
			if len(method.Params) > 0 {
				params = make([]types.AdsMethodParam, len(method.Params))
				for j, param := range method.Params {
					params[j] = types.AdsMethodParam{
						Name: param.Name,
						Type: param.Type,
					}
				}
			}
			result.Methods[i] = types.AdsMethod{
				Name:       method.Name,
				ReturnType: method.ReturnType,
				Params:     params,
			}
		}
	}

	return result
}
