package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jarmoCluyse/ads-go/pkg/ads"
)

func main() {
	settings := ads.ClientSettings{
		TargetNetID:   "127.0.0.1.1.1",   // Replace with your target NetID
		TargetPort:    851,               // Replace with your target port
		RouterAddr:    "127.0.0.1:48898", // Replace with your router address
		Timeout:       5 * time.Second,   // Set a timeout for ADS requests
		AllowHalfOpen: true,              // Set to true to allow connecting to a PLC in non-RUN state
	}

	client := ads.NewClient(settings)

	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	fmt.Println("Connected to ADS router")

	// Example: Read Device Info
	deviceInfo, err := client.ReadDeviceInfo()
	if err != nil {
		log.Fatalf("Failed to read device info: %v", err)
	}

	fmt.Printf("Device Info: %+v\n", deviceInfo)

	// Example: Read ADS State
	state, err := client.ReadState()
	if err != nil {
		log.Fatalf("Failed to read state: %v", err)
	}
	fmt.Printf("ADS State: %+v\n", state)

	// Example: Set to Config Mode (uncomment to use)
	// if _, err := client.SetToConfig(); err != nil {
	// 	log.Fatalf("Failed to set to config: %v", err)
	// }
	// fmt.Println("Set to Config mode")

	// Example: Set to Run Mode (uncomment to use)
	// if _, err := client.SetToRun(); err != nil {
	// 	log.Fatalf("Failed to set to run: %v", err)
	// }
	// fmt.Println("Set to Run mode")

	// Example: Read a variable (replace with your details)
	// readResp, err := client.Read(0x4020, 0, 4)
	// if err != nil {
	// 	log.Fatalf("Failed to read: %v", err)
	// }
	// fmt.Printf("Read response: %+v\n", readResp)

	// Example: Write to a variable (replace with your details)
	// writeData := []byte{1, 2, 3, 4}
	// writeResp, err := client.Write(0x4020, 0, writeData)
	// if err != nil {
	// 	log.Fatalf("Failed to write: %v", err)
	// }
	// fmt.Printf("Write response: %+v\n", writeResp)
}
