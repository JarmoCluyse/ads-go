package ads

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConnectWithLocalNetIDSkipsPortRegistration verifies that with LocalNetID
// set, Connect sends no port registration (a remote PLC router closes the
// connection on one) and uses the configured address as the requests' source.
func TestConnectWithLocalNetIDSkipsPortRegistration(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = ln.Close() }()

	first := make(chan []byte, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		// AMS/TCP header (6 bytes) + target and source address of the AMS header (16 bytes)
		buf := make([]byte, 22)
		if _, err := io.ReadFull(conn, buf); err == nil {
			first <- buf
		}
		_, _ = io.Copy(io.Discard, conn)
	}()

	c := NewClient(ClientSettings{
		TargetNetID: "10.255.4.78.1.1",
		RouterHost:  "127.0.0.1",
		RouterPort:  ln.Addr().(*net.TCPAddr).Port,
		Timeout:     100 * time.Millisecond,
		LocalNetID:  "192.168.3.254.1.1",
	}, nil)
	require.NoError(t, c.Connect())
	defer func() { _ = c.Disconnect() }()

	select {
	case buf := <-first:
		assert.Equal(t, []byte{0, 0}, buf[0:2], "first packet must be an ADS command, not a port registration")
		assert.Equal(t, []byte{192, 168, 3, 254, 1, 1}, buf[14:20], "source AMS Net ID")
		assert.Equal(t, []byte{0xEE, 0x7F}, buf[20:22], "source ADS port 32750")
	case <-time.After(time.Second):
		t.Fatal("no packet received")
	}
}
