package main

import (
	"log/slog"
	"os"

	"github.com/jarmocluyse/ads-go/example/cli"
	"github.com/jarmocluyse/ads-go/pkg/ads"
	"github.com/lmittmann/tint"
)

func main() {
	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.LevelDebug)
	// logLevel.Set(slog.LevelInfo)
	// logLevel.Set(slog.LevelWarn)

	handler := tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel})
	slog.SetDefault(slog.New(handler))
	// logger := slog.Default()

	slog.Info("main: Starting application")

	settings := cli.GetConfig()
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
