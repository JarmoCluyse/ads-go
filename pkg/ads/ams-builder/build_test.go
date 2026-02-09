package amsbuilder

import (
	"encoding/binary"
	"errors"
	"testing"

	"github.com/jarmocluyse/ads-go/pkg/ads/constants"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/jarmocluyse/ads-go/pkg/ads/utils"
	"github.com/stretchr/testify/assert"
)

// TestBuildAmsTcpHeader tests the BuildAmsTcpHeader function.
func TestBuildAmsTcpHeader(t *testing.T) {
	tests := []struct {
		name       string
		command    types.AMSHeaderFlag
		dataLength uint32
	}{
		{
			name:       "Standard ADS command",
			command:    types.AMSTCPPortAMSCommand,
			dataLength: 44,
		},
		{
			name:       "Port close",
			command:    types.AMSTCPPortClose,
			dataLength: 0,
		},
		{
			name:       "Port connect",
			command:    types.AMSTCPPortConnect,
			dataLength: 8,
		},
		{
			name:       "Router note",
			command:    types.AMSTCPPortRouterNote,
			dataLength: 32,
		},
		{
			name:       "Get local NetID",
			command:    types.GetLocalNetID,
			dataLength: 0,
		},
		{
			name:       "Zero data length",
			command:    types.AMSTCPPortAMSCommand,
			dataLength: 0,
		},
		{
			name:       "Large data length",
			command:    types.AMSTCPPortAMSCommand,
			dataLength: 1000000,
		},
		{
			name:       "Max uint32 data length",
			command:    types.AMSTCPPortAMSCommand,
			dataLength: 0xFFFFFFFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildAmsTcpHeader(tt.command, tt.dataLength)

			// Verify length
			if len(result) != constants.AMSTCPHeaderLength {
				t.Fatalf("Expected length %d, got %d", constants.AMSTCPHeaderLength, len(result))
			}

			// Verify command
			command := binary.LittleEndian.Uint16(result[0:2])
			assert.Equal(t, uint16(tt.command), command)

			// Verify data length
			dataLength := binary.LittleEndian.Uint32(result[2:6])
			assert.Equal(t, tt.dataLength, dataLength)
		})
	}
}

