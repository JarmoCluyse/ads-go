package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Client represents an ADS client.
type Client struct {
	conn         net.Conn
	settings     ClientSettings
	mutex        sync.Mutex
	invokeID     uint32
	requests     map[uint32]chan []byte
	notifications map[uint32]func([]byte)
	localAmsAddr AmsAddress
}

// ClientSettings holds the settings for the ADS client.
type ClientSettings struct {
	TargetNetID   string
	TargetPort    uint16
	RouterAddr    string
	Timeout       time.Duration
	AllowHalfOpen bool
}

// NewClient creates a new ADS client.
func NewClient(settings ClientSettings) *Client {
	if settings.Timeout == 0 {
		settings.Timeout = 2 * time.Second
	}
	return &Client{
		settings:     settings,
		requests:     make(map[uint32]chan []byte),
		notifications: make(map[uint32]func([]byte)),
	}
}

// Connect establishes a connection to the ADS router.
func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.settings.RouterAddr, c.settings.Timeout)
	if err != nil {
		return err
	}
	c.conn = conn

	if err := c.registerAdsPort(); err != nil {
		c.conn.Close()
		return err
	}

	go c.receive()

	if err := c.setupPlcConnection(); err != nil {
		if !c.settings.AllowHalfOpen {
			c.conn.Close()
			return fmt.Errorf("failed to setup PLC connection: %w", err)
		}
		fmt.Printf("WARNING: allowHalfOpen is active and PLC connection failed: %v\n", err)
	}

	return nil
}

func (c *Client) setupPlcConnection() error {
	// Read device info to check if we can communicate
	_, err := c.ReadDeviceInfo()
	if err != nil {
		return fmt.Errorf("failed to read device info: %w", err)
	}

	// Check if PLC is in RUN state
	state, err := c.ReadState()
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}

	if ADSState(state.AdsState) != ADSStateRun {
		return fmt.Errorf("PLC not in RUN mode (state: %d)", state.AdsState)
	}

	return nil
}

// Disconnect closes the connection to the ADS router.
func (c *Client) Disconnect() error {
	if c.conn != nil {
		defer c.conn.Close()
		return c.unregisterAdsPort()
	}
	return nil
}

func (c *Client) registerAdsPort() error {
	amsTcpHeader := createAmsTcpHeader(AMSTCPPortConnect, 2)
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, 0) // Let router decide port
	packet := append(amsTcpHeader, data...)

	if _, err := c.conn.Write(packet); err != nil {
		return err
	}

	respAmsTcpHeader := make([]byte, AMSTCPHeaderLength)
	if _, err := c.conn.Read(respAmsTcpHeader); err != nil {
		return err
	}

	length := binary.LittleEndian.Uint32(respAmsTcpHeader[2:6])
	respData := make([]byte, length)
	if _, err := c.conn.Read(respData); err != nil {
		return err
	}

	c.localAmsAddr.NetID = ByteArrayToAmsNetIDStr(respData[0:6])
	c.localAmsAddr.Port = binary.LittleEndian.Uint16(respData[6:8])

	return nil
}

func (c *Client) unregisterAdsPort() error {
	amsTcpHeader := createAmsTcpHeader(AMSTCPPortClose, 2)
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, c.localAmsAddr.Port)
	packet := append(amsTcpHeader, data...)

	_, err := c.conn.Write(packet)
	return err
}

// receive handles incoming data from the ADS router.
func (c *Client) receive() {
	defer c.conn.Close()
	for {
		amsTcpHeader := make([]byte, AMSTCPHeaderLength)
		if _, err := io.ReadFull(c.conn, amsTcpHeader); err != nil {
			// Handle error
			return
		}

		length := binary.LittleEndian.Uint32(amsTcpHeader[2:6])
		amsPacket := make([]byte, length)
		if _, err := io.ReadFull(c.conn, amsPacket); err != nil {
			// Handle error
			return
		}

		amsHeader := amsPacket[:AMSHeaderLength]
		invokeID := binary.LittleEndian.Uint32(amsHeader[28:32])

		c.mutex.Lock()
		ch, ok := c.requests[invokeID]
		c.mutex.Unlock()

		if ok {
			ch <- amsPacket[AMSHeaderLength:]
		} else {
			// Handle notification
		}
	}
}

