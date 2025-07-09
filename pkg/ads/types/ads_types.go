package types

// AmsAddress represents an AMS address with NetId and port.
type AmsAddress struct {
	NetID string
	Port  uint16
}

// AmsTcpHeader represents the AMS/TCP header.
type AmsTcpHeader struct {
	Command uint16
	Length  uint32
}

// AmsHeader represents the AMS header.
type AmsHeader struct {
	TargetAmsAddress AmsAddress
	SourceAmsAddress AmsAddress
	Command          uint16 // ADS command
	StateFlags       uint16
	DataLength       uint32
	ErrorCode        uint32
	InvokeID         uint32
}

// AmsTcpPacket represents a full AMS/TCP packet.
type AmsTcpPacket struct {
	TcpHeader AmsTcpHeader
	AmsHeader AmsHeader
	ADSData   []byte // Payload
}

// AdsCommandRequest represents a request for an ADS command.
type AdsCommandRequest struct {
	Command     ADSCommand
	Data        []byte
	TargetNetID string
	TargetPort  uint16
}

// AdsStateResponse represents the ADS state response payload.
type AdsStateResponse struct {
	ADSState    ADSState
	DeviceState uint16
}

// BaseAdsResponse represents the common fields for all ADS responses.
type BaseAdsResponse struct {
	ErrorCode uint32
}

// AdsReadDeviceInfoResponse represents the response for an ADS ReadDeviceInfo command.
type AdsReadDeviceInfoResponse struct {
	BaseAdsResponse
	MajorVersion uint8
	MinorVersion uint8
	VersionBuild uint16
	DeviceName   string
}

// AdsReadResponse represents the response for an ADS Read command.
type AdsReadResponse struct {
	BaseAdsResponse
	Length uint32
	Data   []byte
}

// AdsWriteResponse represents the response for an ADS Write command.
type AdsWriteResponse struct {
	BaseAdsResponse
}

// AdsReadWriteResponse represents the response for an ADS ReadWrite command.
type AdsReadWriteResponse struct {
	BaseAdsResponse
	Length uint32
	Data   []byte
}

// AdsReadStateResponse represents the response for an ADS ReadState command.
type AdsReadStateResponse struct {
	BaseAdsResponse
	ADSState AdsStateResponse
}

// AdsWriteControlResponse represents the response for an ADS WriteControl command.
type AdsWriteControlResponse struct {
	BaseAdsResponse
}

// AdsTcSystemStateResponse represents the TwinCAT system state response.
type AdsTcSystemStateResponse struct {
	BaseAdsResponse
	AdsState    ADSState
	DeviceState uint16
}
