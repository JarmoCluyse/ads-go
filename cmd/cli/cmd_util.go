package cli

import (
	"fmt"
	"os"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

// handleHelp displays a list of available commands.
// Usage: help
func handleHelp(args []string, client *ads.Client) {
	fmt.Println("\n=== ADS CLI - Available Commands ===")
	fmt.Println()

	fmt.Println("System Commands:")
	fmt.Println("  device_info              - Get device information")
	fmt.Println("  state                    - Read current TwinCAT state")
	fmt.Println("  state_loop               - Continuously monitor TwinCAT state")
	fmt.Println("  monitor                  - Monitor system notifications")
	fmt.Println("  set_state <config|run>   - Switch TwinCAT state")

	fmt.Println("\nRead Commands:")
	fmt.Println("  read_value               - Read GLOBAL.gMyInt")
	fmt.Println("  read_bool                - Read GLOBAL.gMyBool")
	fmt.Println("  read_object              - Read GLOBAL.gMyDUT (struct)")
	fmt.Println("  read_array               - Read GLOBAL.gIntArray")

	fmt.Println("\nWrite Commands:")
	fmt.Println("  write_value <int>        - Write integer to GLOBAL.gMyInt")
	fmt.Println("  write_bool <true|false>  - Write boolean to GLOBAL.gMyBool")
	fmt.Println("  write_object Counter=<int> Ready=<bool> - Write to GLOBAL.gMyDUT")
	fmt.Println("  write_array <i1> <i2> <i3> <i4> <i5>    - Write 5 ints to GLOBAL.gIntArray")

	fmt.Println("\nRaw Commands:")
	fmt.Println("  read_raw                 - Read raw data by index/offset")
	fmt.Println("  write_raw                - Write raw data by index/offset")

	fmt.Println("\nSubscription Commands:")
	fmt.Println("  subscribe [path]         - Subscribe to variable changes (default: GLOBAL.gMyBoolToogle)")
	fmt.Println("  list_subs                - List active subscriptions")
	fmt.Println("  unsubscribe <id>         - Unsubscribe by ID")
	fmt.Println("  unsubscribe_all          - Unsubscribe from all")

	fmt.Println("\nSubscription Shortcuts (Quick subscriptions to example project variables):")
	fmt.Println("  sub_counter              - Subscribe to cycle-based counter (gMyIntCounter)")
	fmt.Println("  sub_toggle               - Subscribe to cycle-based toggle (gMyBoolToogle)")
	fmt.Println("  sub_timed_counter        - Subscribe to time-based counter (gTimedIntCounter)")
	fmt.Println("  sub_timed_toggle         - Subscribe to time-based toggle (gTimedBoolToogle)")
	fmt.Println("  sub_all                  - Subscribe to all counters/toggles at once")

	fmt.Println("\nControl Commands (for example project):")
	fmt.Println("  enable_counter <bool>         - Enable/disable cycle-based counter")
	fmt.Println("  enable_toggle <bool>          - Enable/disable cycle-based toggle")
	fmt.Println("  enable_timed_counter <bool>   - Enable/disable time-based counter")
	fmt.Println("  enable_timed_toggle <bool>    - Enable/disable time-based toggle")
	fmt.Println("  read_counters                 - Read all counter and toggle values")
	fmt.Println("  reset_counters                - Reset all counters to 0")
	fmt.Println("  read_status                   - Read all enable flag states")
	fmt.Println("  set_period <seconds>          - Set cycle period (1-3600 seconds)")
	fmt.Println("  read_period                   - Read current cycle period")

	fmt.Println("\nUtility Commands:")
	fmt.Println("  help                     - Show this help message")
	fmt.Println("  exit, quit               - Exit the CLI")
	fmt.Println()
}

// handleExit exits the CLI gracefully.
// Usage: exit or quit
func handleExit(args []string, client *ads.Client) {
	fmt.Println("[INFO] Exiting CLI...")
	os.Exit(0)
}
