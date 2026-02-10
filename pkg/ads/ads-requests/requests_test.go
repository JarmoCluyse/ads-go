package adsrequests

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test BuildReadRequest
func TestBuildReadRequest(t *testing.T) {
	t.Run("standard read request", func(t *testing.T) {
		payload := BuildReadRequest(0x4020, 0x1234, 100)

		assert.Len(t, payload, 12)

		indexGroup := binary.LittleEndian.Uint32(payload[0:4])
		assert.Equal(t, uint32(0x4020), indexGroup)

		indexOffset := binary.LittleEndian.Uint32(payload[4:8])
		assert.Equal(t, uint32(0x1234), indexOffset)

		readLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(100), readLength)
	})

	t.Run("read request with zero values", func(t *testing.T) {
		payload := BuildReadRequest(0, 0, 0)

		assert.Len(t, payload, 12)
		assert.Equal(t, make([]byte, 12), payload)
	})

	t.Run("read request with max values", func(t *testing.T) {
		payload := BuildReadRequest(0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF)

		assert.Len(t, payload, 12)

		indexGroup := binary.LittleEndian.Uint32(payload[0:4])
		assert.Equal(t, uint32(0xFFFFFFFF), indexGroup)

		indexOffset := binary.LittleEndian.Uint32(payload[4:8])
		assert.Equal(t, uint32(0xFFFFFFFF), indexOffset)

		readLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(0xFFFFFFFF), readLength)
	})

	t.Run("read request for symbol info", func(t *testing.T) {
		// Common pattern: IndexGroup 0xF009 for symbol info
		payload := BuildReadRequest(0xF009, 0, 0xFFFFFFFF)

		assert.Len(t, payload, 12)

		indexGroup := binary.LittleEndian.Uint32(payload[0:4])
		assert.Equal(t, uint32(0xF009), indexGroup)

		readLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(0xFFFFFFFF), readLength)
	})
}

// Test BuildWriteRequest
func TestBuildWriteRequest(t *testing.T) {
	t.Run("write request with data", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03, 0x04}
		payload := BuildWriteRequest(0x4020, 0x5678, data)

		assert.Len(t, payload, 16) // 12 header + 4 data

		indexGroup := binary.LittleEndian.Uint32(payload[0:4])
		assert.Equal(t, uint32(0x4020), indexGroup)

		indexOffset := binary.LittleEndian.Uint32(payload[4:8])
		assert.Equal(t, uint32(0x5678), indexOffset)

		dataLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(4), dataLength)

		assert.Equal(t, data, payload[12:16])
	})

	t.Run("write request with empty data", func(t *testing.T) {
		payload := BuildWriteRequest(0x1000, 0x2000, []byte{})

		assert.Len(t, payload, 12) // Just header, no data

		dataLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(0), dataLength)
	})

	t.Run("write request with large data", func(t *testing.T) {
		data := make([]byte, 1024)
		for i := range data {
			data[i] = byte(i % 256)
		}

		payload := BuildWriteRequest(0xABCD, 0xEF01, data)

		assert.Len(t, payload, 1036) // 12 header + 1024 data

		dataLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(1024), dataLength)

		assert.Equal(t, data, payload[12:])
	})

	t.Run("write request with single byte", func(t *testing.T) {
		payload := BuildWriteRequest(0x1234, 0x5678, []byte{0xFF})

		assert.Len(t, payload, 13)

		dataLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(1), dataLength)

		assert.Equal(t, byte(0xFF), payload[12])
	})
}

