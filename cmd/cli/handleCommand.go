package cli

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

var handlers = map[string]func([]string, *ads.Client){
	"device_info": func(args []string, client *ads.Client) {
		info, err := client.ReadDeviceInfo()
		if err != nil {
			fmt.Printf("[ERROR] Command 'device_info': Failed to read device info: %v\n", err)
			return
		}
		fmt.Printf("[OK] Device version: %d.%d.%d\n", info.MajorVersion, info.MinorVersion, info.VersionBuild)
		fmt.Printf("[OK] Device name: %s\n", info.DeviceName)
	},
	"state": func(args []string, client *ads.Client) {
		state, err := client.ReadTcSystemState()
		if err != nil {
			fmt.Printf("[ERROR] Command 'state': Failed to read system state: %v\n", err)
			return
		}
		fmt.Printf("[OK] System ADS state: %s\n", state.AdsState.String())
	},
	"state_loop": func(args []string, client *ads.Client) {
		// Run state every 550ms until ^C
		fmt.Println("[INFO] Press Ctrl+C to exit system state loop...")
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		done := make(chan struct{})
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		go func() {
			<-sig
			close(done)
		}()

		for {
			select {
			case <-ticker.C:
				state, err := client.ReadTcSystemState()
				if err != nil {
					fmt.Printf("[ERROR] Command 'state_loop': Failed to read system state: %v\n", err)
				} else {
					fmt.Printf("[INFO] System ADS state: %s\n", state.AdsState.String())
				}
			case <-done:
				fmt.Println("[INFO] Exiting system state loop.")
				return
			}
		}
	},
	"toConfig": func(args []string, client *ads.Client) {
		err := client.SetTcSystemToConfig()
		if err != nil {
			fmt.Printf("[ERROR] Command 'toConfig': Failed to set TwinCAT to Config state: %v\n", err)
		}
	},
	"toRun": func(args []string, client *ads.Client) {
		err := client.SetTcSystemToRun()
		if err != nil {
			fmt.Println("Error reading system state:", err)
		}
	},
	"read_value": func(args []string, client *ads.Client) {
		data := "service_interface.input.in_servicetool_serviceint_cmd"
		var port uint16 = 350
		value, err := client.ReadValue(port, data)
		if err != nil {
			fmt.Printf("[ERROR] Command 'read_array': Failed to read array from '%s' (port %d): %v\n", data, port, err)
			fmt.Printf("[ERROR] Command 'read_value': Failed to read value from '%s' (port %d): %v\n", data, port, err)
			fmt.Printf("[ERROR] Command 'read_object': Failed to read object from '%s' (port %d): %v\n", data, port, err)
		} else {
			fmt.Printf("[OK] Object value read from '%s' (port %d): %v\n", data, port, value)
		}
	},
	"write_value": func(args []string, client *ads.Client) {
		data := "service_interface.input.in_servicetool_serviceint_cmd"
		var port uint16 = 350
		if len(args) == 0 {
			fmt.Println("[ERROR] Command 'write_bool': No value provided to write.")
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
		} else {
			fmt.Printf("[OK] Wrote value '%v' to '%s' (port %d) successfully.\n", value, data, port)
		}
	},

	"read_bool": func(args []string, client *ads.Client) {
		data := "Service_interface.Input.IN_MAIN_SERVICEINT_ENABLE"
		var port uint16 = 350
		value, err := client.ReadValue(port, data)
		if err != nil {
			fmt.Printf("[ERROR] Command 'read_bool': Failed to read bool from '%s' (port %d): %v\n", data, port, err)
		} else {
			fmt.Printf("[OK] Bool value read from '%s' (port %d): %v\n", data, port, value)
		}
	},
	"write_bool": func(args []string, client *ads.Client) {
		data := "Service_interface.Input.IN_MAIN_SERVICEINT_ENABLE"
		var port uint16 = 350
		if len(args) == 0 {
			fmt.Println("Error: No value provided to write.")
			return
		}
		var boolValue bool
		var err error
		switch strings.ToLower(args[0]) {
		case "true":
			boolValue = true
		case "false":
			boolValue = false
		default:
			fmt.Println("[ERROR] Command 'write_bool': Value must be 'true' or 'false'.")
			return
		}
		err = client.WriteValue(port, data, boolValue)
		if err != nil {
			fmt.Printf("[ERROR] Command 'write_bool': Failed to write bool value '%v' to '%s' (port %d): %v\n", boolValue, data, port, err)
		} else {
			fmt.Printf("[OK] Wrote bool value '%v' to '%s' (port %d) successfully.\n", boolValue, data, port)
		}
	},
	"read_object": func(args []string, client *ads.Client) {
		data := "Service_interface.Input.IN_busInfo_Main.busPosInit"
		var port uint16 = 350
		value, err := client.ReadValue(port, data)
		if err != nil {
			fmt.Println("Error reading system state:", err)
		}
		fmt.Printf("value %v\n", value)
	},
	"read_array": func(args []string, client *ads.Client) {
		data := "Service_interface.Service_interface_DW_DS_CMDPARAMS.DS_CMDPARAMS.arrParams"
		var port uint16 = 350
		value, err := client.ReadValue(port, data)
		if err != nil {
			fmt.Println("Error reading system state:", err)
		}
		fmt.Printf("value %v\n", value)
	},
	"write_array": func(args []string, client *ads.Client) {
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
		} else {
			fmt.Printf("[OK] Wrote array %v to '%s' (port %d) successfully.\n", arr, data, port)
		}
	},

	"write_object": func(args []string, client *ads.Client) {
		// Usage: write_object Ln=1 Lv=2 PosX=3 PosY=4 PosZ=5 Sct=6 cntrNewTryLock=7 flgForce=true
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
		} else {
			fmt.Printf("[OK] Wrote object %v to '%s' (port %d) successfully.\n", object, data, port)
		}
	},

	"read_raw": func(args []string, client *ads.Client) {
		indexGroup := uint32(0x1010290)
		indexOffset := uint32(0x80000001)
		var port uint16 = 350
		size := uint32(1) // adjust size as needed
		result, err := client.ReadRaw(port, indexGroup, indexOffset, size)
		if err != nil {
			fmt.Printf("[ERROR] Command 'read_raw': Failed to read raw data [IG: 0x%X, IO: 0x%X, port: %d]: %v\n", indexGroup, indexOffset, port, err)
			return
		}
		fmt.Printf("[OK] Raw read [IG: 0x%X, IO: 0x%X, port: %d]: %v\n", indexGroup, indexOffset, port, result)
	},
	"write_raw": func(args []string, client *ads.Client) {
		indexGroup := uint32(0x1010290)
		indexOffset := uint32(0x80000001)
		var port uint16 = 350
		data := []byte{36}
		err := client.WriteRaw(port, indexGroup, indexOffset, data)
		if err != nil {
			fmt.Printf("[ERROR] Command 'write_raw': Failed to write raw data [IG: 0x%X, IO: 0x%X, port: %d]: %v\n", indexGroup, indexOffset, port, err)
			return
		}
		fmt.Printf("[OK] Raw write [IG: 0x%X, IO: 0x%X, port: %d] succeeded.\n", indexGroup, indexOffset, port)
	},
}

func handleCommand(input string, client *ads.Client) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	handler, ok := handlers[parts[0]]
	if ok {
		handler(parts[1:], client)
	} else {
		fmt.Printf("[ERROR] Unknown command: '%s'. Type 'help' for a list of available commands.\n", parts[0])
	}
}
