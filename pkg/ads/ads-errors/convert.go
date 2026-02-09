package adserrors

import "fmt"

// convert the error code to string
func ErrorCodeToString(code uint32) string {
	adsError, ok := ADSError[code]
	if !ok {
		return fmt.Sprintf("Unknown ads error %d", code)
	}
	return adsError
}
