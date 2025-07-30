package cli

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/jarmocluyse/ads-go/pkg/ads"
)

func Commandline(client *ads.Client) {
	// Auto-generate completer based on available commands
	items := []readline.PrefixCompleterInterface{}
	for cmd := range handlers {
		items = append(items, readline.PcItem(cmd))
	}
	completer := readline.NewPrefixCompleter(items...)
	// Use readline to provide command history and up arrow support
	config := &readline.Config{
		Prompt:          "🚀 > ",
		AutoComplete:    completer,
		InterruptPrompt: "^C\n",
		EOFPrompt:       "exit\n",
	}
	rl, err := readline.NewEx(config)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize readline: %v", err))
	}
	defer rl.Close()

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
