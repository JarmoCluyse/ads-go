package types

// ADSTransMode defines ADS notification transmission mode.
// Source: TwinCAT.Ads.dll By Beckhoff
type ADSTransMode uint32

const (
	ADSTransModeNone              ADSTransMode = 0
	ADSTransModeClientCycle       ADSTransMode = 1
	ADSTransModeClientOnChange    ADSTransMode = 2
	ADSTransModeCyclic            ADSTransMode = 3
	ADSTransModeOnChange          ADSTransMode = 4
	ADSTransModeCyclicInContext   ADSTransMode = 5
	ADSTransModeOnChangeInContext ADSTransMode = 6
)

// ADSTransModeToString returns the name of the ADS transmission mode for a given value.
func (mode ADSTransMode) String() string {
	switch mode {
	case ADSTransModeNone:
		return "None"
	case ADSTransModeClientCycle:
		return "ClientCycle"
	case ADSTransModeClientOnChange:
		return "ClientOnChange"
	case ADSTransModeCyclic:
		return "Cyclic"
	case ADSTransModeOnChange:
		return "OnChange"
	case ADSTransModeCyclicInContext:
		return "CyclicInContext"
	case ADSTransModeOnChangeInContext:
		return "OnChangeInContext"
	default:
		return "UNKNOWN"
	}
}
