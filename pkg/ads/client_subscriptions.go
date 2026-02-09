package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// SubscribeValue subscribes to a variable by path (like ReadValue).
// The value is automatically parsed using the symbol's data type.
// The callback is called with parsed data (int32, map[string]any, []any, etc.)
// whenever the value changes or at the specified cycle time.
//
// Example:
//
//	sub, err := client.SubscribeValue(851, "GVL.MyVariable",
//	    func(data ads.SubscriptionData) {
//	        fmt.Printf("Value: %v at %s\n", data.Value, data.Timestamp)
//	    },
//	    ads.SubscriptionSettings{
//	        CycleTime:    100 * time.Millisecond,
//	        SendOnChange: true,
//	    },
//	)
func (c *Client) SubscribeValue(port uint16, path string, callback SubscriptionCallback, settings SubscriptionSettings) (*ActiveSubscription, error) {
	c.logger.Debug("SubscribeValue: Subscribing to value", "port", port, "path", path)

	// Get symbol info (like ReadValue does)
	symbol, err := c.GetSymbol(port, path)
	if err != nil {
		return nil, fmt.Errorf("SubscribeValue: failed to get symbol: %w", err)
	}
	c.logger.Debug("SubscribeValue: Symbol received", "symbol", symbol)

	// Get data type (like ReadValue does)
	dataType, err := c.GetDataType(symbol.Type, port)
	if err != nil {
		return nil, fmt.Errorf("SubscribeValue: failed to get data type: %w", err)
	}

	// Subscribe using raw address with symbol and data type info
	return c.addSubscription(port, symbol.IndexGroup, symbol.IndexOffset, symbol.Size, callback, settings, symbol, &dataType, false)
}

// SubscribeRaw subscribes to a variable by raw ADS address (no parsing).
// The callback is called with raw []byte data whenever the value changes
// or at the specified cycle time.
//
// Example:
//
//	sub, err := client.SubscribeRaw(851, 16448, 414816, 4,
//	    func(data ads.SubscriptionData) {
//	        value := binary.LittleEndian.Uint32(data.RawValue)
//	        fmt.Printf("Raw value: %d\n", value)
//	    },
//	    ads.SubscriptionSettings{
//	        CycleTime:    50 * time.Millisecond,
//	        SendOnChange: true,
//	    },
//	)
func (c *Client) SubscribeRaw(port uint16, indexGroup, indexOffset, size uint32, callback SubscriptionCallback, settings SubscriptionSettings) (*ActiveSubscription, error) {
	c.logger.Debug("SubscribeRaw: Subscribing to raw address", "port", port, "indexGroup", indexGroup, "indexOffset", indexOffset, "size", size)
	return c.addSubscription(port, indexGroup, indexOffset, size, callback, settings, nil, nil, true)
}

// addSubscription is the internal method that sends the AddNotification command.
func (c *Client) addSubscription(port uint16, indexGroup, indexOffset, size uint32, callback SubscriptionCallback, settings SubscriptionSettings, symbol *types.AdsSymbol, dataType *types.AdsDataType, isRaw bool) (*ActiveSubscription, error) {
	c.logger.Debug("addSubscription: Creating subscription", "port", port, "indexGroup", indexGroup, "indexOffset", indexOffset)

	// Apply defaults
	if settings.CycleTime == 0 {
		settings.CycleTime = 200 * time.Millisecond
	}

	// Build 40-byte AddNotification request payload
	payload := make([]byte, 40)
	pos := 0

	// 0..3 IndexGroup
	binary.LittleEndian.PutUint32(payload[pos:], indexGroup)
	pos += 4

	// 4..7 IndexOffset
	binary.LittleEndian.PutUint32(payload[pos:], indexOffset)
	pos += 4

	// 8..11 Data length
	binary.LittleEndian.PutUint32(payload[pos:], size)
	pos += 4

	// 12..15 Transmission mode (3=Cyclic, 4=OnChange)
	transmissionMode := uint32(types.ADSTransModeOnChange) // 4 = OnChange
	if !settings.SendOnChange {
		transmissionMode = uint32(types.ADSTransModeCyclic) // 3 = Cyclic
	}
	binary.LittleEndian.PutUint32(payload[pos:], transmissionMode)
	pos += 4

	// 16..19 Maximum delay in 100ns units (milliseconds * 10000)
	maxDelayUnits := uint32(settings.MaxDelay.Milliseconds() * 10000)
	binary.LittleEndian.PutUint32(payload[pos:], maxDelayUnits)
	pos += 4

	// 20..23 Cycle time in 100ns units (milliseconds * 10000)
	cycleTimeUnits := uint32(settings.CycleTime.Milliseconds() * 10000)
	binary.LittleEndian.PutUint32(payload[pos:], cycleTimeUnits)
	pos += 4

	// 24..39 Reserved (zeros)
	// Already zero-initialized

	// Send AddNotification command
	responseData, err := c.send(AdsCommandRequest{
		Command:    types.ADSCommandAddNotification,
		TargetPort: port,
		Data:       payload,
	})
	if err != nil {
		c.logger.Error("addSubscription: Failed to send AddNotification command", "error", err)
		return nil, fmt.Errorf("addSubscription: failed to send AddNotification command: %w", err)
	}

	// Parse notification handle from response
	// Response format: [0..3] Error code, [4..7] Notification handle
	if len(responseData) < 8 {
		return nil, fmt.Errorf("addSubscription: invalid response length: %d bytes (expected 8)", len(responseData))
	}

	// Check error code (bytes 0-3)
	errorCode := binary.LittleEndian.Uint32(responseData[0:4])
	if errorCode != 0 {
		return nil, fmt.Errorf("addSubscription: ADS error code %d returned", errorCode)
	}

	// Parse notification handle (bytes 4-7)
	notificationHandle := binary.LittleEndian.Uint32(responseData[4:8])

	c.logger.Info("addSubscription: Subscription created", "handle", notificationHandle, "port", port)

	// Create ActiveSubscription
	sub := &ActiveSubscription{
		Handle:   notificationHandle,
		Port:     port,
		Symbol:   symbol,
		DataType: dataType,
		Settings: settings,
		Callback: callback,
		IsRaw:    isRaw,
	}

	// Store in subscriptions map (thread-safe)
	c.subscriptionsMutex.Lock()
	c.subscriptions[notificationHandle] = sub
	c.subscriptionsMutex.Unlock()

	return sub, nil
}