// TestBuildAmsHeader tests the BuildAmsHeader function.
func TestBuildAmsHeader(t *testing.T) {
	tests := []struct {
		name       string
		target     AmsAddress
		source     AmsAddress
		command    types.ADSCommand
		dataLength uint32
		invokeID   uint32
	}{
		{
			name:       "Read command",
			target:     AmsAddress{NetID: "192.168.1.100.1.1", Port: 851},
			source:     AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905},
			command:    types.ADSCommandRead,
			dataLength: 12,
			invokeID:   1,
		},
		{
			name:       "Write command",
			target:     AmsAddress{NetID: "10.0.0.1.1.1", Port: 801},
			source:     AmsAddress{NetID: "10.0.0.2.1.1", Port: 48898},
			command:    types.ADSCommandWrite,
			dataLength: 100,
			invokeID:   42,
		},
		{
			name:       "ReadWrite command",
			target:     AmsAddress{NetID: "172.16.0.10.1.1", Port: 851},
			source:     AmsAddress{NetID: "172.16.0.20.1.1", Port: 40000},
			command:    types.ADSCommandReadWrite,
			dataLength: 256,
			invokeID:   999,
		},
		{
			name:       "Write control",
			target:     AmsAddress{NetID: "5.6.7.8.1.1", Port: 851},
			source:     AmsAddress{NetID: "1.2.3.4.1.1", Port: 50000},
			command:    types.ADSCommandWriteControl,
			dataLength: 8,
			invokeID:   123,
		},
		{
			name:       "Read state",
			target:     AmsAddress{NetID: "255.255.255.255.1.1", Port: 851},
			source:     AmsAddress{NetID: "0.0.0.0.1.1", Port: 1},
			command:    types.ADSCommandReadState,
			dataLength: 0,
			invokeID:   0,
		},
		{
			name:       "Read device info",
			target:     AmsAddress{NetID: "127.0.0.1.1.1", Port: 851},
			source:     AmsAddress{NetID: "127.0.0.2.1.1", Port: 65535},
			command:    types.ADSCommandReadDeviceInfo,
			dataLength: 0,
			invokeID:   0xFFFFFFFF,
		},
		{
			name:       "Zero ports",
			target:     AmsAddress{NetID: "1.1.1.1.1.1", Port: 0},
			source:     AmsAddress{NetID: "2.2.2.2.2.2", Port: 0},
			command:    types.ADSCommandRead,
			dataLength: 50,
			invokeID:   500,
		},
		{
			name:       "Large data length",
			target:     AmsAddress{NetID: "192.168.1.1.1.1", Port: 851},
			source:     AmsAddress{NetID: "192.168.1.2.1.1", Port: 32905},
			command:    types.ADSCommandRead,
			dataLength: 0xFFFFFFFF,
			invokeID:   0xFFFFFFFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildAmsHeader(tt.target, tt.source, tt.command, tt.dataLength, tt.invokeID)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// Verify length
			if len(result) != constants.AMSHeaderLength {
				t.Fatalf("Expected length %d, got %d", constants.AMSHeaderLength, len(result))
			}

			// Verify target NetID and port
			targetNetID := utils.ByteArrayToAmsNetIdStr(result[0:6])
			assert.Equal(t, tt.target.NetID, targetNetID)
			targetPort := binary.LittleEndian.Uint16(result[6:8])
			assert.Equal(t, tt.target.Port, targetPort)

			// Verify source NetID and port
			sourceNetID := utils.ByteArrayToAmsNetIdStr(result[8:14])
			assert.Equal(t, tt.source.NetID, sourceNetID)
			sourcePort := binary.LittleEndian.Uint16(result[14:16])
			assert.Equal(t, tt.source.Port, sourcePort)

			// Verify command
			command := binary.LittleEndian.Uint16(result[16:18])
			assert.Equal(t, uint16(tt.command), command)

			// Verify state flags (always ADSStateFlagAdsCommand for requests)
			stateFlags := binary.LittleEndian.Uint16(result[18:20])
			assert.Equal(t, uint16(types.ADSStateFlagAdsCommand), stateFlags)

			// Verify data length
			dataLength := binary.LittleEndian.Uint32(result[20:24])
			assert.Equal(t, tt.dataLength, dataLength)

			// Verify error code (always 0 for requests)
			errorCode := binary.LittleEndian.Uint32(result[24:28])
			assert.Equal(t, uint32(0), errorCode)

			// Verify invoke ID
			invokeID := binary.LittleEndian.Uint32(result[28:32])
			assert.Equal(t, tt.invokeID, invokeID)
		})
	}
}

// TestBuildAmsHeader_InvalidNetID tests error handling for invalid NetIDs.
func TestBuildAmsHeader_InvalidNetID(t *testing.T) {
	tests := []struct {
		name   string
		target AmsAddress
		source AmsAddress
	}{
		{
			name:   "Invalid target NetID - too short",
			target: AmsAddress{NetID: "192.168.1.100", Port: 851},
			source: AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905},
		},
		{
			name:   "Invalid target NetID - too long",
			target: AmsAddress{NetID: "192.168.1.100.1.1.1", Port: 851},
			source: AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905},
		},
		{
			name:   "Invalid target NetID - invalid format",
			target: AmsAddress{NetID: "invalid.netid", Port: 851},
			source: AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905},
		},
		{
			name:   "Invalid source NetID - too short",
			target: AmsAddress{NetID: "192.168.1.100.1.1", Port: 851},
			source: AmsAddress{NetID: "192.168", Port: 32905},
		},
		{
			name:   "Invalid source NetID - too long",
			target: AmsAddress{NetID: "192.168.1.100.1.1", Port: 851},
			source: AmsAddress{NetID: "192.168.1.50.1.1.1.1", Port: 32905},
		},
		{
			name:   "Invalid source NetID - invalid format",
			target: AmsAddress{NetID: "192.168.1.100.1.1", Port: 851},
			source: AmsAddress{NetID: "abc.def.ghi.jkl.mno.pqr", Port: 32905},
		},
		{
			name:   "Empty target NetID",
			target: AmsAddress{NetID: "", Port: 851},
			source: AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905},
		},
		{
			name:   "Empty source NetID",
			target: AmsAddress{NetID: "192.168.1.100.1.1", Port: 851},
			source: AmsAddress{NetID: "", Port: 32905},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildAmsHeader(tt.target, tt.source, types.ADSCommandRead, 12, 1)
			if err == nil {
				t.Fatal("Expected error for invalid NetID")
			}
		})
	}
}

