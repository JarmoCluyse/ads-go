package ads

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/jarmoCluyse/ads-go/pkg/ads/constants"
	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
	"github.com/jarmoCluyse/ads-go/pkg/ads/utils"
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

	if err := c.setupPlcConnection(); err != nil {
		if !c.settings.AllowHalfOpen {
			c.conn.Close()
			c.logger.Error("Connect: Failed to setup PLC connection and AllowHalfOpen is false", "error", err)
			return fmt.Errorf("failed to setup PLC connection: %w", err)
		}
		c.logger.Warn("Connect: allowHalfOpen is active and PLC connection failed", "error", err)
	}
	c.logger.Debug("Connect: PLC connection setup complete (or half-open allowed).")

	return nil
}

// Disconnect closes the connection to the ADS router.
func (c *Client) Disconnect() error {
	c.logger.Debug("Disconnect: Attempting to disconnect.")
	if c.conn != nil {
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
	amsTcpHeader := createAmsTcpHeader(types.AMSTCPPortConnect, 2)
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
	amsTcpHeader := createAmsTcpHeader(types.AMSTCPPortClose, 2)
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
