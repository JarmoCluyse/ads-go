package cli

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

// handleDeviceInfo reads and displays device information.
// Usage: device_info
func handleDeviceInfo(args []string, client *ads.Client) {
	info, err := client.ReadDeviceInfo()
	if err != nil {
		fmt.Printf("[ERROR] Command 'device_info': Failed to read device info: %v\n", err)
		return
	}
	fmt.Printf("[OK] Device version: %d.%d.%d\n", info.MajorVersion, info.MinorVersion, info.VersionBuild)
	fmt.Printf("[OK] Device name: %s\n", info.DeviceName)
}

// handleState reads and displays the TwinCAT system state.
// Usage: state
func handleState(args []string, client *ads.Client) {
	state, err := client.ReadTcSystemState()
	if err != nil {
		fmt.Printf("[ERROR] Command 'state': Failed to read system state: %v\n", err)
		return
	}
	fmt.Printf("[OK] System ADS state: %s\n", state.AdsState.String())
}

// handleStateLoop continuously reads and displays the system state every 500ms until Ctrl+C.
// Usage: state_loop
func handleStateLoop(args []string, client *ads.Client) {
	// Run state every 500ms until ^C
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
}

// handleSetState sets the TwinCAT system state.
// Usage: set_state <config|run>
func handleSetState(args []string, client *ads.Client) {
	if len(args) == 0 {
		fmt.Println("[ERROR] Command 'set_state': No state provided. Use 'config' or 'run'.")
		return
	}

	state := args[0]

	switch state {
	case "config":
		err := client.SetTcSystemToConfig()
		if err != nil {
			fmt.Printf("[ERROR] Command 'set_state': Failed to set TwinCAT to Config state: %v\n", err)
			return
		}
		fmt.Println("[OK] TwinCAT system set to Config state.")
	case "run":
		err := client.SetTcSystemToRun()
		if err != nil {
			fmt.Printf("[ERROR] Command 'set_state': Failed to set TwinCAT to Run state: %v\n", err)
			return
		}
		fmt.Println("[OK] TwinCAT system set to Run state.")
	default:
		fmt.Printf("[ERROR] Command 'set_state': Invalid state '%s'. Use 'config' or 'run'.\n", state)
	}
}

// handleMonitor displays current TwinCAT state and information about background monitoring.
// Usage: monitor
func handleMonitor(args []string, client *ads.Client) {
	currentState := client.GetCurrentState()
	if currentState == nil {
		fmt.Println("[INFO] TwinCAT state not available yet (still initializing)")
		fmt.Println("[INFO] Background state monitoring is active (2s interval)")
		fmt.Println("[INFO] State changes will be logged automatically")
		return
	}

	fmt.Printf("[INFO] Current TwinCAT state: %s\n", currentState.AdsState.String())
	fmt.Printf("[INFO] Device state: %d\n", currentState.DeviceState)
	fmt.Println()
	fmt.Println("[INFO] Prompt Indicators:")
	fmt.Println("  🟢 > = TwinCAT in Run mode (operations available)")
	fmt.Println("  🔵 > = TwinCAT in Config mode (operations blocked)")
	fmt.Println("  🔴 > = TwinCAT stopped (operations blocked)")
	fmt.Println("  ❌ > = TwinCAT in error state")
	fmt.Println("  ⚪ > = State unknown or disconnected")
	fmt.Println()
	fmt.Println("[INFO] Background state monitoring is active (2s interval)")
	fmt.Println("[INFO] State changes are logged automatically")
	fmt.Println("[INFO] When state leaves Run mode, auto-reconnection is triggered")
}
