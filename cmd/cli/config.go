package cli

import (
	"os"
	"strconv"
	"time"

	"github.com/jarmoCluyse/ads-go/pkg/ads"
)

const (
	DefaultAdsTargetNetID   = "192.168.157.131.1.1"
	DefaultAdsTargetPort    = 350
	DefaultAdsRouterAddr    = "127.0.0.1:48898"
	DefaultAdsTimeout       = 5 * time.Second
	DefaultAdsAllowHalfOpen = true
)

func GetConfig() ads.ClientSettings {
	// Read ADS client settings from environment variables
	adsTargetNetID := os.Getenv("ADS_TARGET_NET_ID")
	adsTargetPortStr := os.Getenv("ADS_TARGET_PORT")
	adsRouterAddr := os.Getenv("ADS_ROUTER_ADDR")
	adsTimeoutStr := os.Getenv("ADS_TIMEOUT")
	adsAllowHalfOpenStr := os.Getenv("ADS_ALLOW_HALF_OPEN")

	// Parse integer port
	adsTargetPort := DefaultAdsTargetPort
	if adsTargetPortStr != "" {
		if v, err := strconv.Atoi(adsTargetPortStr); err == nil {
			adsTargetPort = v
		}
	}
	// Parse timeout
	timeout := DefaultAdsTimeout
	if adsTimeoutStr != "" {
		if d, err := time.ParseDuration(adsTimeoutStr); err == nil {
			timeout = d
		}
	}

	// Parse allow half open
	allowHalfOpen := DefaultAdsAllowHalfOpen
	if adsAllowHalfOpenStr != "" {
		if v, err := strconv.ParseBool(adsAllowHalfOpenStr); err == nil {
			allowHalfOpen = v
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
		TargetPort:        uint16(adsTargetPort),
		RouterAddr:        adsRouterAddr,
		Timeout:           timeout,
		AllowHalfOpen:     allowHalfOpen,
		AdsSymbolsUseUtf8: true,
	}
	return settings
}
