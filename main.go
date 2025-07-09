package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jarmoCluyse/ads-go/pkg/ads"
)

func main() {
	log.Println("main: Starting application")

	settings := ads.ClientSettings{
		TargetNetID: "192.168.157.131.1.1", // Replace with your target NetID
		// TargetNetID:   "127.0.0.1.1.1",   // Replace with your target NetID
		TargetPort:    851,               // Replace with your target port
		RouterAddr:    "127.0.0.1:48898", // Replace with your router address
		Timeout:       5 * time.Second,   // Set a timeout for ADS requests
		AllowHalfOpen: true,              // Set to true to allow connecting to a PLC in non-RUN state
	}

	log.Println("main: Creating new ADS client with settings:", settings)
	client := ads.NewClient(settings)
	log.Println("main: ADS client created.")

	log.Println("main: Attempting to connect to ADS router...")
	if err := client.Connect(); err != nil {
		log.Fatalf("main: Failed to connect: %v", err)
	}
	defer func() {
		log.Println("main: Disconnecting from ADS router...")
		if err := client.Disconnect(); err != nil {
			log.Printf("main: Error during disconnect: %v", err)
		}
		log.Println("main: Disconnected from ADS router.")
	}()

	fmt.Println("Connected to ADS router")
	log.Println("main: Successfully connected to ADS router.")

	// Example: Read Device Info
	log.Println("main: Reading device info...")
	deviceInfo, err := client.ReadDeviceInfo()
	if err != nil {
		log.Fatalf("main: Failed to read device info: %v", err)
	}
	log.Printf("main: Device Info: %+v\n", deviceInfo)
	fmt.Printf("Device Info: %+v\n", deviceInfo)

	// Example: Read ADS State
	log.Println("main: Reading ADS state...")
	state, err := client.ReadTcSystemState()
	if err != nil {
		log.Fatalf("main: Failed to read state: %v", err)
	}
	log.Printf("main: ADS State: %+v\n", state)
	fmt.Printf("ADS State: %+v\n", state)

	// Example: Set TwinCAT to Config mode
	log.Println("main: Setting TwinCAT to Config mode...")
	// err = client.SetTcSystemToConfig()
	err = client.SetTcSystemToRun(true)
	if err != nil {
		log.Fatalf("main: Failed to set TwinCAT to Config mode: %v", err)
	}
	log.Println("main: TwinCAT set to Config mode.")
	fmt.Println("TwinCAT set to Config mode.")

	// Read ADS State again to confirm
	log.Println("main: Reading ADS state after setting to Config...")
	state, err = client.ReadTcSystemState()
	if err != nil {
		log.Fatalf("main: Failed to read state after setting to Config: %v", err)
	}
	log.Printf("main: ADS State after Config: %+v\n", state)
	fmt.Printf("ADS State after Config: %+v\n", state)

	log.Println("main: Application finished.")
}
