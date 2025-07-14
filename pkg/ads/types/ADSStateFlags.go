package types

import "slices"

// ADSStateFlags defines the ADS state flags (from TwinCAT.Ads.dll by Beckhoff)
type ADSStateFlags uint16

const (
	ADSStateFlagResponse       ADSStateFlags = 1     // AMSCMDSF_RESPONSE
	ADSStateFlagNoReturn       ADSStateFlags = 2     // AMSCMDSF_NORETURN
	ADSStateFlagAdsCommand     ADSStateFlags = 4     // AdsCommand
	ADSStateFlagSysCommand     ADSStateFlags = 8     // AMSCMDSF_SYSCMD
	ADSStateFlagHighPriority   ADSStateFlags = 16    // AMSCMDSF_HIGHPRIO
	ADSStateFlagTimeStampAdded ADSStateFlags = 32    // AMSCMDSF_TIMESTAMPADDED
	ADSStateFlagUdp            ADSStateFlags = 64    // AMSCMDSF_UDP
	ADSStateFlagInitCmd        ADSStateFlags = 128   // AMSCMDSF_INITCMD
	ADSStateFlagBroadcast      ADSStateFlags = 32768 // AMSCMDSF_BROADCAST
)

// String returns the flags as a comma separated list, with special cases for Tcp/Request.
func (f ADSStateFlags) String() string {
	flags := make([]string, 0)
	if f&ADSStateFlagResponse == ADSStateFlagResponse {
		flags = append(flags, "Response")
	}
	if f&ADSStateFlagNoReturn == ADSStateFlagNoReturn {
		flags = append(flags, "NoReturn")
	}
	if f&ADSStateFlagAdsCommand == ADSStateFlagAdsCommand {
		flags = append(flags, "AdsCommand")
	}
	if f&ADSStateFlagSysCommand == ADSStateFlagSysCommand {
		flags = append(flags, "SysCommand")
	}
	if f&ADSStateFlagHighPriority == ADSStateFlagHighPriority {
		flags = append(flags, "HighPriority")
	}
	if f&ADSStateFlagTimeStampAdded == ADSStateFlagTimeStampAdded {
		flags = append(flags, "TimeStampAdded")
	}
	if f&ADSStateFlagUdp == ADSStateFlagUdp {
		flags = append(flags, "Udp")
	}
	if f&ADSStateFlagInitCmd == ADSStateFlagInitCmd {
		flags = append(flags, "InitCmd")
	}
	if f&ADSStateFlagBroadcast == ADSStateFlagBroadcast {
		flags = append(flags, "Broadcast")
	}

	present := func(name string) bool {
		return slices.Contains(flags, name)
	}
	if !present("Udp") {
		flags = append(flags, "Tcp")
	}
	if !present("Response") {
		flags = append(flags, "Request")
	}
	if len(flags) == 0 {
		return "0"
	}
	return joinADSStrings(flags, ", ")
}

func joinADSStrings(list []string, sep string) string {
	if len(list) == 0 {
		return ""
	}
	result := list[0]
	for _, s := range list[1:] {
		result += sep + s
	}
	return result
}
