package cli

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/jarmoCluyse/ads-go/pkg/ads"
)

var handlers = map[string]func([]string, *ads.Client){
	"state": func(args []string, client *ads.Client) {
		state, err := client.ReadTcSystemState()
		if err != nil {
			fmt.Println("Error reading system state:", err)
		}
		fmt.Printf("System state %v\n", state)
	},
	"state_loop": func(args []string, client *ads.Client) {
		// Run state every 550ms until ^C
		fmt.Println("Press Ctrl+C to exit state loop...")
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
					fmt.Println("Error reading system state:", err)
				} else {
					fmt.Printf("System state %v\n", state)
				}
			case <-done:
				fmt.Println("Exiting state loop.")
				return
			}
		}
	},
	"toConfig": func(args []string, client *ads.Client) {
		err := client.SetTcSystemToConfig()
		if err != nil {
			fmt.Println("Error reading system state:", err)
		}
	},

	"toRun": func(args []string, client *ads.Client) {
		err := client.SetTcSystemToRun()
		if err != nil {
			fmt.Println("Error reading system state:", err)
		}
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
		fmt.Println("Unknown command:", parts[0])
	}
}
