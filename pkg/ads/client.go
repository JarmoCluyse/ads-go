package ads

import (
	"bytes"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
)

// Response represents a response from an ADS device.
type Response struct {
	Data  []byte // received data
	Error error  // received ads error
}

// Client represents an ADS client.
type Client struct {
	conn          net.Conn                 // tcp connection
	settings      ClientSettings           // client settings
	mutex         sync.Mutex               // mutex for invoke id and request map
	invokeID      uint32                   // last used invoke id
	requests      map[uint32]chan Response // channel map to write the responses to
	localAmsAddr  AmsAddress               // local asigned ams adres
	receiveBuffer bytes.Buffer             // Buffer for incoming data
	logger        *slog.Logger             // logger
}

// ClientSettings holds the settings for the ADS client.
type ClientSettings struct {
	TargetNetID       string        // target ams net id
	RouterAddr        string        // adres of the router (127.0.0.1 asumed if empty)
	Timeout           time.Duration // message timeout
	AdsSymbolsUseUtf8 bool          // bool if names are utf8  TODO: check needed! for our purpose
}

// NewClient creates a new ADS client.
func NewClient(settings ClientSettings, logger *slog.Logger) *Client {
	if logger == nil { // silent logger when not added
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

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
