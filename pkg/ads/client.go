package ads

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/jarmoCluyse/ads-go/pkg/ads/types"
)

// Client represents an ADS client.
type Client struct {
	conn          net.Conn
	settings      ClientSettings
	mutex         sync.Mutex
	invokeID      uint32
	requests      map[uint32]chan []byte
	localAmsAddr  types.AmsAddress
	receiveBuffer bytes.Buffer // Buffer for incoming data
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
	log.Println("NewClient: Initializing new ADS client.")
	if settings.Timeout == 0 {
		settings.Timeout = 2 * time.Second
		log.Println("NewClient: Timeout not set, defaulting to 2 seconds.")
	}
	client := &Client{
		settings: settings,
		requests: make(map[uint32]chan []byte),
	}
	log.Println("NewClient: ADS client initialized.")
	return client
}

// Connect establishes a connection to the ADS router.
func (c *Client) Connect() error {
	log.Printf("Connect: Attempting to connect to router at %s", c.settings.RouterAddr)
	conn, err := net.DialTimeout("tcp", c.settings.RouterAddr, c.settings.Timeout)
	if err != nil {
		log.Printf("Connect: Failed to dial router: %v", err)
		return err
	}
	c.conn = conn
	log.Println("Connect: Successfully dialed router.")

	log.Println("Connect: Registering ADS port...")
	if err := c.registerAdsPort(); err != nil {
		c.conn.Close()
		log.Printf("Connect: Failed to register ADS port: %v", err)
		return err
	}
	log.Println("Connect: ADS port registered.")

	log.Println("Connect: Setting up PLC connection...")
	if err := c.setupPlcConnection(); err != nil {
		if !c.settings.AllowHalfOpen {
			c.conn.Close()
			log.Printf("Connect: Failed to setup PLC connection and AllowHalfOpen is false: %v", err)
			return fmt.Errorf("failed to setup PLC connection: %w", err)
		}
		log.Printf("Connect: WARNING: allowHalfOpen is active and PLC connection failed: %v\n", err)
		fmt.Printf("WARNING: allowHalfOpen is active and PLC connection failed: %v\n", err)
	}
	log.Println("Connect: PLC connection setup complete (or half-open allowed).")

	log.Println("Connect: Starting receive goroutine.")
	go c.receive()

	log.Println("Connect: Connection process finished.")
	return nil
}

func (c *Client) setupPlcConnection() error {
	log.Println("setupPlcConnection: Reading device info to check communication.")
	// Read device info to check if we can communicate
	_, err := c.ReadDeviceInfo()
	if err != nil {
		log.Printf("setupPlcConnection: Failed to read device info: %v", err)
		return fmt.Errorf("failed to read device info: %w", err)
	}
	log.Println("setupPlcConnection: Successfully read device info.")

	log.Println("setupPlcConnection: Checking if PLC is in RUN state.")
	// Check if PLC is in RUN state
	state, err := c.ReadState()
	if err != nil {
		log.Printf("setupPlcConnection: Failed to read state: %v", err)
		return fmt.Errorf("failed to read state: %w", err)
	}
	log.Printf("setupPlcConnection: Current PLC state: %d", state.AdsState)

	if types.ADSState(state.AdsState) != types.ADSStateRun {
		log.Printf("setupPlcConnection: PLC not in RUN mode (state: %d).", state.AdsState)
		return fmt.Errorf("PLC not in RUN mode (state: %d)", state.AdsState)
	}
	log.Println("setupPlcConnection: PLC is in RUN mode.")

	log.Println("setupPlcConnection: PLC connection setup successful.")
	return nil
}

// Disconnect closes the connection to the ADS router.
func (c *Client) Disconnect() error {
	log.Println("Disconnect: Attempting to disconnect.")
	if c.conn != nil {
		log.Println("Disconnect: Unregistering ADS port...")
		err := c.unregisterAdsPort()
		if err != nil {
			log.Printf("Disconnect: Error unregistering ADS port: %v", err)
		}

		log.Println("Disconnect: Closing connection.")
		defer c.conn.Close()
		log.Println("Disconnect: Connection closed.")
		return err
	}
	log.Println("Disconnect: No active connection to disconnect.")
	return nil
}

