/*
Package ads provides a Go client library for Beckhoff TwinCAT ADS protocol.

# Overview

The ads package implements the Automation Device Specification (ADS) protocol
for communicating with Beckhoff TwinCAT automation systems. It supports both
TwinCAT 2 and TwinCAT 3 runtimes and provides automatic type conversion between
PLC and Go data types.

# Features

  - TwinCAT 2 and 3 support
  - Connect to local TwinCAT 3 runtime
  - Connect to remote PLC systems over TCP/IP
  - Read and write any variable type with automatic conversion
  - Support for primitives, structs, arrays, and enums
  - Symbol and data type introspection
  - PLC state control (start, stop, config mode)
  - Device information reading
  - Raw memory operations for advanced use cases
  - Structured logging support (log/slog)

# Quick Start

Basic connection and read/write example:

	package main

	import (
		"fmt"
		"log"

		"github.com/jarmocluyse/ads-go/pkg/ads"
	)

	func main() {
		// Create client
		client := ads.NewClient(ads.ClientSettings{
			TargetAmsNetId: "192.168.1.120.1.1",
			TargetAdsPort:  851,
		}, nil)

		// Connect
		if err := client.Connect(); err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect()

		// Read a value
		value, err := client.ReadValue(851, "GVL_Main.Counter")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Counter value: %v\n", value)

		// Write a value
		err = client.WriteValue(851, "GVL_Main.Counter", 42)
		if err != nil {
			log.Fatal(err)
		}
	}

# Connection Setup

The client supports multiple connection modes depending on your environment.

Setup 1 - Windows with TwinCAT Router installed:

	client := ads.NewClient(ads.ClientSettings{
		TargetAmsNetId: "192.168.1.120.1.1",
		TargetAdsPort:  851,
	}, nil)

Setup 2 - Direct connection from any system (Linux, Raspberry Pi, etc.):

	client := ads.NewClient(ads.ClientSettings{
		LocalAmsNetId:  "192.168.1.10.1.1",  // Your system's AmsNetId
		LocalAdsPort:   32750,                // Any unused port
		TargetAmsNetId: "192.168.1.120.1.1", // PLC's AmsNetId
		TargetAdsPort:  851,
		RouterAddress:  "192.168.1.120",     // PLC's IP address
		RouterTcpPort:  48898,                // ADS router port
	}, nil)

Setup 3 - Localhost connection:

	client := ads.NewClient(ads.ClientSettings{
		TargetAmsNetId: "127.0.0.1.1.1", // or "localhost"
		TargetAdsPort:  851,
	}, nil)

For direct connections, you need to configure a static route on the target PLC.
See the README.md for detailed setup instructions.

# Reading Values

The ReadValue method automatically resolves symbol information and converts
PLC data to appropriate Go types.

Reading primitives:

	// Read an integer
	value, err := client.ReadValue(851, "GVL.IntValue")
	if err != nil {
		log.Fatal(err)
	}
	intValue := value.(int32) // Type assertion

	// Read a boolean
	value, err = client.ReadValue(851, "GVL.BoolValue")
	boolValue := value.(bool)

	// Read a string
	value, err = client.ReadValue(851, "GVL.StringValue")
	stringValue := value.(string)

Reading structs:

	value, err := client.ReadValue(851, "GVL.MyStruct")
	if err != nil {
		log.Fatal(err)
	}

	// Structs are returned as map[string]any
	structMap := value.(map[string]any)
	field1 := structMap["Field1"]
	field2 := structMap["Field2"]

Reading arrays:

	value, err := client.ReadValue(851, "GVL.IntArray")
	if err != nil {
		log.Fatal(err)
	}

	// Arrays are returned as []any
	arrayValue := value.([]any)
	for i, item := range arrayValue {
		fmt.Printf("Array[%d] = %v\n", i, item)
	}

Reading enums:

	value, err := client.ReadValue(851, "GVL.StatusEnum")
	if err != nil {
		log.Fatal(err)
	}

	// Enums are returned as map with "name" and "value" fields
	enumMap := value.(map[string]any)
	enumName := enumMap["name"].(string)
	enumValue := enumMap["value"].(int32)

# Writing Values

The WriteValue method automatically converts Go values to PLC format.

Writing primitives:

	// Write an integer
	err := client.WriteValue(851, "GVL.IntValue", 42)

	// Write a boolean
	err = client.WriteValue(851, "GVL.BoolValue", true)

	// Write a string
	err = client.WriteValue(851, "GVL.StringValue", "Hello PLC")

Writing structs:

	structData := map[string]any{
		"Field1": 100,
		"Field2": true,
		"Field3": 3.14,
	}
	err := client.WriteValue(851, "GVL.MyStruct", structData)

Writing arrays:

	arrayData := []int{1, 2, 3, 4, 5}
	err := client.WriteValue(851, "GVL.IntArray", arrayData)

	// Or with []any for mixed types
	arrayData := []any{1, 2, 3, 4, 5}
	err := client.WriteValue(851, "GVL.IntArray", arrayData)

Writing enums:

	// By name
	err := client.WriteValue(851, "GVL.StatusEnum", "Running")

	// By value
	err := client.WriteValue(851, "GVL.StatusEnum", 100)

# Raw Operations

For performance-critical code or when you need direct memory access:

	// Read raw bytes
	data, err := client.ReadRaw(851, 0x4020, 0x1000, 4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Raw data: %x\n", data)

	// Write raw bytes
	rawData := []byte{0x01, 0x02, 0x03, 0x04}
	err = client.WriteRaw(851, 0x4020, 0x1000, rawData)

	// Combined read-write operation
	writeData := []byte{0x05, 0x06}
	readData, err := client.ReadWriteRaw(851, 0x4020, 0x1000, 4, writeData)

# Symbol and Type Information

Get metadata about PLC variables and types:

	// Get symbol information
	symbol, err := client.GetSymbol(851, "GVL.MyVariable")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("IndexGroup: 0x%x\n", symbol.IndexGroup)
	fmt.Printf("IndexOffset: 0x%x\n", symbol.IndexOffset)
	fmt.Printf("Size: %d bytes\n", symbol.Size)
	fmt.Printf("Type: %s\n", symbol.Type)

	// Get data type definition
	dataType, err := client.GetDataType("ST_MyStruct", 851)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Type name: %s\n", dataType.Name)
	fmt.Printf("Size: %d bytes\n", dataType.Size)
	// Access sub-items for struct fields
	for _, subItem := range dataType.SubItems {
		fmt.Printf("  Field: %s, Type: %s\n", subItem.Name, subItem.Type)
	}

# PLC Control

Control the PLC runtime state:

	// Set system to CONFIG mode (restart in config)
	err := client.SetTcSystemToConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Set system to RUN mode (restart and run)
	err = client.SetTcSystemToRun()
	if err != nil {
		log.Fatal(err)
	}

	// Read current system state
	state, err := client.ReadTcSystemState()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ADS State: %d, Device State: %d\n", state.AdsState, state.DeviceState)

# Device Information

Read information about the target device:

	info, err := client.ReadDeviceInfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Device: %s\n", info.DeviceName)
	fmt.Printf("Version: %d.%d.%d\n",
		info.MajorVersion,
		info.MinorVersion,
		info.BuildVersion)

# Data Types

Supported PLC data types and their Go equivalents:

	BOOL              -> bool
	INT, SINT         -> int8, int16
	DINT, LINT        -> int32, int64
	USINT, UINT       -> uint8, uint16
	UDINT, ULINT      -> uint32, uint64
	REAL              -> float32
	LREAL             -> float64
	STRING            -> string
	WSTRING           -> string (UTF-16LE)
	STRUCT            -> map[string]any
	ARRAY             -> []any
	ENUM              -> map[string]any (with "name" and "value" fields)

# Architecture

The package is organized into submodules for maintainability and testing:

	ads-errors      Parse and validate ADS error codes
	ads-header      Parse 8-byte ADS response headers
	ads-symbol      Parse ADS symbol information
	ads-datatype    Parse complex data type definitions
	ads-stateinfo   Parse system state and device info
	ads-primitives  Read/write primitive types
	ads-requests    Build ADS command payloads
	ads-serializer  Type serialization and deserialization
	ams-header      Parse AMS protocol packet headers
	ams-builder     Build AMS/TCP and AMS headers

This modular design allows each component to be tested independently and
promotes code reusability.

# Logging

The client uses structured logging via the standard log/slog package.
By default, logging is disabled (nil logger).

Enable logging with a custom logger:

	import "log/slog"

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	client := ads.NewClient(settings, logger)

Using JSON output:

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	client := ads.NewClient(settings, logger)

To disable logging (default):

	client := ads.NewClient(settings, nil)

# TwinCAT 2 Differences

When connecting to TwinCAT 2 systems, note these differences:

1. PLC runtime ADS port is 801 instead of 851:

	client := ads.NewClient(ads.ClientSettings{
		TargetAmsNetId: "192.168.1.120.1.1",
		TargetAdsPort:  801, // TwinCAT 2
	}, nil)

2. All variable and data type names are UPPERCASE:

	// TwinCAT 3
	client.ReadValue(851, "GVL_Test.Counter")

	// TwinCAT 2
	client.ReadValue(801, ".COUNTER") // UPPERCASE

3. Global variables use dot prefix without GVL name:

	// TwinCAT 3
	client.ReadValue(851, "GVL_Main.Variable")

	// TwinCAT 2
	client.ReadValue(801, ".VARIABLE") // Dot prefix, no GVL

4. ENUMs return numeric values only (no name strings)

5. Empty structs and function blocks cannot be read

# Examples

Complete working examples can be found in:
  - cmd/main.go - Command-line interface example
  - example/ - Additional usage examples

# API Documentation

For detailed API documentation, see:
https://pkg.go.dev/github.com/jarmocluyse/ads-go/pkg/ads

# Credits

Documentation structure inspired by jisotalo/ads-client (used with permission).
https://github.com/jisotalo/ads-client

Created by Jarmo Cluyse - https://github.com/jarmocluyse
*/
package ads