// Test BuildReadWriteRequest
func TestBuildReadWriteRequest(t *testing.T) {
	t.Run("standard readwrite request", func(t *testing.T) {
		writeData := []byte("MAIN.myVar")
		payload := BuildReadWriteRequest(0xF009, 0, 0xFFFFFFFF, writeData)

		assert.Len(t, payload, 26) // 16 header + 10 data

		indexGroup := binary.LittleEndian.Uint32(payload[0:4])
		assert.Equal(t, uint32(0xF009), indexGroup)

		indexOffset := binary.LittleEndian.Uint32(payload[4:8])
		assert.Equal(t, uint32(0), indexOffset)

		readLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(0xFFFFFFFF), readLength)

		writeLength := binary.LittleEndian.Uint32(payload[12:16])
		assert.Equal(t, uint32(10), writeLength)

		assert.Equal(t, writeData, payload[16:])
	})

	t.Run("readwrite request with empty write data", func(t *testing.T) {
		payload := BuildReadWriteRequest(0x1000, 0x2000, 100, []byte{})

		assert.Len(t, payload, 16) // Just header, no data

		writeLength := binary.LittleEndian.Uint32(payload[12:16])
		assert.Equal(t, uint32(0), writeLength)
	})

	t.Run("readwrite request with zero read length", func(t *testing.T) {
		writeData := []byte{0x01, 0x02}
		payload := BuildReadWriteRequest(0x5000, 0x6000, 0, writeData)

		readLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(0), readLength)

		writeLength := binary.LittleEndian.Uint32(payload[12:16])
		assert.Equal(t, uint32(2), writeLength)
	})

	t.Run("readwrite for data type query", func(t *testing.T) {
		typeName := []byte("ST_MyStruct")
		payload := BuildReadWriteRequest(0xF011, 0, 0xFFFFFFFF, typeName)

		assert.Len(t, payload, 27) // 16 header + 11 data

		indexGroup := binary.LittleEndian.Uint32(payload[0:4])
		assert.Equal(t, uint32(0xF011), indexGroup)

		assert.Equal(t, typeName, payload[16:])
	})
}

// Test BuildReadWriteRequestWithNullTerminator
func TestBuildReadWriteRequestWithNullTerminator(t *testing.T) {
	t.Run("readwrite with null terminator", func(t *testing.T) {
		writeData := []byte("MAIN.Counter")
		payload := BuildReadWriteRequestWithNullTerminator(0xF009, 0, 0xFFFFFFFF, writeData)

		assert.Len(t, payload, 29) // 16 header + 12 data + 1 null

		writeLength := binary.LittleEndian.Uint32(payload[12:16])
		assert.Equal(t, uint32(13), writeLength) // 12 chars + 1 null

		assert.Equal(t, writeData, payload[16:28])
		assert.Equal(t, byte(0x00), payload[28]) // Null terminator
	})

	t.Run("empty string with null terminator", func(t *testing.T) {
		payload := BuildReadWriteRequestWithNullTerminator(0xF009, 0, 100, []byte{})

		assert.Len(t, payload, 17) // 16 header + 1 null

		writeLength := binary.LittleEndian.Uint32(payload[12:16])
		assert.Equal(t, uint32(1), writeLength)

		assert.Equal(t, byte(0x00), payload[16])
	})

	t.Run("single char with null terminator", func(t *testing.T) {
		payload := BuildReadWriteRequestWithNullTerminator(0xF009, 0, 50, []byte("X"))

		assert.Len(t, payload, 18) // 16 header + 1 char + 1 null

		writeLength := binary.LittleEndian.Uint32(payload[12:16])
		assert.Equal(t, uint32(2), writeLength)

		assert.Equal(t, byte('X'), payload[16])
		assert.Equal(t, byte(0x00), payload[17])
	})

	t.Run("comparison with manual null termination", func(t *testing.T) {
		text := []byte("Test")

		// Using helper function
		payload1 := BuildReadWriteRequestWithNullTerminator(0x1000, 0x2000, 100, text)

		// Manual approach
		textWithNull := append(text, 0x00)
		payload2 := BuildReadWriteRequest(0x1000, 0x2000, 100, textWithNull)

		// Should produce identical payloads
		assert.Equal(t, payload2, payload1)
	})
}

