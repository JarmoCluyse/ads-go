package amsheader

import "github.com/jarmocluyse/ads-go/pkg/ads/types"

// Packet represents a complete AMS packet with header and data.
type Packet struct {
	TargetNetID string              // Target AMS Net ID (e.g., "192.168.1.100.1.1")
	TargetPort  uint16              // Target AMS port
	SourceNetID string              // Source AMS Net ID
	SourcePort  uint16              // Source AMS port
	Command     types.ADSCommand    // ADS command
	StateFlags  types.ADSStateFlags // State flags
	DataLength  uint32              // Length of data payload
	ErrorCode   uint32              // AMS header error code (not ADS payload error)
	InvokeID    uint32              // Invoke ID for request/response matching
	Data        []byte              // Data payload
}

// Address represents an AMS address (NetID + Port).
type Address struct {
	NetID string // AMS Net ID (e.g., "192.168.1.100.1.1")
	Port  uint16 // AMS Port
}
