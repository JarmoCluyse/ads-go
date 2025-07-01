package adsgo

import "net"

// AdsClientConnection represents an active client connection.
type AdsClientConnection struct {
	// Connection status of the client, true if connected
	Connected bool
	// True if connected to local TwinCAT system (loopback)
	IsLocal *bool
	// Local AmsNetId of the client
	LocalAmsNetId *string
	// Local ADS port of the client
	LocalAdsPort *int
	// Target AmsNetId
	TargetAmsNetId *string
	// Target ADS port
	TargetAdsPort *int
	// Active Dial connection
	Connection net.Conn
}
