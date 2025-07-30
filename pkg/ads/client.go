package ads

import (
	"bytes"
	"io"
	"log/slog"
	"net"
	"sync"
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

// NewClient creates a new ADS client.
func NewClient(settings ClientSettings, logger *slog.Logger) *Client {
	if logger == nil { // silent logger when not added
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	logger.Info("NewClient: Initializing new ADS client.")
	settings.LoadDefaults()
	client := &Client{
		settings: settings,
		requests: make(map[uint32]chan Response),
		logger:   logger,
	}
	logger.Info("NewClient: ADS client initialized.")
	return client
}