// TestBuildAmsTcpHeader_LittleEndian verifies little-endian encoding.
func TestBuildAmsTcpHeader_LittleEndian(t *testing.T) {
	// Use a distinctive value to verify byte order
	command := types.AMSHeaderFlag(0x1234)
	dataLength := uint32(0x56789ABC)

	result := BuildAmsTcpHeader(command, dataLength)

	// Manually verify little-endian encoding
	assert.Equal(t, byte(0x34), result[0]) // Low byte of command
	assert.Equal(t, byte(0x12), result[1]) // High byte of command
	assert.Equal(t, byte(0xBC), result[2]) // Lowest byte of dataLength
	assert.Equal(t, byte(0x9A), result[3])
	assert.Equal(t, byte(0x78), result[4])
	assert.Equal(t, byte(0x56), result[5]) // Highest byte of dataLength
}

// TestBuildAmsHeader_LittleEndian verifies little-endian encoding.
func TestBuildAmsHeader_LittleEndian(t *testing.T) {
	target := AmsAddress{NetID: "1.2.3.4.5.6", Port: 0x1234}
	source := AmsAddress{NetID: "10.20.30.40.50.60", Port: 0x5678}
	command := types.ADSCommand(0xABCD)
	dataLength := uint32(0x12345678)
	invokeID := uint32(0x9ABCDEF0)

	result, err := BuildAmsHeader(target, source, command, dataLength, invokeID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify NetID bytes are copied directly (not endian-sensitive)
	assert.Equal(t, byte(1), result[0])
	assert.Equal(t, byte(2), result[1])
	assert.Equal(t, byte(3), result[2])
	assert.Equal(t, byte(4), result[3])
	assert.Equal(t, byte(5), result[4])
	assert.Equal(t, byte(6), result[5])

	// Verify little-endian encoding for port
	assert.Equal(t, byte(0x34), result[6]) // Low byte of target port
	assert.Equal(t, byte(0x12), result[7]) // High byte of target port

	// Verify source NetID
	assert.Equal(t, byte(10), result[8])
	assert.Equal(t, byte(20), result[9])
	assert.Equal(t, byte(30), result[10])
	assert.Equal(t, byte(40), result[11])
	assert.Equal(t, byte(50), result[12])
	assert.Equal(t, byte(60), result[13])

	// Verify little-endian encoding for source port
	assert.Equal(t, byte(0x78), result[14]) // Low byte of source port
	assert.Equal(t, byte(0x56), result[15]) // High byte of source port

	// Verify command
	assert.Equal(t, byte(0xCD), result[16]) // Low byte of command
	assert.Equal(t, byte(0xAB), result[17]) // High byte of command

	// Verify state flags (always 0x0004)
	assert.Equal(t, byte(0x04), result[18])
	assert.Equal(t, byte(0x00), result[19])

	// Verify data length
	assert.Equal(t, byte(0x78), result[20])
	assert.Equal(t, byte(0x56), result[21])
	assert.Equal(t, byte(0x34), result[22])
	assert.Equal(t, byte(0x12), result[23])

	// Verify error code (always 0)
	assert.Equal(t, byte(0x00), result[24])
	assert.Equal(t, byte(0x00), result[25])
	assert.Equal(t, byte(0x00), result[26])
	assert.Equal(t, byte(0x00), result[27])

	// Verify invoke ID
	assert.Equal(t, byte(0xF0), result[28])
	assert.Equal(t, byte(0xDE), result[29])
	assert.Equal(t, byte(0xBC), result[30])
	assert.Equal(t, byte(0x9A), result[31])
}

// TestBuildAmsHeader_StateFlags verifies state flags are always set to ADSStateFlagAdsCommand.
func TestBuildAmsHeader_StateFlags(t *testing.T) {
	commands := []types.ADSCommand{
		types.ADSCommandRead,
		types.ADSCommandWrite,
		types.ADSCommandReadWrite,
		types.ADSCommandWriteControl,
		types.ADSCommandReadState,
		types.ADSCommandReadDeviceInfo,
	}

	target := AmsAddress{NetID: "192.168.1.100.1.1", Port: 851}
	source := AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905}

	for _, cmd := range commands {
		t.Run(cmd.String(), func(t *testing.T) {
			result, err := BuildAmsHeader(target, source, cmd, 0, 1)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			stateFlags := binary.LittleEndian.Uint16(result[18:20])
			assert.Equal(t, uint16(types.ADSStateFlagAdsCommand), stateFlags)
		})
	}
}

