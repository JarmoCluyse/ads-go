package cli

import (
	"fmt"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

// handleReadValue reads a generic value from the PLC.
// Usage: read_value
func handleReadValue(args []string, client *ads.Client) {
	data := "service_interface.input.in_servicetool_serviceint_cmd"
	var port uint16 = 350
	value, err := client.ReadValue(port, data)
	if err != nil {
		fmt.Printf("[ERROR] Command 'read_value': Failed to read value from '%s' (port %d): %v\n", data, port, err)
		return
	}
	fmt.Printf("[OK] Value read from '%s' (port %d): %v\n", data, port, value)
}

// handleReadBool reads a boolean value from the PLC.
// Usage: read_bool
func handleReadBool(args []string, client *ads.Client) {
	data := "Service_interface.Input.IN_MAIN_SERVICEINT_ENABLE"
	var port uint16 = 350
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
	data := "Service_interface.Input.IN_busInfo_Main.busPosInit"
	var port uint16 = 350
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
	data := "Service_interface.Service_interface_DW_DS_CMDPARAMS.DS_CMDPARAMS.arrParams"
	var port uint16 = 350
	value, err := client.ReadValue(port, data)
	if err != nil {
		fmt.Printf("[ERROR] Command 'read_array': Failed to read array from '%s' (port %d): %v\n", data, port, err)
		return
	}
	fmt.Printf("[OK] Array value read from '%s' (port %d): %v\n", data, port, value)
}
