package ads

import (
	"io"
	"log/slog"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// newTestClient returns a minimal Client suitable for reconnect tests.
// It bypasses NewClient/Connect so we can control conn directly.
func newTestClient(settings ClientSettings) *Client {
	return &Client{
		settings:      settings,
		requests:      make(map[uint32]chan Response),
		subscriptions: make(map[uint32]*ActiveSubscription),
		logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

// startReceive launches c.receive() in a goroutine and blocks until the
// goroutine has captured c.conn into its local variable. This uses the
// onConnCaptured test hook so there is no data race when the caller
// subsequently writes to c.conn.
func startReceive(c *Client) <-chan struct{} {
	captured := make(chan struct{})
	c.onConnCaptured = func() {
		close(captured)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		c.receive()
	}()
	<-captured // wait until receive() has captured c.conn
	return done
}

// TestStaleGoroutineDoesNotCloseNewConnection verifies that when c.conn is
// replaced (simulating a reconnect) the old receive() goroutine's deferred
// conn.Close() only closes the connection it captured at startup, leaving the
// new connection fully operational.
func TestStaleGoroutineDoesNotCloseNewConnection(t *testing.T) {
	serverA, clientA := net.Pipe()
	serverB, clientB := net.Pipe()
	defer func() { _ = serverB.Close() }()
	defer func() { _ = clientB.Close() }()

	c := newTestClient(ClientSettings{})
	c.conn = clientA
	done := startReceive(c)

	// Simulate reconnect: replace c.conn with the new connection.
	c.conn = clientB

	// Trigger EOF in the stale goroutine.
	_ = serverA.Close()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("old receive() goroutine did not exit in time")
	}

	// clientB must still be usable.
	msg := []byte("still alive")
	writeErr := make(chan error, 1)
	go func() {
		_, err := serverB.Write(msg)
		writeErr <- err
	}()

	buf := make([]byte, len(msg))
	_ = clientB.SetDeadline(time.Now().Add(500 * time.Millisecond))
	_, readErr := clientB.Read(buf)

	assert.NoError(t, <-writeErr, "serverB write should succeed — clientB must still be open")
	assert.NoError(t, readErr, "clientB read should succeed — it must not have been closed by the stale goroutine")
	assert.Equal(t, msg, buf)
}

// TestStaleGoroutineDoesNotFireOnConnectionLost verifies that when the old
// receive() goroutine detects that c.conn has already been replaced it does
// NOT call the OnConnectionLost hook.
func TestStaleGoroutineDoesNotFireOnConnectionLost(t *testing.T) {
	var hookCalled atomic.Bool

	settings := ClientSettings{
		OnConnectionLost: func(_ *Client, _ error) {
			hookCalled.Store(true)
		},
	}

	serverOld, clientOld := net.Pipe()
	c := newTestClient(settings)
	c.conn = clientOld
	done := startReceive(c)

	// Simulate completed reconnect.
	_, clientNew := net.Pipe()
	defer func() { _ = clientNew.Close() }()
	c.conn = clientNew

	// Trigger EOF in the stale goroutine.
	_ = serverOld.Close()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("stale receive() goroutine did not exit in time")
	}

	time.Sleep(50 * time.Millisecond)
	assert.False(t, hookCalled.Load(), "OnConnectionLost must NOT be called by a stale goroutine")
}

// TestReceiveBufferResetOnReconnect verifies that Connect() calling
// receiveBuffer.Reset() clears any bytes left over from the previous session.
func TestReceiveBufferResetOnReconnect(t *testing.T) {
	c := newTestClient(ClientSettings{})

	c.receiveBuffer.Write([]byte{0x01, 0x02, 0x03, 0x04})
	assert.Equal(t, 4, c.receiveBuffer.Len(), "precondition: buffer should contain stale bytes")

	c.receiveBuffer.Reset()

	assert.Equal(t, 0, c.receiveBuffer.Len(), "receive buffer must be empty after Reset()")
}