// TestBuildAmsHeader_ErrorCodeZero verifies error code is always 0 for requests.
func TestBuildAmsHeader_ErrorCodeZero(t *testing.T) {
	target := AmsAddress{NetID: "192.168.1.100.1.1", Port: 851}
	source := AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905}

	result, err := BuildAmsHeader(target, source, types.ADSCommandRead, 100, 999)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	errorCode := binary.LittleEndian.Uint32(result[24:28])
	assert.Equal(t, uint32(0), errorCode)
}

// TestRealWorldScenario tests a complete ADS communication scenario.
func TestRealWorldScenario(t *testing.T) {
	// Scenario: Read variable from PLC
	target := AmsAddress{NetID: "192.168.1.100.1.1", Port: 851}
	source := AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905}
	command := types.ADSCommandRead
	dataLength := uint32(12) // Read request payload
	invokeID := uint32(12345)

	// Build AMS/TCP header
	tcpHeader := BuildAmsTcpHeader(types.AMSTCPPortAMSCommand, constants.AMSHeaderLength+dataLength)

	// Verify AMS/TCP header
	if len(tcpHeader) != constants.AMSTCPHeaderLength {
		t.Fatalf("Expected AMS/TCP header length %d, got %d", constants.AMSTCPHeaderLength, len(tcpHeader))
	}

	tcpCommand := binary.LittleEndian.Uint16(tcpHeader[0:2])
	assert.Equal(t, uint16(types.AMSTCPPortAMSCommand), tcpCommand)

	tcpDataLength := binary.LittleEndian.Uint32(tcpHeader[2:6])
	assert.Equal(t, constants.AMSHeaderLength+dataLength, tcpDataLength)

	// Build AMS header
	amsHeader, err := BuildAmsHeader(target, source, command, dataLength, invokeID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify AMS header
	if len(amsHeader) != constants.AMSHeaderLength {
		t.Fatalf("Expected AMS header length %d, got %d", constants.AMSHeaderLength, len(amsHeader))
	}

	// Verify all fields are correctly set
	targetNetID := utils.ByteArrayToAmsNetIdStr(amsHeader[0:6])
	assert.Equal(t, target.NetID, targetNetID)

	targetPort := binary.LittleEndian.Uint16(amsHeader[6:8])
	assert.Equal(t, target.Port, targetPort)

	sourceNetID := utils.ByteArrayToAmsNetIdStr(amsHeader[8:14])
	assert.Equal(t, source.NetID, sourceNetID)

	sourcePort := binary.LittleEndian.Uint16(amsHeader[14:16])
	assert.Equal(t, source.Port, sourcePort)

	amsCommand := binary.LittleEndian.Uint16(amsHeader[16:18])
	assert.Equal(t, uint16(command), amsCommand)

	stateFlags := binary.LittleEndian.Uint16(amsHeader[18:20])
	assert.Equal(t, uint16(types.ADSStateFlagAdsCommand), stateFlags)

	amsDataLength := binary.LittleEndian.Uint32(amsHeader[20:24])
	assert.Equal(t, dataLength, amsDataLength)

	errorCode := binary.LittleEndian.Uint32(amsHeader[24:28])
	assert.Equal(t, uint32(0), errorCode)

	amsInvokeID := binary.LittleEndian.Uint32(amsHeader[28:32])
	assert.Equal(t, invokeID, amsInvokeID)
}

// TestPortRegistration tests port registration scenario.
func TestPortRegistration(t *testing.T) {
	// Port registration uses AMSTCPPortConnect command
	dataLength := uint32(8) // AmsNetID (6 bytes) + Port (2 bytes)

	tcpHeader := BuildAmsTcpHeader(types.AMSTCPPortConnect, dataLength)

	// Verify header
	if len(tcpHeader) != constants.AMSTCPHeaderLength {
		t.Fatalf("Expected length %d, got %d", constants.AMSTCPHeaderLength, len(tcpHeader))
	}

	command := binary.LittleEndian.Uint16(tcpHeader[0:2])
	assert.Equal(t, uint16(types.AMSTCPPortConnect), command)

	length := binary.LittleEndian.Uint32(tcpHeader[2:6])
	assert.Equal(t, dataLength, length)
}

