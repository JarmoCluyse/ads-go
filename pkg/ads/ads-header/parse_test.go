package adsheader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stripAndCheckLength(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		bytes   []byte
		want    []byte
		wantErr bool
		errMsg  string // expected error message substring
	}{
		{
			name:    "Valid - length matches data exactly",
			bytes:   []byte{0x04, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}, // length=4, data=4 bytes
			want:    []byte{0xAA, 0xBB, 0xCC, 0xDD},
			wantErr: false,
		},
		{
			name:    "Valid - zero length returns empty slice",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00}, // length=0, no data
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "Valid - large data payload",
			bytes:   []byte{0x0A, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A}, // length=10, data=10 bytes
			want:    []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A},
			wantErr: false,
		},
		{
			name:    "Invalid - extra data beyond declared length",
			bytes:   []byte{0x03, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE}, // length=3, data=5 bytes
			want:    nil,
			wantErr: true,
			errMsg:  "expected 3 bytes received len 5",
		},
		{
			name:    "Invalid - less than 4 bytes (no length field)",
			bytes:   []byte{0x00, 0x00, 0x00},
			want:    nil,
			wantErr: true,
			errMsg:  "invalid data length received",
		},
		{
			name:    "Invalid - empty input",
			bytes:   []byte{},
			want:    nil,
			wantErr: true,
			errMsg:  "invalid data length received",
		},
		{
			name:    "Invalid - data shorter than declared length",
			bytes:   []byte{0x0A, 0x00, 0x00, 0x00, 0xAA, 0xBB}, // length=10, data=2 bytes
			want:    nil,
			wantErr: true,
			errMsg:  "expected 10 bytes received len 2",
		},
		{
			name:    "Invalid - only length field, no data",
			bytes:   []byte{0x05, 0x00, 0x00, 0x00}, // length=5, data=0 bytes
			want:    nil,
			wantErr: true,
			errMsg:  "expected 5 bytes received len 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := stripAndCheckLength(tt.bytes)

			if tt.wantErr {
				assert.Error(t, gotErr, "stripAndCheckLength() should return an error")
				if tt.errMsg != "" {
					assert.Contains(t, gotErr.Error(), tt.errMsg, "Error message should contain expected substring")
				}
				assert.Nil(t, got, "stripAndCheckLength() should return nil when error occurs")
				return
			}

			assert.NoError(t, gotErr, "stripAndCheckLength() should not return an error")
			assert.Equal(t, tt.want, got, "stripAndCheckLength() returned incorrect result")
		})
	}
}

func TestCheckAdsHeader(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		bytes   []byte
		wantErr bool
		errMsg  string // expected error message substring
	}{
		{
			name:    "Valid - no error + valid length + data",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}, // error=0, length=4, data=4 bytes
			wantErr: false,
		},
		{
			name:    "Valid - no error + zero length",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // error=0, length=0, no data
			wantErr: false,
		},
		{
			name:    "Valid - no error + large payload",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A}, // error=0, length=10, data=10 bytes
			wantErr: false,
		},
		{
			name:    "Invalid - less than 8 bytes total",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00}, // only 7 bytes
			wantErr: true,
			errMsg:  "invalid header length received",
		},
		{
			name:    "Invalid - empty input",
			bytes:   []byte{},
			wantErr: true,
			errMsg:  "invalid header length received",
		},
		{
			name:    "Invalid - only 4 bytes (just error code)",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00}, // only error code, no length
			wantErr: true,
			errMsg:  "invalid header length received",
		},
		{
			name:    "Invalid - ADS error code 1",
			bytes:   []byte{0x01, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}, // error=1
			wantErr: true,
			errMsg:  "ads error",
		},
		{
			name:    "Invalid - ADS error code 1792",
			bytes:   []byte{0x00, 0x07, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}, // error=1792
			wantErr: true,
			errMsg:  "ads error",
		},
		{
			name:    "Invalid - valid error but data length mismatch",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00, 0xAA, 0xBB}, // error=0, length=10, data=2 bytes
			wantErr: true,
			errMsg:  "expected 10 bytes received len 2",
		},
		{
			name:    "Invalid - valid error but no data after length",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00}, // error=0, length=5, data=0 bytes
			wantErr: true,
			errMsg:  "expected 5 bytes received len 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := CheckAdsHeader(tt.bytes)

			if tt.wantErr {
				assert.Error(t, gotErr, "CheckAdsHeader() should return an error")
				if tt.errMsg != "" {
					assert.Contains(t, gotErr.Error(), tt.errMsg, "Error message should contain expected substring")
				}
				return
			}

			assert.NoError(t, gotErr, "CheckAdsHeader() should not return an error")
		})
	}
}

func TestStripAdsHeader(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		bytes   []byte
		want    []byte
		wantErr bool
		errMsg  string // expected error message substring
	}{
		{
			name:    "Valid - strips 8 bytes (4 error + 4 length) and returns data",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}, // error=0, length=4, data=4 bytes
			want:    []byte{0xAA, 0xBB, 0xCC, 0xDD},
			wantErr: false,
		},
		{
			name:    "Valid - no error + zero length returns empty slice",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // error=0, length=0, no data
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "Valid - large payload correctly stripped",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A}, // error=0, length=10, data=10 bytes
			want:    []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A},
			wantErr: false,
		},
		{
			name:    "Invalid - extra data beyond declared length",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE}, // error=0, length=3, data=5 bytes
			want:    nil,
			wantErr: true,
			errMsg:  "expected 3 bytes received len 5",
		},
		{
			name:    "Invalid - less than 8 bytes",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00}, // only 7 bytes
			want:    nil,
			wantErr: true,
			errMsg:  "invalid header length received",
		},
		{
			name:    "Invalid - empty input",
			bytes:   []byte{},
			want:    nil,
			wantErr: true,
			errMsg:  "invalid header length received",
		},
		{
			name:    "Invalid - ADS error code 1 present",
			bytes:   []byte{0x01, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}, // error=1
			want:    nil,
			wantErr: true,
			errMsg:  "ads error",
		},
		{
			name:    "Invalid - ADS error code 1792 present",
			bytes:   []byte{0x00, 0x07, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}, // error=1792
			want:    nil,
			wantErr: true,
			errMsg:  "ads error",
		},
		{
			name:    "Invalid - declared length 10, actual data 5",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE}, // error=0, length=10, data=5 bytes
			want:    nil,
			wantErr: true,
			errMsg:  "expected 10 bytes received len 5",
		},
		{
			name:    "Invalid - length field says 5 but no data",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00}, // error=0, length=5, data=0 bytes
			want:    nil,
			wantErr: true,
			errMsg:  "expected 5 bytes received len 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := StripAdsHeader(tt.bytes)

			if tt.wantErr {
				assert.Error(t, gotErr, "StripAdsHeader() should return an error")
				if tt.errMsg != "" {
					assert.Contains(t, gotErr.Error(), tt.errMsg, "Error message should contain expected substring")
				}
				assert.Nil(t, got, "StripAdsHeader() should return nil when error occurs")
				return
			}

			assert.NoError(t, gotErr, "StripAdsHeader() should not return an error")
			assert.Equal(t, tt.want, got, "StripAdsHeader() returned incorrect result")
		})
	}
}
