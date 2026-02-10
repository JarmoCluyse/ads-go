package main

import (
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/jarmocluyse/ads-go/example/cli"
	"github.com/jarmocluyse/ads-go/pkg/ads"
	adsstateinfo "github.com/jarmocluyse/ads-go/pkg/ads/ads-stateinfo"
	"github.com/lmittmann/tint"
)

func main() {
	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.LevelDebug)
	// logLevel.Set(slog.LevelInfo)
	// logLevel.Set(slog.LevelWarn)

	handler := tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel})
	slog.SetDefault(slog.New(handler))

	slog.Info("main: Starting application")

	settings := cli.GetConfig()
	slog.Info("main: Creating new ADS client with settings", "settings", settings)

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
			slog.Info("EVENT: Initial TwinCAT state read", "state", newState.AdsState.String())
		} else {
			slog.Info("EVENT: TwinCAT state changed",
				"from", oldState.AdsState.String(),
				"to", newState.AdsState.String())
		}
	}

	// Create client with nil logger (silent internal logs)
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
