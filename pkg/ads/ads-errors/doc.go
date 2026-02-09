// Package adserrors provides error handling for the ADS protocol.
//
// # Error Code Format
//
// ADS error codes follow this binary format:
//   - 4 bytes (uint32) in little-endian byte order
//   - Error code 0 indicates success (no error)
//   - Non-zero values indicate specific error conditions
//   - Error codes are defined by Beckhoff InfoSys documentation
//
// # Core Functions
//
// CheckAdsError validates a 4-byte slice and returns an error if a non-zero error code is found:
//
//	errorBytes := []byte{0x00, 0x07, 0x00, 0x00} // Error code 1792
//	if err := adserrors.CheckAdsError(errorBytes); err != nil {
//	    log.Printf("ADS error: %v", err) // Output: ads error: General device error
//	}
//
// StripAdsError validates and removes the first 4 bytes from a response payload:
//
//	response := []byte{0x00, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}
//	data, err := adserrors.StripAdsError(response)
//	if err != nil {
//	    return err
//	}
//	// data now contains: []byte{0xAA, 0xBB, 0xCC, 0xDD}
//
// ErrorCodeToString converts numeric error codes to human-readable messages:
//
//	message := adserrors.ErrorCodeToString(1792)
//	fmt.Println(message) // Output: General device error
//
// # Error Types
//
// The package exports two sentinel errors for error inspection:
//   - ErrAdsError: Returned when a non-zero ADS error code is detected
//   - ErrInvalidLength: Returned when the byte slice length is invalid
//
// Use errors.Is() for error type checking:
//
//	if errors.Is(err, adserrors.ErrAdsError) {
//	    // Handle ADS protocol error
//	}
//
// # Error Code Reference
//
// The ADSError map contains the complete mapping of error codes to messages.
// Error codes are sourced from Beckhoff InfoSys documentation and include:
//   - 0-28: General ADS errors (internal errors, communication errors)
//   - 1280-1290: Router errors
//   - 1792-1877: Device and service errors
//   - 4096-4125: TwinCAT real-time system errors
//
// For unknown error codes, ErrorCodeToString returns a formatted string with the numeric code.
package adserrors
