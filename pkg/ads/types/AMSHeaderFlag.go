package types

// AMSHeaderFlag defines the AMS command flags.
type AMSHeaderFlag uint16

const (
	AMSTCPPortAMSCommand AMSHeaderFlag = 0    // 0x0000 - Used for ADS commands
	AMSTCPPortClose      AMSHeaderFlag = 1    // 0x0001 - Port close command
	AMSTCPPortConnect    AMSHeaderFlag = 4096 // 0x1000 - Port connect command
	AMSTCPPortRouterNote AMSHeaderFlag = 4097 // 0x1001 - Router notification
	GetLocalNetID        AMSHeaderFlag = 4098 // 0x1002 - Requests local AmsNetId
)

// Returns the corresponding key as string by given value (number)
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