// Test BuildWriteControlRequest
func TestBuildWriteControlRequest(t *testing.T) {
	t.Run("write control run state", func(t *testing.T) {
		payload := BuildWriteControlRequest(5, 0) // Run state

		assert.Len(t, payload, 8)

		adsState := binary.LittleEndian.Uint16(payload[0:2])
		assert.Equal(t, uint16(5), adsState)

		deviceState := binary.LittleEndian.Uint16(payload[2:4])
		assert.Equal(t, uint16(0), deviceState)

		dataLength := binary.LittleEndian.Uint32(payload[4:8])
		assert.Equal(t, uint32(0), dataLength)
	})

	t.Run("write control config state", func(t *testing.T) {
		payload := BuildWriteControlRequest(15, 0) // Config state

		adsState := binary.LittleEndian.Uint16(payload[0:2])
		assert.Equal(t, uint16(15), adsState)
	})

	t.Run("write control reset state", func(t *testing.T) {
		payload := BuildWriteControlRequest(2, 0) // Reset state

		adsState := binary.LittleEndian.Uint16(payload[0:2])
		assert.Equal(t, uint16(2), adsState)
	})

	t.Run("write control reconfig state", func(t *testing.T) {
		payload := BuildWriteControlRequest(16, 0) // Reconfig state

		adsState := binary.LittleEndian.Uint16(payload[0:2])
		assert.Equal(t, uint16(16), adsState)
	})

	t.Run("write control with device state", func(t *testing.T) {
		payload := BuildWriteControlRequest(5, 100)

		adsState := binary.LittleEndian.Uint16(payload[0:2])
		assert.Equal(t, uint16(5), adsState)

		deviceState := binary.LittleEndian.Uint16(payload[2:4])
		assert.Equal(t, uint16(100), deviceState)
	})

	t.Run("write control with max values", func(t *testing.T) {
		payload := BuildWriteControlRequest(0xFFFF, 0xFFFF)

		assert.Len(t, payload, 8)

		adsState := binary.LittleEndian.Uint16(payload[0:2])
		assert.Equal(t, uint16(0xFFFF), adsState)

		deviceState := binary.LittleEndian.Uint16(payload[2:4])
		assert.Equal(t, uint16(0xFFFF), deviceState)
	})
}

// Test BuildReadStateRequest
func TestBuildReadStateRequest(t *testing.T) {
	t.Run("read state request is empty", func(t *testing.T) {
		payload := BuildReadStateRequest()

		assert.NotNil(t, payload)
		assert.Len(t, payload, 0)
		assert.Equal(t, []byte{}, payload)
	})
}

// Test BuildReadDeviceInfoRequest
func TestBuildReadDeviceInfoRequest(t *testing.T) {
	t.Run("read device info request is empty", func(t *testing.T) {
		payload := BuildReadDeviceInfoRequest()

		assert.NotNil(t, payload)
		assert.Len(t, payload, 0)
		assert.Equal(t, []byte{}, payload)
	})
}

// Test Little-Endian Encoding
func TestLittleEndianEncoding(t *testing.T) {
	t.Run("verify little-endian for read request", func(t *testing.T) {
		payload := BuildReadRequest(0x12345678, 0x9ABCDEF0, 0x11223344)

		// 0x12345678 in little-endian = [0x78, 0x56, 0x34, 0x12]
		assert.Equal(t, byte(0x78), payload[0])
		assert.Equal(t, byte(0x56), payload[1])
		assert.Equal(t, byte(0x34), payload[2])
		assert.Equal(t, byte(0x12), payload[3])

		// 0x9ABCDEF0 in little-endian = [0xF0, 0xDE, 0xBC, 0x9A]
		assert.Equal(t, byte(0xF0), payload[4])
		assert.Equal(t, byte(0xDE), payload[5])
		assert.Equal(t, byte(0xBC), payload[6])
		assert.Equal(t, byte(0x9A), payload[7])

		// 0x11223344 in little-endian = [0x44, 0x33, 0x22, 0x11]
		assert.Equal(t, byte(0x44), payload[8])
		assert.Equal(t, byte(0x33), payload[9])
		assert.Equal(t, byte(0x22), payload[10])
		assert.Equal(t, byte(0x11), payload[11])
	})

	t.Run("verify little-endian for write control", func(t *testing.T) {
		payload := BuildWriteControlRequest(0x1234, 0x5678)

		// 0x1234 in little-endian = [0x34, 0x12]
		assert.Equal(t, byte(0x34), payload[0])
		assert.Equal(t, byte(0x12), payload[1])

		// 0x5678 in little-endian = [0x78, 0x56]
		assert.Equal(t, byte(0x78), payload[2])
		assert.Equal(t, byte(0x56), payload[3])
	})
}

