package main

import (
	"unicode"

	"golang.org/x/text/width"
)

func newGetWidth(ambiguousIsWide bool) func(rune) int {
	var cache = map[rune]int{}
	return func(r rune) (result int) {
		if w, ok := cache[r]; ok {
			return w
		}
		if unicode.IsPrint(r) {
			switch width.LookupRune(r).Kind() {
			case width.Neutral, width.EastAsianNarrow, width.EastAsianHalfwidth:
				result = 1
			case width.EastAsianWide, width.EastAsianFullwidth:
				result = 2
			case width.EastAsianAmbiguous:
				if ambiguousIsWide {
					result = 2
				} else {
					result = 1
				}
			}
		}
		cache[r] = result
		return
	}
}
