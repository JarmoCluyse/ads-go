package adsheader

import (
	"encoding/binary"
	"errors"
	"fmt"

	adserrors "github.com/jarmocluyse/ads-go/pkg/ads/ads-errors"
)

var (
	ErrInvalidHeaderLength = errors.New("invalid header length received")
	ErrInvalidDataLength   = errors.New("invalid data length received")
)

// check the described length in the header and the received length
func stripAndCheckLength(bytes []byte) ([]byte, error) {
	if len(bytes) < 4 {
		return nil, fmt.Errorf("%w : expected at least 4 bytes received len %d", ErrInvalidDataLength, len(bytes))
	}
	dataLen := binary.LittleEndian.Uint32(bytes[0:4])
	data := bytes[4:]
	if len(data) != int(dataLen) {
		return nil, fmt.Errorf("%w : expected %d bytes received len %d", ErrInvalidDataLength, dataLen, len(data))
	}
	return data, nil

}

// check the ads header, length + ads error
func CheckAdsHeader(bytes []byte) error {
	if len(bytes) < 8 {
		return fmt.Errorf("%w : received len %d", ErrInvalidHeaderLength, len(bytes))
	}
	data, err := adserrors.StripAdsError(bytes)
	if err != nil {
		return err
	}
	_, err = stripAndCheckLength(data)
	if err != nil {
		return err
	}
	return nil
}

// check the ads header, length + ads error and strip it
func StripAdsHeader(bytes []byte) ([]byte, error) {
	if len(bytes) < 8 {
		return nil, fmt.Errorf("%w : received len %d", ErrInvalidHeaderLength, len(bytes))
	}
	data, err := adserrors.StripAdsError(bytes)
	if err != nil {
		return nil, err
	}
	data, err = stripAndCheckLength(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
