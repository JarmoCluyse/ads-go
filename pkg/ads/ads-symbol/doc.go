// Package adssymbol provides parsing and validation for ADS symbol information.
//
// This package handles parsing of ADS symbol responses from Beckhoff TwinCAT devices.
// ADS symbols contain metadata about PLC variables including their names, types,
// memory locations, sizes, and comments.
//
// # ADS Symbol Format
//
// ADS symbol responses follow this binary format (all fields in little-endian):
//
//   - Bytes 0-3:   Data Length (uint32) - redundant length field
//   - Bytes 4-7:   IndexGroup (uint32) - ADS index group for the symbol
//   - Bytes 8-11:  IndexOffset (uint32) - ADS index offset for the symbol
//   - Bytes 12-15: Size (uint32) - Size of the variable in bytes
//   - Bytes 16-19: DataType (uint32) - ADS data type identifier
//   - Bytes 20-23: Flags (uint32) - Symbol flags (persistent, static, etc.)
//   - Bytes 24-25: NameLength (uint16) - Length of the variable name
//   - Bytes 26-27: TypeLength (uint16) - Length of the type name
//   - Bytes 28-29: CommentLength (uint16) - Length of the comment
//   - Bytes 30+:   Name - Variable name (null-terminated, length = NameLength + 1)
//   - Bytes ...:   TypeName - Type name (null-terminated, length = TypeLength + 1)
//   - Bytes ...:   Comment - Comment text (NOT null-terminated, length = CommentLength)
//
// # Redundant Length Field
//
// Note: Beckhoff includes a redundant data length field at the beginning of the
// symbol data (bytes 0-3). This duplicates the length already provided in the
// ADS response header. The ParseSymbol function reads and discards this field
// for compatibility with the Beckhoff protocol.
//
// # Core Functions
//
// ParseSymbol parses complete symbol information from binary data:
//
//	data := []byte{ /* symbol data from ADS response */ }
//	symbol, err := adssymbol.ParseSymbol(data)
//	if err != nil {
//	    log.Printf("Failed to parse symbol: %v", err)
//	    return err
//	}
//	fmt.Printf("Symbol: %s (Type: %s, Size: %d bytes)\n",
//	    symbol.Name, symbol.Type, symbol.Size)
//
// CheckSymbol validates symbol data without parsing it:
//
//	if err := adssymbol.CheckSymbol(data); err != nil {
//	    log.Printf("Invalid symbol data: %v", err)
//	    return err
//	}
//
// # String Encoding
//
// Symbol names and type names are null-terminated strings in the PLC character set.
// The Comment field is NOT null-terminated. The package uses utils.DecodePlcStringBuffer
// to properly decode these strings from the PLC encoding to Go strings.
//
// # Error Types
//
// The package exports sentinel errors for error inspection:
//   - ErrInvalidSymbolLength: Response is less than 30 bytes (too short for header)
//   - ErrInsufficientData: Declared string lengths exceed available data
//
// Use errors.Is() for error type checking:
//
//	if errors.Is(err, adssymbol.ErrInvalidSymbolLength) {
//	    // Handle invalid length
//	}
//
// # Common Use Cases
//
// Reading symbol information by name:
//
//	// Query symbol by name using ReadWriteRaw
//	response, err := client.ReadWriteRaw(
//	    port,
//	    uint32(types.ADSReservedIndexGroupSymbolInfoByNameEx),
//	    0,
//	    0xFFFFFFFF,
//	    utils.EncodeStringToPlcStringBuffer("MAIN.MyVariable"),
//	)
//	if err != nil {
//	    return nil, err
//	}
//
//	// Parse the symbol information
//	symbol, err := adssymbol.ParseSymbol(response)
//	if err != nil {
//	    return nil, fmt.Errorf("failed to parse symbol: %w", err)
//	}
//
//	// Use symbol information to read the variable's value
//	value, err := client.ReadRaw(port, symbol.IndexGroup, symbol.IndexOffset, symbol.Size)
//
// # Symbol Flags
//
// The Flags field contains bitwise flags indicating symbol properties:
//   - Persistent: Variable retains value across PLC restarts
//   - Static: Variable has static storage duration
//   - ReadOnly: Variable cannot be written
//   - And other flags defined in types.ADSSymbolFlags
//
// # Data Types
//
// The DataType field uses the ADSDataType enumeration from the types package.
// Common types include:
//   - ADST_INT16, ADST_INT32, ADST_INT64 - Signed integers
//   - ADST_UINT16, ADST_UINT32, ADST_UINT64 - Unsigned integers
//   - ADST_REAL32, ADST_REAL64 - Floating point
//   - ADST_STRING - String type
//   - ADST_BIT - Boolean/bit type
//
// Use types.ADSDataTypeToString() to convert data type codes to readable names.
package adssymbol
