package ads

import "errors"

// ErrNotConnected is returned when an operation is attempted without an
// active connection. Callers can match on this with errors.Is.
var ErrNotConnected = errors.New("not connected")
