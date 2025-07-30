package utils

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unicode/utf16"
)

func TestTrimPlcString(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"foo\x00bar", "foo"},
		{"foobar", "foobar"},
		{"hello\x00", "hello"},
		{"\x00", ""},
		{"", ""},
	}
	for _, c := range cases {
		result := TrimPlcString(c.input)
		if result != c.expected {
			t.Errorf("TrimPlcString(%q) = %q, want %q", c.input, result, c.expected)
		}
	}
}

func TestDecodePlcStringBuffer(t *testing.T) {
	cases := []struct {
		input    []byte
		expected string
	}{
		// TODO: do I still need this, since I always want to look at length instead
		{[]byte("hello\x00world"), "hello"},
		{[]byte("hello world"), "hello world"},
		{[]byte{'f', 'o', 'o', 0}, "foo"},
		{[]byte{0}, ""},
		{[]byte{}, ""},
	}
	for _, c := range cases {
		result := DecodePlcStringBuffer(c.input)
		if result != c.expected {
			t.Errorf("DecodePlcStringBuffer(%v) = %q, want %q", c.input, result, c.expected)
		}
	}
}

func TestEncodeStringToPlcStringBuffer(t *testing.T) {
	cases := []struct {
		input    string
		expected []byte
	}{
		{"hello", []byte("hello")},
		{"", []byte{}},
		{"foo\x00bar", []byte{'f', 'o', 'o', 0, 'b', 'a', 'r'}},
	}
	for _, c := range cases {
		result := EncodeStringToPlcStringBuffer(c.input)
		if !bytes.Equal(result, c.expected) {
			t.Errorf("EncodeStringToPlcStringBuffer(%q) = %v, want %v", c.input, result, c.expected)
		}
	}
}

func TestDecodePlcWstringBuffer(t *testing.T) {
	cases := []struct {
		input    []byte
		expected string
	}{
		// "h\x00e\x00l\x00l\x00o\x00\x00\x00" is "hello\x00" in UTF-16LE
		{[]byte{'h', 0, 'e', 0, 'l', 0, 'l', 0, 'o', 0, 0, 0}, "hello"},
		{[]byte{'f', 0, 'o', 0, 'o', 0, 0, 0, 'b', 0, 'a', 0, 'r', 0}, "foo"},
		{[]byte{0, 0}, ""},
		{[]byte{}, ""},
	}
	for _, c := range cases {
		result := DecodePlcWstringBuffer(c.input)
		if result != c.expected {
			t.Errorf("DecodePlcWstringBuffer(%v) = %q, want %q", c.input, result, c.expected)
		}
	}
}

func TestEncodeStringToPlcWstringBuffer(t *testing.T) {
	cases := []struct {
		input    string
		expected []byte
	}{
		{"hello", encodeUtf16LE("hello")},
		{"", encodeUtf16LE("")},
		{"foo", encodeUtf16LE("foo")},
	}
	for _, c := range cases {
		result := EncodeStringToPlcWstringBuffer(c.input)
		if !bytes.Equal(result, c.expected) {
			t.Errorf("EncodeStringToPlcWstringBuffer(%q) = %v, want %v", c.input, result, c.expected)
		}
	}
}

// Helper to encode to UTF-16LE (like PLC WSTRING)
func encodeUtf16LE(s string) []byte {
	u16 := utf16.Encode([]rune(s))
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, u16)
	return buf.Bytes()
}
