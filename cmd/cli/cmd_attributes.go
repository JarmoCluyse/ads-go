package cli

import (
	"fmt"

	"github.com/jarmocluyse/ads-go/pkg/ads"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// handleReadAttributes reads a symbol's type info, pragma attribute declarations, and current value.
// Usage: read_attributes [symbol_path] [port]
func handleReadAttributes(args []string, client *ads.Client) {
	if len(args) == 0 {
		fmt.Println("[ERROR] read_attributes: symbol path required")
		fmt.Println("  Usage: read_attributes <symbol_path> [port]")
		fmt.Println("  Tip:   use list_symbols to discover available symbol paths")
		return
	}

	path := args[0]
	var port uint16 = 852

	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &port)
	}

	sym, err := client.GetSymbol(port, path)
	if err != nil {
		fmt.Printf("[ERROR] read_attributes: failed to get symbol '%s': %v\n", path, err)
		return
	}

	dt, err := client.GetDataType(sym.Type, port)
	if err != nil {
		fmt.Printf("[ERROR] read_attributes: failed to get data type '%s': %v\n", sym.Type, err)
		return
	}

	value, err := client.ReadValue(port, path)
	if err != nil {
		fmt.Printf("[ERROR] read_attributes: failed to read value of '%s': %v\n", path, err)
		return
	}

	fmt.Printf("\n--- Symbol: %s ---\n", sym.Name)
	fmt.Printf("  Type:    %s (%s, %d bytes)\n", sym.Type, types.ADSDataTypeToString(dt.DataType), sym.Size)
	if sym.Comment != "" {
		fmt.Printf("  Comment: %s\n", sym.Comment)
	}
	fmt.Printf("  Value:   %v\n", value)

	fmt.Printf("\n  Symbol Attributes (%d):\n", len(sym.Attributes))
	if len(sym.Attributes) == 0 {
		fmt.Printf("    (none)\n")
	}
	for _, attr := range sym.Attributes {
		fmt.Printf("    %-30s = %q\n", attr.Name, attr.Value)
	}

	fmt.Printf("\n  Type-level Attributes (%d):\n", len(dt.Attributes))
	if len(dt.Attributes) == 0 {
		fmt.Printf("    (none)\n")
	}
	for _, attr := range dt.Attributes {
		fmt.Printf("    %-30s = %q\n", attr.Name, attr.Value)
	}

	if len(dt.SubItems) > 0 {
		fmt.Printf("\n  Sub-items (%d):\n", len(dt.SubItems))
		printAttrSubItems(dt.SubItems, "    ")
	}

	if len(dt.ArrayInfo) > 0 {
		fmt.Printf("\n  Array dimensions (%d):\n", len(dt.ArrayInfo))
		for i, ai := range dt.ArrayInfo {
			fmt.Printf("    [%d] start=%d length=%d\n", i, ai.StartIndex, ai.Length)
		}
	}

	fmt.Println()
}

func printAttrSubItems(items []types.AdsDataType, indent string) {
	for _, item := range items {
		fmt.Printf("%s%-20s : %-20s (offset %d, %d bytes)\n",
			indent, item.Name, item.Type, item.Offset, item.Size)
		for _, attr := range item.Attributes {
			fmt.Printf("%s  {attr} %-28s = %q\n", indent, attr.Name, attr.Value)
		}
		if len(item.SubItems) > 0 {
			printAttrSubItems(item.SubItems, indent+"  ")
		}
	}
}
