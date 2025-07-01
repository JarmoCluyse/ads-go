
package ads

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

const (
	AMSTCPHeaderLength   = 6
	AMSHeaderLength      = 32
	AMSNetIDLength       = 6
	ADSIndexOffsetLength = 4
	ADSIndexGroupLength  = 4
	ADSInvokeIDMaxValue  = 4294967295
	ADSDefaultTCPPort    = 48898
	LoopbackAmsNetID     = "127.0.0.1.1.1"
)

// AMSHeaderFlag defines the AMS command flags.
type AMSHeaderFlag uint16

const (
	AMSTCPPortAMSCommand AMSHeaderFlag = 0
	AMSTCPPortClose      AMSHeaderFlag = 1
	AMSTCPPortConnect    AMSHeaderFlag = 4096
	AMSTCPPortRouterNote AMSHeaderFlag = 4097
	GetLocalNetID        AMSHeaderFlag = 4098
)

func (f AMSHeaderFlag) String() string {
	switch f {
	case AMSTCPPortAMSCommand:
		return "AMS_TCP_PORT_AMS_CMD"
	case AMSTCPPortClose:
		return "AMS_TCP_PORT_CLOSE"
	case AMSTCPPortConnect:
		return "AMS_TCP_PORT_CONNECT"
	case AMSTCPPortRouterNote:
		return "AMS_TCP_PORT_ROUTER_NOTE"
	case GetLocalNetID:
		return "GET_LOCAL_NETID"
	default:
		return "UNKNOWN"
	}
}

// ADSCommand defines the ADS commands.
type ADSCommand uint16

const (
	ADSCommandInvalid        ADSCommand = 0
	ADSCommandNone           ADSCommand = 0
	ADSCommandReadDeviceInfo ADSCommand = 1
	ADSCommandRead           ADSCommand = 2
	ADSCommandWrite          ADSCommand = 3
	ADSCommandReadState      ADSCommand = 4
	ADSCommandWriteControl   ADSCommand = 5
	ADSCommandAddNotification ADSCommand = 6
	ADSCommandDeleteNotification ADSCommand = 7
	ADSCommandNotification   ADSCommand = 8
	ADSCommandReadWrite      ADSCommand = 9
)

func (c ADSCommand) String() string {
	switch c {
	case ADSCommandReadDeviceInfo:
		return "ReadDeviceInfo"
	case ADSCommandRead:
		return "Read"
	case ADSCommandWrite:
		return "Write"
	case ADSCommandReadState:
		return "ReadState"
	case ADSCommandWriteControl:
		return "WriteControl"
	case ADSCommandAddNotification:
		return "AddNotification"
	case ADSCommandDeleteNotification:
		return "DeleteNotification"
	case ADSCommandNotification:
		return "Notification"
	case ADSCommandReadWrite:
		return "ReadWrite"
	default:
		return "UNKNOWN"
	}
}

// AmsAddress represents the AMS address.
type AmsAddress struct {
	NetID string
	Port  uint16
}

func (a AmsAddress) String() string {
	return fmt.Sprintf("%s:%d", a.NetID, a.Port)
}

// AmsNetIDStrToByteArray converts an AmsNetId string to a byte array.
func AmsNetIDStrToByteArray(s string) ([]byte, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid AmsNetId: %s", s)
	}
	bytes := make([]byte, 6)
	for i, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid part in AmsNetId: %s", part)
		}
		bytes[i] = byte(val)
	}
	return bytes, nil
}

// ByteArrayToAmsNetIDStr converts a byte array to an AmsNetId string.
func ByteArrayToAmsNetIDStr(b []byte) string {
	parts := make([]string, len(b))
	for i, byte := range b {
		parts[i] = fmt.Sprintf("%d", byte)
	}
	return strings.Join(parts, ".")
}

// AmsTcpHeader represents the AMS/TCP header.
type AmsTcpHeader struct {
	Command AMSHeaderFlag
	Length  uint32
}

// AmsHeader represents the AMS header.
type AmsHeader struct {
	Target     AmsAddress
	Source     AmsAddress
	Command    ADSCommand
	StateFlags uint16
	DataLength uint32
	ErrorCode  uint32
	InvokeID   uint32
}

