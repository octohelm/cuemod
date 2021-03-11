package core

import (
	"strings"
	"unicode"
)

func SafeIdentifierFromImportPath(s string) string {
	parts := strings.Split(s, "/")

	lastIdx := len(parts)

	//
	for {
		lastIdx = lastIdx - 1

		if lastIdx < 0 {
			continue
		}

		last := parts[lastIdx]

		// use parent when v1 v2
		if len(last) > 2 && last[0] == 'v' && unicode.IsNumber(rune(last[1])) {
			continue
		}

		// use parent when number only
		if len(last) > 0 && unicode.IsNumber(rune(last[0])) {
			continue
		}

		runes := []rune(last)

		for i, r := range runes {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				runes[i] = '_'
			}
		}

		return string(runes)
	}
}
