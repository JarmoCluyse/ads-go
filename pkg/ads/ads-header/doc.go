// Package adsheader provides parsing and validation for ADS response headers.
//
// This package handles the common ADS response pattern consisting of an error code
// and a length field, followed by the actual data payload. ADS responses from Beckhoff
// TwinCAT devices follow this 8-byte header format before the payload data.
//
// # ADS Response Header Format
//
// ADS responses follow this binary format:
//   - Bytes 0-3: Error code (uint32) in little-endian byte order
//   - Bytes 4-7: Data length (uint32) in little-endian byte order
//   - Bytes 8+: Actual data payload (length specified by bytes 4-7)
//
// The error code (first 4 bytes) indicates success (0) or a specific ADS error.
// The length field (next 4 bytes) declares how many bytes of valid payload data follow.
//
// # Core Functions
//
// CheckAdsHeader validates an ADS response header without extracting data:
//
//	response := []byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}
//	if err := adsheader.CheckAdsHeader(response); err != nil {
//	    log.Printf("Invalid ADS header: %v", err)
//	}
//
// StripAdsHeader validates and removes the 8-byte header, returning only the payload:
//
//	response := []byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}
//	payload, err := adsheader.StripAdsHeader(response)
//	if err != nil {
//	    return err
//	}
//	// payload now contains: []byte{0xAA, 0xBB, 0xCC, 0xDD}
//
// # Exact Length Validation
//
// The StripAdsHeader function performs exact length validation. The actual data
// length must exactly match the declared length in the header. If the response
// contains extra bytes or fewer bytes than declared, an error is returned:
//
//	// Response with 3-byte length but 5 bytes of data (INVALID - too many bytes)
//	response := []byte{0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE}
//	payload, err := adsheader.StripAdsHeader(response)
//	// err = ErrInvalidDataLength (data length mismatch)
//
// # Error Types
//
// The package exports sentinel errors for error inspection:
//   - ErrInvalidHeaderLength: Response is less than 8 bytes (too short for header)
//   - ErrInvalidDataLength: Actual data doesn't match the declared length
//
// Use errors.Is() for error type checking:
//
//	if errors.Is(err, adsheader.ErrInvalidHeaderLength) {
//	    // Handle invalid header length
//	}
//
// # Integration with ads-errors
//
// This package uses the ads-errors package to validate the error code portion
// of the header. If a non-zero ADS error code is detected, the error will be
// of type adserrors.ErrAdsError:
//
//	payload, err := adsheader.StripAdsHeader(response)
//	if errors.Is(err, adserrors.ErrAdsError) {
//	    log.Printf("ADS protocol error: %v", err)
//	}
//
// # Common Use Cases
//
// Most ADS Read commands return responses with this header format:
//   - ADSCommandRead (read data from device)
//   - ADSCommandReadWrite (combined read/write operation)
//   - ADSCommandReadInfo (read device information)
//
// Example - Processing a read response:
//
//	// Send ADS read command
//	response, err := client.send(readRequest)
//	if err != nil {
//	    return nil, err
//	}
//
//	// Strip header and get payload
//	data, err := adsheader.StripAdsHeader(response)
//	if err != nil {
//	    return nil, fmt.Errorf("invalid ADS response: %w", err)
//	}
//
//	// Process data...
//	return parseData(data), nil
package adsheader
