package cli

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

// Subscription tracking state
var (
	subscriptions      = make(map[int]*ads.ActiveSubscription)
	subscriptionsMutex sync.RWMutex
	nextSubID          = 1
)

// SubscriptionStats tracks statistics for each subscription
type SubscriptionStats struct {
	LastValue         interface{}
	LastUpdate        time.Time
	NotificationCount int64
}

var (
	subscriptionStats      = make(map[int]*SubscriptionStats)
	subscriptionStatsMutex sync.RWMutex
)

// createSubscriptionCallback creates a callback function for a subscription
func createSubscriptionCallback(id int, path string) ads.SubscriptionCallback {
	return func(data ads.SubscriptionData) {
		fmt.Printf("[NOTIFICATION #%d] %s: %v (at %s)\n",
			id, path, data.Value, data.Timestamp.Format("15:04:05.000"))

		// Update statistics
		subscriptionStatsMutex.Lock()
		if subscriptionStats[id] == nil {
			subscriptionStats[id] = &SubscriptionStats{}
		}
		subscriptionStats[id].LastValue = data.Value
		subscriptionStats[id].LastUpdate = data.Timestamp
		subscriptionStats[id].NotificationCount++
		subscriptionStatsMutex.Unlock()
	}
}

// handleSubscribe subscribes to a variable and starts receiving notifications.
// Usage: subscribe [path]
// If path is not provided, uses a hardcoded test variable.
func handleSubscribe(args []string, client *ads.Client) {
	// Hardcoded test configuration
	var port uint16 = 852
	path := "GLOBAL.gMyBoolToogle"

	// Allow optional path override
	if len(args) > 0 {
		path = args[0]
	}

	// Create subscription settings
	settings := ads.SubscriptionSettings{
		CycleTime:    100 * time.Millisecond,
		SendOnChange: true,
	}

	// Get next subscription ID
	subscriptionsMutex.Lock()
	id := nextSubID
	nextSubID++
	subscriptionsMutex.Unlock()

	// Create callback
	callback := createSubscriptionCallback(id, path)

	// Subscribe
	sub, err := client.SubscribeValue(port, path, callback, settings)
	if err != nil {
		fmt.Printf("[ERROR] Command 'subscribe': Failed to subscribe to '%s' (port %d): %v\n", path, port, err)
		return
	}

	// Store subscription
	subscriptionsMutex.Lock()
	subscriptions[id] = sub
	subscriptionsMutex.Unlock()

	fmt.Printf("[OK] Subscription #%d created for '%s' (port %d)\n", id, path, port)
}

// handleListSubs lists all active subscriptions.
// Usage: list_subs
func handleListSubs(args []string, client *ads.Client) {
	subscriptionsMutex.RLock()
	defer subscriptionsMutex.RUnlock()

	if len(subscriptions) == 0 {
		fmt.Println("[INFO] No active subscriptions")
		return
	}

	fmt.Println("[INFO] Active subscriptions:")
	for id, sub := range subscriptions {
		onChangeStr := "false"
		if sub.Settings.SendOnChange {
			onChangeStr = "true"
		}

		// Get path from symbol if available, otherwise show "raw"
		path := "raw"
		if sub.Symbol != nil {
			path = sub.Symbol.Name
		}

		fmt.Printf("  #%d: %s (port %d)\n", id, path, sub.Port)
		fmt.Printf("      CycleTime: %dms | OnChange: %s\n",
			sub.Settings.CycleTime.Milliseconds(), onChangeStr)

		// Add statistics if available
		subscriptionStatsMutex.RLock()
		if stats, ok := subscriptionStats[id]; ok && stats.NotificationCount > 0 {
			timeSince := time.Since(stats.LastUpdate)
			timeStr := formatDuration(timeSince)
			fmt.Printf("      Last Value: %v | Last Update: %s ago | Notifications: %d\n",
				stats.LastValue, timeStr, stats.NotificationCount)
		} else {
			fmt.Println("      No notifications received yet")
		}
		subscriptionStatsMutex.RUnlock()
	}
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	} else if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

