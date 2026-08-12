package main

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jarmocluyse/ads-go/cmd/cli"
	"github.com/jarmocluyse/ads-go/pkg/ads"
	adsstateinfo "github.com/jarmocluyse/ads-go/pkg/ads/ads-stateinfo"
	adsconstants "github.com/jarmocluyse/ads-go/pkg/ads/constants"
	"github.com/lmittmann/tint"
)

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.LevelDebug)

	handler := tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel})
	slog.SetDefault(slog.New(handler))
	slog.Info("main: Starting application")

	// Resolve env var defaults before defining flags so that flag --help shows the effective default.
	defaultTargetNetID := envOrDefault("ADS_TARGET_NET_ID", "127.0.0.1.1.1")
	defaultRouterHost := envOrDefault("ADS_ROUTER_HOST", "127.0.0.1")

	defaultRouterPort := adsconstants.ADSDefaultTCPPort
	if v := os.Getenv("ADS_ROUTER_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			defaultRouterPort = p
		} else {
			slog.Warn("main: invalid ADS_ROUTER_PORT, using default", "value", v, "default", defaultRouterPort)
		}
	}

	defaultTimeout := 2 * time.Second
	if v := os.Getenv("ADS_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			defaultTimeout = d
		} else {
			slog.Warn("main: invalid ADS_TIMEOUT, using default", "value", v, "default", defaultTimeout)
		}
	}

	targetNetID := flag.String("target-net-id", defaultTargetNetID, "AMS NetID of the target TwinCAT runtime (env: ADS_TARGET_NET_ID)")
	routerHost := flag.String("router-host", defaultRouterHost, "Hostname or IP of the AMS router (env: ADS_ROUTER_HOST)")
	routerPort := flag.Int("router-port", defaultRouterPort, "TCP port of the AMS router (env: ADS_ROUTER_PORT)")
	timeout := flag.Duration("timeout", defaultTimeout, "Per-message read/write timeout (env: ADS_TIMEOUT)")
	flag.Parse()

	settings := ads.ClientSettings{
		TargetNetID: *targetNetID,
		RouterHost:  *routerHost,
		RouterPort:  *routerPort,
		Timeout:     *timeout,
	}

	// Synchronization for reconnection logic
	var reconnecting sync.Mutex
	var gracefulDisconnect bool
	var gracefulDisconnectMutex sync.Mutex

	// Configure connection event hooks
	settings.OnConnect = func(client *ads.Client, addr ads.AmsAddress) error {
		slog.Info("EVENT: ADS client connected", "localAMS", addr.NetID, "port", addr.Port)
		return nil
	}

	settings.OnDisconnect = func(client *ads.Client) {
		slog.Info("EVENT: ADS client disconnected gracefully")
		gracefulDisconnectMutex.Lock()
		gracefulDisconnect = true
		gracefulDisconnectMutex.Unlock()
	}

	settings.OnConnectionLost = func(client *ads.Client, err error) {
		slog.Error("EVENT: ADS connection lost unexpectedly", "error", err)

		// Check if this was a graceful disconnect
		gracefulDisconnectMutex.Lock()
		wasGraceful := gracefulDisconnect
		gracefulDisconnectMutex.Unlock()

		if wasGraceful {
			slog.Debug("Skipping reconnection for graceful disconnect")
			return
		}

		// Prevent concurrent reconnection attempts
		if !reconnecting.TryLock() {
			slog.Debug("Reconnection already in progress, skipping")
			return
		}

		// Start reconnection loop
		go func() {
			defer reconnecting.Unlock()

			slog.Info("Starting reconnection loop...")
			attemptNum := 1
			reconnectInterval := 5 * time.Second

			for {
				slog.Info("Attempting to reconnect...", "attempt", attemptNum)

				// Wait before attempting reconnection
				time.Sleep(reconnectInterval)

				// Try to reconnect
				if err := client.Connect(); err != nil {
					slog.Warn("Reconnection attempt failed", "attempt", attemptNum, "error", err)
					attemptNum++
					continue
				}

				slog.Info("Successfully reconnected to ADS router!", "attempts", attemptNum)
				return
			}
		}()
	}

	settings.OnStateChange = func(client *ads.Client, newState, oldState *adsstateinfo.SystemState) {
		if oldState == nil {
			// Initial state read
			slog.Info("EVENT: Initial TwinCAT state read",
				"state", newState.AdsState.String(),
				"deviceState", newState.DeviceState)
		} else {
			slog.Info("EVENT: TwinCAT system state changed",
				"fromState", oldState.AdsState.String(),
				"toState", newState.AdsState.String(),
				"fromDeviceState", oldState.DeviceState,
				"toDeviceState", newState.DeviceState)
		}
	}

	// Create client with nil logger (silent internal logs)
	slog.Info("main: Creating new ADS client with settings", "settings", settings)

	client := ads.NewClient(settings, nil)
	slog.Debug("main: ADS client created.")

	slog.Info("main: Attempting to connect to ADS router...")
	if err := client.Connect(); err != nil {
		slog.Error("main: Failed to connect", "error", err)
		os.Exit(1)
	}

	defer func() {
		slog.Info("main: Disconnecting from ADS router...")
		if err := client.Disconnect(); err != nil {
			slog.Error("main: Error during disconnect", "error", err)
		}
		slog.Info("main: Disconnected from ADS router.")
	}()
	cli.Commandline(client)
}
