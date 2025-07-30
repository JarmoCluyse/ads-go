package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	errInvalidAmsNetId     = errors.New("invalid AmsNetId")
	errInvalidAmsNetIdPart = errors.New("invalid part in AmsNetId")
)

// Converts given AmsAddress struct to string "amsNetId:port"
func AmsAddressToString(amsNetId string, adsPort uint16) string {
	return amsNetId + ":" + strconv.Itoa(int(adsPort))
}

// AmsNetIDStrToByteArray converts an AmsNetId string to a byte array.
func AmsNetIdStrToByteArray(s string) ([]byte, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 6 {
		return nil, fmt.Errorf("%w: %s", errInvalidAmsNetId, s)
	}
	bytes := make([]byte, 6)
	for i, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", errInvalidAmsNetIdPart, part)
		}
		bytes[i] = byte(val)
	}
	return bytes, nil
}

// ByteArrayToAmsNetIdStr converts a byte array to an AmsNetId string.
func ByteArrayToAmsNetIdStr(b []byte) string {
	parts := make([]string, len(b))
	for i, byte := range b {
		parts[i] = fmt.Sprintf("%d", byte)
	}
	return strings.Join(parts, ".")
}
