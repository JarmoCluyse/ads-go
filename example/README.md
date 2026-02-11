# Example TwinCAT Project

This directory contains a complete TwinCAT 3 project that demonstrates the capabilities of the ads-go library and CLI tool.

## Project Structure

```
example/
├── example.sln              # Visual Studio solution file
└── example/
    ├── example.tsproj       # TwinCAT project file
    └── smallproject/        # PLC project
        ├── POUs/
        │   └── MAIN.TcPOU   # Main program with cycle-based and time-based logic
        ├── DUTs/
        │   └── DUTSample.TcDUT  # Sample structured data type
        └── GVLs/
            └── GLOBAL.TcGVL # Global variable list with test variables
```

## Opening the Project

1. Open TwinCAT XAE (Extended Automation Engineering)
2. Open `example/example.sln`
3. Build the project: **TwinCAT → Activate Configuration**
4. Set TwinCAT to RUN mode

**Note:** The `_Boot/` directory (containing compiled boot files) is excluded from Git per TwinCAT best practices. It will be automatically generated when you activate the configuration in step 3.

## Available Variables

The project exposes the following global variables for ADS access:

### Basic Types (For Testing)
- `GLOBAL.gMyInt: INT` - Simple integer value (default: 0)
- `GLOBAL.gMyBool: BOOL` - Simple boolean value (default: FALSE)
- `GLOBAL.gMyDINT: DINT` - Simple double integer value (default: 0)

### Cycle-Based Variables (Updates Every PLC Scan)
- `GLOBAL.gMyIntCounter: INT` - Increments every PLC cycle when enabled (default: 0)
- `GLOBAL.gMyBoolToogle: BOOL` - Toggles every PLC cycle when enabled (default: FALSE)

### Time-Based Variables (Updates Every Cycle Period)
- `GLOBAL.gTimedIntCounter: INT` - Increments every cycle period when enabled (default: 0)
- `GLOBAL.gTimedBoolToogle: BOOL` - Toggles every cycle period when enabled (default: FALSE)

### Control Flags
- `GLOBAL.gIntCounterActive: BOOL` - Enable/disable cycle-based counter (default: TRUE)
- `GLOBAL.gBoolToggleActive: BOOL` - Enable/disable cycle-based toggle (default: TRUE)
- `GLOBAL.gTimedCounterActive: BOOL` - Enable/disable time-based counter (default: TRUE)
- `GLOBAL.gTimedToggleActive: BOOL` - Enable/disable time-based toggle (default: TRUE)

### Configuration
- `GLOBAL.gCyclePeriod: TIME` - Period for time-based operations (default: T#2S = 2 seconds)

### Complex Types
- `GLOBAL.gIntArray: ARRAY[0..100] OF INT` - Array of 101 integers
- `GLOBAL.gMyDUT: DUTSample` - Structured data type with:
  - `Counter: INT` - Integer counter field
  - `Ready: BOOL` - Boolean ready flag
  - `gIntArray: ARRAY[0..50] OF INT` - Nested integer array

## PLC Program Logic

The `MAIN` program implements two types of behavior:

### 1. Cycle-Based Behavior (Every PLC Scan)
Executes on every PLC cycle (typically 1-10ms):

```iecst
// Increment counter if enabled
IF GLOBAL.gIntCounterActive THEN
    GLOBAL.gMyIntCounter := GLOBAL.gMyIntCounter + 1;
END_IF

// Toggle boolean if enabled
IF GLOBAL.gBoolToggleActive THEN
    GLOBAL.gMyBoolToogle := NOT GLOBAL.gMyBoolToogle;
END_IF
```

**Use case:** Fast-changing values for testing rapid subscriptions and notifications.

### 2. Time-Based Behavior (Every Cycle Period)
Executes periodically based on `gCyclePeriod` (default 2 seconds):

```iecst
// Timer triggers periodic actions
fbTimer(IN := TRUE, PT := GLOBAL.gCyclePeriod);

IF fbTimer.Q THEN
    // Increment timed counter if enabled
    IF GLOBAL.gTimedCounterActive THEN
        GLOBAL.gTimedIntCounter := GLOBAL.gTimedIntCounter + 1;
    END_IF
    
    // Toggle timed boolean if enabled
    IF GLOBAL.gTimedToggleActive THEN
        GLOBAL.gTimedBoolToogle := NOT GLOBAL.gTimedBoolToogle;
    END_IF
    
    // Retrigger timer for next cycle
    fbTimer(IN := FALSE);
END_IF
```

**Use case:** Slower-changing values for testing different subscription cycle times and demonstrating configurable periods.

## Using with ads-go CLI

The CLI tool is pre-configured to work with these variables. See the main [README.md](../README.md#example-twincat-project) for CLI commands.

### Quick Start Examples

**Monitor fast-changing counter:**
```bash
./ads-cli
> sub_counter
# Watch counter increment every PLC cycle
```

**Control the toggles:**
```bash
> enable_toggle false     # Disable cycle-based toggle
> enable_timed_toggle true  # Enable time-based toggle
> read_status             # Check current states
```

**Change timer period:**
```bash
> read_period             # Check current period (default 2s)
> set_period 5            # Change to 5 seconds
> sub_timed_counter       # Watch it increment every 5s
```

**Test all subscriptions:**
```bash
> sub_all                 # Subscribe to all 4 counters/toggles
> list_subs               # View statistics
> unsubscribe_all         # Clean up
```

## Variable Access Paths

When using ads-go library programmatically:

```go
// Read basic values
value, _ := client.ReadValue(851, "GLOBAL.gMyInt")
value, _ := client.ReadValue(851, "GLOBAL.gMyBool")

// Read counters
counter, _ := client.ReadValue(851, "GLOBAL.gMyIntCounter")
timedCounter, _ := client.ReadValue(851, "GLOBAL.gTimedIntCounter")

// Write control flags
client.WriteValue(851, "GLOBAL.gIntCounterActive", true)
client.WriteValue(851, "GLOBAL.gCyclePeriod", int32(5000)) // 5 seconds in ms

// Read array
array, _ := client.ReadValue(851, "GLOBAL.gIntArray")

// Read/write struct
dut, _ := client.ReadValue(851, "GLOBAL.gMyDUT")
client.WriteValue(851, "GLOBAL.gMyDUT", map[string]any{
    "Counter": int16(100),
    "Ready":   true,
})

// Subscribe to changes
callback := func(data ads.SubscriptionData) {
    fmt.Printf("Value changed: %v\n", data.Value)
}
settings := ads.SubscriptionSettings{
    CycleTime:    100 * time.Millisecond,
    SendOnChange: true,
}
sub, _ := client.SubscribeValue(851, "GLOBAL.gMyBoolToogle", callback, settings)
```

## Requirements

- TwinCAT 3 XAE or XAR (Build 4024 or newer recommended)
- Windows OS with TwinCAT runtime
- ADS route configured between client and PLC (see main README for setup instructions)

## Port Configuration

- **Default PLC Port:** 851 (TwinCAT 3 first runtime)
- For TwinCAT 2, use port 801

## Troubleshooting

**Variables not found:**
- Ensure TwinCAT is in RUN mode (not CONFIG)
- Verify the PLC program is activated
- Check ADS route is configured correctly

**Counters not incrementing:**
- Check enable flags: `read_status` in CLI
- Verify PLC is running: `state` in CLI
- For timed counters, verify cycle period: `read_period` in CLI

**Subscription notifications not received:**
- Ensure PLC is in RUN mode
- Verify the variable is actually changing (check enable flags)
- Check subscription settings (cycle time, on-change mode)

## License

This example project is part of ads-go and is licensed under the MIT License.
