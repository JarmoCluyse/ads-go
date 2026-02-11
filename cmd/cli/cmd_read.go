package cli

import (
	"fmt"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

// handleReadValue reads a generic value from the PLC.
// Usage: read_value [symbol_path] [port]
func handleReadValue(args []string, client *ads.Client) {
	data := "GLOBAL.gMyInt"
	var port uint16 = 852

	// Parse arguments if provided
	if len(args) > 0 {
		data = args[0]
	}
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &port)
	}

	value, err := client.ReadValue(port, data)
	if err != nil {
		fmt.Printf("[ERROR] Command 'read_value': Failed to read value from '%s' (port %d): %v\n", data, port, err)
		return
	}
	fmt.Printf("[OK] Value read from '%s' (port %d): %v\n", data, port, value)
}

// handleReadBool reads a boolean value from the PLC.
// Usage: read_bool [symbol_path] [port]
func handleReadBool(args []string, client *ads.Client) {
	data := "GLOBAL.gMyBool"
	var port uint16 = 852

	// Parse arguments if provided
	if len(args) > 0 {
		data = args[0]
	}
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &port)
	}

	value, err := client.ReadValue(port, data)
	if err != nil {
		fmt.Printf("[ERROR] Command 'read_bool': Failed to read bool from '%s' (port %d): %v\n", data, port, err)
		return
	}
	fmt.Printf("[OK] Bool value read from '%s' (port %d): %v\n", data, port, value)
}

// handleReadObject reads a structured object from the PLC.
// Usage: read_object
func handleReadObject(args []string, client *ads.Client) {
	data := "GLOBAL.gMyDUT"
	var port uint16 = 852
	value, err := client.ReadValue(port, data)
	if err != nil {
		fmt.Printf("[ERROR] Command 'read_object': Failed to read object from '%s' (port %d): %v\n", data, port, err)
		return
	}
	fmt.Printf("[OK] Object value read from '%s' (port %d): %v\n", data, port, value)
}

// handleReadArray reads an array from the PLC.
// Usage: read_array
func handleReadArray(args []string, client *ads.Client) {
	data := "GLOBAL.gIntArray"
	var port uint16 = 852
	value, err := client.ReadValue(port, data)
	if err != nil {
		fmt.Printf("[ERROR] Command 'read_array': Failed to read array from '%s' (port %d): %v\n", data, port, err)
		return
	}
	fmt.Printf("[OK] Array value read from '%s' (port %d): %v\n", data, port, value)
}

// handleListSymbols lists all available symbols from the PLC.
// Usage: list_symbols [port]
func handleListSymbols(args []string, client *ads.Client) {
	var port uint16 = 852

	// Parse port if provided
	if len(args) > 0 {
		fmt.Sscanf(args[0], "%d", &port)
	}

	// Use ReadRaw to get symbol upload info
	data, err := client.ReadRaw(port, 0xF00C, 0, 24) // ADSReservedIndexGroupSymbolUploadInfo2
	if err != nil {
		fmt.Printf("[ERROR] Command 'list_symbols': Failed to get symbol info (port %d): %v\n", port, err)
		return
	}

	if len(data) < 24 {
		fmt.Printf("[ERROR] Command 'list_symbols': Insufficient data returned: %d bytes\n", len(data))
		return
	}

	// Parse symbol upload info (first 3 uint32s are what we need)
	symbolCount := uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
	symbolsSize := uint32(data[4]) | uint32(data[5])<<8 | uint32(data[6])<<16 | uint32(data[7])<<24

	fmt.Printf("[INFO] Symbol Count: %d, Symbols Size: %d bytes\n", symbolCount, symbolsSize)

	// Now read the full symbol table
	symbols, err := client.ReadRaw(port, 0xF00B, 0, symbolsSize) // ADSReservedIndexGroupSymbolUpload
	if err != nil {
		fmt.Printf("[ERROR] Command 'list_symbols': Failed to read symbol table (port %d): %v\n", port, err)
		return
	}

	fmt.Printf("[OK] Retrieved %d bytes of symbol data\n", len(symbols))
	fmt.Printf("First 100 symbols:\n")

	// Parse symbol entries - this is a simplified parser
	offset := uint32(0)
	count := 0
	for offset < uint32(len(symbols)) && count < 100 {
		if offset+36 > uint32(len(symbols)) {
			break
		}

		// Read entry length (first uint32)
		entryLength := uint32(symbols[offset]) | uint32(symbols[offset+1])<<8 | uint32(symbols[offset+2])<<16 | uint32(symbols[offset+3])<<24

		if entryLength == 0 || offset+entryLength > uint32(len(symbols)) {
			break
		}

		// Skip to name offset (at byte 24)
		nameOffset := offset + 24

		// Read name length (uint16 at nameOffset)
		if nameOffset+2 > uint32(len(symbols)) {
			break
		}
		nameLength := uint16(symbols[nameOffset]) | uint16(symbols[nameOffset+1])<<8

		// Read name string (starts at nameOffset+2)
		nameStart := nameOffset + 2
		nameEnd := nameStart + uint32(nameLength)

		if nameEnd > uint32(len(symbols)) {
			break
		}

		name := string(symbols[nameStart:nameEnd])
		if len(name) > 0 && name[len(name)-1] == 0 {
			name = name[:len(name)-1] // Remove null terminator
		}

		fmt.Printf("  %3d: %s\n", count+1, name)

		count++
		offset += entryLength
	}

	if symbolCount > 100 {
		fmt.Printf("... and %d more symbols (showing first 100)\n", symbolCount-100)
	}
}
