# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **CLI Enhancements**: Major improvements to the command-line interface
  - **Intelligent Autocomplete System**:
    - Nested command completion with TAB key
    - Argument suggestions for boolean commands (`true`/`false`)
    - Variable path suggestions for `subscribe` command (14 common paths from example project)
    - Dynamic subscription ID completions for `unsubscribe` command
    - Object field suggestions for `write_object` (Counter=, Ready=)
  - **Subscription Shortcut Commands**: Quick access to example project variables
    - `sub_counter` - Subscribe to cycle-based counter (`GLOBAL.gMyIntCounter`)
    - `sub_toggle` - Subscribe to cycle-based toggle (`GLOBAL.gMyBoolToogle`)
    - `sub_timed_counter` - Subscribe to time-based counter (`GLOBAL.gTimedIntCounter`)
    - `sub_timed_toggle` - Subscribe to time-based toggle (`GLOBAL.gTimedBoolToogle`)
    - `sub_all` - Subscribe to all 4 counters/toggles simultaneously
  - **Control Commands**: Interactive control for example project toggles and timers
    - `enable_counter`, `enable_toggle`, `enable_timed_counter`, `enable_timed_toggle` - Control enable flags
    - `read_counters` - Read all counter and toggle values at once
    - `reset_counters` - Reset all counters and toggles to initial values
    - `read_status` - Display status of all enable flags
    - `set_period <seconds>` - Configure cycle period for timed operations (1-3600s)
    - `read_period` - Display current cycle period
  - **Enhanced Subscription Display**:
    - `list_subs` now shows last received value, time since last update, and notification count
    - Statistics tracking for each subscription (value, timestamp, count)
    - Human-readable time formatting (ms/s/m/h)
    - Multi-line display format for better readability
- **Example TwinCAT Project**: Added complete working TwinCAT 3 project
  - Location: `example/example/`
  - Includes cycle-based and time-based counters and toggles
  - Configurable timer period via ADS
  - Enable/disable flags for all automatic behaviors
  - Demonstrates real-time subscriptions and control operations
  - 14 global variables available for testing (basic types, arrays, structs)

### Changed
- Updated CLI read/write commands to use example project variables
  - `read_value`/`write_value` now use `GLOBAL.gMyInt`
  - `read_bool`/`write_bool` now use `GLOBAL.gMyBool`
  - `read_object`/`write_object` now use `GLOBAL.gMyDUT` (simplified to Counter/Ready fields)
  - `read_array`/`write_array` now use `GLOBAL.gIntArray` (changed from 10 to 5 elements for usability)
  - `subscribe` default path changed to `GLOBAL.gMyBoolToogle`
- **Renamed state commands** for consistency: `toConfig`/`toRun` → `set_state config`/`set_state run`
- Enhanced CLI help command with detailed command descriptions and usage examples
- Improved subscription callback to track statistics automatically

### Fixed
- Removed redundant newline in help command output (go vet warning)

## [0.2.0] - 2026-02-10

### Added
- **ADS Notifications & Subscriptions**: Full support for ADS device notifications with automatic change detection
  - `SubscribeValue()` - Subscribe to variable changes with configurable transmission modes (OnChange, Cyclic, CyclicInContext)
  - `Unsubscribe()` - Remove individual subscriptions
  - `UnsubscribeAll()` - Remove all active subscriptions
  - Automatic subscription lifecycle management with TwinCAT restart detection and re-registration
- **State Monitoring & Event Handling**: Real-time PLC state tracking with connection lifecycle hooks
  - `OnConnect()` - Callback when ADS connection is established
  - `OnDisconnect()` - Callback when client disconnects cleanly
  - `OnConnectionLost()` - Callback when connection drops unexpectedly
  - `MonitorState()` - Continuous ADS state monitoring with configurable interval
  - `GetState()` - Query current PLC state (Run, Stop, Config, etc.)
  - Extended state support for TwinCAT 4022+ with graceful fallback for older versions
- **System Information Commands**: New CLI commands for PLC diagnostics
  - `system state` - Display current PLC state
  - `system version` - Display TwinCAT version information
  - `subscribe` - Subscribe to variable changes via CLI
  - `list_subs` - List all active subscriptions
  - `unsubscribe` - Remove subscriptions by ID
  - `unsubscribe_all` - Clear all subscriptions
- **Enhanced State Parsing**: Comprehensive PLC state interpretation
  - Parse ADS state codes with detailed descriptions
  - Support for TwinCAT 2/3 state flags (AFSERVICELOAD, AFADSIOCREATION, AFEVENTPROCESSING, etc.)
  - System Service state parsing (SystemServiceRunning, ServiceInternalError, etc.)
  - Router state parsing with detailed descriptions
  - Extensive test coverage for state parsing functions
- **Documentation Improvements**: 
  - Added 500+ lines of subscriptions documentation with working examples
  - Updated README with subscription lifecycle management guide
  - Added troubleshooting guide for common scenarios
  - Updated project status to reflect completed features

### Changed
- Enhanced ADS client connection handling with automatic reconnection support
- Improved logging throughout the client with structured logging
- Updated CLI to enable logger output for better debugging experience

### Fixed
- Import casing issues in subscription types
- Subscription command registration in CLI handlers
- Buffer handling for ADS notification messages
- Type conversions for subscription-related ADS commands
- Code formatting issues identified by Go Report Card (gofmt)

## [0.1.4] - 2026-02-10

### Fixed
- Default configuration loading functionality

## [0.1.3] - 2026-02-09

### Fixed
- Import casing issues in package imports

## [0.1.2] - 2026-02-09

### Fixed
- Module casing issue in go.mod

## [0.1.1] - 2026-02-09

### Fixed
- Casing in go.mod file

## [0.1.0] - 2026-02-09

### Added
- Initial release of ads-go
- Core ADS protocol client implementation
- Read operations for all basic types (bool, integers, floats, strings)
- Write operations for all basic types
- Symbol resolution by name
- Connection management to TwinCAT AMS Router
- CLI tool for testing ADS operations
- Basic documentation and examples

### Features
- Read/Write BOOL, BYTE, WORD, DWORD, LWORD
- Read/Write INT, UINT, SINT, USINT, DINT, UDINT, LINT, ULINT
- Read/Write REAL, LREAL
- Read/Write STRING, WSTRING
- Symbol lookup by name
- Raw read/write with index group/offset
- Port-independent implementation (works on Linux/macOS/Windows)
- Configurable AMS NetID and port
- Integration with TwinCAT AMS Router
