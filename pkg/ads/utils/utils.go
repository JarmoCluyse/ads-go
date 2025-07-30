package utils

import (
	"bytes"
	"encoding/binary"
	"strings"
	"unicode/utf16"
	// "golang.org/x/text/encoding/charmap"
)

// Trims given PLC string until '\\0' is found (removes empty bytes from the end)
func TrimPlcString(str string) string {
	if idx := strings.IndexRune(str, 0); idx >= 0 {
		return str[:idx]
	}
	return str
}

// Decodes provided []byte to PLC STRING
func DecodePlcStringBuffer(data []byte) string {
	var decoded string
	decoded = string(data)
	// this will be needed for older versions of twincat
	// decoded, _ = charmap.Windows1252.NewDecoder().String(string(data))
	return TrimPlcString(decoded)
}

// Encodes provided string to []byte as PLC STRING
func EncodeStringToPlcStringBuffer(str string) []byte {
	return []byte(str)
	// this will be needed for older versions of twincat
	// encoded, _ := charmap.Windows1252.NewEncoder().String(str)
	// return []byte(encoded)
}

// Decodes provided []byte to PLC WSTRING using ucs2 encoding, also trims zeroes
func DecodePlcWstringBuffer(data []byte) string {
	u16 := make([]uint16, len(data)/2)
	_ = binary.Read(bytes.NewReader(data), binary.LittleEndian, &u16)
	str := string(utf16.Decode(u16))
	return TrimPlcString(str)
}

// Encodes provided string to []byte as PLC WSTRING using ucs2 encoding
func EncodeStringToPlcWstringBuffer(str string) []byte {
	u16 := utf16.Encode([]rune(str))
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, u16)
	return buf.Bytes()
}
