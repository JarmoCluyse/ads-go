// Package amsbuilder provides functions for building AMS/TCP and AMS protocol headers.
//
// This package is symmetric with the ams-header package: while ams-header parses
// incoming AMS packets, ams-builder constructs outgoing AMS packets. Together, they
// handle the complete AMS protocol layer for ADS communication.
//
// # AMS/TCP Protocol
//
// AMS/TCP is the transport layer protocol that wraps AMS packets. It provides a
// simple 6-byte header containing a command flag and the length of the following data.
//
// Binary format (6 bytes, little-endian):
//
//	Bytes 0-1: Command/Flag (uint16) - Type of AMS/TCP packet
//	Bytes 2-5: Data Length (uint32) - Length of following data
//
// Common AMS/TCP commands:
//   - 0x0000 (AMSTCPPortAMSCommand) - Standard ADS command packet
//   - 0x0001 (AMSTCPPortClose) - Close port
//   - 0x1000 (AMSTCPPortConnect) - Port registration
//   - 0x1001 (AMSTCPPortRouterNote) - Router notification
//   - 0x1002 (GetLocalNetID) - Request local NetID
//
// # AMS Protocol
//
// The AMS (Automation Message Specification) header provides routing information
// and command details for ADS communication.
//
// Binary format (32 bytes, little-endian):
//
//	Bytes 0-5:   Target NetID (6 bytes)
//	Bytes 6-7:   Target Port (uint16)
//	Bytes 8-13:  Source NetID (6 bytes)
//	Bytes 14-15: Source Port (uint16)
//	Bytes 16-17: ADS Command (uint16)
//	Bytes 18-19: State Flags (uint16) - Always 4 (ADSStateFlagAdsCommand) for requests
//	Bytes 20-23: Data Length (uint32) - Length of ADS command payload only
//	Bytes 24-27: Error Code (uint32) - Always 0 when sending requests
//	Bytes 28-31: Invoke ID (uint32) - Unique request identifier
//
// # NetID Format
//
// AMS NetIDs identify ADS devices on the network. They consist of six decimal
// numbers separated by dots, similar to IPv4 addresses but with two extra octets.
//
// Format: "a.b.c.d.e.f" where each value is 0-255
//
// Common NetID patterns:
//   - PLC: "192.168.1.100.1.1" (IP address + ".1.1")
//   - PC: "192.168.1.50.1.1" (IP address + ".1.1")
//   - Local: "127.0.0.1.1.1"
//
// # AMS Ports
//
// AMS ports identify services on an ADS device:
//   - 851: TwinCAT PLC Runtime (most common)
//   - 801: TwinCAT System Service
//   - 10000: First TwinCAT 2 Runtime
//   - 32768-65535: Dynamic client ports
//
// # Usage Example
//
// Building a complete ADS Read request packet:
//
//	import (
//	    amsbuilder "github.com/jarmocluyse/ads-go/pkg/ads/ams-builder"
//	    adsrequests "github.com/jarmocluyse/ads-go/pkg/ads/ads-requests"
//	    "github.com/jarmocluyse/ads-go/pkg/ads/types"
//	)
//
//	// Define target and source addresses
//	target := amsbuilder.AmsAddress{
//	    NetID: "192.168.1.100.1.1",  // PLC address
//	    Port:  851,                    // PLC runtime
//	}
//	source := amsbuilder.AmsAddress{
//	    NetID: "192.168.1.50.1.1",   // PC address
//	    Port:  32905,                  // Dynamic client port
//	}
//
//	// Build ADS Read request payload
//	payload := adsrequests.BuildReadRequest(
//	    0x4020,  // Index group
//	    0x0,     // Index offset
//	    4,       // Read 4 bytes
//	)
//
//	// Build AMS header
//	amsHeader, err := amsbuilder.BuildAmsHeader(
//	    target,
//	    source,
//	    types.ADSCommandRead,
//	    uint32(len(payload)),
//	    12345,  // Invoke ID
//	)
//	if err != nil {
//	    // Handle error (invalid NetID)
//	    return err
//	}
//
//	// Build AMS/TCP header
//	tcpHeader := amsbuilder.BuildAmsTcpHeader(
//	    types.AMSTCPPortAMSCommand,
//	    uint32(len(amsHeader) + len(payload)),
//	)
//
//	// Combine into complete packet
//	packet := append(tcpHeader, amsHeader...)
//	packet = append(packet, payload...)
//
//	// Send packet over TCP connection...
//
// # Port Registration Example
//
// Registering a dynamic client port with the AMS router:
//
//	import (
//	    "encoding/binary"
//	    amsbuilder "github.com/jarmocluyse/ads-go/pkg/ads/ams-builder"
//	    "github.com/jarmocluyse/ads-go/pkg/ads/types"
//	    "github.com/jarmocluyse/ads-go/pkg/ads/utils"
//	)
//
//	// Build registration payload (NetID + Port)
//	payload := make([]byte, 8)
//	netID, _ := utils.AmsNetIdStrToByteArray("192.168.1.50.1.1")
//	copy(payload[0:6], netID)
//	binary.LittleEndian.PutUint16(payload[6:8], 32905)
//
//	// Build AMS/TCP header for port registration
//	tcpHeader := amsbuilder.BuildAmsTcpHeader(
//	    types.AMSTCPPortConnect,
//	    uint32(len(payload)),
//	)
//
//	// Combine and send
//	packet := append(tcpHeader, payload...)
//	// Send packet...
//
// # Length Calculations
//
// When building packets, pay attention to length fields:
//
//   - AMS/TCP Data Length: Length of AMS header + ADS payload
//     (typically 32 + payload length for standard commands)
//
//   - AMS Data Length: Length of ADS command payload only
//     (does NOT include AMS header)
//
// Example for Read command with 12-byte payload:
//   - AMS/TCP Data Length: 32 + 12 = 44 bytes
//   - AMS Data Length: 12 bytes
//
// # Error Handling
//
// The only function that returns an error is BuildAmsHeader, which can fail if
// the target or source NetID is invalid. NetID validation checks:
//   - Must have exactly 6 parts separated by dots
//   - Each part must be a valid integer (parsing errors are caught)
//
// BuildAmsTcpHeader never returns an error and accepts any uint16/uint32 values.
//
// # Integration with Other Modules
//
// This package is designed to work with:
//   - ads-requests: Builds ADS command payloads
//   - ams-header: Parses response packets (symmetric operation)
//   - Client code: Combines headers and payloads, manages TCP connections
//
// # Protocol References
//
// For complete protocol specification, see:
//   - Beckhoff ADS Protocol Documentation
//   - TwinCAT AMS/TCP Specification
//   - https://infosys.beckhoff.com/
package amsbuilder
