package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

// handleEnableCounter enables or disables the cycle-based integer counter.
// Usage: enable_counter <true|false>
func handleEnableCounter(args []string, client *ads.Client) {
	data := "GLOBAL.gIntCounterActive"
	var port uint16 = 852
	if len(args) == 0 {
		fmt.Println("[ERROR] Command 'enable_counter': No value provided. Use 'true' or 'false'.")
		return
	}
	var boolValue bool
	switch strings.ToLower(args[0]) {
	case "true":
		boolValue = true
	case "false":
		boolValue = false
	default:
		fmt.Println("[ERROR] Command 'enable_counter': Value must be 'true' or 'false'.")
		return
	}
	err := client.WriteValue(port, data, boolValue)
	if err != nil {
		fmt.Printf("[ERROR] Command 'enable_counter': Failed to write value '%v' to '%s' (port %d): %v\n", boolValue, data, port, err)
		return
	}
	status := "disabled"
	if boolValue {
		status = "enabled"
	}
	fmt.Printf("[OK] Cycle-based counter %s (wrote '%v' to '%s')\n", status, boolValue, data)
}

// handleEnableToggle enables or disables the cycle-based boolean toggle.
// Usage: enable_toggle <true|false>
func handleEnableToggle(args []string, client *ads.Client) {
	data := "GLOBAL.gBoolToggleActive"
	var port uint16 = 852
	if len(args) == 0 {
		fmt.Println("[ERROR] Command 'enable_toggle': No value provided. Use 'true' or 'false'.")
		return
	}
	var boolValue bool
	switch strings.ToLower(args[0]) {
	case "true":
		boolValue = true
	case "false":
		boolValue = false
	default:
		fmt.Println("[ERROR] Command 'enable_toggle': Value must be 'true' or 'false'.")
		return
	}
	err := client.WriteValue(port, data, boolValue)
	if err != nil {
		fmt.Printf("[ERROR] Command 'enable_toggle': Failed to write value '%v' to '%s' (port %d): %v\n", boolValue, data, port, err)
		return
	}
	status := "disabled"
	if boolValue {
		status = "enabled"
	}
	fmt.Printf("[OK] Cycle-based toggle %s (wrote '%v' to '%s')\n", status, boolValue, data)
}

// handleEnableTimedCounter enables or disables the time-based integer counter.
// Usage: enable_timed_counter <true|false>
func handleEnableTimedCounter(args []string, client *ads.Client) {
	data := "GLOBAL.gTimedCounterActive"
	var port uint16 = 852
	if len(args) == 0 {
		fmt.Println("[ERROR] Command 'enable_timed_counter': No value provided. Use 'true' or 'false'.")
		return
	}
	var boolValue bool
	switch strings.ToLower(args[0]) {
	case "true":
		boolValue = true
	case "false":
		boolValue = false
	default:
		fmt.Println("[ERROR] Command 'enable_timed_counter': Value must be 'true' or 'false'.")
		return
	}
	err := client.WriteValue(port, data, boolValue)
	if err != nil {
		fmt.Printf("[ERROR] Command 'enable_timed_counter': Failed to write value '%v' to '%s' (port %d): %v\n", boolValue, data, port, err)
		return
	}
	status := "disabled"
	if boolValue {
		status = "enabled"
	}
	fmt.Printf("[OK] Time-based counter %s (wrote '%v' to '%s')\n", status, boolValue, data)
}

// handleEnableTimedToggle enables or disables the time-based boolean toggle.
// Usage: enable_timed_toggle <true|false>
func handleEnableTimedToggle(args []string, client *ads.Client) {
	data := "GLOBAL.gTimedToggleActive"
	var port uint16 = 852
	if len(args) == 0 {
		fmt.Println("[ERROR] Command 'enable_timed_toggle': No value provided. Use 'true' or 'false'.")
		return
	}
	var boolValue bool
	switch strings.ToLower(args[0]) {
	case "true":
		boolValue = true
	case "false":
		boolValue = false
	default:
		fmt.Println("[ERROR] Command 'enable_timed_toggle': Value must be 'true' or 'false'.")
		return
	}
	err := client.WriteValue(port, data, boolValue)
	if err != nil {
		fmt.Printf("[ERROR] Command 'enable_timed_toggle': Failed to write value '%v' to '%s' (port %d): %v\n", boolValue, data, port, err)
		return
	}
	status := "disabled"
	if boolValue {
		status = "enabled"
	}
	fmt.Printf("[OK] Time-based toggle %s (wrote '%v' to '%s')\n", status, boolValue, data)
}

