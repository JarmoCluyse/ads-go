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

func Commandline(client *ads.Client) {
	// Auto-generate completer based on available commands
	items := []readline.PrefixCompleterInterface{}
	for cmd := range handlers {
		items = append(items, readline.PcItem(cmd))
	}
	completer := readline.NewPrefixCompleter(items...)
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
	defer rl.Close()

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