// TestPortUnregistration tests port unregistration scenario.
func TestPortUnregistration(t *testing.T) {
	// Port unregistration uses AMSTCPPortClose command
	dataLength := uint32(8) // AmsNetID (6 bytes) + Port (2 bytes)

	tcpHeader := BuildAmsTcpHeader(types.AMSTCPPortClose, dataLength)

	// Verify header
	if len(tcpHeader) != constants.AMSTCPHeaderLength {
		t.Fatalf("Expected length %d, got %d", constants.AMSTCPHeaderLength, len(tcpHeader))
	}

	command := binary.LittleEndian.Uint16(tcpHeader[0:2])
	assert.Equal(t, uint16(types.AMSTCPPortClose), command)

	length := binary.LittleEndian.Uint32(tcpHeader[2:6])
	assert.Equal(t, dataLength, length)
}

// TestBuildAmsHeader_AllADSCommands tests building headers for all ADS command types.
func TestBuildAmsHeader_AllADSCommands(t *testing.T) {
	target := AmsAddress{NetID: "192.168.1.100.1.1", Port: 851}
	source := AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905}

	commands := []struct {
		command    types.ADSCommand
		dataLength uint32
	}{
		{types.ADSCommandRead, 12},
		{types.ADSCommandWrite, 100},
		{types.ADSCommandReadState, 0},
		{types.ADSCommandWriteControl, 8},
		{types.ADSCommandReadDeviceInfo, 0},
		{types.ADSCommandReadWrite, 256},
	}

	for _, tc := range commands {
		t.Run(tc.command.String(), func(t *testing.T) {
			result, err := BuildAmsHeader(target, source, tc.command, tc.dataLength, 1)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(result) != constants.AMSHeaderLength {
				t.Fatalf("Expected length %d, got %d", constants.AMSHeaderLength, len(result))
			}

			command := binary.LittleEndian.Uint16(result[16:18])
			assert.Equal(t, uint16(tc.command), command)

			dataLength := binary.LittleEndian.Uint32(result[20:24])
			assert.Equal(t, tc.dataLength, dataLength)
		})
	}
}

// BenchmarkBuildAmsTcpHeader benchmarks the BuildAmsTcpHeader function.
func BenchmarkBuildAmsTcpHeader(b *testing.B) {
	command := types.AMSTCPPortAMSCommand
	dataLength := uint32(44)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildAmsTcpHeader(command, dataLength)
	}
}

// BenchmarkBuildAmsHeader benchmarks the BuildAmsHeader function.
func BenchmarkBuildAmsHeader(b *testing.B) {
	target := AmsAddress{NetID: "192.168.1.100.1.1", Port: 851}
	source := AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905}
	command := types.ADSCommandRead
	dataLength := uint32(12)
	invokeID := uint32(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BuildAmsHeader(target, source, command, dataLength, invokeID)
	}
}

// BenchmarkBuildAmsHeader_InvalidNetID benchmarks error handling.
func BenchmarkBuildAmsHeader_InvalidNetID(b *testing.B) {
	target := AmsAddress{NetID: "invalid", Port: 851}
	source := AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905}
	command := types.ADSCommandRead
	dataLength := uint32(12)
	invokeID := uint32(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BuildAmsHeader(target, source, command, dataLength, invokeID)
	}
}

