package cli

import (
	"os"
	"time"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

const (
	DefaultAdsTargetNetID = "192.168.157.131.1.1"
	DefaultAdsTargetPort  = 350
	DefaultAdsRouterAddr  = "127.0.0.1:48898"
	DefaultAdsTimeout     = 5 * time.Second
)

func GetConfig() ads.ClientSettings {
	// Read ADS client settings from environment variables
	adsTargetNetID := os.Getenv("ADS_TARGET_NET_ID")
	adsRouterAddr := os.Getenv("ADS_ROUTER_ADDR")
	adsTimeoutStr := os.Getenv("ADS_TIMEOUT")

	// Parse timeout
	timeout := DefaultAdsTimeout
	if adsTimeoutStr != "" {
		if d, err := time.ParseDuration(adsTimeoutStr); err == nil {
			timeout = d
		}
	}

	// Use defaults if vars not set
	if adsTargetNetID == "" {
		adsTargetNetID = DefaultAdsTargetNetID
	}
	if adsRouterAddr == "" {
		adsRouterAddr = DefaultAdsRouterAddr
	}

	settings := ads.ClientSettings{
		TargetNetID:       adsTargetNetID,
		RouterAddr:        adsRouterAddr,
		Timeout:           timeout,
		AdsSymbolsUseUtf8: true,
	}
	return settings
}