func (c *Client) registerAdsPort() error {
	log.Println("registerAdsPort: Creating AMS TCP header for port connection.")
	amsTcpHeader := createAmsTcpHeader(types.AMSTCPPortConnect, 2)
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, 0) // Let router decide port
	packet := append(amsTcpHeader, data...)

	log.Printf("registerAdsPort: Sending registration packet (length: %d).", len(packet))
	if _, err := c.conn.Write(packet); err != nil {
		log.Printf("registerAdsPort: Failed to write registration packet: %v", err)
		return err
	}
	log.Println("registerAdsPort: Registration packet sent. Waiting for response.")

	respAmsTcpHeader := make([]byte, AMSTCPHeaderLength)
	if _, err := c.conn.Read(respAmsTcpHeader); err != nil {
		log.Printf("registerAdsPort: Failed to read response AMS TCP header: %v", err)
		return err
	}
	log.Println("registerAdsPort: Received response AMS TCP header.")

	length := binary.LittleEndian.Uint32(respAmsTcpHeader[2:6])
	respData := make([]byte, length)
	log.Printf("registerAdsPort: Reading response data (length: %d).", length)
	if _, err := c.conn.Read(respData); err != nil {
		log.Printf("registerAdsPort: Failed to read response data: %v", err)
		return err
	}
	log.Println("registerAdsPort: Received response data.")

	c.localAmsAddr.NetID = ByteArrayToAmsNetIDStr(respData[0:6])
	c.localAmsAddr.Port = binary.LittleEndian.Uint16(respData[6:8])
	log.Printf("registerAdsPort: Local AMS Address set to NetID: %s, Port: %d", c.localAmsAddr.NetID, c.localAmsAddr.Port)

	log.Println("registerAdsPort: ADS port registration successful.")
	return nil
}

func (c *Client) unregisterAdsPort() error {
	log.Println("unregisterAdsPort: Creating AMS TCP header for port close.")
	amsTcpHeader := createAmsTcpHeader(types.AMSTCPPortClose, 2)
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, c.localAmsAddr.Port)
	packet := append(amsTcpHeader, data...)

	log.Printf("unregisterAdsPort: Sending unregistration packet (length: %d).", len(packet))
	_, err := c.conn.Write(packet)
	if err != nil {
		log.Printf("unregisterAdsPort: Failed to write unregistration packet: %v", err)
		return err
	}
	log.Println("unregisterAdsPort: Unregistration packet sent.")
	return nil
}

// receive handles incoming data from the ADS router.
func (c *Client) receive() {
	log.Println("receive: Starting receive goroutine.")
	defer func() {
		log.Println("receive: Closing connection from receive goroutine.")
		c.conn.Close()
		log.Println("receive: Receive goroutine terminated.")
	}()

	// Temporary buffer for reading from connection
	tempBuf := make([]byte, 4096) // Read in chunks

	for {
		n, err := c.conn.Read(tempBuf)
		if err != nil {
			if err == io.EOF {
				log.Println("receive: Connection closed by remote.")
			} else {
				log.Printf("receive: Error reading from connection: %v", err)
			}
			return // Exit goroutine on error or EOF
		}

		// Write read data to the receive buffer
		c.receiveBuffer.Write(tempBuf[:n])

		// Process packets from the receive buffer
		c.processReceiveBuffer()
	}
}

