package amsheader

import (
	"encoding/binary"
	"errors"
	"testing"

	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/stretchr/testify/assert"
)

// buildAmsPacket creates a valid AMS packet for testing.
func buildAmsPacket(targetNetID, sourceNetID []byte, targetPort, sourcePort uint16, command, stateFlags uint16, errorCode, invokeID uint32, data []byte) []byte {
	dataLen := uint32(len(data))
	packetLen := AMSHeaderLength + dataLen

	packet := make([]byte, AMSTCPHeaderLength+AMSHeaderLength+dataLen)

	// AMS/TCP Header
	binary.LittleEndian.PutUint16(packet[0:2], 0) // Reserved
	binary.LittleEndian.PutUint32(packet[2:6], packetLen)

	// AMS Header
	copy(packet[6:12], targetNetID)
	binary.LittleEndian.PutUint16(packet[12:14], targetPort)
	copy(packet[14:20], sourceNetID)
	binary.LittleEndian.PutUint16(packet[20:22], sourcePort)
	binary.LittleEndian.PutUint16(packet[22:24], command)
	binary.LittleEndian.PutUint16(packet[24:26], stateFlags)
	binary.LittleEndian.PutUint32(packet[26:30], dataLen)
	binary.LittleEndian.PutUint32(packet[30:34], errorCode)
	binary.LittleEndian.PutUint32(packet[34:38], invokeID)

	// Data
	copy(packet[38:], data)

	return packet
}

// netIDFromString converts a string representation to 6-byte NetID.
func netIDFromInts(a, b, c, d, e, f byte) []byte {
	return []byte{a, b, c, d, e, f}
}

func TestParsePacket_Valid(t *testing.T) {
	targetNetID := netIDFromInts(192, 168, 1, 100, 1, 1)
	sourceNetID := netIDFromInts(192, 168, 1, 101, 1, 1)
	data := []byte{0x01, 0x02, 0x03, 0x04}

	packet := buildAmsPacket(
		targetNetID, sourceNetID,
		851, 852,
		uint16(types.ADSCommandRead), 0x0004,
		0, 12345,
		data,
	)

	result, err := ParsePacket(packet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, "192.168.1.100.1.1", result.TargetNetID)
	assert.Equal(t, uint16(851), result.TargetPort)
	assert.Equal(t, "192.168.1.101.1.1", result.SourceNetID)
	assert.Equal(t, uint16(852), result.SourcePort)
	assert.Equal(t, types.ADSCommandRead, result.Command)
	assert.Equal(t, types.ADSStateFlags(0x0004), result.StateFlags)
	assert.Equal(t, uint32(4), result.DataLength)
	assert.Equal(t, uint32(0), result.ErrorCode)
	assert.Equal(t, uint32(12345), result.InvokeID)
	assert.Equal(t, data, result.Data)
}

func TestParsePacket_NoData(t *testing.T) {
	targetNetID := netIDFromInts(1, 2, 3, 4, 5, 6)
	sourceNetID := netIDFromInts(6, 5, 4, 3, 2, 1)

	packet := buildAmsPacket(
		targetNetID, sourceNetID,
		100, 200,
		uint16(types.ADSCommandWrite), 0,
		0, 999,
		nil,
	)

	result, err := ParsePacket(packet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, "1.2.3.4.5.6", result.TargetNetID)
	assert.Equal(t, uint16(100), result.TargetPort)
	assert.Equal(t, "6.5.4.3.2.1", result.SourceNetID)
	assert.Equal(t, uint16(200), result.SourcePort)
	assert.Equal(t, types.ADSCommandWrite, result.Command)
	assert.Equal(t, uint32(0), result.DataLength)
	assert.Empty(t, result.Data)
}

func TestParsePacket_WithError(t *testing.T) {
	targetNetID := netIDFromInts(10, 20, 30, 40, 50, 60)
	sourceNetID := netIDFromInts(60, 50, 40, 30, 20, 10)

	packet := buildAmsPacket(
		targetNetID, sourceNetID,
		500, 600,
		uint16(types.ADSCommandReadWrite), 0,
		1792, // Error code
		555,
		nil,
	)

	result, err := ParsePacket(packet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, uint32(1792), result.ErrorCode)
	assert.Equal(t, uint32(555), result.InvokeID)
}

func TestParsePacket_InsufficientData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"Empty", []byte{}},
		{"Only 1 byte", []byte{0x01}},
		{"Only AMS/TCP header", make([]byte, AMSTCPHeaderLength)},
		{"Partial AMS header", make([]byte, AMSTCPHeaderLength+10)},
		{"Missing data payload", func() []byte {
			// Create packet claiming 10 bytes of data but don't include it
			packet := make([]byte, AMSTCPHeaderLength+AMSHeaderLength)
			binary.LittleEndian.PutUint32(packet[2:6], AMSHeaderLength+10) // Claim 10 bytes
			return packet
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParsePacket(tt.data)
			if err == nil {
				t.Fatal("Expected error for insufficient data")
			}
			assert.True(t, errors.Is(err, ErrInsufficientData))
		})
	}
}

