package ads

const (
	AMSTCPHeaderLength   = 6               // AMS/TCP header length
	AMSHeaderLength      = 32              // AMS header length
	AMSNetIDLength       = 6               // AmsNetId length
	ADSIndexOffsetLength = 4               // ADS index offset length
	ADSIndexGroupLength  = 4               // ADS index group length
	ADSInvokeIDMaxValue  = 4294967295      // ADS invoke ID maximum value (32bit unsigned integer)
	ADSDefaultTCPPort    = 48898           // Default ADS server TCP port for incoming connections
	LoopbackAmsNetID     = "127.0.0.1.1.1" // Loopback (localhost) AmsNetId
)