// handleUnsubscribe unsubscribes from a specific subscription by ID.
// Usage: unsubscribe <id>
func handleUnsubscribe(args []string, client *ads.Client) {
	if len(args) == 0 {
		fmt.Println("[ERROR] Command 'unsubscribe': Usage: unsubscribe <id>")
		return
	}

	// Parse subscription ID
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("[ERROR] Command 'unsubscribe': Invalid subscription ID '%s'\n", args[0])
		return
	}

	// Look up subscription
	subscriptionsMutex.Lock()
	sub, ok := subscriptions[id]
	if !ok {
		subscriptionsMutex.Unlock()
		fmt.Printf("[ERROR] Command 'unsubscribe': Subscription #%d not found\n", id)
		return
	}
	delete(subscriptions, id)
	subscriptionsMutex.Unlock()

	// Clean up stats
	subscriptionStatsMutex.Lock()
	delete(subscriptionStats, id)
	subscriptionStatsMutex.Unlock()

	// Unsubscribe from client
	err = client.Unsubscribe(sub)
	if err != nil {
		fmt.Printf("[ERROR] Command 'unsubscribe': Failed to unsubscribe #%d: %v\n", id, err)
		return
	}

	fmt.Printf("[OK] Unsubscribed from subscription #%d\n", id)
}

// handleUnsubscribeAll unsubscribes from all active subscriptions.
// Usage: unsubscribe_all
func handleUnsubscribeAll(args []string, client *ads.Client) {
	// Clear local maps first
	subscriptionsMutex.Lock()
	count := len(subscriptions)
	subscriptions = make(map[int]*ads.ActiveSubscription)
	subscriptionsMutex.Unlock()

	subscriptionStatsMutex.Lock()
	subscriptionStats = make(map[int]*SubscriptionStats)
	subscriptionStatsMutex.Unlock()

	// Unsubscribe from client
	err := client.UnsubscribeAll()
	if err != nil {
		fmt.Printf("[ERROR] Command 'unsubscribe_all': Failed to unsubscribe all: %v\n", err)
		return
	}

	fmt.Printf("[OK] Unsubscribed from all subscriptions (%d total)\n", count)
}

// handleSubCounter subscribes to the cycle-based counter (gMyIntCounter).
// Usage: sub_counter
func handleSubCounter(args []string, client *ads.Client) {
	var port uint16 = 852
	path := "GLOBAL.gMyIntCounter"

	settings := ads.SubscriptionSettings{
		CycleTime:    100 * time.Millisecond,
		SendOnChange: true,
	}

	subscriptionsMutex.Lock()
	id := nextSubID
	nextSubID++
	subscriptionsMutex.Unlock()

	callback := createSubscriptionCallback(id, path)

	sub, err := client.SubscribeValue(port, path, callback, settings)
	if err != nil {
		fmt.Printf("[ERROR] Command 'sub_counter': Failed to subscribe to '%s' (port %d): %v\n", path, port, err)
		return
	}

	subscriptionsMutex.Lock()
	subscriptions[id] = sub
	subscriptionsMutex.Unlock()

	fmt.Printf("[OK] Subscription #%d created for cycle-based counter '%s' (port %d)\n", id, path, port)
}

// handleSubToggle subscribes to the cycle-based toggle (gMyBoolToogle).
// Usage: sub_toggle
func handleSubToggle(args []string, client *ads.Client) {
	var port uint16 = 852
	path := "GLOBAL.gMyBoolToogle"

	settings := ads.SubscriptionSettings{
		CycleTime:    100 * time.Millisecond,
		SendOnChange: true,
	}

	subscriptionsMutex.Lock()
	id := nextSubID
	nextSubID++
	subscriptionsMutex.Unlock()

	callback := createSubscriptionCallback(id, path)

	sub, err := client.SubscribeValue(port, path, callback, settings)
	if err != nil {
		fmt.Printf("[ERROR] Command 'sub_toggle': Failed to subscribe to '%s' (port %d): %v\n", path, port, err)
		return
	}

	subscriptionsMutex.Lock()
	subscriptions[id] = sub
	subscriptionsMutex.Unlock()

	fmt.Printf("[OK] Subscription #%d created for cycle-based toggle '%s' (port %d)\n", id, path, port)
}

