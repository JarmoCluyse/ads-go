package ads

import (
	"encoding/binary"
	"fmt"
	"net"

	amsbuilder "github.com/jarmocluyse/ads-go/pkg/ads/ams-builder"
	"github.com/jarmocluyse/ads-go/pkg/ads/constants"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/jarmocluyse/ads-go/pkg/ads/utils"
)

// Connect establishes a connection to the ADS router.
func (c *Client) Connect() error {
	c.logger.Debug("Connect: Attempting to connect to router", "routerAddr", c.settings.RouterAddr)
	conn, err := net.DialTimeout("tcp", c.settings.RouterAddr, c.settings.Timeout)
	if err != nil {
		c.logger.Error("Connect: Failed to dial router", "error", err)
		return err
	}
	c.conn = conn

	if err := c.registerAdsPort(); err != nil {
		c.conn.Close()
		c.logger.Error("Connect: Failed to register ADS port", "error", err)
		return err
	}
	c.logger.Debug("Connect: ADS port registered.")

	// Start receiving
	go c.receive()

	// TODO: check if this is needed
	if err := c.setupPlcConnection(); err != nil {
		c.logger.Warn("Connect: PLC setup not complete", "error", err)
	}

	c.logger.Info("Connect: Successfully connected to ADS router", "localAMS", c.localAmsAddr.NetID, "port", c.localAmsAddr.Port)

	// Invoke OnConnect hook (synchronous)
	if err := c.invokeConnectHook(c.localAmsAddr); err != nil {
		c.logger.Error("Connect: OnConnect hook failed, disconnecting", "error", err)
		_ = c.Disconnect() // Clean up connection
		return fmt.Errorf("connection hook failed: %w", err)
	}

	// Read initial state and start state monitoring
	if initialState, err := c.ReadTcSystemState(); err != nil {
		c.logger.Warn("Connect: Failed to read initial TwinCAT state", "error", err)
	} else {
		c.stateMutex.Lock()
		c.currentState = initialState
		c.stateMutex.Unlock()
		c.logger.Info("Connect: Initial TwinCAT state", "state", initialState.AdsState.String())
	}

	// Start state poller if enabled
	if c.settings.StatePollingInterval > 0 {
		c.startStatePoller()
	}

	return nil
}

// Disconnect closes the connection to the ADS router.
func (c *Client) Disconnect() error {
	c.logger.Debug("Disconnect: Attempting to disconnect.")
	if c.conn != nil {
		// Stop state monitoring
		c.stopStatePoller()

		// Clear cached state
		c.stateMutex.Lock()
		c.currentState = nil
		c.stateMutex.Unlock()

		// Unsubscribe from all active subscriptions before disconnecting
		if err := c.UnsubscribeAll(); err != nil {
			c.logger.Warn("Disconnect: Error unsubscribing from all subscriptions", "error", err)
		}
		c.logger.Info("Disconnect: Unsubscribed from all active subscriptions.")

		// Invoke OnDisconnect hook asynchronously (fire-and-forget)
		go c.invokeHook("OnDisconnect", func() {
			c.settings.OnDisconnect(c)
		})

		err := c.unregisterAdsPort()
		if err != nil {
			c.logger.Error("Disconnect: Error unregistering ADS port", "error", err)
		}

		defer c.conn.Close()
		c.logger.Info("Disconnect: Connection closed.")
		return err
	}
	c.logger.Warn("Disconnect: No active connection to disconnect.")
	return nil
}

// Connect to the PLC
func (c *Client) setupPlcConnection() error {
	c.logger.Debug("setupPlcConnection: Reading device info to check communication.")
	// Read device info to check if we can communicate
	_, err := c.ReadDeviceInfo()
	if err != nil {
		c.logger.Error("setupPlcConnection: Failed to read device info", "error", err)
		return fmt.Errorf("failed to read device info: %w", err)
	}

	// Check if PLC is in RUN state
	state, err := c.ReadTcSystemState()
	if err != nil {
		c.logger.Error("setupPlcConnection: Failed to read state", "error", err)
		return fmt.Errorf("failed to read state: %w", err)
	}
	c.logger.Info("setupPlcConnection: Current PLC state", "state", state.AdsState)

	if types.ADSState(state.AdsState) != types.ADSStateRun {
		c.logger.Warn("setupPlcConnection: PLC not in RUN mode", "state", state.AdsState)
		return fmt.Errorf("PLC not in RUN mode (state: %d)", state.AdsState)
	}

	c.logger.Debug("setupPlcConnection: PLC is in RUN mode.")
	return nil
}

// registerAdsPort
func (c *Client) registerAdsPort() error {
	c.logger.Debug("registerAdsPort: Creating AMS TCP header for port connection.")
	amsTcpHeader := amsbuilder.BuildAmsTcpHeader(types.AMSTCPPortConnect, 2)
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, 0) // Let router decide port
	packet := append(amsTcpHeader, data...)

	c.logger.Debug("registerAdsPort: Sending registration packet", "length", len(packet), "packet", packet)
	if _, err := c.conn.Write(packet); err != nil {
		c.logger.Error("registerAdsPort: Failed to write registration packet", "error", err)
		return err
	}

	respAmsTcpHeader := make([]byte, constants.AMSTCPHeaderLength)
	if _, err := c.conn.Read(respAmsTcpHeader); err != nil {
		c.logger.Error("registerAdsPort: Failed to read response AMS TCP header", "error", err)
		return err
	}

	c.logger.Debug("registerAdsPort: respAmsTcpHeader", "length", len(respAmsTcpHeader), "packet", respAmsTcpHeader)

	length := binary.LittleEndian.Uint32(respAmsTcpHeader[2:6])
	respData := make([]byte, length)
	if _, err := c.conn.Read(respData); err != nil {
		c.logger.Error("registerAdsPort: Failed to read response data", "error", err)
		return err
	}

	c.logger.Debug("registerAdsPort: respData", "length", len(respData), "packet", respData)
	c.localAmsAddr.NetID = utils.ByteArrayToAmsNetIdStr(respData[0:6])
	c.localAmsAddr.Port = binary.LittleEndian.Uint16(respData[6:8])

	c.logger.Debug("registerAdsPort: Local AMS Address set", "netID", c.localAmsAddr.NetID, "port", c.localAmsAddr.Port)
	c.logger.Info("registerAdsPort: ADS port registration successful.")
	return nil
}

// unregisterAdsPort
func (c *Client) unregisterAdsPort() error {
	c.logger.Debug("unregisterAdsPort: Creating AMS TCP header for port close.")
	amsTcpHeader := amsbuilder.BuildAmsTcpHeader(types.AMSTCPPortClose, 2)
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, c.localAmsAddr.Port)
	packet := append(amsTcpHeader, data...)

	_, err := c.conn.Write(packet)
	if err != nil {
		c.logger.Error("unregisterAdsPort: Failed to write unregistration packet", "error", err)
		return err
	}
	c.logger.Info("unregisterAdsPort: Unregistration packet sent.")
	return nil
}
