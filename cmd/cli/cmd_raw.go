package cli

import (
	"fmt"

	"github.com/jarmocluyse/ads-go/pkg/ads"
)

// handleReadRaw reads raw data from the PLC using index group and offset.
// Usage: read_raw
func handleReadRaw(args []string, client *ads.Client) {
	indexGroup := uint32(0x1010290)
	indexOffset := uint32(0x80000001)
	var port uint16 = 350
	size := uint32(1) // adjust size as needed
	result, err := client.ReadRaw(port, indexGroup, indexOffset, size)
	if err != nil {
		fmt.Printf("[ERROR] Command 'read_raw': Failed to read raw data [IG: 0x%X, IO: 0x%X, port: %d]: %v\n", indexGroup, indexOffset, port, err)
		return
	}
	fmt.Printf("[OK] Raw read [IG: 0x%X, IO: 0x%X, port: %d]: %v\n", indexGroup, indexOffset, port, result)
}

// handleWriteRaw writes raw data to the PLC using index group and offset.
// Usage: write_raw
func handleWriteRaw(args []string, client *ads.Client) {
	indexGroup := uint32(0x1010290)
	indexOffset := uint32(0x80000001)
	var port uint16 = 350
	data := []byte{36}
	err := client.WriteRaw(port, indexGroup, indexOffset, data)
	if err != nil {
		fmt.Printf("[ERROR] Command 'write_raw': Failed to write raw data [IG: 0x%X, IO: 0x%X, port: %d]: %v\n", indexGroup, indexOffset, port, err)
		return
	}
	fmt.Printf("[OK] Raw write [IG: 0x%X, IO: 0x%X, port: %d] succeeded.\n", indexGroup, indexOffset, port)
}