func TestParsePacket_InvalidLength(t *testing.T) {
	tests := []struct {
		name      string
		packetLen uint32
	}{
		{"Zero length", 0},
		{"Less than AMS header", AMSHeaderLength - 1},
		{"Much too small", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := make([]byte, AMSTCPHeaderLength+100)
			binary.LittleEndian.PutUint32(packet[2:6], tt.packetLen)

			_, err := ParsePacket(packet)
			if err == nil {
				t.Fatal("Expected error for invalid length")
			}
			assert.True(t, errors.Is(err, ErrInvalidLength))
		})
	}
}

func TestParsePacket_LargeData(t *testing.T) {
	targetNetID := netIDFromInts(1, 1, 1, 1, 1, 1)
	sourceNetID := netIDFromInts(2, 2, 2, 2, 2, 2)
	data := make([]byte, 1000)
	for i := range data {
		data[i] = byte(i % 256)
	}

	packet := buildAmsPacket(
		targetNetID, sourceNetID,
		800, 900,
		uint16(types.ADSCommandRead), 0,
		0, 7777,
		data,
	)

	result, err := ParsePacket(packet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, uint32(1000), result.DataLength)
	assert.Equal(t, data, result.Data)
}

func TestCheckPacket_Valid(t *testing.T) {
	targetNetID := netIDFromInts(1, 2, 3, 4, 5, 6)
	sourceNetID := netIDFromInts(6, 5, 4, 3, 2, 1)
	data := []byte{0xAA, 0xBB}

	packet := buildAmsPacket(
		targetNetID, sourceNetID,
		100, 200,
		uint16(types.ADSCommandRead), 0,
		0, 111,
		data,
	)

	err := CheckPacket(packet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCheckPacket_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		buildPacket func() []byte
		expectError error
	}{
		{
			"Empty data",
			func() []byte { return []byte{} },
			ErrInsufficientData,
		},
		{
			"Only TCP header",
			func() []byte { return make([]byte, AMSTCPHeaderLength) },
			ErrInsufficientData,
		},
		{
			"Invalid packet length",
			func() []byte {
				packet := make([]byte, 100)
				binary.LittleEndian.PutUint32(packet[2:6], 10) // Too small
				return packet
			},
			ErrInvalidLength,
		},
		{
			"Length mismatch",
			func() []byte {
				packet := make([]byte, AMSTCPHeaderLength+AMSHeaderLength)
				binary.LittleEndian.PutUint32(packet[2:6], AMSHeaderLength+50) // Claims more
				return packet
			},
			ErrInsufficientData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := tt.buildPacket()
			err := CheckPacket(packet)
			if err == nil {
				t.Fatal("Expected error")
			}
			assert.True(t, errors.Is(err, tt.expectError))
		})
	}
}

func TestCheckTCPPacketLength_Valid(t *testing.T) {
	targetNetID := netIDFromInts(1, 2, 3, 4, 5, 6)
	sourceNetID := netIDFromInts(6, 5, 4, 3, 2, 1)
	data := []byte{0x01, 0x02, 0x03}

	packet := buildAmsPacket(
		targetNetID, sourceNetID,
		100, 200,
		uint16(types.ADSCommandRead), 0,
		0, 123,
		data,
	)

	expectedLen := uint32(AMSTCPHeaderLength + AMSHeaderLength + len(data))

	length, err := CheckTCPPacketLength(packet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, expectedLen, length)
}

func TestCheckTCPPacketLength_NoData(t *testing.T) {
	targetNetID := netIDFromInts(1, 2, 3, 4, 5, 6)
	sourceNetID := netIDFromInts(6, 5, 4, 3, 2, 1)

	packet := buildAmsPacket(
		targetNetID, sourceNetID,
		100, 200,
		uint16(types.ADSCommandRead), 0,
		0, 123,
		nil,
	)

	expectedLen := uint32(AMSTCPHeaderLength + AMSHeaderLength)

	length, err := CheckTCPPacketLength(packet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, expectedLen, length)
}

func TestCheckTCPPacketLength_InsufficientData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"Empty", []byte{}},
		{"Only 3 bytes", []byte{0x01, 0x02, 0x03}},
		{"Only 5 bytes", make([]byte, 5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CheckTCPPacketLength(tt.data)
			if err == nil {
				t.Fatal("Expected error")
			}
			assert.True(t, errors.Is(err, ErrInsufficientData))
		})
	}
}

func TestCheckTCPPacketLength_InvalidLength(t *testing.T) {
	packet := make([]byte, 100)
	binary.LittleEndian.PutUint32(packet[2:6], 20) // Less than AMS header

	_, err := CheckTCPPacketLength(packet)
	if err == nil {
		t.Fatal("Expected error for invalid length")
	}
	assert.True(t, errors.Is(err, ErrInvalidLength))
}

func TestCheckTCPPacketLength_PartialPacket(t *testing.T) {
	packet := make([]byte, AMSTCPHeaderLength+10)
	binary.LittleEndian.PutUint32(packet[2:6], AMSHeaderLength+50) // Claims more

	_, err := CheckTCPPacketLength(packet)
	if err == nil {
		t.Fatal("Expected error for partial packet")
	}
	assert.True(t, errors.Is(err, ErrInsufficientData))
}
