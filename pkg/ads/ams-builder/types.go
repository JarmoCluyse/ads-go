package amsbuilder

// AmsAddress represents an AMS address consisting of a NetID and port.
type AmsAddress struct {
	NetID string // AMS Net ID (e.g., "192.168.1.100.1.1")
	Port  uint16 // AMS port number
}
