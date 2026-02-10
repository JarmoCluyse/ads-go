package adserrors

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrAdsError      = errors.New("ads error")
	ErrInvalidLength = errors.New("invalid length received")
)

// parse in bytes and check the ads error
func CheckAdsError(bytes []byte) error {
	if len(bytes) != 4 {
		return fmt.Errorf("%w : received len %d", ErrInvalidLength, len(bytes))
	}
	errorCode := binary.LittleEndian.Uint32(bytes)
	if errorCode != 0 {
		code := ErrorCodeToString(errorCode)
		return fmt.Errorf("%w: %s", ErrAdsError, code)
	}
	return nil

}

// strip the first 4 bytes and check the ads error
// from the perspective that the first 4 bytes are the ads error
func StripAdsError(bytes []byte) ([]byte, error) {
	if len(bytes) < 4 {
		return nil, fmt.Errorf("%w : received len %d", ErrInvalidLength, len(bytes))
	}
	err := CheckAdsError(bytes[0:4])
	if err != nil {
		return nil, err
	}
	return bytes[4:], nil

}