// Unsubscribe removes a subscription by sending DeleteNotification command.
//
// Example:
//
//	err := client.Unsubscribe(sub)
func (c *Client) Unsubscribe(sub *ActiveSubscription) error {
	c.logger.Debug("Unsubscribe: Unsubscribing", "handle", sub.Handle, "port", sub.Port)

	// Build 4-byte DeleteNotification request
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload[0:4], sub.Handle)

	// Send DeleteNotification command
	_, err := c.send(AdsCommandRequest{
		Command:    types.ADSCommandDeleteNotification,
		TargetPort: sub.Port,
		Data:       payload,
	})
	if err != nil {
		c.logger.Error("Unsubscribe: Failed to send DeleteNotification command", "error", err, "handle", sub.Handle)
		return fmt.Errorf("Unsubscribe: failed to send DeleteNotification command: %w", err)
	}

	// Remove from subscriptions map (thread-safe)
	c.subscriptionsMutex.Lock()
	delete(c.subscriptions, sub.Handle)
	c.subscriptionsMutex.Unlock()

	c.logger.Info("Unsubscribe: Subscription removed", "handle", sub.Handle, "port", sub.Port)
	return nil
}

// UnsubscribeAll removes all active subscriptions.
//
// Example:
//
//	err := client.UnsubscribeAll()
func (c *Client) UnsubscribeAll() error {
	c.logger.Debug("UnsubscribeAll: Unsubscribing from all subscriptions")

	// Get copy of all subscriptions (thread-safe)
	c.subscriptionsMutex.RLock()
	subs := make([]*ActiveSubscription, 0, len(c.subscriptions))
	for _, sub := range c.subscriptions {
		subs = append(subs, sub)
	}
	c.subscriptionsMutex.RUnlock()

	// Unsubscribe each
	var firstError error
	successCount := 0
	for _, sub := range subs {
		if err := c.Unsubscribe(sub); err != nil {
			if firstError == nil {
				firstError = err
			}
			c.logger.Error("UnsubscribeAll: Failed to unsubscribe", "handle", sub.Handle, "error", err)
		} else {
			successCount++
		}
	}

	c.logger.Info("UnsubscribeAll: Unsubscribed from subscriptions", "total", len(subs), "success", successCount)

	if firstError != nil {
		return fmt.Errorf("UnsubscribeAll: failed to unsubscribe some subscriptions: %w", firstError)
	}
	return nil
}