// send sends a command to the ADS router.
func (c *Client) send(command ADSCommand, data []byte) ([]byte, error) {
	c.mutex.Lock()
	c.invokeID++
	invokeID := c.invokeID
	ch := make(chan []byte)
	c.requests[invokeID] = ch
	c.mutex.Unlock()

	defer func() {
		c.mutex.Lock()
		delete(c.requests, invokeID)
		c.mutex.Unlock()
	}()

	target := AmsAddress{NetID: c.settings.TargetNetID, Port: c.settings.TargetPort}

	amsHeader, err := createAmsHeader(target, c.localAmsAddr, command, uint32(len(data)), invokeID)
	if err != nil {
		return nil, err
	}
	amsTcpHeader := createAmsTcpHeader(AMSTCPPortAMSCommand, uint32(len(amsHeader)+len(data)))

	packet := append(amsTcpHeader, amsHeader...)
	packet = append(packet, data...)

	_, err = c.conn.Write(packet)
	if err != nil {
		return nil, err
	}

	select {
	case response := <-ch:
		return response, nil
	case <-time.After(c.settings.Timeout):
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

// ReadDeviceInfo reads the device information.
func (c *Client) ReadDeviceInfo() (*AdsReadDeviceInfoResponse, error) {
	data, err := c.send(ADSCommandReadDeviceInfo, []byte{})
	if err != nil {
		return nil, err
	}

	if len(data) < 4 {
		return nil, fmt.Errorf("invalid response length: %d", len(data))
	}

	errorCode := binary.LittleEndian.Uint32(data[0:4])
	if errorCode != 0 {
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	if len(data) < 24 {
		return nil, fmt.Errorf("invalid response length for device info: %d", len(data))
	}

	resp := &AdsReadDeviceInfoResponse{
		ErrorCode:    errorCode,
		MajorVersion: data[4],
		MinorVersion: data[5],
		VersionBuild: binary.LittleEndian.Uint16(data[6:8]),
		DeviceName:   string(bytes.Trim(data[8:24], "\x00")),
	}

	return resp, nil
}

// Read reads data from an ADS device.
func (c *Client) Read(indexGroup, indexOffset, length uint32) (*AdsReadResponse, error) {
	data := make([]byte, 12)
	binary.LittleEndian.PutUint32(data[0:4], indexGroup)
	binary.LittleEndian.PutUint32(data[4:8], indexOffset)
	binary.LittleEndian.PutUint32(data[8:12], length)

	respData, err := c.send(ADSCommandRead, data)
	if err != nil {
		return nil, err
	}

	if len(respData) < 8 {
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	errorCode := binary.LittleEndian.Uint32(respData[0:4])
	if errorCode != 0 {
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	readLength := binary.LittleEndian.Uint32(respData[4:8])

	resp := &AdsReadResponse{
		ErrorCode: errorCode,
		Length:    readLength,
		Data:      respData[8 : 8+readLength],
	}

	return resp, nil
}

// Write writes data to an ADS device.
func (c *Client) Write(indexGroup, indexOffset uint32, data []byte) (*AdsWriteResponse, error) {
	reqData := make([]byte, 12+len(data))
	binary.LittleEndian.PutUint32(reqData[0:4], indexGroup)
	binary.LittleEndian.PutUint32(reqData[4:8], indexOffset)
	binary.LittleEndian.PutUint32(reqData[8:12], uint32(len(data)))
	copy(reqData[12:], data)

	respData, err := c.send(ADSCommandWrite, reqData)
	if err != nil {
		return nil, err
	}

	if len(respData) < 4 {
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	resp := &AdsWriteResponse{
		ErrorCode: binary.LittleEndian.Uint32(respData[0:4]),
	}

	return resp, nil
}

// ReadWrite reads and writes data to an ADS device.
func (c *Client) ReadWrite(indexGroup, indexOffset, readLength uint32, dataToWrite []byte) (*AdsReadWriteResponse, error) {
	reqData := make([]byte, 16+len(dataToWrite))
	binary.LittleEndian.PutUint32(reqData[0:4], indexGroup)
	binary.LittleEndian.PutUint32(reqData[4:8], indexOffset)
	binary.LittleEndian.PutUint32(reqData[8:12], readLength)
	binary.LittleEndian.PutUint32(reqData[12:16], uint32(len(dataToWrite)))
	copy(reqData[16:], dataToWrite)

	respData, err := c.send(ADSCommandReadWrite, reqData)
	if err != nil {
		return nil, err
	}

	if len(respData) < 8 {
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	errorCode := binary.LittleEndian.Uint32(respData[0:4])
	if errorCode != 0 {
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	length := binary.LittleEndian.Uint32(respData[4:8])

	resp := &AdsReadWriteResponse{
		ErrorCode: errorCode,
		Length:    length,
		Data:      respData[8 : 8+length],
	}

	return resp, nil
}

// ReadState reads the ADS state of the device.
func (c *Client) ReadState() (*AdsReadStateResponse, error) {
	data, err := c.send(ADSCommandReadState, []byte{})
	if err != nil {
		return nil, err
	}

	if len(data) < 8 {
		return nil, fmt.Errorf("invalid response length: %d", len(data))
	}

	errorCode := binary.LittleEndian.Uint32(data[0:4])
	if errorCode != 0 {
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	resp := &AdsReadStateResponse{
		ErrorCode:   errorCode,
		AdsState:    binary.LittleEndian.Uint16(data[4:6]),
		DeviceState: binary.LittleEndian.Uint16(data[6:8]),
	}

	return resp, nil
}

// SetToConfig sets the TwinCAT system to config mode.
func (c *Client) SetToConfig() (*AdsWriteControlResponse, error) {
    return c.setSystemState(ADSStateConfig)
}

// SetToRun sets the TwinCAT system to run mode.
func (c *Client) SetToRun() (*AdsWriteControlResponse, error) {
    return c.setSystemState(ADSStateRun)
}

func (c *Client) setSystemState(state ADSState) (*AdsWriteControlResponse, error) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint16(data[0:2], uint16(state))
	binary.LittleEndian.PutUint16(data[2:4], 0) // DeviceState, not used

	respData, err := c.send(ADSCommandWriteControl, data)
	if err != nil {
		return nil, err
	}

	if len(respData) < 4 {
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	resp := &AdsWriteControlResponse{
		ErrorCode: binary.LittleEndian.Uint32(respData[0:4]),
	}

	return resp, nil
}