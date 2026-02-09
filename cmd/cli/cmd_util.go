package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

// handleHelp displays a list of available commands.
// Usage: help
func handleHelp(args []string, client *ads.Client) {
	fmt.Println("\nAvailable commands:")

	// Get sorted command names from getHandlers() to avoid initialization cycle
	handlerMap := getHandlers()
	commands := make([]string, 0, len(handlerMap))
	for cmd := range handlerMap {
		commands = append(commands, cmd)
	}
	sort.Strings(commands)

	// Print commands in columns
	for _, cmd := range commands {
		fmt.Printf("  %s\n", cmd)
	}
	fmt.Println()
}

// handleExit exits the CLI gracefully.
// Usage: exit or quit
func handleExit(args []string, client *ads.Client) {
	fmt.Println("[INFO] Exiting CLI...")
	os.Exit(0)
}
