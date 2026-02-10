package cli

import (
	"fmt"
	"strings"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

// CommandHandler is the function signature for all CLI command handlers.
type CommandHandler func(args []string, client *ads.Client)

// getHandlers returns the map of command names to their handler functions.
// This function is used to avoid initialization cycles.
func getHandlers() map[string]CommandHandler {
	return map[string]CommandHandler{
		// System commands (cmd_system.go)
		"device_info": handleDeviceInfo,
		"state":       handleState,
		"state_loop":  handleStateLoop,
		"toConfig":    handleToConfig,
		"toRun":       handleToRun,

		// Read commands (cmd_read.go)
		"read_value":  handleReadValue,
		"read_bool":   handleReadBool,
		"read_object": handleReadObject,
		"read_array":  handleReadArray,

		// Write commands (cmd_write.go)
		"write_value":  handleWriteValue,
		"write_bool":   handleWriteBool,
		"write_object": handleWriteObject,
		"write_array":  handleWriteArray,

		// Raw commands (cmd_raw.go)
		"read_raw":  handleReadRaw,
		"write_raw": handleWriteRaw,

		// Utility commands (cmd_util.go)
		"help": handleHelp,
		"exit": handleExit,
		"quit": handleExit,
	}
}

// handlers maps command names to their handler functions.
var handlers = getHandlers()

// handleCommand routes user input to the appropriate handler.
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