// Test Real-World Scenarios
func TestRealWorldScenarios(t *testing.T) {
	t.Run("read PLC memory", func(t *testing.T) {
		// Read 4 bytes from memory area MB (0x4020) offset 100
		payload := BuildReadRequest(0x4020, 100, 4)

		assert.Len(t, payload, 12)

		indexGroup := binary.LittleEndian.Uint32(payload[0:4])
		assert.Equal(t, uint32(0x4020), indexGroup)
	})

	t.Run("write PLC output", func(t *testing.T) {
		// Write boolean TRUE to output bit
		data := []byte{0x01}
		payload := BuildWriteRequest(0x4010, 0, data)

		assert.Len(t, payload, 13)

		dataLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(1), dataLength)

		assert.Equal(t, byte(0x01), payload[12])
	})

	t.Run("query symbol by name", func(t *testing.T) {
		// Get symbol info for "MAIN.Counter"
		symbolName := []byte("MAIN.Counter")
		payload := BuildReadWriteRequestWithNullTerminator(0xF009, 0, 0xFFFFFFFF, symbolName)

		indexGroup := binary.LittleEndian.Uint32(payload[0:4])
		assert.Equal(t, uint32(0xF009), indexGroup) // SymbolInfoByNameEx

		readLength := binary.LittleEndian.Uint32(payload[8:12])
		assert.Equal(t, uint32(0xFFFFFFFF), readLength) // Read all

		assert.Equal(t, symbolName, payload[16:28])
		assert.Equal(t, byte(0x00), payload[28]) // Null terminator
	})

	t.Run("query data type by name", func(t *testing.T) {
		// Get data type info for "INT"
		typeName := []byte("INT")
		payload := BuildReadWriteRequestWithNullTerminator(0xF011, 0, 0xFFFFFFFF, typeName)

		indexGroup := binary.LittleEndian.Uint32(payload[0:4])
		assert.Equal(t, uint32(0xF011), indexGroup) // DataTypeInfoByNameEx

		assert.Equal(t, typeName, payload[16:19])
		assert.Equal(t, byte(0x00), payload[19])
	})

	t.Run("set PLC to run mode", func(t *testing.T) {
		payload := BuildWriteControlRequest(5, 0) // ADSStateRun = 5

		adsState := binary.LittleEndian.Uint16(payload[0:2])
		assert.Equal(t, uint16(5), adsState)
	})

	t.Run("set PLC to config mode", func(t *testing.T) {
		payload := BuildWriteControlRequest(15, 0) // ADSStateConfig = 15

		adsState := binary.LittleEndian.Uint16(payload[0:2])
		assert.Equal(t, uint16(15), adsState)
	})
}

// Benchmark tests
func BenchmarkBuildReadRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildReadRequest(0x4020, 0x1234, 100)
	}
}

func BenchmarkBuildWriteRequest(b *testing.B) {
	data := make([]byte, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildWriteRequest(0x4020, 0x5678, data)
	}
}

func BenchmarkBuildReadWriteRequest(b *testing.B) {
	writeData := []byte("MAIN.myVariable")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildReadWriteRequest(0xF009, 0, 0xFFFFFFFF, writeData)
	}
}

func BenchmarkBuildReadWriteRequestWithNullTerminator(b *testing.B) {
	writeData := []byte("MAIN.myVariable")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildReadWriteRequestWithNullTerminator(0xF009, 0, 0xFFFFFFFF, writeData)
	}
}
