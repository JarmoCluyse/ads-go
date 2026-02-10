# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
