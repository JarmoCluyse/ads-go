package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

// handleWriteValue writes a numeric value to the PLC.
// Usage: write_value <number>
func handleWriteValue(args []string, client *ads.Client) {
	data := "service_interface.input.in_servicetool_serviceint_cmd"
	var port uint16 = 350
	if len(args) == 0 {
		fmt.Println("[ERROR] Command 'write_value': No value provided to write.")
		return
	}
	// Try to parse the argument as an integer first
	var value any
	if intVal, err := strconv.Atoi(args[0]); err == nil {
		value = intVal
	} else if floatVal, err := strconv.ParseFloat(args[0], 64); err == nil {
		value = floatVal
	} else {
		fmt.Printf("[ERROR] Command 'write_value': Provided value '%s' is not a valid number.\n", args[0])
		return
	}
	err := client.WriteValue(port, data, value)
	if err != nil {
		fmt.Printf("[ERROR] Command 'write_value': Failed to write value '%v' to '%s' (port %d): %v\n", value, data, port, err)
		return
	}
	fmt.Printf("[OK] Wrote value '%v' to '%s' (port %d) successfully.\n", value, data, port)
}

// handleWriteBool writes a boolean value to the PLC.
// Usage: write_bool <true|false>
func handleWriteBool(args []string, client *ads.Client) {
	data := "Service_interface.Input.IN_MAIN_SERVICEINT_ENABLE"
	var port uint16 = 350
	if len(args) == 0 {
		fmt.Println("[ERROR] Command 'write_bool': No value provided to write.")
		return
	}
	var boolValue bool
	switch strings.ToLower(args[0]) {
	case "true":
		boolValue = true
	case "false":
		boolValue = false
	default:
		fmt.Println("[ERROR] Command 'write_bool': Value must be 'true' or 'false'.")
		return
	}
	err := client.WriteValue(port, data, boolValue)
	if err != nil {
		fmt.Printf("[ERROR] Command 'write_bool': Failed to write bool value '%v' to '%s' (port %d): %v\n", boolValue, data, port, err)
		return
	}
	fmt.Printf("[OK] Wrote bool value '%v' to '%s' (port %d) successfully.\n", boolValue, data, port)
}

// handleWriteObject writes a structured object to the PLC.
// Usage: write_object Ln=<int> Lv=<int> PosX=<int> PosY=<int> PosZ=<int> Sct=<int> cntrNewTryLock=<int> flgForce=<true|false>
func handleWriteObject(args []string, client *ads.Client) {
	data := "Service_interface.Input.IN_busInfo_Main.busPosInit"
	var port uint16 = 350
	fields := map[string]string{}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("[ERROR] Command 'write_object': Argument '%s' must be in key=value format.\n", arg)
			return
		}
		fields[parts[0]] = parts[1]
	}
	// Define expected fields and types
	object := map[string]any{}
	for _, name := range []string{"Ln", "Lv", "PosX", "PosY", "PosZ", "Sct", "cntrNewTryLock"} {
		valStr, ok := fields[name]
		if !ok {
			fmt.Printf("[ERROR] Command 'write_object': Missing required field '%s'.\n", name)
			return
		}
		valInt, err := strconv.Atoi(valStr)
		if err != nil {
			fmt.Printf("[ERROR] Command 'write_object': Field '%s' must be an integer, got '%s'.\n", name, valStr)
			return
		}
		object[name] = valInt
	}
	// Handle flgForce (bool)
	flgForceStr, ok := fields["flgForce"]
	if !ok {
		fmt.Println("[ERROR] Command 'write_object': Missing required field 'flgForce'.")
		return
	}
	var flgForce bool
	switch flgForceStr {
	case "true":
		flgForce = true
	case "false":
		flgForce = false
	default:
		fmt.Println("[ERROR] Command 'write_object': flgForce must be 'true' or 'false'.")
		return
	}
	object["flgForce"] = flgForce
	err := client.WriteValue(port, data, object)
	if err != nil {
		fmt.Printf("[ERROR] Command 'write_object': Failed to write object to '%s' (port %d): %v\n", data, port, err)
		return
	}
	fmt.Printf("[OK] Wrote object %v to '%s' (port %d) successfully.\n", object, data, port)
}

// handleWriteArray writes an array of integers to the PLC.
// Usage: write_array <int1> <int2> <int3> <int4> <int5> <int6> <int7> <int8> <int9> <int10>
func handleWriteArray(args []string, client *ads.Client) {
	data := "Service_interface.Service_interface_DW_DS_CMDPARAMS.DS_CMDPARAMS.arrParams"
	var port uint16 = 350
	if len(args) != 10 {
		fmt.Printf("[ERROR] Command 'write_array': You must provide exactly 10 elements to write to the array. Got %d.\n", len(args))
		return
	}
	arr := make([]int, 10)
	for i := range 10 {
		val, err := strconv.Atoi(args[i])
		if err != nil {
			fmt.Printf("[ERROR] Command 'write_array': Argument %d ('%s') is not a valid integer.\n", i+1, args[i])
			return
		}
		arr[i] = val
	}
	err := client.WriteValue(port, data, arr)
	if err != nil {
		fmt.Printf("[ERROR] Command 'write_array': Failed to write array to '%s' (port %d): %v\n", data, port, err)
		return
	}
	fmt.Printf("[OK] Wrote array %v to '%s' (port %d) successfully.\n", arr, data, port)
}
