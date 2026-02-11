package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/chzyer/readline"
	"github.com/jarmocluyse/ads-go/pkg/ads"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// getPrompt returns the appropriate prompt based on client state
func getPrompt(client *ads.Client) string {
	state := client.GetCurrentState()

	if state == nil {
		// State not available yet (initializing or disconnected)
		return "⚪ > "
	}

	switch state.AdsState {
	case types.ADSStateRun:
		return "🟢 > " // Green - Running
	case types.ADSStateConfig:
		return "🔵 > " // Blue - Config mode
	case types.ADSStateStop:
		return "🔴 > " // Red - Stopped
	case types.ADSStateError:
		return "❌ > " // Error
	default:
		return "⚪ > " // Unknown state
	}
}

// getUnsubscribeCompletions returns dynamic completions for unsubscribe command
func getUnsubscribeCompletions(line string) []string {
	subscriptionsMutex.RLock()
	defer subscriptionsMutex.RUnlock()

	completions := []string{}
	for id := range subscriptions {
		// Format: "ID"
		completions = append(completions, fmt.Sprintf("%d", id))
	}
	return completions
}

func Commandline(client *ads.Client) {
	// Enhanced autocomplete with nested completions
	completer := readline.NewPrefixCompleter(
		// System commands
		readline.PcItem("device_info"),
		readline.PcItem("state"),
		readline.PcItem("state_loop"),
		readline.PcItem("monitor"),
		readline.PcItem("set_state",
			readline.PcItem("config"),
			readline.PcItem("run"),
		),

		// Read commands
		readline.PcItem("read_value"),
		readline.PcItem("read_bool"),
		readline.PcItem("read_object"),
		readline.PcItem("read_array"),

		// Write commands with argument completions
		readline.PcItem("write_value"),
		readline.PcItem("write_bool",
			readline.PcItem("true"),
			readline.PcItem("false"),
		),
		readline.PcItem("write_object",
			readline.PcItem("Counter="),
			readline.PcItem("Ready="),
		),
		readline.PcItem("write_array"),

		// Raw commands
		readline.PcItem("read_raw"),
		readline.PcItem("write_raw"),

		// Subscription commands with variable path completions
		readline.PcItem("subscribe",
			readline.PcItem("GLOBAL.gMyIntCounter"),
			readline.PcItem("GLOBAL.gMyBoolToogle"),
			readline.PcItem("GLOBAL.gTimedIntCounter"),
			readline.PcItem("GLOBAL.gTimedBoolToogle"),
			readline.PcItem("GLOBAL.gMyInt"),
			readline.PcItem("GLOBAL.gMyBool"),
			readline.PcItem("GLOBAL.gMyDINT"),
			readline.PcItem("GLOBAL.gIntCounterActive"),
			readline.PcItem("GLOBAL.gBoolToggleActive"),
			readline.PcItem("GLOBAL.gTimedCounterActive"),
			readline.PcItem("GLOBAL.gTimedToggleActive"),
			readline.PcItem("GLOBAL.gCyclePeriod"),
			readline.PcItem("GLOBAL.gIntArray"),
			readline.PcItem("GLOBAL.gMyDUT"),
		),
		readline.PcItem("list_subs"),
		readline.PcItemDynamic(func(line string) []string {
			return getUnsubscribeCompletions(line)
		},
			readline.PcItem("unsubscribe"),
		),
		readline.PcItem("unsubscribe_all"),

		// Subscription shortcuts
		readline.PcItem("sub_counter"),
		readline.PcItem("sub_toggle"),
		readline.PcItem("sub_timed_counter"),
		readline.PcItem("sub_timed_toggle"),
		readline.PcItem("sub_all"),

		// Control commands with boolean argument completions
		readline.PcItem("enable_counter",
			readline.PcItem("true"),
			readline.PcItem("false"),
		),
		readline.PcItem("enable_toggle",
			readline.PcItem("true"),
			readline.PcItem("false"),
		),
		readline.PcItem("enable_timed_counter",
			readline.PcItem("true"),
			readline.PcItem("false"),
		),
		readline.PcItem("enable_timed_toggle",
			readline.PcItem("true"),
			readline.PcItem("false"),
		),
		readline.PcItem("read_counters"),
		readline.PcItem("reset_counters"),
		readline.PcItem("read_status"),
		readline.PcItem("set_period"),
		readline.PcItem("read_period"),

		// Utility commands
		readline.PcItem("help"),
		readline.PcItem("exit"),
		readline.PcItem("quit"),
	)
	// Use readline to provide command history and up arrow support
	config := &readline.Config{
		Prompt:          getPrompt(client),
		AutoComplete:    completer,
		InterruptPrompt: "^C\n",
		EOFPrompt:       "exit\n",
	}
	rl, err := readline.NewEx(config)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize readline: %v", err))
	}
	defer func() {
		if err := rl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing readline: %v\n", err)
		}
	}()

	// Start a goroutine to update the prompt based on state changes
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		lastPrompt := getPrompt(client)
		for range ticker.C {
			newPrompt := getPrompt(client)
			if newPrompt != lastPrompt {
				rl.SetPrompt(newPrompt)
				rl.Refresh()
				lastPrompt = newPrompt
			}
		}
	}()

	for {
		line, err := rl.Readline()
		if err != nil {
			os.Exit(0)
		}

		if len(line) > 0 {
			handleCommand(line, client)
		}
	}
}
