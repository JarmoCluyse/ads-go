// Package adsdatatype provides parsing functionality for ADS data type information responses.
//
// ADS data types are complex structures that describe the metadata of PLC variables and types,
// including primitive types, structures, enums, arrays, and their nested relationships.
//
// # Binary Format
//
// The ADS data type response has a complex binary format with a base header followed by
// variable-length fields and optional sections based on flag values:
//
//	Offset  Size  Field
//	------  ----  -----
//	0       4     Data length (uint32)
//	4       4     Version (uint32)
//	8       4     Hash value (uint32)
//	12      4     Type hash (uint32)
//	16      4     Size in bytes (uint32)
//	20      4     Offset (uint32)
//	24      4     ADS data type enum (uint32)
//	28      4     Flags (uint32)
//	32      2     Name length (uint16)
//	34      2     Type length (uint16)
//	36      2     Comment length (uint16)
//	38      2     Array dimension (uint16)
//	40      2     Subitem count (uint16)
//	42      N+1   Name string (null-terminated)
//	...     M+1   Type string (null-terminated)
//	...     P+1   Comment string (null-terminated)
//	...     8*D   Array info entries (StartIndex int32, Length uint32) * ArrayDim
//	...     var   Subitems (recursive, each prefixed with uint32 entry length)
//	...     var   Optional fields based on flags
//
// # Optional Fields
//
// The following optional fields may be present based on the Flags value:
//
//   - TypeGuid (0x80): 16-byte GUID
//   - CopyMask (0x200): Variable-length mask data (Size bytes)
//   - MethodInfos (0x800): Method information entries
//   - Attributes (0x1000): Name-value attribute pairs
//   - EnumInfos (0x2000): Enumeration value definitions
//   - ExtendedFlags (0x80000000): 4-byte extended flags
//   - DeRefTypeItem (0x10000000): GUID references
//   - ExtendedEnumInfos (0x20000000): Extended enum information with comments and attributes
//
// # Recursive Parsing
//
// Data types can contain nested subitems, where each subitem is itself a complete data type
// definition. The parser handles this recursion automatically.
//
// # Usage
//
// Parse a complete data type definition:
//
//	dataType, err := adsdatatype.ParseDataType(responseData)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Type: %s, Size: %d bytes\n", dataType.Name, dataType.Size)
//	fmt.Printf("Subitems: %d\n", len(dataType.SubItems))
//
// Validate data before full parsing:
//
//	if err := adsdatatype.CheckDataType(responseData); err != nil {
//	    log.Fatal("Invalid data type response")
//	}
//
// # Error Handling
//
// The package defines sentinel errors that can be checked with errors.Is():
//
//   - ErrInvalidData: The data is malformed or contains invalid values
//   - ErrInsufficientData: Not enough data to parse the response
//   - ErrInvalidEntryLength: Subitem entry length is invalid
//
// Example:
//
//	dataType, err := adsdatatype.ParseDataType(data)
//	if errors.Is(err, adsdatatype.ErrInsufficientData) {
//	    // Handle insufficient data
//	}
package adsdatatype
