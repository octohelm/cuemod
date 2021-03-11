package jsonnet

import (
	"bytes"
	"regexp"
	"text/scanner"
)

var reDocUtilImport = regexp.MustCompile(`local d = [(]?import ['"]doc-util/main\.libsonnet['"][)]?([,;]\n)?`)
var reDocField = regexp.MustCompile(`([ \t]+)?'#(.+)?'[+:]?: 'ignore'([,;]\n)?`)

func dropDocUtil(data []byte) []byte {
	data = reDocUtilImport.ReplaceAll(data, []byte{})

	s := scanner.Scanner{}
	s.Init(bytes.NewReader(data))
	s.Error = func(s *scanner.Scanner, msg string) {
	}

	ranges := make([][2]int, 0)
	scopeDepth := 0
	from := -1
	fromScopeDepth := 0

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		c := s.TokenText()

		switch c {
		case "(":
			scopeDepth++
		case ")":
			scopeDepth--
			if from > 0 && scopeDepth == fromScopeDepth {
				ranges = append(ranges, [2]int{from, s.Position.Offset + 1})
				from = -1
			}
		case "d":
			offset := s.Position.Offset

			_ = s.Scan()
			c = s.TokenText()

			// ignore nested
			if c == "." && from < 0 {
				from = offset
				fromScopeDepth = scopeDepth
			}
		}
	}

	results := make([]byte, 0)

	start := 0

	for _, r := range ranges {
		results = append(results, data[start:r[0]]...)
		start = r[1]

		results = append(results, []byte("'ignore'")...)
	}

	results = append(results, data[start:]...)

	return reDocField.ReplaceAll(results, []byte{})
}
