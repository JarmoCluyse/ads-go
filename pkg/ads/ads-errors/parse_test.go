package adserrors_test

import (
	"testing"

	"github.com/jarmocluyse/ads-go/pkg/ads/ads-errors"
	"github.com/stretchr/testify/assert"
)

func TestCheckAdsError(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		bytes   []byte
		wantErr bool
		errMsg  string // expected error message substring
	}{
		{
			name:    "No error - error code 0",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00}, // 0 in little endian
			wantErr: false,
		},
		{
			name:    "ADS error - error code 1",
			bytes:   []byte{0x01, 0x00, 0x00, 0x00}, // 1 in little endian
			wantErr: true,
			errMsg:  "ads error",
		},
		{
			name:    "ADS error - error code 1792",
			bytes:   []byte{0x00, 0x07, 0x00, 0x00}, // 1792 in little endian
			wantErr: true,
			errMsg:  "ads error",
		},
		{
			name:    "Invalid length - too short",
			bytes:   []byte{0x00, 0x00, 0x00},
			wantErr: true,
			errMsg:  "invalid length received",
		},
		{
			name:    "Invalid length - too long",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x00},
			wantErr: true,
			errMsg:  "invalid length received",
		},
		{
			name:    "Invalid length - empty",
			bytes:   []byte{},
			wantErr: true,
			errMsg:  "invalid length received",
		},
		{
			name:    "Unknown error code",
			bytes:   []byte{0xFF, 0xFF, 0xFF, 0xFF}, // 4294967295 in little endian
			wantErr: true,
			errMsg:  "ads error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := adserrors.CheckAdsError(tt.bytes)
			if tt.wantErr {
				assert.Error(t, gotErr, "CheckAdsError() should return an error")
				if tt.errMsg != "" {
					assert.Contains(t, gotErr.Error(), tt.errMsg, "Error message should contain expected substring")
				}
			} else {
				assert.NoError(t, gotErr, "CheckAdsError() should not return an error")
			}
		})
	}
}

func TestStripAdsError(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		bytes   []byte
		want    []byte
		wantErr bool
		errMsg  string // expected error message substring
	}{
		{
			name:    "No error - strips first 4 bytes",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04},
			want:    []byte{0x01, 0x02, 0x03, 0x04},
			wantErr: false,
		},
		{
			name:    "No error - exactly 4 bytes returns empty slice",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00},
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "ADS error - error code 1",
			bytes:   []byte{0x01, 0x00, 0x00, 0x00, 0x01, 0x02},
			want:    nil,
			wantErr: true,
			errMsg:  "ads error",
		},
		{
			name:    "ADS error - error code 1792",
			bytes:   []byte{0x00, 0x07, 0x00, 0x00, 0xFF, 0xFF},
			want:    nil,
			wantErr: true,
			errMsg:  "ads error",
		},
		{
			name:    "Invalid length - too short",
			bytes:   []byte{0x00, 0x00, 0x00},
			want:    nil,
			wantErr: true,
			errMsg:  "invalid length received",
		},
		{
			name:    "Invalid length - empty",
			bytes:   []byte{},
			want:    nil,
			wantErr: true,
			errMsg:  "invalid length received",
		},
		{
			name:    "No error - strips first 4 bytes from longer payload",
			bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
			want:    []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
			wantErr: false,
		},
		{
			name:    "Unknown error code",
			bytes:   []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x01, 0x02},
			want:    nil,
			wantErr: true,
			errMsg:  "ads error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := adserrors.StripAdsError(tt.bytes)
			if tt.wantErr {
				assert.Error(t, gotErr, "StripAdsError() should return an error")
				if tt.errMsg != "" {
					assert.Contains(t, gotErr.Error(), tt.errMsg, "Error message should contain expected substring")
				}
				assert.Nil(t, got, "StripAdsError() should return nil when error occurs")
			} else {
				assert.NoError(t, gotErr, "StripAdsError() should not return an error")
				assert.Equal(t, tt.want, got, "StripAdsError() returned incorrect result")
			}
		})
	}
}
