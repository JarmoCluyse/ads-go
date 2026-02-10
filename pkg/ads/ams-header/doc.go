// Package amsheader provides parsing functionality for AMS (Automation Message Specification) packet headers.
//
// The AMS protocol is the transport layer used by the ADS (Automation Device Specification) protocol
// in TwinCAT systems. Each AMS packet consists of an AMS/TCP header, an AMS header, and a data payload.
//
// # Binary Format
//
// An AMS packet has the following structure:
//
//	AMS/TCP Header (6 bytes):
//	  Offset  Size  Field
//	  ------  ----  -----
//	  0       2     Reserved (uint16)
//	  2       4     Packet length (uint32) - excludes the AMS/TCP header itself
//
//	AMS Header (32 bytes):
//	  Offset  Size  Field
//	  ------  ----  -----
//	  0       6     Target AMS Net ID (6 bytes, e.g., 192.168.1.100.1.1)
//	  6       2     Target Port (uint16)
//	  8       6     Source AMS Net ID (6 bytes)
//	  14      2     Source Port (uint16)
//	  16      2     ADS Command (uint16)
//	  18      2     State Flags (uint16)
//	  20      4     Data Length (uint32)
//	  24      4     Error Code (uint32) - AMS header error, not ADS error
//	  28      4     Invoke ID (uint32) - for matching requests/responses
//
//	Data Payload (variable length):
//	  32      N     ADS data
//
// # Total Packet Length
//
// The total packet length is: 6 (AMS/TCP header) + packet_length field.
// The minimum packet length is 38 bytes (6 + 32 with no data).
//
// # Usage
//
// Parse a complete AMS packet:
//
//	packet, err := amsheader.ParsePacket(data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("From: %s:%d\n", packet.SourceNetID, packet.SourcePort)
//	fmt.Printf("Command: %v, InvokeID: %d\n", packet.Command, packet.InvokeID)
//
// Check if enough data is available (useful for buffering):
//
//	totalLength, err := amsheader.CheckTCPPacketLength(buffer)
//	if err != nil {
//	    // Not enough data yet, wait for more
//	    return
//	}
//	// We have a complete packet of totalLength bytes
//
// Validate packet structure before parsing:
//
//	if err := amsheader.CheckPacket(data); err != nil {
//	    log.Fatal("Invalid packet")
//	}
//
// # Error Handling
//
// The package defines sentinel errors that can be checked with errors.Is():
//
//   - ErrInsufficientData: Not enough data to parse the packet
//   - ErrInvalidData: The data is malformed or contains invalid values
//   - ErrInvalidLength: The packet length field is invalid
//
// Example:
//
//	packet, err := amsheader.ParsePacket(data)
//	if errors.Is(err, amsheader.ErrInsufficientData) {
//	    // Need to read more data from the connection
//	}
package amsheader
