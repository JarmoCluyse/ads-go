package utils

import (
	"testing"
)

func TestAmsNetIdStrToByteArray(t *testing.T) {
	tests := []struct {
		input    string
		expect   []byte
		wantErr  bool
	}{
		{"5.33.160.43.1.1", []byte{5, 33, 160, 43, 1, 1}, false},
		{"192.168.0.1.2.3", []byte{192, 168, 0, 1, 2, 3}, false},
		{"1.2.3.4.5.6", []byte{1, 2, 3, 4, 5, 6}, false},
		{"1.2.3.4.5", nil, true}, // too short
		{"1.2.3.4.5.6.7", nil, true}, // too long
		{"a.b.c.d.e.f", nil, true}, // non-numeric
		{"255.255.255.255.255.255", []byte{255, 255, 255, 255, 255, 255}, false},
	}

	for _, test := range tests {
		got, err := AmsNetIdStrToByteArray(test.input)
		if test.wantErr {
			if err == nil {
				t.Errorf("AmsNetIdStrToByteArray(%q): expected error, got none", test.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("AmsNetIdStrToByteArray(%q): unexpected error: %v", test.input, err)
			continue
		}
		for i := range test.expect {
			if got[i] != test.expect[i] {
				t.Errorf("AmsNetIdStrToByteArray(%q): got %v, want %v", test.input, got, test.expect)
				break
			}
		}
	}
}

func TestByteArrayToAmsNetIdStr(t *testing.T) {
	tests := []struct {
		input []byte
		expect string
	}{
		{[]byte{5, 33, 160, 43, 1, 1}, "5.33.160.43.1.1"},
		{[]byte{192, 168, 0, 1, 2, 3}, "192.168.0.1.2.3"},
		{[]byte{1, 2, 3, 4, 5, 6}, "1.2.3.4.5.6"},
		{[]byte{255, 255, 255, 255, 255, 255}, "255.255.255.255.255.255"},
	}
	for _, test := range tests {
		got := ByteArrayToAmsNetIdStr(test.input)
		if got != test.expect {
			t.Errorf("ByteArrayToAmsNetIdStr(%v): got %q, want %q", test.input, got, test.expect)
		}
	}
}

func TestAmsAddressToString(t *testing.T) {
	tests := []struct {
		amsNetId string
		adsPort  uint16
		expect    string
	}{
		{"5.33.160.43.1.1", 851, "5.33.160.43.1.1:851"},
		{"192.168.0.1.2.3", 48898, "192.168.0.1.2.3:48898"},
	}
	for _, test := range tests {
		got := AmsAddressToString(test.amsNetId, test.adsPort)
		if got != test.expect {
			t.Errorf("AmsAddressToString(%q, %d): got %q, want %q", test.amsNetId, test.adsPort, got, test.expect)
		}
	}
}

