package main

import (
	"unicode"
)

func IsValidIdentifier(s string) bool {
	if len(s) > 64 || len(s) == 0 {
		return false
	}
	r0 := rune(s[0])
	if !(unicode.IsLetter(r0) && (r0 <= unicode.MaxASCII) || r0 == '_') {
		return false
	}
	for i := 1; i < len(s); i++ {
		r := rune(s[i])
		if !(unicode.IsLetter(r) && (r <= unicode.MaxASCII) || r == '_' || unicode.IsDigit(r)) {
			return false
		}

	}
	return true
}
