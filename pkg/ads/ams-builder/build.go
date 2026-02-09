package amsbuilder

import (
	"encoding/binary"

	"github.com/jarmocluyse/ads-go/pkg/ads/constants"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/jarmocluyse/ads-go/pkg/ads/utils"
)

// BuildAmsTcpHeader builds the AMS/TCP header (6 bytes).
//
// The AMS/TCP header is the transport layer header that wraps AMS packets.
//
// Binary format (6 bytes, little-endian):
//
//	Bytes 0-1: Command/Flag (uint16) - Type of AMS/TCP packet
//	Bytes 2-5: Data Length (uint32) - Length of following data (AMS header + payload)
//
// Common commands:
//   - 0x0000 (AMSTCPPortAMSCommand) - Standard ADS command
//   - 0x0001 (AMSTCPPortClose) - Close port
//   - 0x1000 (AMSTCPPortConnect) - Port registration
//   - 0x1001 (AMSTCPPortRouterNote) - Router notification
//   - 0x1002 (GetLocalNetID) - Request local NetID
//
// The dataLength parameter should be the combined length of the AMS header
// (32 bytes) plus the ADS command payload.
//
// Example:
//
//	// For a Read command with 12-byte payload:
//	// AMS header = 32 bytes, payload = 12 bytes, total = 44 bytes
//	header := BuildAmsTcpHeader(types.AMSTCPPortAMSCommand, 44)
func BuildAmsTcpHeader(command types.AMSHeaderFlag, dataLength uint32) []byte {
	buf := make([]byte, constants.AMSTCPHeaderLength)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(command)) // AMS command
	binary.LittleEndian.PutUint32(buf[2:6], dataLength)      // Length of the data
	return buf
}

// BuildAmsHeader builds the AMS header (32 bytes).
//
// The AMS header provides routing information and command details for ADS communication.
//
// Binary format (32 bytes, little-endian):
//
//	Bytes 0-5:   Target NetID (6 bytes)
//	Bytes 6-7:   Target Port (uint16)
//	Bytes 8-13:  Source NetID (6 bytes)
//	Bytes 14-15: Source Port (uint16)
//	Bytes 16-17: ADS Command (uint16)
//	Bytes 18-19: State Flags (uint16) - Always ADSStateFlagAdsCommand (4) for requests
//	Bytes 20-23: Data Length (uint32) - Length of ADS command payload only
//	Bytes 24-27: Error Code (uint32) - Always 0 when sending requests
//	Bytes 28-31: Invoke ID (uint32) - Unique request identifier
//
// Parameters:
//   - target: Destination AMS address (NetID and port)
//   - source: Source AMS address (NetID and port)
//   - command: ADS command to execute (Read, Write, ReadWrite, etc.)
//   - dataLength: Length of the ADS command payload (not including headers)
//   - invokeID: Unique request identifier for matching responses
//
// The NetID is converted from string format (e.g., "192.168.1.100.1.1") to
// 6-byte binary format.
//
// Example:
//
//	target := AmsAddress{NetID: "192.168.1.100.1.1", Port: 851}
//	source := AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905}
//	header, err := BuildAmsHeader(target, source, types.ADSCommandRead, 12, 1)
func BuildAmsHeader(target AmsAddress, source AmsAddress, command types.ADSCommand, dataLength uint32, invokeID uint32) ([]byte, error) {
	buf := make([]byte, constants.AMSHeaderLength)

	// Target NetID and port
	targetNetID, err := utils.AmsNetIdStrToByteArray(target.NetID)
	if err != nil {
		return nil, err
	}
	copy(buf[0:6], targetNetID)
	binary.LittleEndian.PutUint16(buf[6:8], target.Port)

	// Source NetID and port
	sourceNetID, err := utils.AmsNetIdStrToByteArray(source.NetID)
	if err != nil {
		return nil, err
	}
	copy(buf[8:14], sourceNetID)
	binary.LittleEndian.PutUint16(buf[14:16], source.Port)

	// Command and flags
	binary.LittleEndian.PutUint16(buf[16:18], uint16(command))                      // ADS command
	binary.LittleEndian.PutUint16(buf[18:20], uint16(types.ADSStateFlagAdsCommand)) // State flags (always 4 for commands)
	binary.LittleEndian.PutUint32(buf[20:24], dataLength)                           // Payload length
	binary.LittleEndian.PutUint32(buf[24:28], 0)                                    // Error code (0 when sending)
	binary.LittleEndian.PutUint32(buf[28:32], invokeID)                             // Invoke ID

	return buf, nil
}
