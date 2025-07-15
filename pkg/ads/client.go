package ads

import (
	"bytes"
	"log/slog"
	"net"
	"sync"
	"time"
)

// Response represents a response from an ADS device.
type Response struct {
	Data  []byte
	Error error
}

// Client represents an ADS client.
type Client struct {
	conn          net.Conn
	settings      ClientSettings
	mutex         sync.Mutex
	invokeID      uint32
	requests      map[uint32]chan Response
	localAmsAddr  AmsAddress
	receiveBuffer bytes.Buffer // Buffer for incoming data
	logger        *slog.Logger
}

// ClientSettings holds the settings for the ADS client.
type ClientSettings struct {
	TargetNetID       string
	RouterAddr        string
	Timeout           time.Duration
	AdsSymbolsUseUtf8 bool
}

// NewClient creates a new ADS client.
func NewClient(settings ClientSettings, logger *slog.Logger) *Client {
	logger.Info("NewClient: Initializing new ADS client.")
	if settings.Timeout == 0 {
		settings.Timeout = 2 * time.Second
		logger.Info("NewClient: Timeout not set, defaulting to 2 seconds.")
	}
	client := &Client{
		settings: settings,
		requests: make(map[uint32]chan Response),
		logger:   logger,
	}
	logger.Info("NewClient: ADS client initialized.")
	return client
}
