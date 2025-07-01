package adsgo

// AdsClientSettings defines the configuration settings for the ADS client.
type AdsClientSettings struct {
	// Default target AmsNetId address (REQUIRED)
	// Examples:
	//   - "localhost" or "127.0.0.1.1.1": Local (same machine)
	//   - "192.168.1.5.1.1": PLC example
	//   - "192.168.1.5.2.1": EtherCAT I/O device example
	TargetAmsNetId string

	// Default target ADS port (REQUIRED)
	// Examples:
	//   - 851: TwinCAT 3 PLC runtime 1
	//   - 852: TwinCAT 3 PLC runtime 2
	//   - 801: TwinCAT 2 PLC runtime 1
	//   - 10000: TwinCAT system service
	// NOTE: Not a TCP port.
	TargetAdsPort int

	// Target ADS router TCP port (Optional, default: 48898)
	// Usually 48898, unless using port forwarding or separate router.
	// NOTE: Firewall must allow this port.
	RouterTcpPort *int

	// Target ADS router IP address/hostname (Optional, default: 127.0.0.1)
	// Usually local machine, change if connecting from a system without router.
	RouterAddress *string

	// Local IP address to use (Optional, default: "")
	// Can be used to force using another network interface.
	LocalAddress *string

	// Local TCP port to use for outgoing connection (Optional, default: 0)
	// Can be used to force using specific local TCP port.
	LocalTcpPort *int

	// Local AmsNetId to use (Optional, default: "")
	// Can be set manually or received from target router.
	LocalAmsNetId *string

	// Local ADS port to use (Optional, default: 0)
	// Can be set manually or received from target router.
	LocalAdsPort *int

	// Time (ms) after a command is timeouted if no response received (Optional, default: 2000)
	TimeoutDelay *int

	// If set, the client tries to reconnect automatically after a connection loss (Optional, default: true)
	AutoReconnect *bool

	// Interval (ms) how often the lost connection is tried to re-establish (Optional, default: 2000)
	ReconnectInterval *int

	// If set, ENUM data types are converted to objects (Optional, default: true)
	ObjectifyEnumerations *bool

	// If set, PLC date types are converted to Go time.Time (Optional, default: true)
	ConvertDatesToTime *bool

	// If set, all symbols from target are read and cached after connecting (Optional, default: false)
	ReadAndCacheSymbols *bool

	// If set, all data types from target are read and cached after connecting (Optional, default: false)
	ReadAndCacheDataTypes *bool

	// If set, client detects PLC symbol version changes and reloads symbols and data types (Optional, default: true)
	MonitorPlcSymbolVersion *bool

	// If set, no warnings are written to console (Optional, default: false)
	HideConsoleWarnings *bool

	// Interval (ms) how often the client checks if the connection is working (Optional, default: 1000)
	ConnectionCheckInterval *int

	// Time (ms) how long after the target TwinCAT system state is not available, the connection is determined to be lost (Optional, default: 5000)
	ConnectionDownDelay *int

	// If set, connecting to the target will succeed, even if there is no PLC runtime or system is in CONFIG mode (Optional, default: false)
	AllowHalfOpen *bool

	// If set, only a direct raw ADS connection is established (Optional, default: false)
	RawClient *bool

	// If set, the client never caches symbols and data types (Optional, default: false)
	DisableCaching *bool

	// If set, the client automatically deletes ADS notifications for unknown subscriptions (Optional, default: true)
	DeleteUnknownSubscriptions *bool

	// If set, the client always uses UTF-8 for encoding/decoding ADS symbols (Optional, default: false)
	ForceUtf8ForAdsSymbols *bool
}