// parseNotification parses the raw notification packet received from the PLC.
// Returns a slice of notification stamps, each containing one or more samples.
func parseNotification(data []byte) ([]notificationStamp, error) {
	reader := bytes.NewReader(data)

	// 0..3 Length
	var length uint32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return nil, fmt.Errorf("parseNotification: failed to read length: %w", err)
	}

	// 4..7 Stamp count
	var stampCount uint32
	if err := binary.Read(reader, binary.LittleEndian, &stampCount); err != nil {
		return nil, fmt.Errorf("parseNotification: failed to read stamp count: %w", err)
	}

	stamps := make([]notificationStamp, 0, stampCount)

	// Parse each stamp
	for i := uint32(0); i < stampCount; i++ {
		stamp := notificationStamp{}

		// 0..7 Timestamp (Windows FILETIME - 100ns intervals since Jan 1, 1601)
		var filetime uint64
		if err := binary.Read(reader, binary.LittleEndian, &filetime); err != nil {
			return nil, fmt.Errorf("parseNotification: failed to read timestamp: %w", err)
		}
		stamp.Timestamp = filetimeToTime(filetime)

		// 8..11 Sample count
		var sampleCount uint32
		if err := binary.Read(reader, binary.LittleEndian, &sampleCount); err != nil {
			return nil, fmt.Errorf("parseNotification: failed to read sample count: %w", err)
		}

		stamp.Samples = make([]notificationSample, 0, sampleCount)

		// Parse each sample
		for j := uint32(0); j < sampleCount; j++ {
			sample := notificationSample{}

			// 0..3 Notification handle
			if err := binary.Read(reader, binary.LittleEndian, &sample.Handle); err != nil {
				return nil, fmt.Errorf("parseNotification: failed to read notification handle: %w", err)
			}

			// 4..7 Data length
			var dataLength uint32
			if err := binary.Read(reader, binary.LittleEndian, &dataLength); err != nil {
				return nil, fmt.Errorf("parseNotification: failed to read data length: %w", err)
			}

			// 8..n Data
			sample.Payload = make([]byte, dataLength)
			if _, err := reader.Read(sample.Payload); err != nil {
				return nil, fmt.Errorf("parseNotification: failed to read payload: %w", err)
			}

			stamp.Samples = append(stamp.Samples, sample)
		}

		stamps = append(stamps, stamp)
	}

	return stamps, nil
}

// filetimeToTime converts Windows FILETIME to Go time.Time.
// FILETIME = 100-nanosecond intervals since Jan 1, 1601.
// Unix epoch = Jan 1, 1970.
func filetimeToTime(filetime uint64) time.Time {
	// Windows epoch (1601) to Unix epoch (1970) difference in seconds
	const windowsEpochDiff = 11644473600

	// Convert to seconds and subtract epoch difference
	seconds := int64(filetime/10000000 - windowsEpochDiff)
	nanos := int64((filetime % 10000000) * 100)
	return time.Unix(seconds, nanos)
}

// handleNotification processes received notification packets and routes them
// to the appropriate subscription callbacks.
func (c *Client) handleNotification(data []byte) {
	c.logger.Debug("handleNotification: Processing notification packet", "dataLen", len(data))

	// Parse notification packet
	stamps, err := parseNotification(data)
	if err != nil {
		c.logger.Error("handleNotification: Failed to parse notification", "error", err)
		return
	}

	c.logger.Debug("handleNotification: Parsed notification", "stampCount", len(stamps))

	// Process each stamp
	for _, stamp := range stamps {
		for _, sample := range stamp.Samples {
			// Look up subscription (thread-safe read)
			c.subscriptionsMutex.RLock()
			sub := c.subscriptions[sample.Handle]
			c.subscriptionsMutex.RUnlock()

			if sub == nil {
				c.logger.Warn("handleNotification: Unknown notification handle", "handle", sample.Handle)
				continue
			}

			c.logger.Debug("handleNotification: Processing notification for subscription", "handle", sample.Handle, "port", sub.Port)

			// Process notification in a goroutine (don't block)
			go c.processNotification(sub, sample.Payload, stamp.Timestamp)
		}
	}
}

// processNotification parses the notification data and calls the user callback.
// Runs in a separate goroutine to prevent blocking notification processing.
func (c *Client) processNotification(sub *ActiveSubscription, rawData []byte, timestamp time.Time) {
	// Recover from panics in user callback
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("processNotification: Subscription callback panic", "handle", sub.Handle, "panic", r)
		}
	}()

	var value any
	var err error

	if sub.IsRaw || sub.DataType == nil {
		// Raw subscription - don't parse, return raw bytes
		value = rawData
	} else {
		// Parse using existing convertBufferToValue (like ReadValue does)
		value, err = c.convertBufferToValue(rawData, *sub.DataType)
		if err != nil {
			c.logger.Error("processNotification: Failed to parse notification value", "error", err, "handle", sub.Handle)
			return
		}
	}

	c.logger.Debug("processNotification: Calling user callback", "handle", sub.Handle, "timestamp", timestamp)

	// Call user callback
	sub.Callback(SubscriptionData{
		Value:     value,
		RawValue:  rawData,
		Timestamp: timestamp,
	})
}