func (c *Client) processReceiveBuffer() {
	for {
		// Check if we have enough data for AMS/TCP header (6 bytes)
		if c.receiveBuffer.Len() < AMSTCPHeaderLength {
			return // Not enough data for header
		}

		// Read packet length from AMS/TCP header (bytes 2-5)
		// We need to peek without advancing the buffer's read pointer
		headerBytes := c.receiveBuffer.Bytes()[:AMSTCPHeaderLength]
		packetLength := binary.LittleEndian.Uint32(headerBytes[2:6])

		// Total length of the full packet (AMS/TCP header + AMS header + ADS data)
		totalPacketLength := AMSTCPHeaderLength + packetLength

		// Check if we have the full packet
		if c.receiveBuffer.Len() < int(totalPacketLength) {
			return // Not enough data for full packet
		}

		// Extract the full packet
		fullPacket := make([]byte, totalPacketLength)
		c.receiveBuffer.Read(fullPacket)

		// Now process the full packet
		amsPacket := fullPacket[AMSTCPHeaderLength:]
		amsHeader := amsPacket[:AMSHeaderLength]
		invokeID := binary.LittleEndian.Uint32(amsHeader[28:32])

		log.Printf("receive: Received packet with InvokeID: %d", invokeID)

		c.mutex.Lock()
		ch, ok := c.requests[invokeID]
		c.mutex.Unlock()

		if ok {
			log.Printf("receive: Found channel for InvokeID %d, sending response.", invokeID)
			ch <- amsPacket[AMSHeaderLength:]
		} else {
			log.Printf("receive: No channel found for InvokeID %d, discarding packet.", invokeID)
		}
	}
}

