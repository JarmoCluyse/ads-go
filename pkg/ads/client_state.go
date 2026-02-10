package ads

import (
	"fmt"
	"time"

	"github.com/jarmocluyse/ads-go/pkg/ads/ads-stateinfo"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

// GetCurrentState returns the cached TwinCAT system state.
// Returns nil if state has not been read yet or connection is not established.
func (c *Client) GetCurrentState() *adsstateinfo.SystemState {
	c.stateMutex.RLock()
	defer c.stateMutex.RUnlock()
	return c.currentState
}

// startStatePoller starts the background state polling timer.
// This is called automatically after successful connection if StatePollingInterval > 0.
func (c *Client) startStatePoller() {
	if c.settings.StatePollingInterval <= 0 {
		c.logger.Debug("startStatePoller: State polling disabled (interval <= 0)")
		return
	}

	c.statePollerMutex.Lock()
	defer c.statePollerMutex.Unlock()

	// Stop any existing timer
	if c.statePollerTimer != nil {
		c.statePollerTimer.Stop()
	}

	// Increment poller ID to invalidate any in-flight checks
	c.statePollerID++
	pollerID := c.statePollerID

	c.logger.Info("startStatePoller: Starting state monitoring", "interval", c.settings.StatePollingInterval, "pollerID", pollerID)

	// Start the timer
	c.statePollerTimer = time.AfterFunc(c.settings.StatePollingInterval, func() {
		c.checkState(pollerID)
	})
}

// stopStatePoller stops the background state polling timer.
// This is called automatically before disconnection.
func (c *Client) stopStatePoller() {
	c.statePollerMutex.Lock()
	defer c.statePollerMutex.Unlock()

	if c.statePollerTimer != nil {
		c.statePollerTimer.Stop()
		c.statePollerTimer = nil
		c.logger.Debug("stopStatePoller: State monitoring stopped")
	}

	// Increment poller ID to invalidate any in-flight checks
	c.statePollerID++
}

// checkState reads the current state, detects changes, and schedules the next check.
// This runs in the background as part of the state polling system.
func (c *Client) checkState(pollerID int) {
	c.logger.Debug("checkState: Checking TwinCAT system state", "pollerID", pollerID)

	// Check if this timer is still valid
	c.statePollerMutex.Lock()
	if pollerID != c.statePollerID {
		c.logger.Debug("checkState: Timer invalidated, skipping check", "pollerID", pollerID, "currentID", c.statePollerID)
		c.statePollerMutex.Unlock()
		return
	}
	c.statePollerMutex.Unlock()

	// Get the old state
	c.stateMutex.RLock()
	oldState := c.currentState
	c.stateMutex.RUnlock()

	// Read the current state
	newState, err := c.ReadTcSystemState()
	if err != nil {
		c.logger.Warn("checkState: Failed to read system state", "error", err)
		// Schedule next check even if this one failed
		c.scheduleNextStateCheck(pollerID)
		return
	}

	// Update the cached state
	c.stateMutex.Lock()
	c.currentState = newState
	c.stateMutex.Unlock()

	// Detect state changes
	if oldState == nil {
		// First state read
		c.logger.Info("checkState: Initial state read", "state", newState.AdsState.String())
		c.invokeStateChangeHook(newState, nil)
	} else if newState.AdsState != oldState.AdsState {
		// State changed
		c.logger.Info("checkState: TwinCAT state changed",
			"from", oldState.AdsState.String(),
			"to", newState.AdsState.String())

		c.invokeStateChangeHook(newState, oldState)

		// If state is not Run, trigger connection lost (which triggers reconnection)
		if newState.AdsState != types.ADSStateRun {
			c.logger.Warn("checkState: TwinCAT not in Run mode, triggering connection lost",
				"state", newState.AdsState.String())

			// Don't schedule next check - let reconnection handle it
			go c.invokeHook("OnConnectionLost", func() {
				c.settings.OnConnectionLost(c, fmt.Errorf("TwinCAT state changed to %s (not Run)", newState.AdsState.String()))
			})
			return
		}
	}

	// Schedule next check
	c.scheduleNextStateCheck(pollerID)
}

// scheduleNextStateCheck schedules the next state check if the poller is still valid.
func (c *Client) scheduleNextStateCheck(pollerID int) {
	c.statePollerMutex.Lock()
	defer c.statePollerMutex.Unlock()

	// Only schedule if this poller is still valid
	if pollerID != c.statePollerID {
		c.logger.Debug("scheduleNextStateCheck: Timer invalidated, not scheduling next check",
			"pollerID", pollerID, "currentID", c.statePollerID)
		return
	}

	c.logger.Debug("scheduleNextStateCheck: Scheduling next state check", "interval", c.settings.StatePollingInterval)

	c.statePollerTimer = time.AfterFunc(c.settings.StatePollingInterval, func() {
		c.checkState(pollerID)
	})
}

// invokeStateChangeHook safely calls the OnStateChange hook with panic recovery.
func (c *Client) invokeStateChangeHook(newState, oldState *adsstateinfo.SystemState) {
	if c.settings.OnStateChange == nil {
		return
	}

	go c.invokeHook("OnStateChange", func() {
		c.settings.OnStateChange(c, newState, oldState)
	})
}

// checkStateForOperation verifies that the system is in Run mode before performing read/write.
// Returns error if not in Run mode. Logs warning if state is unknown.
func (c *Client) checkStateForOperation(operationName string) error {
	c.stateMutex.RLock()
	state := c.currentState
	c.stateMutex.RUnlock()

	if state == nil {
		c.logger.Debug(operationName + ": System state unknown, allowing operation")
		return nil // Allow operation - state might not be cached yet
	}

	if state.AdsState != types.ADSStateRun {
		return fmt.Errorf("%s: Operation not allowed - TwinCAT is in %s mode (expected Run)",
			operationName, state.AdsState.String())
	}

	return nil
}