// handleSubTimedCounter subscribes to the time-based counter (gTimedIntCounter).
// Usage: sub_timed_counter
func handleSubTimedCounter(args []string, client *ads.Client) {
	var port uint16 = 852
	path := "GLOBAL.gTimedIntCounter"

	settings := ads.SubscriptionSettings{
		CycleTime:    500 * time.Millisecond,
		SendOnChange: true,
	}

	subscriptionsMutex.Lock()
	id := nextSubID
	nextSubID++
	subscriptionsMutex.Unlock()

	callback := createSubscriptionCallback(id, path)

	sub, err := client.SubscribeValue(port, path, callback, settings)
	if err != nil {
		fmt.Printf("[ERROR] Command 'sub_timed_counter': Failed to subscribe to '%s' (port %d): %v\n", path, port, err)
		return
	}

	subscriptionsMutex.Lock()
	subscriptions[id] = sub
	subscriptionsMutex.Unlock()

	fmt.Printf("[OK] Subscription #%d created for time-based counter '%s' (port %d)\n", id, path, port)
}

// handleSubTimedToggle subscribes to the time-based toggle (gTimedBoolToogle).
// Usage: sub_timed_toggle
func handleSubTimedToggle(args []string, client *ads.Client) {
	var port uint16 = 852
	path := "GLOBAL.gTimedBoolToogle"

	settings := ads.SubscriptionSettings{
		CycleTime:    500 * time.Millisecond,
		SendOnChange: true,
	}

	subscriptionsMutex.Lock()
	id := nextSubID
	nextSubID++
	subscriptionsMutex.Unlock()

	callback := createSubscriptionCallback(id, path)

	sub, err := client.SubscribeValue(port, path, callback, settings)
	if err != nil {
		fmt.Printf("[ERROR] Command 'sub_timed_toggle': Failed to subscribe to '%s' (port %d): %v\n", path, port, err)
		return
	}

	subscriptionsMutex.Lock()
	subscriptions[id] = sub
	subscriptionsMutex.Unlock()

	fmt.Printf("[OK] Subscription #%d created for time-based toggle '%s' (port %d)\n", id, path, port)
}

// handleSubAll subscribes to all counters and toggles at once.
// Usage: sub_all
func handleSubAll(args []string, client *ads.Client) {
	var port uint16 = 852

	// Variables to subscribe to
	variables := []struct {
		path      string
		cycleTime time.Duration
		name      string
	}{
		{"GLOBAL.gMyIntCounter", 100 * time.Millisecond, "cycle-based counter"},
		{"GLOBAL.gMyBoolToogle", 100 * time.Millisecond, "cycle-based toggle"},
		{"GLOBAL.gTimedIntCounter", 500 * time.Millisecond, "time-based counter"},
		{"GLOBAL.gTimedBoolToogle", 500 * time.Millisecond, "time-based toggle"},
	}

	createdIDs := []int{}

	for _, v := range variables {
		settings := ads.SubscriptionSettings{
			CycleTime:    v.cycleTime,
			SendOnChange: true,
		}

		subscriptionsMutex.Lock()
		id := nextSubID
		nextSubID++
		subscriptionsMutex.Unlock()

		callback := createSubscriptionCallback(id, v.path)

		sub, err := client.SubscribeValue(port, v.path, callback, settings)
		if err != nil {
			fmt.Printf("[ERROR] Failed to subscribe to '%s' (%s): %v\n", v.path, v.name, err)
			continue
		}

		subscriptionsMutex.Lock()
		subscriptions[id] = sub
		subscriptionsMutex.Unlock()

		createdIDs = append(createdIDs, id)
		fmt.Printf("[OK] Subscription #%d created for %s '%s'\n", id, v.name, v.path)
	}

	if len(createdIDs) > 0 {
		fmt.Printf("[OK] Created %d subscriptions: IDs %v\n", len(createdIDs), createdIDs)
	} else {
		fmt.Println("[ERROR] Failed to create any subscriptions")
	}
}
