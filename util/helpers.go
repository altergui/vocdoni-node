package util

import (
	"encoding/hex"
	"regexp"
	"strconv"
)

var validHexRegex = regexp.MustCompile("^([0-9a-fA-F])+$")

// IsHex checks if the given string contains only valid hex symbols
func IsHex(str string) bool { return validHexRegex.MatchString(str) }

// IsHexEncodedStringWithLength checks if the given string contains only valid hex symbols and have the desired length
func IsHexEncodedStringWithLength(str string, length int) bool {
	return hex.DecodedLen(len(str)) == length && IsHex(str)
}

func Hex2int64(s string) int64 {
	// base 16 for hexadecimal
	result, err := strconv.ParseUint(TrimHex(s), 16, 64)
	if err != nil {
		panic(err)
	}
	return int64(result)
}

func TrimHex(s string) string {
	if len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		return s[2:]
	}
	return s
}
