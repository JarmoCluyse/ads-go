package ads

import (
	"bytes"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// Client represents an ADS client.
type Client struct {
	conn          net.Conn
	settings      ClientSettings
	mutex         sync.Mutex
	invokeID      uint32
	requests      map[uint32]chan []byte
	localAmsAddr  AmsAddress
	receiveBuffer bytes.Buffer // Buffer for incoming data
	logger        *slog.Logger
	plcSymbols    map[string]types.AdsSymbol
	plcDataTypes  map[string]types.AdsDataType
}

// ClientSettings holds the settings for the ADS client.
type ClientSettings struct {
	TargetNetID   string
	TargetPort    uint16
	RouterAddr    string
	Timeout       time.Duration
	AllowHalfOpen bool
}

// NewClient creates a new ADS client.
func NewClient(settings ClientSettings, logger *slog.Logger) *Client {
	logger.Info("NewClient: Initializing new ADS client.")
	if settings.Timeout == 0 {
		settings.Timeout = 2 * time.Second
		logger.Info("NewClient: Timeout not set, defaulting to 2 seconds.")
	}
	client := &Client{
		settings:     settings,
		requests:     make(map[uint32]chan []byte),
		logger:       logger,
		plcSymbols:   make(map[string]types.AdsSymbol),
		plcDataTypes: make(map[string]types.AdsDataType),
	}
	logger.Info("NewClient: ADS client initialized.")
	return client
}

