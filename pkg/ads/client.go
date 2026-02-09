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
	conn               net.Conn                       // tcp connection
	settings           ClientSettings                 // client settings
	mutex              sync.Mutex                     // mutex for invoke id and request map
	invokeID           uint32                         // last used invoke id
	requests           map[uint32]chan Response       // channel map to write the responses to
	localAmsAddr       AmsAddress                     // local asigned ams adres
	receiveBuffer      bytes.Buffer                   // Buffer for incoming data
	logger             *slog.Logger                   // logger
	subscriptions      map[uint32]*ActiveSubscription // active subscriptions map[notificationHandle]subscription
	subscriptionsMutex sync.RWMutex                   // mutex for subscriptions map
}

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

// NewClient creates a new ADS client.
func NewClient(settings ClientSettings, logger *slog.Logger) *Client {
	if logger == nil { // silent logger when not added
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	logger.Info("NewClient: Initializing new ADS client.")
	settings.LoadDefaults()
	client := &Client{
		settings:      settings,
		requests:      make(map[uint32]chan Response),
		subscriptions: make(map[uint32]*ActiveSubscription),
		logger:        logger,
	}
	logger.Info("NewClient: ADS client initialized.")
	return client
}