// AdsReadDeviceInfoResponse represents the response of a ReadDeviceInfo command.
type AdsReadDeviceInfoResponse struct {
	ErrorCode    uint32
	MajorVersion uint8
	MinorVersion uint8
	VersionBuild uint16
	DeviceName   string
}

// AdsReadStateResponse represents the response of a ReadState command.
type AdsReadStateResponse struct {
	ErrorCode uint32
	AdsState  uint16
	DeviceState uint16
}

// AdsWriteResponse represents the response of a Write command.
type AdsWriteResponse struct {
	ErrorCode uint32
}

// AdsReadResponse represents the response of a Read command.
type AdsReadResponse struct {
	ErrorCode uint32
	Length    uint32
	Data      []byte
}

// AdsReadWriteResponse represents the response of a ReadWrite command.
type AdsReadWriteResponse struct {
	ErrorCode uint32
	Length    uint32
	Data      []byte
}

// AdsAddNotificationResponse represents the response of an AddNotification command.
type AdsAddNotificationResponse struct {
	ErrorCode        uint32
	NotificationHandle uint32
}

// AdsDeleteNotificationResponse represents the response of a DeleteNotification command.
type AdsDeleteNotificationResponse struct {
	ErrorCode uint32
}

// AdsNotificationSample represents a single notification sample.
type AdsNotificationSample struct {
	NotificationHandle uint32
	SampleSize         uint32
	Data               []byte
}

// AdsNotificationStamp represents a notification stamp with multiple samples.
type AdsNotificationStamp struct {
	Timestamp uint64
	Samples   []AdsNotificationSample
}

// AdsNotificationResponse represents a notification response.
type AdsNotificationResponse struct {
	Length uint32
	Stamps []AdsNotificationStamp
}

// ADSState defines the ADS state of a device.
type ADSState uint16

const (
	ADSStateInvalid    ADSState = 0
	ADSStateIdle       ADSState = 1
	ADSStateReset      ADSState = 2
	ADSStateInitialize ADSState = 3
	ADSStateStart      ADSState = 4
	ADSStateRun        ADSState = 5
	ADSStateStop       ADSState = 6
	ADSStateSaveConfig ADSState = 7
	ADSStateLoadConfig ADSState = 8
	ADSStatePowerFailure ADSState = 9
	ADSStatePowerGood  ADSState = 10
	ADSStateError      ADSState = 11
	ADSStateShutdown   ADSState = 12
	ADSStateSuspend    ADSState = 13
	ADSStateResume     ADSState = 14
	ADSStateConfig     ADSState = 15
	ADSStateReconfig   ADSState = 16
	ADSStateStopping   ADSState = 17
)

// AdsWriteControlResponse represents the response of a WriteControl command.
type AdsWriteControlResponse struct {
	ErrorCode uint32
}

// createAmsTcpHeader creates the AMS/TCP header.
func createAmsTcpHeader(command AMSHeaderFlag, dataLength uint32) []byte {
	buf := make([]byte, AMSTCPHeaderLength)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(command))
	binary.LittleEndian.PutUint32(buf[2:6], dataLength)
	return buf
}

// createAmsHeader creates the AMS header.
func createAmsHeader(target AmsAddress, source AmsAddress, command ADSCommand, dataLength uint32, invokeID uint32) ([]byte, error) {
	buf := make([]byte, AMSHeaderLength)
	targetNetID, err := AmsNetIDStrToByteArray(target.NetID)
	if err != nil {
		return nil, err
	}
	copy(buf[0:6], targetNetID)
	binary.LittleEndian.PutUint16(buf[6:8], target.Port)

	sourceNetID, err := AmsNetIDStrToByteArray(source.NetID)
	if err != nil {
		return nil, err
	}
	copy(buf[8:14], sourceNetID)
	binary.LittleEndian.PutUint16(buf[14:16], source.Port)
	binary.LittleEndian.PutUint16(buf[16:18], uint16(command))
	binary.LittleEndian.PutUint16(buf[18:20], 0x0004) // StateFlags
	binary.LittleEndian.PutUint32(buf[20:24], dataLength)
	binary.LittleEndian.PutUint32(buf[24:28], 0) // ErrorCode
	binary.LittleEndian.PutUint32(buf[28:32], invokeID)
	return buf, nil
}
