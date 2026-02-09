// Package adsrequests provides functions for building binary payloads for
// ADS (Automation Device Specification) commands.
//
// This package handles the construction of request payloads that are sent
// to Beckhoff TwinCAT PLCs via the ADS protocol. All multi-byte values are
// encoded in little-endian byte order as required by the ADS specification.
//
// # Request Types
//
// The package supports building payloads for these ADS commands:
//
//   - Read: Read data from PLC memory or symbols
//   - Write: Write data to PLC memory or symbols
//   - ReadWrite: Combined read/write operation (used for symbol/type queries)
//   - WriteControl: Change PLC state (Run, Config, Reset, etc.)
//   - ReadState: Query current PLC state
//   - ReadDeviceInfo: Query device version and name
//
// # Read Requests
//
// Read requests access data by index group and offset. Common index groups:
//
//	0x4020 - Memory bits (MB)
//	0xF009 - Symbol information by name
//	0xF011 - Data type information by name
//
// Example:
//
//	// Read 4 bytes from memory area MB offset 100
//	payload := BuildReadRequest(0x4020, 100, 4)
//
// # Write Requests
//
// Write requests send data to specific memory locations:
//
// Example:
//
//	// Write boolean TRUE to output bit at offset 0
//	data := []byte{0x01}
//	payload := BuildWriteRequest(0x4010, 0, data)
//
// # ReadWrite Requests
//
// ReadWrite requests are used for queries where you send a name/query and
// receive structured data in response. This is commonly used for:
//
//   - Symbol lookups (send variable name, receive symbol metadata)
//   - Data type queries (send type name, receive type structure)
//
// Example:
//
//	// Query symbol information by name
//	symbolName := []byte("MAIN.Counter")
//	payload := BuildReadWriteRequestWithNullTerminator(
//	    0xF009,      // SymbolInfoByNameEx
//	    0,           // Offset
//	    0xFFFFFFFF,  // Read all available data
//	    symbolName,
//	)
//
// Note: Many PLC queries expect null-terminated strings. Use
// BuildReadWriteRequestWithNullTerminator for these cases.
//
// # WriteControl Requests
//
// WriteControl requests change the PLC's ADS state:
//
// Common ADS states:
//
//	2  - Reset (transition to Run)
//	5  - Run (normal operation)
//	15 - Config (configuration mode)
//	16 - Reconfig (reconfiguration)
//
// Example:
//
//	// Set PLC to Run mode
//	payload := BuildWriteControlRequest(5, 0)
//
// The deviceState parameter is device-specific and typically set to 0.
//
// # Binary Format
//
// All payloads use little-endian encoding for multi-byte values.
//
// Read Request (12 bytes):
//
//	Bytes 0-3:  Index Group (uint32)
//	Bytes 4-7:  Index Offset (uint32)
//	Bytes 8-11: Read Length (uint32)
//
// Write Request (12 + N bytes):
//
//	Bytes 0-3:   Index Group (uint32)
//	Bytes 4-7:   Index Offset (uint32)
//	Bytes 8-11:  Data Length (uint32)
//	Bytes 12+:   Data
//
// ReadWrite Request (16 + N bytes):
//
//	Bytes 0-3:   Index Group (uint32)
//	Bytes 4-7:   Index Offset (uint32)
//	Bytes 8-11:  Read Length (uint32)
//	Bytes 12-15: Write Length (uint32)
//	Bytes 16+:   Write Data
//
// WriteControl Request (8 bytes):
//
//	Bytes 0-1: ADS State (uint16)
//	Bytes 2-3: Device State (uint16)
//	Bytes 4-7: Data Length (uint32) - always 0 for state changes
//
// # Usage with Client
//
// These functions build only the command payload. To send a complete request,
// the payload must be wrapped in:
//
//  1. AMS Header (32 bytes) - routing and command information
//  2. AMS/TCP Header (6 bytes) - transport layer
//
// The Client handles this wrapping automatically. This package is useful for:
//
//   - Testing payload construction
//   - Custom protocol implementations
//   - Understanding the binary format
//
// # Index Groups Reference
//
// Common reserved index groups:
//
//	0x4000 - Input bits (IB)
//	0x4010 - Output bits (OB)
//	0x4020 - Memory bits (MB)
//	0xF000 - Symbol table
//	0xF003 - Symbol handle by name
//	0xF004 - Symbol value by name
//	0xF009 - Symbol info by name (extended)
//	0xF011 - Data type info by name (extended)
//
// # Reserved Ports
//
// Common ADS ports for different services:
//
//	10000 - System Service (state control, device info)
//	851   - TwinCAT 3 PLC Runtime 1
//	852   - TwinCAT 3 PLC Runtime 2
//	801   - TwinCAT 2 PLC Runtime 1
//
// Note: Ports are not part of the payload - they're specified in the AMS header.
package adsrequests
