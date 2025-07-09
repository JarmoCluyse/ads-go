package types

// ADSState represents the state of ADS (from TwinCAT.Ads.dll by Beckhoff)
type ADSState int

const (
	ADSStateInvalid      ADSState = iota
	ADSStateIdle         ADSState = 1
	ADSStateReset        ADSState = 2
	ADSStateInitialize   ADSState = 3
	ADSStateStart        ADSState = 4
	ADSStateRun          ADSState = 5
	ADSStateStop         ADSState = 6
	ADSStateSaveConfig   ADSState = 7
	ADSStateLoadConfig   ADSState = 8
	ADSStatePowerFailure ADSState = 9
	ADSStatePowerGood    ADSState = 10
	ADSStateError        ADSState = 11
	ADSStateShutdown     ADSState = 12
	ADSStateSuspend      ADSState = 13
	ADSStateResume       ADSState = 14
	ADSStateConfig       ADSState = 15
	ADSStateReconfig     ADSState = 16
	ADSStateStopping     ADSState = 17
	ADSStateIncompatible ADSState = 18
	ADSStateException    ADSState = 19
)

// String returns the string representation of ADSState (like toString in TypeScript)
func (s ADSState) String() string {
	switch s {
	case ADSStateInvalid:
		return "Invalid"
	case ADSStateIdle:
		return "Idle"
	case ADSStateReset:
		return "Reset"
	case ADSStateInitialize:
		return "Initialize"
	case ADSStateStart:
		return "Start"
	case ADSStateRun:
		return "Run"
	case ADSStateStop:
		return "Stop"
	case ADSStateSaveConfig:
		return "SaveConfig"
	case ADSStateLoadConfig:
		return "LoadConfig"
	case ADSStatePowerFailure:
		return "PowerFailure"
	case ADSStatePowerGood:
		return "PowerGood"
	case ADSStateError:
		return "Error"
	case ADSStateShutdown:
		return "Shutdown"
	case ADSStateSuspend:
		return "Suspend"
	case ADSStateResume:
		return "Resume"
	case ADSStateConfig:
		return "Config"
	case ADSStateReconfig:
		return "Reconfig"
	case ADSStateStopping:
		return "Stopping"
	case ADSStateIncompatible:
		return "Incompatible"
	case ADSStateException:
		return "Exception"
	default:
		return "UNKNOWN"
	}
}