// send sends a command to the ADS router.
func (c *Client) send(req types.AdsCommandRequest) ([]byte, error) {
	log.Printf("send: Preparing to send command %s (data length: %d).", req.Command.String(), len(req.Data))
	c.mutex.Lock()
	c.invokeID++
	invokeID := c.invokeID
	ch := make(chan []byte)
	c.requests[invokeID] = ch
	c.mutex.Unlock()
	log.Printf("send: Assigned InvokeID: %d", invokeID)

	defer func() {
		c.mutex.Lock()
		delete(c.requests, invokeID)
		c.mutex.Unlock()
		log.Printf("send: Cleaned up request for InvokeID: %d", invokeID)
	}()

	target := types.AmsAddress{NetID: req.TargetNetID, Port: req.TargetPort}
	log.Printf("send: Target AMS Address: %s:%d", target.NetID, target.Port)

	amsHeader, err := createAmsHeader(target, c.localAmsAddr, req.Command, uint32(len(req.Data)), invokeID)
	if err != nil {
		log.Printf("send: Failed to create AMS header: %v", err)
		return nil, err
	}
	amsTcpHeader := createAmsTcpHeader(types.AMSTCPPortAMSCommand, uint32(len(amsHeader)+len(req.Data)))

	packet := append(amsTcpHeader, amsHeader...)
	packet = append(packet, req.Data...)
	log.Printf("send: Constructed packet (total length: %d). Raw packet: %x", len(packet), packet)

	_, err = c.conn.Write(packet)
	if err != nil {
		log.Printf("send: Failed to write packet to connection: %v", err)
		return nil, err
	}
	log.Println("send: Packet sent. Waiting for response or timeout.")

	select {
	case response := <-ch:
		log.Printf("send: Received response for InvokeID: %d", invokeID)
		return response, nil
	case <-time.After(c.settings.Timeout):
		log.Printf("send: Timeout waiting for response for InvokeID: %d", invokeID)
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

// ReadDeviceInfo reads the device information.
func (c *Client) ReadDeviceInfo() (*types.AdsReadDeviceInfoResponse, error) {
	log.Println("ReadDeviceInfo: Sending ReadDeviceInfo command.")
	req := types.AdsCommandRequest{
		Command:     types.ADSCommandReadDeviceInfo,
		Data:        []byte{},
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}
	data, err := c.send(req)
	if err != nil {
		log.Printf("ReadDeviceInfo: Failed to send ReadDeviceInfo command: %v", err)
		return nil, err
	}
	log.Printf("ReadDeviceInfo: Received raw response data (length: %d). Data: %x", len(data), data)

	if len(data) < 4 {
		log.Printf("ReadDeviceInfo: Invalid response length: %d, expected at least 4.", len(data))
		return nil, fmt.Errorf("invalid response length: %d", len(data))
	}

	errorCode := binary.LittleEndian.Uint32(data[0:4])
	if errorCode != 0 {
		log.Printf("ReadDeviceInfo: ADS error received: 0x%x", errorCode)
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	// Handle cases where only the error code is returned (e.g., older/embedded runtimes)
	if len(data) == 4 && errorCode == 0 {
		log.Println("ReadDeviceInfo: Received only error code, returning default device info.")
		return &types.AdsReadDeviceInfoResponse{
			BaseAdsResponse: types.BaseAdsResponse{ErrorCode: 0},
			MajorVersion:    0,
			MinorVersion:    0,
			VersionBuild:    0,
			DeviceName:      "",
		}, nil
	}

	if len(data) < 24 {
		log.Printf("ReadDeviceInfo: Invalid response length for device info: %d, expected at least 24.", len(data))
		return nil, fmt.Errorf("invalid response length for device info: %d", len(data))
	}

	resp := &types.AdsReadDeviceInfoResponse{
		BaseAdsResponse: types.BaseAdsResponse{ErrorCode: errorCode},
		MajorVersion:    data[4],
		MinorVersion:    data[5],
		VersionBuild:    binary.LittleEndian.Uint16(data[6:8]),
		DeviceName:      string(bytes.Trim(data[8:24], "\x00")),
	}
	log.Printf("ReadDeviceInfo: Successfully parsed device info: %+v", resp)

	return resp, nil
}

// Read reads data from an ADS device.
func (c *Client) Read(indexGroup, indexOffset, length uint32) (*types.AdsReadResponse, error) {
	log.Printf("Read: Reading data from indexGroup 0x%x, indexOffset 0x%x, length %d.", indexGroup, indexOffset, length)
	data := make([]byte, 12)
	binary.LittleEndian.PutUint32(data[0:4], indexGroup)
	binary.LittleEndian.PutUint32(data[4:8], indexOffset)
	binary.LittleEndian.PutUint32(data[8:12], length)

	req := types.AdsCommandRequest{
		Command:     types.ADSCommandRead,
		Data:        data,
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}
	respData, err := c.send(req)
	if err != nil {
		log.Printf("Read: Failed to send command: %v", err)
		return nil, err
	}
	log.Printf("Read: Received raw response data (length: %d).", len(respData))

	if len(respData) < 8 {
		log.Printf("Read: Invalid response length: %d, expected at least 8.", len(respData))
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	errorCode := binary.LittleEndian.Uint32(respData[0:4])
	if errorCode != 0 {
		log.Printf("Read: ADS error received: 0x%x", errorCode)
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	readLength := binary.LittleEndian.Uint32(respData[4:8])
	log.Printf("Read: Reported read length: %d.", readLength)

	resp := &types.AdsReadResponse{
		BaseAdsResponse: types.BaseAdsResponse{ErrorCode: errorCode},
		Length:          readLength,
		Data:            respData[8 : 8+readLength],
	}
	log.Printf("Read: Successfully parsed read response (data length: %d).", len(resp.Data))

	return resp, nil
}

// Write writes data to an ADS device.
func (c *Client) Write(indexGroup, indexOffset uint32, data []byte) (*types.AdsWriteResponse, error) {
	log.Printf("Write: Writing data to indexGroup 0x%x, indexOffset 0x%x (data length: %d).", indexGroup, indexOffset, len(data))
	reqData := make([]byte, 12+len(data))
	binary.LittleEndian.PutUint32(reqData[0:4], indexGroup)
	binary.LittleEndian.PutUint32(reqData[4:8], indexOffset)
	binary.LittleEndian.PutUint32(reqData[8:12], uint32(len(data)))
	copy(reqData[12:], data)

	req := types.AdsCommandRequest{
		Command:     types.ADSCommandWrite,
		Data:        reqData,
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}
	respData, err := c.send(req)
	if err != nil {
		log.Printf("Write: Failed to send command: %v", err)
		return nil, err
	}
	log.Printf("Write: Received raw response data (length: %d).", len(respData))

	if len(respData) < 4 {
		log.Printf("Write: Invalid response length: %d, expected at least 4.", len(respData))
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	resp := &types.AdsWriteResponse{
		BaseAdsResponse: types.BaseAdsResponse{ErrorCode: binary.LittleEndian.Uint32(respData[0:4])},
	}
	log.Printf("Write: Successfully parsed write response (ErrorCode: 0x%x).", resp.ErrorCode)

	return resp, nil
}

// ReadWrite reads and writes data to an ADS device.
func (c *Client) ReadWrite(indexGroup, indexOffset, readLength uint32, dataToWrite []byte) (*types.AdsReadWriteResponse, error) {
	log.Printf("ReadWrite: Reading from indexGroup 0x%x, indexOffset 0x%x (readLength: %d), writing data (length: %d).", indexGroup, indexOffset, readLength, len(dataToWrite))
	reqData := make([]byte, 16+len(dataToWrite))
	binary.LittleEndian.PutUint32(reqData[0:4], indexGroup)
	binary.LittleEndian.PutUint32(reqData[4:8], indexOffset)
	binary.LittleEndian.PutUint32(reqData[8:12], readLength)
	binary.LittleEndian.PutUint32(reqData[12:16], uint32(len(dataToWrite)))
	copy(reqData[16:], dataToWrite)

	req := types.AdsCommandRequest{
		Command:     types.ADSCommandReadWrite,
		Data:        reqData,
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}
	respData, err := c.send(req)
	if err != nil {
		log.Printf("ReadWrite: Failed to send command: %v", err)
		return nil, err
	}
	log.Printf("ReadWrite: Received raw response data (length: %d).", len(respData))

	if len(respData) < 8 {
		log.Printf("ReadWrite: Invalid response length: %d, expected at least 8.", len(respData))
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	errorCode := binary.LittleEndian.Uint32(respData[0:4])
	if errorCode != 0 {
		log.Printf("ReadWrite: ADS error received: 0x%x", errorCode)
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	length := binary.LittleEndian.Uint32(respData[4:8])
	log.Printf("ReadWrite: Reported response length: %d.", length)

	resp := &types.AdsReadWriteResponse{
		BaseAdsResponse: types.BaseAdsResponse{ErrorCode: errorCode},
		Length:          length,
		Data:            respData[8 : 8+length],
	}
	log.Printf("ReadWrite: Successfully parsed read/write response (data length: %d).", len(resp.Data))

	return resp, nil
}

// ReadState sends an AdsCommandReadState (command ID = 4) and returns the PLC’s ADS
// and device states.  Older/embedded runtimes may reply with **only the 4-byte
// result code** (0 = OK).  We accept that short frame and return zeros for the
// state fields so the behaviour matches the Bun implementation.
func (c *Client) ReadState() (*AdsReadStateResponse, error) {
	log.Println("ReadState: Reading ADS state.")
	req := types.AdsCommandRequest{
		Command:     types.ADSCommandReadState,
		Data:        []byte{},
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  c.settings.TargetPort,
	}
	data, err := c.send(req)
	if err != nil {
		log.Printf("ReadState: Failed to send ReadState command: %v", err)
		return nil, err
	}
	log.Printf("ReadState: Received raw response data (length: %d). Data: %x", len(data), data)

	if len(data) < 4 {
		log.Printf("ReadState: Invalid response length: %d, expected at least 4.", len(data))
		return nil, fmt.Errorf("invalid response length: %d", len(data))
	}

	/* ────────────────────────────────
	   Legacy runtimes: reply is only
	   the 4-byte result code (OK = 0)
	   ──────────────────────────────── */
	if len(data) == 4 && binary.LittleEndian.Uint32(data[:4]) == 0 {
		return &AdsReadStateResponse{
			ErrorCode:   0,
			AdsState:    0, // unknown
			DeviceState: 0, // unknown
		}, nil
	}

	/* ────────────────────────────────
	   Normal case – need 8 bytes total
	   ──────────────────────────────── */
	if len(data) < 8 {
		return nil, fmt.Errorf("state reply length %d (want 4 or ≥8)", len(data))
	}

	errorCode := binary.LittleEndian.Uint32(data[0:4])
	if errorCode != 0 {
		log.Printf("ReadState: ADS error received: 0x%x", errorCode)
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	var adsState, deviceState uint16
	// The actual ADS state and device state are after the 4-byte error code
	if len(data) >= 8 { // ErrorCode (4 bytes) + ADS State (2 bytes) + Device State (2 bytes)
		adsState = binary.LittleEndian.Uint16(data[4:6])
		deviceState = binary.LittleEndian.Uint16(data[6:8])
		log.Printf("ReadState: Parsed ADS State: %d, Device State: %d.", adsState, deviceState)
	} else {
		log.Printf("ReadState: Not enough data in response to parse ADS and Device State (length: %d).", len(data))
	}

	respState := &AdsReadStateResponse{
		ErrorCode:   errorCode,
		AdsState:    types.ADSState(adsState),
		DeviceState: deviceState,
	}
	log.Printf("ReadState: Successfully parsed ADS state response: %+v", respState)

	return respState, nil
}

// ReadTcSystemState reads the TwinCAT system state.
func (c *Client) ReadTcSystemState() (*types.AdsTcSystemStateResponse, error) {
	log.Println("ReadTcSystemState: Reading TwinCAT system state.")
	req := types.AdsCommandRequest{
		Command:     types.ADSCommandReadState,
		Data:        []byte{},
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  types.ADSReservedPortSystemService, // Explicitly target SystemService port
	}
	data, err := c.send(req)
	if err != nil {
		log.Printf("ReadTcSystemState: Failed to send ReadState command: %v", err)
		return nil, err
	}
	log.Printf("ReadTcSystemState: Received raw response data (length: %d). Data: %x", len(data), data)

	if len(data) < 4 {
		log.Printf("ReadTcSystemState: Invalid response length: %d, expected at least 4.", len(data))
		return nil, fmt.Errorf("invalid response length: %d", len(data))
	}

	errorCode := binary.LittleEndian.Uint32(data[0:4])
	if errorCode != 0 {
		log.Printf("ReadTcSystemState: ADS error received: 0x%x", errorCode)
		return nil, fmt.Errorf("ADS error: 0x%x", errorCode)
	}

	var adsState, deviceState uint16
	// The actual ADS state and device state are after the 4-byte error code
	if len(data) >= 8 { // ErrorCode (4 bytes) + ADS State (2 bytes) + Device State (2 bytes)
		adsState = binary.LittleEndian.Uint16(data[4:6])
		deviceState = binary.LittleEndian.Uint16(data[6:8])
		log.Printf("ReadTcSystemState: Parsed ADS State: %d, Device State: %d.", adsState, deviceState)
	} else {
		log.Printf("ReadTcSystemState: Not enough data in response to parse ADS and Device State (length: %d).", len(data))
	}

	respState := &types.AdsTcSystemStateResponse{
		BaseAdsResponse: types.BaseAdsResponse{ErrorCode: errorCode},
		AdsState:        types.ADSState(adsState),
		DeviceState:     deviceState,
	}
	log.Printf("ReadTcSystemState: Successfully parsed TwinCAT system state response: %+v", respState)

	return respState, nil
}

// WriteControl writes control data to an ADS device.
// It can optionally target a specific ADS port. If targetPort is 0, c.settings.TargetPort is used.
func (c *Client) WriteControl(adsState types.ADSState, deviceState uint16, targetPort ...uint16) (*types.AdsWriteControlResponse, error) {
	port := c.settings.TargetPort
	if len(targetPort) > 0 && targetPort[0] != 0 {
		port = targetPort[0]
	}

	log.Printf("WriteControl: Setting ADS state to %s (0x%x), device state to 0x%x, target port: %d.", adsState.String(), adsState, deviceState, port)
	data := make([]byte, 8)
	binary.LittleEndian.PutUint16(data[0:2], uint16(adsState))
	binary.LittleEndian.PutUint16(data[2:4], deviceState)
	binary.LittleEndian.PutUint32(data[4:8], 0) // DataLength, 0 for state change

	req := types.AdsCommandRequest{
		Command:     types.ADSCommandWriteControl,
		Data:        data,
		TargetNetID: c.settings.TargetNetID,
		TargetPort:  port,
	}
	respData, err := c.send(req)
	if err != nil {
		log.Printf("WriteControl: Failed to send WriteControl command: %v", err)
		return nil, err
	}
	log.Printf("WriteControl: Received raw response data (length: %d).", len(respData))

	if len(respData) < 4 {
		log.Printf("WriteControl: Invalid response length: %d, expected at least 4.", len(respData))
		return nil, fmt.Errorf("invalid response length: %d", len(respData))
	}

	resp := &types.AdsWriteControlResponse{
		BaseAdsResponse: types.BaseAdsResponse{ErrorCode: binary.LittleEndian.Uint32(respData[0:4])},
	}
	log.Printf("WriteControl: Successfully parsed write control response (ErrorCode: 0x%x).", resp.ErrorCode)

	return resp, nil
}

// ReadPlcRuntimeState reads the PLC runtime state.
func (c *Client) ReadPlcRuntimeState() (*types.AdsStateResponse, error) {
	log.Println("ReadPlcRuntimeState: Reading PLC runtime state.")
	resp, err := c.ReadTcSystemState()
	if err != nil {
		log.Printf("ReadPlcRuntimeState: Failed to read state: %v", err)
		return nil, err
	}
	log.Printf("ReadPlcRuntimeState: Successfully read PLC runtime state: %+v", resp.AdsState)
	return &types.AdsStateResponse{ADSState: types.ADSState(resp.AdsState), DeviceState: resp.DeviceState}, nil
}

// SetTcSystemToConfig sets the TwinCAT system to config mode.
func (c *Client) SetTcSystemToConfig() error {
	log.Println("SetTcSystemToConfig: Setting TwinCAT system to config mode.")
	// Read current state to ensure it's not already in config mode or other unexpected state
	currentState, err := c.ReadPlcRuntimeState()
	if err != nil {
		log.Printf("SetTcSystemToConfig: Failed to read current PLC runtime state: %v", err)
		return fmt.Errorf("failed to read current PLC runtime state: %w", err)
	}

	if currentState.ADSState == types.ADSStateConfig {
		log.Println("SetTcSystemToConfig: TwinCAT system is already in config mode.")
		return nil
	}

	// Set ADS state to Config (15) and device state to 2 (reboot)
	resp, err := c.WriteControl(types.ADSStateReconfig, 2, types.ADSReservedPortSystemService)
	if err != nil {
		log.Printf("SetTcSystemToConfig: Failed to send WriteControl command: %v", err)
		return err
	}

	if resp.ErrorCode != 0 {
		log.Printf("SetTcSystemToConfig: WriteControl command returned ADS error: 0x%x", resp.ErrorCode)
		return fmt.Errorf("WriteControl command returned ADS error: 0x%x", resp.ErrorCode)
	}

	log.Println("SetTcSystemToConfig: TwinCAT system successfully set to config mode.")
	return nil
}

// SetTcSystemToRun sets the TwinCAT system to run mode.
func (c *Client) SetTcSystemToRun(reconnect bool) error {
	log.Printf("SetTcSystemToRun: Setting TwinCAT system to run mode (reconnect: %t).", reconnect)
	if c.conn == nil {
		return fmt.Errorf("SetTcSystemToRun: Client is not connected. Use Connect() to connect to the target first.")
	}

	// Reading device state first as we don't want to change it (even though it's most probably 0)
	state, err := c.ReadPlcRuntimeState()
	if err != nil {
		log.Printf("SetTcSystemToRun: Failed to read current PLC runtime state: %v", err)
		return fmt.Errorf("failed to read current PLC runtime state: %w", err)
	}

	resp, err := c.WriteControl(types.ADSStateReset, state.DeviceState, types.ADSReservedPortSystemService)
	if err != nil {
		log.Printf("SetTcSystemToRun: Failed to send WriteControl command: %v", err)
		return err
	}

	if resp.ErrorCode != 0 {
		log.Printf("SetTcSystemToRun: WriteControl command returned ADS error: 0x%x", resp.ErrorCode)
		return fmt.Errorf("WriteControl command returned ADS error: 0x%x", resp.ErrorCode)
	}

	log.Println("SetTcSystemToRun: TwinCAT system successfully set to run mode.")

	if reconnect {
		log.Println("SetTcSystemToRun: Reconnecting after TwinCAT system restart.")
		if err := c.Disconnect(); err != nil {
			log.Printf("SetTcSystemToRun: Error during disconnect before reconnect: %v", err)
			// Continue with connect even if disconnect fails, as connection might already be broken
		}
		if err := c.Connect(); err != nil {
			log.Printf("SetTcSystemToRun: Error during reconnect: %v", err)
			return fmt.Errorf("reconnect failed after setting system to run: %w", err)
		}
		log.Println("SetTcSystemToRun: Reconnected successfully.")
	}

	return nil
}
