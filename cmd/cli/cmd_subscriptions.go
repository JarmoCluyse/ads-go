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

// createSubscriptionCallback creates a callback function for a subscription
func createSubscriptionCallback(id int, path string) ads.SubscriptionCallback {
	return func(data ads.SubscriptionData) {
		fmt.Printf("[NOTIFICATION #%d] %s: %v (at %s)\n",
			id, path, data.Value, data.Timestamp.Format("15:04:05.000"))
	}
}

// handleSubscribe subscribes to a variable and starts receiving notifications.
// Usage: subscribe [path]
// If path is not provided, uses a hardcoded test variable.
func handleSubscribe(args []string, client *ads.Client) {
	// Hardcoded test configuration
	var port uint16 = 350
	path := "Service_interface.Input.IN_MAIN_SERVICEINT_ENABLE"

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

		fmt.Printf("  #%d: %s (port %d) - CycleTime=%dms, OnChange=%s\n",
			id, path, sub.Port, sub.Settings.CycleTime.Milliseconds(), onChangeStr)
	}
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
	// Clear local map first
	subscriptionsMutex.Lock()
	count := len(subscriptions)
	subscriptions = make(map[int]*ads.ActiveSubscription)
	subscriptionsMutex.Unlock()

	// Unsubscribe from client
	err := client.UnsubscribeAll()
	if err != nil {
		fmt.Printf("[ERROR] Command 'unsubscribe_all': Failed to unsubscribe all: %v\n", err)
		return
	}

	fmt.Printf("[OK] Unsubscribed from all subscriptions (%d total)\n", count)
}
