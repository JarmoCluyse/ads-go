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
	"device_info": func(args []string, client *ads.Client) {
		info, err := client.ReadDeviceInfo()
		if err != nil {
			fmt.Println("Error reading system info:", err)
		}
		fmt.Printf("System version: %d.%d.%d\n", info.MajorVersion, info.MinorVersion, info.VersionBuild)
		fmt.Printf("System name: %v\n", info.DeviceName)
	},
	"state": func(args []string, client *ads.Client) {
		state, err := client.ReadTcSystemState()
		if err != nil {
			fmt.Println("Error reading system state:", err)
		}
		fmt.Printf("System state %v\n", state.AdsState.String())
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
					fmt.Printf("System state %v\n", state.AdsState.String())
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
	"read_value": func(args []string, client *ads.Client) {
		data := "service_interface.input.in_servicetool_serviceint_cmd"
		var port uint16 = 350
		value, err := client.ReadValue(port, data)
		if err != nil {
			fmt.Println("Error reading system state:", err)
		}
		fmt.Printf("value %v\n", value)
	},

	"read_bool": func(args []string, client *ads.Client) {
		data := "Service_interface.Input.IN_busInfo_Main.flgReadyCmd"
		var port uint16 = 350
		value, err := client.ReadValue(port, data)
		if err != nil {
			fmt.Println("Error reading bool:", err)
		}
		fmt.Printf("value %v\n", value)
	},
	"read_object": func(args []string, client *ads.Client) {
		data := "Service_interface.Input.IN_busInfo_Main"
		var port uint16 = 350
		value, err := client.ReadValue(port, data)
		if err != nil {
			fmt.Println("Error reading system state:", err)
		}
		fmt.Printf("value %v\n", value)
	},
	"read_raw": func(args []string, client *ads.Client) {
		indexGroup := uint32(0x1010290)
		indexOffset := uint32(0x80000001)
		var port uint16 = 350
		size := uint32(1) // adjust size as needed
		result, err := client.ReadRaw(port, indexGroup, indexOffset, size)
		if err != nil {
			fmt.Println("Error reading raw data:", err)
			return
		}
		// response[0:3] error
		// response[4:7] length
		// response[8:x] data
		fmt.Printf("Raw read [IG: 0x%X, IO: 0x%X]: %v\n", indexGroup, indexOffset, result)
	},
	"write_raw": func(args []string, client *ads.Client) {
		indexGroup := uint32(0x1010290)
		indexOffset := uint32(0x80000001)
		var port uint16 = 350
		data := []byte{36}
		err := client.WriteRaw(port, indexGroup, indexOffset, data)
		if err != nil {
			fmt.Println("Error writing raw data:", err)
			return
		}
		// response[0:3] error
		fmt.Printf("Raw write [IG: 0x%X, IO: 0x%X] succeeded%v\n", indexGroup, indexOffset)
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