// handleReadCounters reads all counter and toggle values.
// Usage: read_counters
func handleReadCounters(args []string, client *ads.Client) {
	var port uint16 = 852

	// List of variables to read
	vars := []string{
		"GLOBAL.gMyIntCounter",
		"GLOBAL.gMyBoolToogle",
		"GLOBAL.gTimedIntCounter",
		"GLOBAL.gTimedBoolToogle",
	}

	fmt.Println("[INFO] Reading counter and toggle values:")
	for _, varName := range vars {
		value, err := client.ReadValue(port, varName)
		if err != nil {
			fmt.Printf("  ❌ %s: ERROR - %v\n", varName, err)
		} else {
			fmt.Printf("  ✓ %s: %v\n", varName, value)
		}
	}
}

// handleResetCounters resets all counters and toggles to their initial values.
// Usage: reset_counters
func handleResetCounters(args []string, client *ads.Client) {
	var port uint16 = 852

	// Reset cycle-based counter
	if err := client.WriteValue(port, "GLOBAL.gMyIntCounter", 0); err != nil {
		fmt.Printf("[ERROR] Failed to reset gMyIntCounter: %v\n", err)
		return
	}

	// Reset cycle-based toggle
	if err := client.WriteValue(port, "GLOBAL.gMyBoolToogle", false); err != nil {
		fmt.Printf("[ERROR] Failed to reset gMyBoolToogle: %v\n", err)
		return
	}

	// Reset timed counter
	if err := client.WriteValue(port, "GLOBAL.gTimedIntCounter", 0); err != nil {
		fmt.Printf("[ERROR] Failed to reset gTimedIntCounter: %v\n", err)
		return
	}

	// Reset timed toggle
	if err := client.WriteValue(port, "GLOBAL.gTimedBoolToogle", false); err != nil {
		fmt.Printf("[ERROR] Failed to reset gTimedBoolToogle: %v\n", err)
		return
	}

	fmt.Println("[OK] All counters and toggles reset to initial values (0/false)")
}

// handleReadStatus reads the status of all enable flags.
// Usage: read_status
func handleReadStatus(args []string, client *ads.Client) {
	var port uint16 = 852

	// List of enable flags to read
	flags := []string{
		"GLOBAL.gIntCounterActive",
		"GLOBAL.gBoolToggleActive",
		"GLOBAL.gTimedCounterActive",
		"GLOBAL.gTimedToggleActive",
	}

	fmt.Println("[INFO] Reading enable flag status:")
	for _, flagName := range flags {
		value, err := client.ReadValue(port, flagName)
		if err != nil {
			fmt.Printf("  ❌ %s: ERROR - %v\n", flagName, err)
		} else {
			status := "disabled"
			if boolVal, ok := value.(bool); ok && boolVal {
				status = "enabled"
			}
			fmt.Printf("  %s: %s (%v)\n", flagName, status, value)
		}
	}
}

// handleSetCyclePeriod sets the cycle period for timed operations.
// Usage: set_period <seconds>
func handleSetCyclePeriod(args []string, client *ads.Client) {
	data := "GLOBAL.gCyclePeriod"
	var port uint16 = 852

	if len(args) == 0 {
		fmt.Println("[ERROR] Command 'set_period': No value provided. Specify period in seconds (e.g., '2' for 2 seconds).")
		return
	}

	seconds, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("[ERROR] Command 'set_period': Invalid value '%s'. Must be an integer (seconds).\n", args[0])
		return
	}

	if seconds < 1 || seconds > 3600 {
		fmt.Println("[ERROR] Command 'set_period': Period must be between 1 and 3600 seconds.")
		return
	}

	// TwinCAT TIME type is in milliseconds (DINT)
	milliseconds := int32(seconds * 1000)

	err = client.WriteValue(port, data, milliseconds)
	if err != nil {
		fmt.Printf("[ERROR] Command 'set_period': Failed to write value to '%s' (port %d): %v\n", data, port, err)
		return
	}
	fmt.Printf("[OK] Set cycle period to %d seconds (%d ms) at '%s'\n", seconds, milliseconds, data)
}

// handleReadCyclePeriod reads the current cycle period.
// Usage: read_period
func handleReadCyclePeriod(args []string, client *ads.Client) {
	data := "GLOBAL.gCyclePeriod"
	var port uint16 = 852

	value, err := client.ReadValue(port, data)
	if err != nil {
		fmt.Printf("[ERROR] Command 'read_period': Failed to read value from '%s' (port %d): %v\n", data, port, err)
		return
	}

	// Convert milliseconds to seconds for display
	if ms, ok := value.(int32); ok {
		seconds := float64(ms) / 1000.0
		fmt.Printf("[OK] Current cycle period: %.1f seconds (%d ms)\n", seconds, ms)
	} else {
		fmt.Printf("[OK] Current cycle period: %v\n", value)
	}
}
