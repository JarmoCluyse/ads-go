# ads-go

[![Go Reference](https://pkg.go.dev/badge/github.com/jarmocluyse/ads-go.svg)](https://pkg.go.dev/github.com/jarmocluyse/ads-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/jarmocluyse/ads-go)](https://goreportcard.com/report/github.com/jarmocluyse/ads-go)
[![Test](https://github.com/JarmoCluyse/ads-go/actions/workflows/test.yml/badge.svg)](https://github.com/JarmoCluyse/ads-go/actions/workflows/test.yml)
[![Lint](https://github.com/JarmoCluyse/ads-go/actions/workflows/lint.yml/badge.svg)](https://github.com/JarmoCluyse/ads-go/actions/workflows/lint.yml)
[![Build](https://github.com/JarmoCluyse/ads-go/actions/workflows/build.yml/badge.svg)](https://github.com/JarmoCluyse/ads-go/actions/workflows/build.yml)
[![Coverage](https://codecov.io/gh/JarmoCluyse/ads-go/branch/main/graph/badge.svg)](https://codecov.io/gh/JarmoCluyse/ads-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/jarmocluyse/ads-go)](https://github.com/jarmocluyse/ads-go)
[![Release](https://img.shields.io/github/v/release/jarmocluyse/ads-go?include_prereleases)](https://github.com/jarmocluyse/ads-go/releases)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Beckhoff TwinCAT ADS client library for Go (unofficial).

Connect to a Beckhoff TwinCAT automation system using the ADS protocol from a Go application.

> **Note:** Documentation structure inspired by [jisotalo/ads-client](https://github.com/jisotalo/ads-client) (used with permission).

# Project Status

Active development. Core features are stable and tested.

**Implemented:**
- ✅ Connection management (connect, disconnect, port registration)
- ✅ Read/write operations with automatic type conversion
- ✅ Raw memory operations (ReadRaw, WriteRaw, ReadWriteRaw)
- ✅ Symbol and data type introspection
- ✅ PLC state control (config/run modes)
- ✅ Device information reading
- ✅ Full type support (primitives, structs, arrays, enums, strings)
- ✅ ADS notifications (subscriptions) with automatic change detection
- ✅ State monitoring with restart detection
- ✅ Connection lifecycle hooks (OnConnect, OnDisconnect, OnConnectionLost)

**Roadmap:**
- ⏳ Variable handle management
- ⏳ RPC method invocation
- ⏳ Batch operations (sum commands)

# Features

- Supports TwinCAT 2 and 3
- Supports connecting to local TwinCAT 3 runtime
- Supports any ADS-enabled target system (local runtime, remote PLC, I/O devices)
- Multiple connections from same host
- Reading and writing any variable type
- Automatic conversion between PLC and Go types
- Symbol and data type introspection
- PLC state control (start, stop, config mode)
- Device information reading
- Raw memory operations for advanced use cases
- Automatic 32/64-bit variable support (XINT, ULINT, etc.)
- Automatic byte alignment support (all pack-modes)
- ADS notifications/subscriptions with configurable cycle times
- Automatic TwinCAT state monitoring and restart detection
- Connection lifecycle hooks for robust error handling
- Structured logging support (log/slog)

# Table of Contents

- [Support](#support)
- [Installing](#installing)
- [Minimal Example (TLDR)](#minimal-example-tldr)
- [Connection Setup](#connection-setup)
  - [Setup 1 - Connect from Windows](#setup-1---connect-from-windows)
  - [Setup 2 - Connect from Linux/Windows with .NET Router](#setup-2---connect-from-linuxwindows-with-net-router)
  - [Setup 3 - Connect from any system (direct)](#setup-3---connect-from-any-system-direct)
  - [Setup 4 - Connect from local system](#setup-4---connect-from-local-system)
  - [Setup 5 - Docker container](#setup-5---docker-container)
- [Important](#important)
  - [Enabling localhost support on TwinCAT 3](#enabling-localhost-support-on-twincat-3)
  - [Structured variables](#structured-variables)
  - [Differences when using with TwinCAT 2](#differences-when-using-with-twincat-2)
- [Getting Started](#getting-started)
  - [Documentation](#documentation)
  - [Available Methods](#available-methods)
  - [Creating a Client](#creating-a-client)
  - [Connecting](#connecting)
  - [Reading Values](#reading-values)
  - [Writing Values](#writing-values)
  - [Raw Operations](#raw-operations)
  - [Symbol and Type Information](#symbol-and-type-information)
  - [PLC Control](#plc-control)
  - [State Monitoring & Event Handling](#state-monitoring--event-handling)
  - [Subscriptions & Notifications](#subscriptions--notifications)
  - [Device Information](#device-information)
  - [Logging](#logging)
  - [Disconnecting](#disconnecting)
- [Common Issues and Questions](#common-issues-and-questions)
- [Architecture](#architecture)
- [Roadmap](#roadmap)
- [Testing](#testing)
- [Examples](#examples)
- [License](#license)

# Support

- **Issues & bugs:** [GitHub Issues](https://github.com/jarmocluyse/ads-go/issues)
- **Discussions & help:** [GitHub Discussions](https://github.com/jarmocluyse/ads-go/discussions)

# Installing

```bash
go get github.com/jarmocluyse/ads-go@latest
```

Import in your code:

```go
import "github.com/jarmocluyse/ads-go/pkg/ads"
```

# Minimal Example (TLDR)

This connects to a local PLC runtime, reads a value, writes a value, reads it again and then disconnects. The value is a string located at `GVL_Global.StringValue`.

```go
package main

import (
	"fmt"
	"log"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

func main() {
	// Create client
	client := ads.NewClient(ads.ClientSettings{
		TargetAmsNetId: "localhost",
		TargetAdsPort:  851,
	}, nil)

	// Connect
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	fmt.Println("Connected to PLC")

	// Read a value
	value, err := client.ReadValue(851, "GVL_Global.StringValue")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Value read (before): %v\n", value)

	// Write a value
	err = client.WriteValue(851, "GVL_Global.StringValue", "New value from Go!")
	if err != nil {
		log.Fatal(err)
	}

	// Read again to verify
	value, err = client.ReadValue(851, "GVL_Global.StringValue")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Value read (after): %v\n", value)

	fmt.Println("Done!")
}
```

# Connection Setup

The ads-go client can be used with multiple system configurations.

![Connection Setup Diagram](./img/connection_setup.png)

## Setup 1 - Connect from Windows

This is the most common scenario. The client is running on a Windows PC that has TwinCAT Router installed (such as development laptop, Beckhoff IPC/PC, Beckhoff PLC).

**Requirements:**
- Client has one of the following installed:
  - TwinCAT XAE (development environment)
  - TwinCAT XAR (runtime)
  - [TwinCAT ADS](https://www.beckhoff.com/en-en/products/automation/twincat/tc1xxx-twincat-3-base/tc1000.html)
- An ADS route is created between the client and the PLC using TwinCAT router

**Client settings:**

```go
client := ads.NewClient(ads.ClientSettings{
	TargetAmsNetId: "192.168.1.120.1.1", // AmsNetId of the target PLC
	TargetAdsPort:  851,
}, nil)
```

## Setup 2 - Connect from Linux/Windows with .NET Router

In this scenario, the client is running on Linux or Windows without TwinCAT Router. The .NET based router can be run separately on the same machine.

**Requirements:**
- Client has .NET runtime installed
- Client has [AdsRouterConsoleApp](https://github.com/Beckhoff/TF6000_ADS_DOTNET_V5_Samples/tree/main/Sources/RouterSamples/AdsRouterConsoleApp) or similar running
- An ADS route is created between the client and the PLC (see AdsRouterConsoleApp docs)

**Client settings:**

```go
client := ads.NewClient(ads.ClientSettings{
	TargetAmsNetId: "192.168.1.120.1.1", // AmsNetId of the target PLC
	TargetAdsPort:  851,
}, nil)
```

## Setup 3 - Connect from any system (direct)

In this scenario, the client is running on a machine that has no router running (no TwinCAT router and no 3rd party router). For example, Raspberry Pi without any additional installations.

In this setup, the client directly connects to the PLC and uses its TwinCAT router for communication. Only one simultaneous connection from the client is possible.

**Requirements:**
- Target system (PLC) firewall has TCP port 48898 open
  - Windows Firewall might block, make sure Ethernet connection is handled as "private"
- Local AmsNetId and ADS port are set manually
  - Used `LocalAmsNetId` is not already in use
  - Used `LocalAdsPort` is not already in use
- An ADS route is configured to the PLC (see below)

**Setting up the route:**

1. At the PLC, open `C:\TwinCAT\3.1\Target\StaticRoutes.xml`
2. Copy paste the following under `<RemoteConnections>`:

```xml
<Route>
  <Name>GoClient</Name>
  <Address>192.168.1.10</Address>
  <NetId>192.168.1.10.1.1</NetId>
  <Type>TCP_IP</Type>
  <Flags>64</Flags>
</Route>
```

3. Edit `Address` to IP address of the client (which runs the Go app), such as `192.168.1.10`
4. Edit `NetId` to any unused AmsNetId address, such as `192.168.1.10.1.1`
5. Restart the PLC

**Client settings:**

```go
client := ads.NewClient(ads.ClientSettings{
	LocalAmsNetId:  "192.168.1.10.1.1",  // Same as NetId in StaticRoutes.xml
	LocalAdsPort:   32750,                // Can be anything that is not used
	TargetAmsNetId: "192.168.1.120.1.1", // AmsNetId of the target PLC
	TargetAdsPort:  851,
	RouterAddress:  "192.168.1.120",     // PLC IP address
	RouterTcpPort:  48898,
}, nil)
```

## Setup 4 - Connect from local system

In this scenario, the PLC is running the Go app locally. For example, the development PC or Beckhoff PLC with a screen for HMI.

**Requirements:**
- AMS router TCP loopback enabled (see [Enabling localhost support](#enabling-localhost-support-on-twincat-3))
  - Should be already enabled in TwinCAT versions >= 4024.5

**Client settings:**

```go
client := ads.NewClient(ads.ClientSettings{
	TargetAmsNetId: "127.0.0.1.1.1", // or "localhost"
	TargetAdsPort:  851,
}, nil)
```

## Setup 5 - Docker container

It's also possible to run the client in Docker containers, also with a separate router (Linux systems).

Contact me if you need help with Docker setup.

# Important

## Enabling localhost support on TwinCAT 3

If connecting to the local TwinCAT runtime (Go app and PLC on the same machine), the ADS router TCP loopback feature has to be enabled.

TwinCAT 4024.5 and newer already have this enabled as default.

1. Open registry editor (`regedit`)
2. Navigate to:

```
32-bit operating system:
  HKEY_LOCAL_MACHINE\SOFTWARE\Beckhoff\TwinCAT3\System\

64-bit operating system:
  HKEY_LOCAL_MACHINE\SOFTWARE\WOW6432Node\Beckhoff\TwinCAT3\System\
```

3. Create new DWORD registry entry named `EnableAmsTcpLoopback` with value of `1`
4. Restart the system

![Registry Setting](https://user-images.githubusercontent.com/13457157/82748398-2640bf00-9daa-11ea-98e5-0032b3537969.png)

Now you can connect to localhost using `TargetAmsNetId` address of `127.0.0.1.1.1` or `localhost`.

## Structured variables

When writing structured variables, the object properties are handled case-insensitively. This is because TwinCAT is case-insensitive.

In practice, it means that the following objects are equal when passed to `WriteValue()`:

```go
// These are equivalent in TwinCAT
map[string]any{
	"sometext": "hello",
	"somereal": 3.14,
}

map[string]any{
	"SOmeTEXT": "hello",
	"SOMEreal": 3.14,
}
```

If there are multiple properties with the same name (case-insensitive), the behavior is undefined.

## Differences when using with TwinCAT 2

**ADS port for the first PLC runtime is 801 instead of 851:**

```go
client := ads.NewClient(ads.ClientSettings{
	TargetAmsNetId: "192.168.1.120.1.1",
	TargetAdsPort:  801, // TwinCAT 2
}, nil)
```

**All variable and data type names are in UPPERCASE:**

This might cause problems if your app is used with both TC2 & TC3 systems.

![TwinCAT 2 Variables](https://user-images.githubusercontent.com/13457157/86540055-96df0d80-bf0a-11ea-8f94-7e04515213c2.png)

**Global variables are accessed with dot (`.`) prefix (without the GVL name):**

```go
// TwinCAT 3
client.ReadValue(851, "GVL_Test.ExampleSTRUCT")

// TwinCAT 2
client.ReadValue(801, ".EXAMPLESTRUCT")
```

**ENUMs are always numeric values only** (no name strings).

**Empty structs and function blocks (without members) can't be read.**

# Getting Started

## Documentation

Full API documentation is available at [https://pkg.go.dev/github.com/jarmocluyse/ads-go/pkg/ads](https://pkg.go.dev/github.com/jarmocluyse/ads-go/pkg/ads)

Complete working examples can be found in `cmd/main.go` and `example/` directory.

## Available Methods

| Method | Description |
|--------|-------------|
| `Connect()` | Establishes connection to target system |
| `Disconnect()` | Closes connection and cleans up resources |
| `ReadValue(port, path)` | Reads variable value by path with auto type conversion |
| `WriteValue(port, path, value)` | Writes variable value by path with auto type conversion |
| `ReadRaw(port, indexGroup, indexOffset, size)` | Reads raw bytes from memory |
| `WriteRaw(port, indexGroup, indexOffset, data)` | Writes raw bytes to memory |
| `ReadWriteRaw(port, indexGroup, indexOffset, readLength, writeData)` | Combined read-write operation |
| `GetSymbol(port, path)` | Retrieves symbol metadata (IndexGroup, IndexOffset, Size, Type) |
| `GetDataType(name, port)` | Retrieves complete data type definition |
| `BuildDataType(name, port)` | Recursively builds complex data type structures |
| `ReadDeviceInfo()` | Reads device name and version information |
| `ReadTcSystemState()` | Reads current TwinCAT system state |
| `ReadTcSystemExtendedState()` | Reads extended system state including restart index (TwinCAT 4022+) |
| `GetCurrentState()` | Returns cached current system state (updated by state monitoring) |
| `SetTcSystemToConfig()` | Sets TwinCAT system to CONFIG mode |
| `SetTcSystemToRun()` | Sets TwinCAT system to RUN mode |
| `WriteControl(adsState, deviceState, targetPort)` | Low-level state control |
| `SubscribeValue(port, path, callback, settings)` | Subscribe to variable value changes with automatic notifications |
| `Unsubscribe(subscription)` | Unsubscribe from a specific subscription |
| `UnsubscribeAll()` | Unsubscribe from all active subscriptions |

## Creating a Client

Settings are passed via the `ClientSettings` struct. The following settings are mandatory:
- `TargetAmsNetId` - Target runtime AmsNetId
- `TargetAdsPort` - Target runtime ADS port

```go
client := ads.NewClient(ads.ClientSettings{
	TargetAmsNetId: "localhost",
	TargetAdsPort:  851,
}, nil)
```

**With custom logger:**

```go
import "log/slog"

logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

client := ads.NewClient(ads.ClientSettings{
	TargetAmsNetId: "localhost",
	TargetAdsPort:  851,
}, logger)
```

## Connecting

It's good practice to start a connection at startup and keep it open until the app is closed.

```go
package main

import (
	"fmt"
	"log"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

func main() {
	client := ads.NewClient(ads.ClientSettings{
		TargetAmsNetId: "localhost",
		TargetAdsPort:  851,
	}, nil)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	fmt.Println("Connected to PLC")

	// Your code here...
}
```

## Reading Values

### Reading Primitives

Use `ReadValue()` to read any PLC value. The method automatically resolves the symbol and converts the value to an appropriate Go type.

**Reading INT:**

```go
value, err := client.ReadValue(851, "GVL_Read.StandardTypes.INT_")
if err != nil {
	log.Fatal(err)
}

// Type assertion to get specific type
intValue := value.(int16)
fmt.Printf("INT value: %d\n", intValue)
// Output: 32767
```

**Reading BOOL:**

```go
value, err := client.ReadValue(851, "GVL_Read.StandardTypes.BOOL_")
if err != nil {
	log.Fatal(err)
}

boolValue := value.(bool)
fmt.Printf("BOOL value: %v\n", boolValue)
// Output: true
```

**Reading REAL:**

```go
value, err := client.ReadValue(851, "GVL_Read.StandardTypes.REAL_")
if err != nil {
	log.Fatal(err)
}

realValue := value.(float32)
fmt.Printf("REAL value: %.2f\n", realValue)
// Output: 3.14
```

**Reading STRING:**

```go
value, err := client.ReadValue(851, "GVL_Read.StandardTypes.STRING_")
if err != nil {
	log.Fatal(err)
}

stringValue := value.(string)
fmt.Printf("STRING value: %s\n", stringValue)
// Output: Hello from PLC
```

### Reading Structs

Structs are returned as `map[string]any`:

```go
value, err := client.ReadValue(851, "GVL_Read.ComplexTypes.STRUCT_")
if err != nil {
	log.Fatal(err)
}

// Type assertion to map
structMap := value.(map[string]any)

// Access fields
boolField := structMap["BOOL_"].(bool)
intField := structMap["INT_"].(int16)
realField := structMap["REAL_"].(float32)

fmt.Printf("Struct fields: BOOL=%v, INT=%d, REAL=%.2f\n",
	boolField, intField, realField)

// Or print entire struct
fmt.Printf("Entire struct: %+v\n", structMap)
/* Output:
map[BOOL_:true BOOL_2:false BYTE_:255 WORD_:65535 ...]
*/
```

### Reading Arrays

Arrays are returned as `[]any`:

```go
value, err := client.ReadValue(851, "GVL_Read.StandardArrays.INT_5")
if err != nil {
	log.Fatal(err)
}

// Type assertion to slice
arrayValue := value.([]any)

fmt.Printf("Array length: %d\n", len(arrayValue))

// Access individual elements
for i, item := range arrayValue {
	intItem := item.(int16)
	fmt.Printf("Array[%d] = %d\n", i, intItem)
}
/* Output:
Array[0] = 10
Array[1] = 20
Array[2] = 30
Array[3] = 40
Array[4] = 50
*/
```

**Multidimensional arrays:**

```go
value, err := client.ReadValue(851, "GVL_Read.ComplexArrays.INT_2x3")
if err != nil {
	log.Fatal(err)
}

// Outer array
outerArray := value.([]any)

for i, row := range outerArray {
	// Inner array
	innerArray := row.([]any)
	fmt.Printf("Row %d: ", i)
	for _, item := range innerArray {
		fmt.Printf("%d ", item.(int16))
	}
	fmt.Println()
}
/* Output:
Row 0: 1 2 3
Row 1: 4 5 6
*/
```

### Reading Enums

Enums are returned as `map[string]any` with "name" and "value" fields:

```go
value, err := client.ReadValue(851, "GVL_Read.ComplexTypes.ENUM_")
if err != nil {
	log.Fatal(err)
}

enumMap := value.(map[string]any)
enumName := enumMap["name"].(string)
enumValue := enumMap["value"].(int32)

fmt.Printf("Enum: %s = %d\n", enumName, enumValue)
// Output: Running = 100
```

### Safe Type Assertions

Always use the comma-ok idiom for safe type assertions:

```go
value, err := client.ReadValue(851, "GVL.SomeValue")
if err != nil {
	log.Fatal(err)
}

// Safe type assertion
if intValue, ok := value.(int32); ok {
	fmt.Printf("Integer value: %d\n", intValue)
} else {
	fmt.Printf("Unexpected type: %T\n", value)
}
```

## Writing Values

### Writing Primitives

Use `WriteValue()` to write any PLC value.

**Writing INT:**

```go
err := client.WriteValue(851, "GVL_Write.StandardTypes.INT_", 42)
if err != nil {
	log.Fatal(err)
}
```

**Writing BOOL:**

```go
err := client.WriteValue(851, "GVL_Write.StandardTypes.BOOL_", true)
if err != nil {
	log.Fatal(err)
}
```

**Writing REAL:**

```go
err := client.WriteValue(851, "GVL_Write.StandardTypes.REAL_", 3.14)
if err != nil {
	log.Fatal(err)
}
```

**Writing STRING:**

```go
err := client.WriteValue(851, "GVL_Write.StandardTypes.STRING_", "Hello from Go!")
if err != nil {
	log.Fatal(err)
}
```

### Writing Structs

Write structs using `map[string]any`:

```go
structData := map[string]any{
	"BOOL_":  true,
	"INT_":   int16(100),
	"REAL_":  float32(2.71),
	"STRING": "Test",
}

err := client.WriteValue(851, "GVL_Write.ComplexTypes.STRUCT_", structData)
if err != nil {
	log.Fatal(err)
}
```

**Note:** Currently, partial struct updates require reading the existing value first, modifying it, then writing back:

```go
// Read existing value
value, err := client.ReadValue(851, "GVL_Write.ComplexTypes.STRUCT_")
if err != nil {
	log.Fatal(err)
}

// Modify specific field
structMap := value.(map[string]any)
structMap["INT_"] = int16(200)

// Write back
err = client.WriteValue(851, "GVL_Write.ComplexTypes.STRUCT_", structMap)
if err != nil {
	log.Fatal(err)
}
```

### Writing Arrays

Write arrays using slices:

```go
// Using []int
intArray := []int{1, 2, 3, 4, 5}
err := client.WriteValue(851, "GVL_Write.StandardArrays.INT_5", intArray)
if err != nil {
	log.Fatal(err)
}

// Or using []any
anyArray := []any{1, 2, 3, 4, 5}
err = client.WriteValue(851, "GVL_Write.StandardArrays.INT_5", anyArray)
if err != nil {
	log.Fatal(err)
}
```

**Multidimensional arrays:**

```go
// 2D array (2x3)
array2D := []any{
	[]any{1, 2, 3},
	[]any{4, 5, 6},
}

err := client.WriteValue(851, "GVL_Write.ComplexArrays.INT_2x3", array2D)
if err != nil {
	log.Fatal(err)
}
```

### Writing Enums

Write enums by name (string) or value (integer):

```go
// By name
err := client.WriteValue(851, "GVL_Write.ComplexTypes.ENUM_", "Running")
if err != nil {
	log.Fatal(err)
}

// By value
err = client.WriteValue(851, "GVL_Write.ComplexTypes.ENUM_", 100)
if err != nil {
	log.Fatal(err)
}
```

## Raw Operations

For performance-critical code or when you need direct memory access, use raw operations.

### ReadRaw

Read raw bytes from PLC memory:

```go
// Read 4 bytes from IndexGroup 0x4020, IndexOffset 0x1000
data, err := client.ReadRaw(851, 0x4020, 0x1000, 4)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Raw data: %x\n", data)
// Output: Raw data: 01020304
```

**Getting IndexGroup and IndexOffset from symbol:**

```go
// Get symbol info first
symbol, err := client.GetSymbol(851, "GVL.MyVariable")
if err != nil {
	log.Fatal(err)
}

// Use symbol info for raw read
data, err := client.ReadRaw(851, symbol.IndexGroup, symbol.IndexOffset, symbol.Size)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Read %d bytes: %x\n", len(data), data)
```

### WriteRaw

Write raw bytes to PLC memory:

```go
rawData := []byte{0x01, 0x02, 0x03, 0x04}

err := client.WriteRaw(851, 0x4020, 0x1000, rawData)
if err != nil {
	log.Fatal(err)
}

fmt.Println("Raw data written successfully")
```

### ReadWriteRaw

Combined read-write operation (useful for commands that require both):

```go
writeData := []byte{0x05, 0x06}

// Write 2 bytes and read 4 bytes in one operation
readData, err := client.ReadWriteRaw(851, 0x4020, 0x1000, 4, writeData)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Read data after write: %x\n", readData)
```

## Symbol and Type Information

Get metadata about PLC variables and data types.

### GetSymbol

Retrieve symbol information (IndexGroup, IndexOffset, Size, Type):

```go
symbol, err := client.GetSymbol(851, "GVL.MyVariable")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Symbol: %s\n", symbol.Name)
fmt.Printf("  Type: %s\n", symbol.Type)
fmt.Printf("  Size: %d bytes\n", symbol.Size)
fmt.Printf("  IndexGroup: 0x%x\n", symbol.IndexGroup)
fmt.Printf("  IndexOffset: 0x%x\n", symbol.IndexOffset)
fmt.Printf("  Comment: %s\n", symbol.Comment)

/* Output:
Symbol: GVL.MyVariable
  Type: INT
  Size: 2 bytes
  IndexGroup: 0x4020
  IndexOffset: 0x1000
  Comment: Counter variable
*/
```

### GetDataType

Retrieve complete data type definition:

```go
dataType, err := client.GetDataType("ST_MyStruct", 851)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Type: %s\n", dataType.Name)
fmt.Printf("Size: %d bytes\n", dataType.Size)
fmt.Printf("Offset: %d\n", dataType.Offset)

// Access struct fields (SubItems)
fmt.Println("Fields:")
for _, subItem := range dataType.SubItems {
	fmt.Printf("  %s: %s (offset %d, size %d)\n",
		subItem.Name,
		subItem.Type,
		subItem.Offset,
		subItem.Size)
}

/* Output:
Type: ST_MyStruct
Size: 16 bytes
Offset: 0
Fields:
  Field1: INT (offset 0, size 2)
  Field2: BOOL (offset 2, size 1)
  Field3: REAL (offset 4, size 4)
  Field4: STRING(10) (offset 8, size 11)
*/
```

**For arrays:**

```go
dataType, err := client.GetDataType("ARRAY[0..4] OF INT", 851)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Array info: %d dimensions\n", len(dataType.ArrayInfo))
for i, arrInfo := range dataType.ArrayInfo {
	fmt.Printf("  Dimension %d: Length=%d, LowerBound=%d, UpperBound=%d\n",
		i, arrInfo.Length, arrInfo.LowerBound, arrInfo.UpperBound)
}

/* Output:
Array info: 1 dimensions
  Dimension 0: Length=5, LowerBound=0, UpperBound=4
*/
```

## PLC Control

Control the PLC runtime state.

### SetTcSystemToConfig

Set TwinCAT system to CONFIG mode (restart in config):

```go
err := client.SetTcSystemToConfig()
if err != nil {
	log.Fatal(err)
}

fmt.Println("TwinCAT system set to CONFIG mode")
```

### SetTcSystemToRun

Set TwinCAT system to RUN mode (restart and run):

```go
err := client.SetTcSystemToRun()
if err != nil {
	log.Fatal(err)
}

fmt.Println("TwinCAT system set to RUN mode")
```

### ReadTcSystemState

Read current TwinCAT system state:

```go
state, err := client.ReadTcSystemState()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("ADS State: %d\n", state.AdsState)
fmt.Printf("Device State: %d\n", state.DeviceState)

// Common ADS states:
// 0 = Invalid
// 5 = Run
// 6 = Stop

/* Output:
ADS State: 5
Device State: 0
*/
```

### WriteControl (Low-level)

For advanced use cases, you can use the low-level `WriteControl` method:

```go
// Set to RUN state (AdsState=5, DeviceState=0)
err := client.WriteControl(5, 0, 851)
if err != nil {
	log.Fatal(err)
}

// Set to STOP state (AdsState=6, DeviceState=0)
err = client.WriteControl(6, 0, 851)
if err != nil {
	log.Fatal(err)
}
```

## State Monitoring & Event Handling

The client can automatically monitor TwinCAT system state changes and detect restarts. This is useful for handling connection issues, state transitions, and TwinCAT system restarts.

### Automatic State Monitoring

By default, the client checks the system state every 2 seconds and triggers event handlers when changes are detected.

**Key Features:**
- Detects state changes (Run ↔ Config ↔ Stop)
- Detects TwinCAT system restarts (even when state stays "Run")
- Auto-detects extended state support (TwinCAT 4022+)
- Thread-safe with automatic cleanup

### OnStateChange Hook

Called whenever the TwinCAT system state changes:

```go
settings := ads.ClientSettings{
	TargetAmsNetId: "localhost",
	TargetAdsPort:  851,
	
	// Called on state changes
	OnStateChange: func(client *ads.Client, newState, oldState *adsstateinfo.SystemState) {
		if oldState == nil {
			// Initial state after connection
			fmt.Printf("Initial state: %s\n", newState.AdsState.String())
		} else {
			// State changed
			fmt.Printf("State changed: %s → %s\n", 
				oldState.AdsState.String(), 
				newState.AdsState.String())
		}
	},
}

client := ads.NewClient(settings, nil)
```

**Common State Transitions:**
- `Run` → `Config`: PLC stopped for configuration
- `Config` → `Run`: PLC started after configuration
- `Run` → `Stop`: PLC stopped
- `Stop` → `Run`: PLC started

### OnConnectionLost Hook

Called when the connection is lost unexpectedly or TwinCAT restarts:

```go
settings := ads.ClientSettings{
	TargetAmsNetId: "localhost",
	TargetAdsPort:  851,
	
	// Called when connection drops or TwinCAT restarts
	OnConnectionLost: func(client *ads.Client, err error) {
		fmt.Printf("Connection lost: %v\n", err)
		
		// Re-read values and re-subscribe to notifications here
		// The ADS connection is still alive, but TwinCAT restarted
	},
}

client := ads.NewClient(settings, nil)
```

**When This Is Triggered:**
- TwinCAT state leaves "Run" mode (→ Config, Stop, Error, etc.)
- TwinCAT system restarts (detected via restart index change)
- Physical network connection drops

### TwinCAT Restart Detection

When TwinCAT restarts using `set_state run` command, the ADS state may remain "Run" but subscriptions are cleared. The client detects this by monitoring the **restart index** from extended system state.

**How It Works:**
1. On first state check, auto-detects extended state support
2. Monitors both `AdsState` AND `RestartIndex` on each poll
3. When `RestartIndex` changes → triggers `OnConnectionLost`
4. Works with TwinCAT 4022 and newer (gracefully falls back on older versions)

**Example Log Output:**

```
TwinCAT system restarted (restart index: 44 → 48)
EVENT: TwinCAT system state changed fromState=Run toState=Run
EVENT: ADS connection lost unexpectedly
```

### Reading Extended System State

For TwinCAT 4022 and newer, you can read extended system information including the restart index:

```go
extState, err := client.ReadTcSystemExtendedState()
if err != nil {
	// Extended state not supported or error
	log.Printf("Extended state not available: %v", err)
} else {
	fmt.Printf("Restart Index: %d\n", extState.RestartIndex)
	fmt.Printf("TwinCAT Version: %d.%d.%d\n", 
		extState.Version, extState.Revision, extState.Build)
	fmt.Printf("Platform: %d, OS Type: %d\n", 
		extState.Platform, extState.OsType)
}

/* Output:
Restart Index: 48
TwinCAT Version: 3.1.4024
Platform: 1, OS Type: 2
*/
```

**Extended State Fields:**
- `RestartIndex` (uint16): Increments on every TwinCAT restart
- `Version`, `Revision`, `Build`: TwinCAT version information
- `Platform`: Platform identifier (1=PC, 5=ARM, etc.)
- `OsType`: Operating system type (2=Windows, 10=Linux, etc.)
- `Flags`: System service state flags

### Getting Current State

Retrieve the cached current state (updated by background monitoring):

```go
currentState := client.GetCurrentState()
if currentState == nil {
	fmt.Println("State not available yet (still initializing)")
} else {
	fmt.Printf("Current state: %s\n", currentState.AdsState.String())
	
	// Check if PLC is running
	if currentState.AdsState == types.ADSStateRun {
		fmt.Println("PLC is running - operations available")
	}
}
```

### Customizing State Polling

Change the polling interval (default is 2 seconds):

```go
settings := ads.ClientSettings{
	TargetAmsNetId: "localhost",
	TargetAdsPort:  851,
	
	// Check state every 5 seconds
	StatePollingInterval: 5 * time.Second,
}

client := ads.NewClient(settings, nil)
```

**Disable state monitoring:**

```go
settings := ads.ClientSettings{
	TargetAmsNetId: "localhost",
	TargetAdsPort:  851,
	
	// Disable automatic state monitoring
	StatePollingInterval: 0,
}

client := ads.NewClient(settings, nil)
```

### Complete Example

Here's a complete example with state monitoring and reconnection logic:

```go
package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/jarmocluyse/ads-go/pkg/ads"
	"github.com/jarmocluyse/ads-go/pkg/ads/ads-stateinfo"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

func main() {
	settings := ads.ClientSettings{
		TargetAmsNetId: "localhost",
		TargetAdsPort:  851,
		
		// Monitor state changes
		OnStateChange: func(client *ads.Client, newState, oldState *adsstateinfo.SystemState) {
			if oldState == nil {
				fmt.Printf("Initial state: %s\n", newState.AdsState.String())
				return
			}
			
			fmt.Printf("State changed: %s → %s\n",
				oldState.AdsState.String(),
				newState.AdsState.String())
			
			// Detect Run mode entry
			if newState.AdsState == types.ADSStateRun && 
			   oldState.AdsState != types.ADSStateRun {
				fmt.Println("TwinCAT entered RUN mode")
				// Re-initialize your application logic here
			}
		},
		
		// Handle connection loss / restart
		OnConnectionLost: func(client *ads.Client, err error) {
			fmt.Printf("Connection lost: %v\n", err)
			
			// Wait for TwinCAT to come back to Run mode
			fmt.Println("Waiting for TwinCAT to return to Run mode...")
			
			for {
				time.Sleep(1 * time.Second)
				state := client.GetCurrentState()
				
				if state != nil && state.AdsState == types.ADSStateRun {
					fmt.Println("TwinCAT back in Run mode!")
					
					// Re-read values and re-subscribe here
					// Example: resubscribeToNotifications(client)
					break
				}
			}
		},
	}
	
	client := ads.NewClient(settings, nil)
	
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	
	fmt.Println("Connected - monitoring state changes...")
	
	// Your application logic here
	select {} // Keep running
}
```

**Key Points:**
- State monitoring runs automatically in the background
- Hooks are called asynchronously (don't block)
- ADS connection stays alive during TwinCAT restarts
- User must re-read values and re-subscribe after restart
- `GetCurrentState()` returns cached state (no network call)

## Subscriptions & Notifications

The client supports ADS notifications (subscriptions) for monitoring variable value changes in real-time. Instead of polling variables, you can subscribe to them and receive automatic notifications when values change.

### Key Features

- **Event-driven monitoring** - Get notified only when values change
- **Configurable cycle times** - Control how often values are checked (default: 100ms)
- **Change detection** - Option to send notifications only on value changes
- **Multiple subscriptions** - Subscribe to many variables simultaneously
- **Thread-safe** - Safe for concurrent access
- **Automatic cleanup** - Subscriptions are cleared on disconnect

### Basic Subscription

Subscribe to a variable and receive notifications when it changes:

```go
package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/jarmocluyse/ads-go/pkg/ads"
)

func main() {
	client := ads.NewClient(ads.ClientSettings{
		TargetAmsNetId: "localhost",
		TargetAdsPort:  851,
	}, nil)
	
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	
	// Define callback function
	callback := func(data ads.SubscriptionData) {
		fmt.Printf("Value changed: %v (at %s)\n", 
			data.Value, 
			data.Timestamp.Format("15:04:05.000"))
	}
	
	// Subscribe to a variable
	settings := ads.SubscriptionSettings{
		CycleTime:    100 * time.Millisecond,
		SendOnChange: true,
	}
	
	sub, err := client.SubscribeValue(851, "GVL.Counter", callback, settings)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Subscribed! Waiting for notifications...")
	
	// Keep running to receive notifications
	time.Sleep(30 * time.Second)
	
	// Unsubscribe when done
	if err := client.Unsubscribe(sub); err != nil {
		log.Printf("Error unsubscribing: %v", err)
	}
}

/* Output:
Subscribed! Waiting for notifications...
Value changed: 10 (at 14:23:15.123)
Value changed: 11 (at 14:23:15.223)
Value changed: 12 (at 14:23:15.323)
...
*/
```

### Subscription Settings

Control how notifications are sent using `SubscriptionSettings`:

```go
settings := ads.SubscriptionSettings{
	// How often to check the variable (required)
	CycleTime: 100 * time.Millisecond,
	
	// Only send notifications when value changes (default: false)
	// If false, notifications are sent every CycleTime
	SendOnChange: true,
}
```

**Recommended settings:**

```go
// Fast-changing values (motors, sensors)
fastSettings := ads.SubscriptionSettings{
	CycleTime:    50 * time.Millisecond,
	SendOnChange: true,
}

// Slow-changing values (temperature, status)
slowSettings := ads.SubscriptionSettings{
	CycleTime:    1 * time.Second,
	SendOnChange: true,
}

// Always notify (regardless of change)
alwaysSettings := ads.SubscriptionSettings{
	CycleTime:    100 * time.Millisecond,
	SendOnChange: false, // Sends every 100ms
}
```

### Subscription Data

The callback receives `SubscriptionData` with the following fields:

```go
type SubscriptionData struct {
	Value     any       // The variable value (with type conversion)
	Timestamp time.Time // When the notification was received
}
```

**Example callback with type assertion:**

```go
callback := func(data ads.SubscriptionData) {
	// Type assert to expected type
	if intValue, ok := data.Value.(int32); ok {
		fmt.Printf("Counter: %d\n", intValue)
	}
	
	// Or handle multiple types
	switch v := data.Value.(type) {
	case int32:
		fmt.Printf("Integer: %d\n", v)
	case bool:
		fmt.Printf("Boolean: %v\n", v)
	case float32:
		fmt.Printf("Float: %.2f\n", v)
	default:
		fmt.Printf("Unknown type: %v\n", v)
	}
}
```

### Multiple Subscriptions

Subscribe to multiple variables at once:

```go
// Track subscriptions
var subscriptions []*ads.ActiveSubscription

// Subscribe to multiple variables
variables := []string{
	"GVL.Counter",
	"GVL.Temperature",
	"GVL.IsRunning",
	"GVL.ErrorCode",
}

for _, varName := range variables {
	// Create callback for this variable
	callback := func(name string) ads.SubscriptionCallback {
		return func(data ads.SubscriptionData) {
			fmt.Printf("[%s] = %v\n", name, data.Value)
		}
	}(varName)
	
	// Subscribe
	settings := ads.SubscriptionSettings{
		CycleTime:    100 * time.Millisecond,
		SendOnChange: true,
	}
	
	sub, err := client.SubscribeValue(851, varName, callback, settings)
	if err != nil {
		log.Printf("Failed to subscribe to %s: %v", varName, err)
		continue
	}
	
	subscriptions = append(subscriptions, sub)
	fmt.Printf("Subscribed to %s\n", varName)
}

// Later: unsubscribe from all
for _, sub := range subscriptions {
	if err := client.Unsubscribe(sub); err != nil {
		log.Printf("Error unsubscribing: %v", err)
	}
}
```

### Unsubscribing

**Unsubscribe from a specific subscription:**

```go
sub, err := client.SubscribeValue(851, "GVL.Counter", callback, settings)
if err != nil {
	log.Fatal(err)
}

// ... later ...

if err := client.Unsubscribe(sub); err != nil {
	log.Printf("Error unsubscribing: %v", err)
}
```

**Unsubscribe from all active subscriptions:**

```go
if err := client.UnsubscribeAll(); err != nil {
	log.Printf("Error unsubscribing from all: %v", err)
}
```

**Note:** All subscriptions are automatically cleared when:
- `Disconnect()` is called
- TwinCAT system restarts (use `OnConnectionLost` hook to re-subscribe)

### Handling TwinCAT Restarts

When TwinCAT restarts, all subscriptions are cleared. Use the `OnConnectionLost` hook to automatically re-subscribe:

```go
// Track active subscriptions for re-subscription
var activeVars = []string{"GVL.Counter", "GVL.Temperature"}

settings := ads.ClientSettings{
	TargetAmsNetId: "localhost",
	TargetAdsPort:  851,
	
	// Re-subscribe after TwinCAT restart
	OnConnectionLost: func(client *ads.Client, err error) {
		fmt.Printf("Connection lost: %v\n", err)
		fmt.Println("Waiting for TwinCAT to return to Run mode...")
		
		// Wait for Run state
		for {
			time.Sleep(1 * time.Second)
			state := client.GetCurrentState()
			
			if state != nil && state.AdsState == types.ADSStateRun {
				fmt.Println("TwinCAT back in Run mode - re-subscribing...")
				
				// Re-subscribe to all variables
				for _, varName := range activeVars {
					callback := func(data ads.SubscriptionData) {
						fmt.Printf("[%s] = %v\n", varName, data.Value)
					}
					
					subSettings := ads.SubscriptionSettings{
						CycleTime:    100 * time.Millisecond,
						SendOnChange: true,
					}
					
					if _, err := client.SubscribeValue(851, varName, callback, subSettings); err != nil {
						log.Printf("Failed to re-subscribe to %s: %v", varName, err)
					} else {
						fmt.Printf("Re-subscribed to %s\n", varName)
					}
				}
				
				break
			}
		}
	},
}

client := ads.NewClient(settings, nil)
```

### Complete Working Example

Here's a complete example with subscriptions and proper lifecycle management:

```go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/jarmocluyse/ads-go/pkg/ads"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

func main() {
	// Track subscriptions for cleanup
	var subscriptions []*ads.ActiveSubscription
	
	settings := ads.ClientSettings{
		TargetAmsNetId: "localhost",
		TargetAdsPort:  851,
		
		// Handle TwinCAT restarts
		OnConnectionLost: func(client *ads.Client, err error) {
			fmt.Printf("Connection lost: %v\n", err)
			
			// Wait for Run state and re-subscribe
			for {
				time.Sleep(1 * time.Second)
				state := client.GetCurrentState()
				
				if state != nil && state.AdsState == types.ADSStateRun {
					fmt.Println("Re-subscribing...")
					subscribeToVariables(client, &subscriptions)
					break
				}
			}
		},
	}
	
	client := ads.NewClient(settings, nil)
	
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	
	fmt.Println("Connected! Creating subscriptions...")
	
	// Initial subscriptions
	subscribeToVariables(client, &subscriptions)
	
	fmt.Println("\nMonitoring variables. Press Ctrl+C to exit...")
	
	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	
	fmt.Println("\nShutting down...")
}

func subscribeToVariables(client *ads.Client, subscriptions *[]*ads.ActiveSubscription) {
	// Clear old subscriptions
	*subscriptions = nil
	
	variables := map[string]ads.SubscriptionSettings{
		"GVL.Counter": {
			CycleTime:    100 * time.Millisecond,
			SendOnChange: true,
		},
		"GVL.Temperature": {
			CycleTime:    500 * time.Millisecond,
			SendOnChange: true,
		},
		"GVL.IsRunning": {
			CycleTime:    200 * time.Millisecond,
			SendOnChange: true,
		},
	}
	
	for varName, settings := range variables {
		// Create callback for this variable
		callback := func(name string) ads.SubscriptionCallback {
			return func(data ads.SubscriptionData) {
				fmt.Printf("[%s] %s = %v\n",
					data.Timestamp.Format("15:04:05.000"),
					name,
					data.Value)
			}
		}(varName)
		
		// Subscribe
		sub, err := client.SubscribeValue(851, varName, callback, settings)
		if err != nil {
			log.Printf("Failed to subscribe to %s: %v", varName, err)
			continue
		}
		
		*subscriptions = append(*subscriptions, sub)
		fmt.Printf("✓ Subscribed to %s\n", varName)
	}
}
```

### Subscription Lifecycle

```
1. Connect to PLC
   ↓
2. Subscribe to variables
   ↓
3. Receive notifications automatically
   ↓
4. [TwinCAT restarts] → OnConnectionLost triggered
   ↓
5. Wait for Run state
   ↓
6. Re-subscribe to variables
   ↓
7. Continue receiving notifications
   ↓
8. Disconnect (automatic cleanup)
```

### Performance Considerations

**Cycle Time:**
- Shorter cycle times = more frequent checks = higher CPU usage
- Recommended minimum: 50ms
- Default: 100ms
- For slow-changing values: 500ms - 1s

**Send On Change:**
- Always enable `SendOnChange: true` when possible
- Reduces network traffic significantly
- Only use `SendOnChange: false` when you need guaranteed periodic updates

**Number of Subscriptions:**
- The client can handle many simultaneous subscriptions
- Each subscription is managed independently
- TwinCAT may have limits (typically hundreds of subscriptions)

### Common Patterns

**Subscribe to struct fields:**

```go
// Subscribe to individual fields
callback := func(data ads.SubscriptionData) {
	structValue := data.Value.(map[string]any)
	field1 := structValue["Field1"].(int32)
	field2 := structValue["Field2"].(bool)
	
	fmt.Printf("Field1=%d, Field2=%v\n", field1, field2)
}

settings := ads.SubscriptionSettings{
	CycleTime:    100 * time.Millisecond,
	SendOnChange: true,
}

sub, err := client.SubscribeValue(851, "GVL.MyStruct", callback, settings)
```

**Subscribe to array elements:**

```go
// Subscribe to entire array
callback := func(data ads.SubscriptionData) {
	arrayValue := data.Value.([]any)
	fmt.Printf("Array length: %d\n", len(arrayValue))
	
	for i, item := range arrayValue {
		fmt.Printf("  [%d] = %v\n", i, item)
	}
}

sub, err := client.SubscribeValue(851, "GVL.MyArray", callback, settings)
```

**Conditional notifications:**

```go
// Only log when value exceeds threshold
callback := func(data ads.SubscriptionData) {
	if temperature, ok := data.Value.(float32); ok {
		if temperature > 80.0 {
			fmt.Printf("⚠️ High temperature: %.1f°C\n", temperature)
		}
	}
}
```

### Troubleshooting

**Notifications not received:**
- Verify PLC is in RUN mode (`GetCurrentState()`)
- Check that the variable path is correct
- Ensure `CycleTime` is not too long
- Verify the variable value is actually changing

**Too many notifications:**
- Increase `CycleTime` to reduce frequency
- Enable `SendOnChange: true` to filter unchanged values
- Consider if you really need such frequent updates

**Subscriptions lost after restart:**
- This is expected behavior when TwinCAT restarts
- Use `OnConnectionLost` hook to re-subscribe automatically
- See "Handling TwinCAT Restarts" section above

## Device Information

Read information about the target device:

```go
info, err := client.ReadDeviceInfo()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Device Name: %s\n", info.DeviceName)
fmt.Printf("Version: %d.%d (Build %d)\n",
	info.MajorVersion,
	info.MinorVersion,
	info.BuildVersion)

/* Output:
Device Name: PLC-1
Version: 3.1 (Build 4024)
*/
```

## Logging

The client uses structured logging via Go's standard `log/slog` package. By default, logging is disabled.

### Enable Logging

**Text output to console:**

```go
import "log/slog"

logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

client := ads.NewClient(settings, logger)
```

**JSON output:**

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

client := ads.NewClient(settings, logger)
```

**Custom log levels:**

```go
logLevel := &slog.LevelVar{}
logLevel.Set(slog.LevelWarn) // Only warnings and errors

handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: logLevel,
})

logger := slog.New(handler)
client := ads.NewClient(settings, logger)
```

### Disable Logging (default)

```go
client := ads.NewClient(settings, nil) // No logging
```

## Disconnecting

Always disconnect when done to clean up resources:

```go
if err := client.Disconnect(); err != nil {
	log.Printf("Error during disconnect: %v", err)
}
```

**Using defer (recommended):**

```go
func main() {
	client := ads.NewClient(settings, nil)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	// Your code here...
	// Disconnect will be called automatically on exit
}
```

# Common Issues and Questions

## Connection timeouts or failures

**Symptoms:**
- Connection fails immediately
- Timeout errors after 2 minutes
- "Connection refused" errors

**Solutions:**
1. Verify the target PLC is reachable (ping the IP address)
2. Check that TwinCAT is running on the target
3. Verify firewall allows TCP port 48898 (ADS router port)
4. On Windows, ensure Ethernet connection is set to "Private" network
5. For direct connections (Setup 3), verify StaticRoutes.xml is configured correctly
6. Check that the AmsNetId and ADS port are correct

## Symbol not found errors

**Symptoms:**
- Error message: "symbol not found" or similar
- ReadValue/WriteValue fails

**Solutions:**
1. Verify the variable name and path are correct (case-sensitive in TC3, UPPERCASE in TC2)
2. Check that the PLC runtime is in RUN mode (some symbols unavailable in CONFIG)
3. Verify the variable is not optimized away by the compiler
4. For TwinCAT 2, ensure you're using the correct syntax (dot prefix for globals)
5. Try using ReadRaw with GetSymbol to get more details

## Connecting to localhost not working

**Symptoms:**
- Cannot connect when using `TargetAmsNetId: "localhost"` or `"127.0.0.1.1.1"`
- Connection refused on local machine

**Solutions:**
1. Enable TCP loopback in registry (see [Enabling localhost support](#enabling-localhost-support-on-twincat-3))
2. TwinCAT versions < 4024.5 require manual registry edit
3. Restart Windows after changing registry
4. Verify TwinCAT is running locally

## Config mode connections

**Symptoms:**
- Cannot read/write values when PLC is in CONFIG mode
- "Target port not found" errors

**Solutions:**
1. This is expected behavior - most PLC runtime features require RUN mode
2. Use `SetTcSystemToRun()` to start the PLC
3. For system-level operations, you can still read device info and state

## TwinCAT 2 variable names

**Symptoms:**
- Symbol not found when using TwinCAT 2
- Variables not accessible

**Solutions:**
1. All variable names must be UPPERCASE in TwinCAT 2
2. Global variables need dot prefix: `.VARIABLENAME`
3. Use port 801 instead of 851
4. See [Differences when using with TwinCAT 2](#differences-when-using-with-twincat-2)

## Connection from Raspberry Pi or Linux

**Symptoms:**
- Cannot connect from Linux system
- No router available

**Solutions:**
1. Use Setup 3 (direct connection) - see [Setup 3](#setup-3---connect-from-any-system-direct)
2. Configure StaticRoutes.xml on the target PLC
3. Ensure your LocalAmsNetId is unique and not used by other devices
4. No TwinCAT installation needed on the Linux system

## Port already in use

**Symptoms:**
- Error about port being in use
- Cannot start second client

**Solutions:**
1. Use a different `LocalAdsPort` for each client instance
2. Ensure previous client disconnected properly
3. Wait a few seconds for the OS to release the port

# Architecture

The ads-go library uses a modular architecture for maintainability and testing.

## Modular Design

The package is organized into submodules, each handling a specific aspect of the ADS protocol:

| Module | Purpose | Test Coverage |
|--------|---------|---------------|
| **ads-errors** | Parse and validate 4-byte ADS error codes | 100% |
| **ads-header** | Parse 8-byte ADS response headers | 100% |
| **ads-symbol** | Parse ADS symbol information | 100% |
| **ads-datatype** | Parse complex data type definitions | 100% |
| **ads-stateinfo** | Parse system state and device info | 100% |
| **ads-primitives** | Read/write primitive types | 84.8% |
| **ads-requests** | Build ADS command payloads | 100% |
| **ads-serializer** | Type serialization and deserialization | 57.9% |
| **ams-header** | Parse AMS protocol packet headers | 100% |
| **ams-builder** | Build AMS/TCP and AMS headers | 100% |

## Design Patterns

**Invoke ID Management:**
- Each request gets a unique invoke ID
- Responses are matched to requests via invoke ID
- Ensures correct handling of concurrent operations

**Goroutine Receive Loop:**
- Dedicated goroutine for receiving AMS packets
- Channel-based communication with request handlers
- Automatic buffer management and packet reassembly

**Modular Parsing:**
- Each protocol layer has dedicated parser
- Easy to test and maintain
- Clear separation of concerns

# Roadmap

The following features are planned for future releases:

## ✅ ADS Notifications (Subscriptions) - IMPLEMENTED

Event-driven value monitoring:
- ✅ Subscribe to variable value changes
- ✅ Automatic notification handling
- ✅ Multiple simultaneous subscriptions
- ✅ Configurable cycle times and change thresholds

**Status:** ✅ Complete - See [Subscriptions & Notifications](#subscriptions--notifications) section

## Variable Handle Management

Improve performance for repeated reads/writes:
- Create/delete variable handles
- Read/write using handles (faster than by path)
- Automatic handle caching
- Handle lifecycle management

**Status:** Index groups defined, not actively used

## RPC Method Invocation

Call PLC function block methods:
- Invoke FB methods with parameters
- Support for input/output parameters
- Return value handling
- Method metadata parsing

**Status:** Method metadata is parsed, invocation not implemented

## Batch Operations (Sum Commands)

Improve performance for multiple operations:
- Read multiple values in one packet
- Write multiple values in one packet
- Reduced network overhead
- Single round-trip for many operations

**Status:** Index groups defined, not implemented

## Contributing

Contributions, issues, and feature requests are welcome! Please see our [Contributing Guide](CONTRIBUTING.md) for details on how to get started.

**Quick Links:**
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute to this project
- **[Code of Conduct](CODE_OF_CONDUCT.md)** - Our community standards
- **[Security Policy](SECURITY.md)** - How to report security vulnerabilities
- **[Report bugs](https://github.com/jarmocluyse/ads-go/issues)** - GitHub Issues
- **[Suggest features](https://github.com/jarmocluyse/ads-go/discussions)** - GitHub Discussions

# Testing

## Running Tests

Run all tests:
```bash
go test ./pkg/ads/... -v
```

Run tests with coverage:
```bash
go test ./pkg/ads/... -cover
```

Generate coverage report:
```bash
go test ./pkg/ads/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Run specific test:
```bash
go test ./pkg/ads/ads-serializer/... -v -run TestSerialize
```

## Test Structure

The project uses table-driven tests with clear test cases:
- Unit tests for each module
- Integration tests for client operations
- Guard clause style (early returns)
- `github.com/stretchr/testify/assert` for assertions

# Examples

Complete working examples can be found in:

## Command-Line Interface (CLI)

**Location:** `cmd/main.go`

The CLI provides an interactive interface for testing and demonstrating the ads-go library features.

### Running the CLI

```bash
cd cmd
go run main.go
```

Or use the pre-built binary:
```bash
./cmd/ads-cli
```

### CLI Features

**Visual Status Indicators:**
- 🟢 Green prompt = PLC running (operations available)
- 🔵 Blue prompt = PLC in config mode
- 🔴 Red prompt = PLC stopped
- ⚪ White prompt = Disconnected or initializing

**Intelligent Autocomplete:**
- Command completion with TAB key
- Argument suggestions for commands (e.g., `write_bool <TAB>` → `true`, `false`)
- Variable path suggestions for `subscribe` command (14 common paths)
- Dynamic subscription ID completions for `unsubscribe`
- Object field suggestions for `write_object` (Counter=, Ready=)

**Enhanced Subscription Management:**
- Real-time notifications with timestamps
- Subscription statistics (last value, update time, notification count)
- Quick subscription shortcuts for common variables
- Multiple simultaneous subscriptions
- Enhanced list view with detailed information

**Interactive Features:**
- Command history navigation (use arrow keys)
- Auto-reconnection on connection loss
- Automatic state change detection
- Connection lifecycle hooks

### Available Commands

#### System Commands
- `device_info` - Get device information
- `state` - Read current TwinCAT state
- `state_loop` - Continuously monitor TwinCAT state
- `monitor` - Monitor system notifications
- `set_state <config|run>` - Switch TwinCAT state

#### Read/Write Commands
- `read_value` - Read `GLOBAL.gMyInt`
- `read_bool` - Read `GLOBAL.gMyBool`
- `read_object` - Read `GLOBAL.gMyDUT` (struct)
- `read_array` - Read `GLOBAL.gIntArray`
- `write_value <int>` - Write integer to `GLOBAL.gMyInt`
- `write_bool <true|false>` - Write boolean to `GLOBAL.gMyBool`
- `write_object Counter=<int> Ready=<bool>` - Write to `GLOBAL.gMyDUT`
- `write_array <i1> <i2> <i3> <i4> <i5>` - Write 5 ints to `GLOBAL.gIntArray`

#### Subscription Commands
- `subscribe [path]` - Subscribe to variable changes (default: `GLOBAL.gMyBoolToogle`)
- `list_subs` - List active subscriptions with statistics
- `unsubscribe <id>` - Remove specific subscription
- `unsubscribe_all` - Remove all subscriptions

#### Subscription Shortcuts (Quick subscriptions for example project)
- `sub_counter` - Subscribe to cycle-based counter (`GLOBAL.gMyIntCounter`)
- `sub_toggle` - Subscribe to cycle-based toggle (`GLOBAL.gMyBoolToogle`)
- `sub_timed_counter` - Subscribe to time-based counter (`GLOBAL.gTimedIntCounter`)
- `sub_timed_toggle` - Subscribe to time-based toggle (`GLOBAL.gTimedBoolToogle`)
- `sub_all` - Subscribe to all 4 counters/toggles at once

#### Control Commands (For example project)
- `enable_counter <bool>` - Enable/disable cycle-based counter
- `enable_toggle <bool>` - Enable/disable cycle-based toggle
- `enable_timed_counter <bool>` - Enable/disable time-based counter
- `enable_timed_toggle <bool>` - Enable/disable time-based toggle
- `read_counters` - Read all counter and toggle values
- `reset_counters` - Reset all counters to zero
- `read_status` - Show enable flag states
- `set_period <seconds>` - Set cycle period (1-3600s, default 2s)
- `read_period` - Read current cycle period

### Example TwinCAT Project

**Location:** `example/example/`

The CLI works with an included TwinCAT 3 project that demonstrates various features:

**Available Variables:**
- `GLOBAL.gMyInt`, `GLOBAL.gMyBool`, `GLOBAL.gMyDINT` - Basic types for testing
- `GLOBAL.gMyIntCounter` - Counter increments every PLC scan
- `GLOBAL.gMyBoolToogle` - Boolean toggles every PLC scan
- `GLOBAL.gTimedIntCounter` - Counter increments every cycle period (default 2s)
- `GLOBAL.gTimedBoolToogle` - Boolean toggles every cycle period
- `GLOBAL.gIntArray` - Array of 101 integers (CLI writes to first 5)
- `GLOBAL.gMyDUT` - Structured data (Counter: INT, Ready: BOOL, gIntArray: ARRAY[0..50] OF INT)
- `GLOBAL.gCyclePeriod` - Configurable timer period (TIME type, default T#2S)

**Control Flags:**
- `GLOBAL.gIntCounterActive` - Enable/disable cycle-based counter (default TRUE)
- `GLOBAL.gBoolToggleActive` - Enable/disable cycle-based toggle (default TRUE)
- `GLOBAL.gTimedCounterActive` - Enable/disable time-based counter (default TRUE)
- `GLOBAL.gTimedToggleActive` - Enable/disable time-based toggle (default TRUE)

The project includes:
- Cycle-based logic that runs every PLC scan
- Time-based logic triggered by configurable timer
- All variables accessible via ADS for read/write operations
- Perfect for testing subscriptions and real-time updates

### Example Usage

**Quick start with subscriptions:**
```bash
# Start the CLI
./ads-cli

# Subscribe to a fast-changing counter
sub_counter

# Subscribe to all counters/toggles at once
sub_all

# View subscription statistics with last values
list_subs

# Disable the cycle-based toggle
enable_toggle false

# Change timer period to 5 seconds
set_period 5

# Remove a specific subscription
unsubscribe 1

# Remove all subscriptions
unsubscribe_all
```

**Testing variable operations:**
```bash
# Write a value
write_value 42

# Read it back
read_value

# Write a boolean
write_bool true

# Write a structured object
write_object Counter=100 Ready=true

# Write an array (first 5 elements)
write_array 10 20 30 40 50

# Read the array back
read_array
```

**Monitoring system state:**
```bash
# Check current state
state

# Monitor state continuously (Ctrl+C to stop)
state_loop

# Switch to config mode
set_state config

# Return to run mode
set_state run

# Get device information
device_info
```

**Using autocomplete:**
```bash
# Type 'sub_' and press TAB to see subscription shortcuts
sub_<TAB>

# Type 'write_bool ' and press TAB to see options
write_bool <TAB>

# Type 'subscribe ' and press TAB to see common variable paths
subscribe <TAB>

# Type 'unsubscribe ' and press TAB to see active subscription IDs
unsubscribe <TAB>
```

## Additional Examples

- **example/** - TwinCAT 3 example project with PLC program

# License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

# Credits

- **Author:** Jarmo Cluyse (jarmo_cluyse@hotmail.com)
- **GitHub:** https://github.com/jarmocluyse/ads-go
- **Documentation inspired by:** [jisotalo/ads-client](https://github.com/jisotalo/ads-client) by Jussi Isotalo (used with permission)

---

**Made with ❤️ for the Beckhoff automation community**