// TestRoundTrip_WithAmsHeaderParser tests building and parsing headers together.
// This verifies that our builder creates headers that the parser can correctly read.
func TestRoundTrip_WithAmsHeaderParser(t *testing.T) {
	// This test would require importing ams-header module
	// We'll simulate the parsing logic here to verify round-trip

	target := AmsAddress{NetID: "192.168.1.100.1.1", Port: 851}
	source := AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905}
	command := types.ADSCommandRead
	dataLength := uint32(12)
	invokeID := uint32(12345)
	testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C}

	// Build headers
	tcpHeader := BuildAmsTcpHeader(types.AMSTCPPortAMSCommand, constants.AMSHeaderLength+dataLength)
	amsHeader, err := BuildAmsHeader(target, source, command, dataLength, invokeID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Combine into complete packet
	packet := make([]byte, 0, len(tcpHeader)+len(amsHeader)+len(testData))
	packet = append(packet, tcpHeader...)
	packet = append(packet, amsHeader...)
	packet = append(packet, testData...)

	// Parse back (simulating ams-header parser logic)
	if len(packet) < constants.AMSTCPHeaderLength+constants.AMSHeaderLength {
		t.Fatal("Packet too short")
	}

	// Parse AMS/TCP header
	tcpCommand := binary.LittleEndian.Uint16(packet[0:2])
	assert.Equal(t, uint16(types.AMSTCPPortAMSCommand), tcpCommand)

	tcpDataLength := binary.LittleEndian.Uint32(packet[2:6])
	assert.Equal(t, constants.AMSHeaderLength+dataLength, tcpDataLength)

	// Parse AMS header
	amsPacket := packet[constants.AMSTCPHeaderLength:]

	targetNetID := utils.ByteArrayToAmsNetIdStr(amsPacket[0:6])
	assert.Equal(t, target.NetID, targetNetID)

	targetPort := binary.LittleEndian.Uint16(amsPacket[6:8])
	assert.Equal(t, target.Port, targetPort)

	sourceNetID := utils.ByteArrayToAmsNetIdStr(amsPacket[8:14])
	assert.Equal(t, source.NetID, sourceNetID)

	sourcePort := binary.LittleEndian.Uint16(amsPacket[14:16])
	assert.Equal(t, source.Port, sourcePort)

	parsedCommand := types.ADSCommand(binary.LittleEndian.Uint16(amsPacket[16:18]))
	assert.Equal(t, command, parsedCommand)

	stateFlags := binary.LittleEndian.Uint16(amsPacket[18:20])
	assert.Equal(t, uint16(types.ADSStateFlagAdsCommand), stateFlags)

	parsedDataLength := binary.LittleEndian.Uint32(amsPacket[20:24])
	assert.Equal(t, dataLength, parsedDataLength)

	errorCode := binary.LittleEndian.Uint32(amsPacket[24:28])
	assert.Equal(t, uint32(0), errorCode)

	parsedInvokeID := binary.LittleEndian.Uint32(amsPacket[28:32])
	assert.Equal(t, invokeID, parsedInvokeID)

	// Parse data
	parsedData := amsPacket[constants.AMSHeaderLength : constants.AMSHeaderLength+dataLength]
	assert.Equal(t, testData, parsedData)
}

// TestBuildAmsHeader_NetIDEdgeCases tests edge cases for NetID values.
func TestBuildAmsHeader_NetIDEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		target AmsAddress
		source AmsAddress
	}{
		{
			name:   "All zeros",
			target: AmsAddress{NetID: "0.0.0.0.0.0", Port: 851},
			source: AmsAddress{NetID: "0.0.0.0.0.0", Port: 32905},
		},
		{
			name:   "All 255s",
			target: AmsAddress{NetID: "255.255.255.255.255.255", Port: 851},
			source: AmsAddress{NetID: "255.255.255.255.255.255", Port: 32905},
		},
		{
			name:   "Mixed values",
			target: AmsAddress{NetID: "0.255.0.255.0.255", Port: 0},
			source: AmsAddress{NetID: "255.0.255.0.255.0", Port: 65535},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildAmsHeader(tt.target, tt.source, types.ADSCommandRead, 0, 0)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// Verify NetIDs can be parsed back
			targetNetID := utils.ByteArrayToAmsNetIdStr(result[0:6])
			assert.Equal(t, tt.target.NetID, targetNetID)

			sourceNetID := utils.ByteArrayToAmsNetIdStr(result[8:14])
			assert.Equal(t, tt.source.NetID, sourceNetID)
		})
	}
}

// TestBuildAmsHeader_UtilsError tests that errors from utils are properly propagated.
func TestBuildAmsHeader_UtilsError(t *testing.T) {
	// Test that we properly return errors from utils.AmsNetIdStrToByteArray
	invalidNetIDs := []string{
		"",
		"invalid",
		"1.2.3",
		"1.2.3.4.5.6.7",
		"a.b.c.d.e.f",
	}

	for _, netID := range invalidNetIDs {
		t.Run("Invalid_"+netID, func(t *testing.T) {
			target := AmsAddress{NetID: netID, Port: 851}
			source := AmsAddress{NetID: "192.168.1.50.1.1", Port: 32905}

			_, err := BuildAmsHeader(target, source, types.ADSCommandRead, 0, 0)
			if err == nil {
				t.Fatal("Expected error for invalid NetID")
			}

			// Verify it's not a nil error
			if errors.Is(err, nil) {
				t.Fatal("Expected non-nil error")
			}
		})
	}
}
