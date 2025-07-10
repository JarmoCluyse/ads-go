package types

// AMSRouterState defines possible AMS router states.
// Source: TwinCAT.Ads.dll By Beckhoff
type AMSRouterState uint32

const (
	AMSRouterStateStop    AMSRouterState = 0 // Router is stopped
	AMSRouterStateStart   AMSRouterState = 1 // Router is started
	AMSRouterStateRemoved AMSRouterState = 2 // Router is removed (unavailable?)
)

// AMSRouterStateToString returns the string key for the given router state value.
func AMSRouterStateToString(value AMSRouterState) string {
	switch value {
	case AMSRouterStateStop:
		return "STOP"
	case AMSRouterStateStart:
		return "START"
	case AMSRouterStateRemoved:
		return "REMOVED"
	default:
		return "UNKNOWN"
	}
}

