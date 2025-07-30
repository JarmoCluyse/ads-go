package ads

import "time"

// ClientSettings holds the settings for the ADS client.
type ClientSettings struct {
	TargetNetID string        // target ams net id (127.0.0.1.1.1 asumed if empty)
	RouterAddr  string        // adres of the router (127.0.0.1 asumed if empty)
	RouterPort  int           // port of the router (48898 asumed if empty)
	Timeout     time.Duration // message timeout (2s assumed if empty)
}

// LoadDefaults sets the default values for any unset ClientSettings fields.
func (cs *ClientSettings) LoadDefaults() {
	if cs.TargetNetID == "" {
		cs.TargetNetID = "127.0.0.1.1.1"
	}
	if cs.RouterAddr == "" {
		cs.RouterAddr = "127.0.0.1"
	}
	if cs.RouterPort == 0 {
		cs.RouterPort = 48898
	}
	if cs.Timeout == 0 {
		cs.Timeout = 2 * time.Second
	}
}
