package ads

import (
	"time"

	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// SubscriptionSettings configures a subscription.
type SubscriptionSettings struct {
	// CycleTime is how often the PLC checks for value changes (default: 200ms).
	// If SendOnChange is true, PLC checks if value has changed with CycleTime interval.
	// If SendOnChange is false, PLC constantly sends the value with CycleTime interval.
	CycleTime time.Duration

	// SendOnChange determines if notifications are sent only when the value changes (default: true).
	// If true, the value is checked every CycleTime and sent only if it has changed.
	// If false, the value is sent cyclically every CycleTime (even if unchanged).
	// NOTE: When subscribing, the value is always sent once.
	SendOnChange bool

	// MaxDelay is how long the PLC waits before sending values at maximum (default: 0 = off).
	// If value is not changing, first notification after subscribing is sent after MaxDelay.
	// If the value is changing, PLC sends one or more notifications every MaxDelay.
	// This can be useful for throttling high-frequency changes.
	MaxDelay time.Duration
}

// SubscriptionData contains the parsed value and timestamp from a subscription notification.
type SubscriptionData struct {
	// Value is the parsed value (int32, map[string]any, []any, etc.) for SubscribeValue.
	// For SubscribeRaw, this will be the raw []byte data (same as RawValue).
	Value any

	// RawValue is the original raw bytes received from the PLC.
	RawValue []byte

	// Timestamp is when the PLC sampled the value.
	Timestamp time.Time
}

// SubscriptionCallback is called when a subscribed value changes or when cycle time elapses.
type SubscriptionCallback func(data SubscriptionData)

// ActiveSubscription represents an active subscription.
type ActiveSubscription struct {
	// Handle is the notification handle assigned by the PLC.
	Handle uint32

	// Port is the target ADS port.
	Port uint16

	// Symbol contains symbol information (if subscribed by path via SubscribeValue).
	Symbol *types.AdsSymbol

	// DataType contains full data type information (if subscribed by path via SubscribeValue).
	DataType *types.AdsDataType

	// Settings are the subscription settings.
	Settings SubscriptionSettings

	// Callback is the user-provided callback function.
	Callback SubscriptionCallback

	// IsRaw indicates if this is a raw subscription (SubscribeRaw).
	// If true, data is not parsed and Value = RawValue in callback.
	IsRaw bool
}

// notificationStamp represents a timestamp with multiple notification samples.
// Internal type used when parsing notification packets from the PLC.
type notificationStamp struct {
	Timestamp time.Time
	Samples   []notificationSample
}

// notificationSample represents a single notification sample with handle and data.
// Internal type used when parsing notification packets from the PLC.
type notificationSample struct {
	Handle  uint32
	Payload []byte
}
