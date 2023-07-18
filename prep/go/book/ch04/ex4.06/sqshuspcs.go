/*
Exercise 4.6: Write an in-place function that squashes each run of adjacent Unicode spaces (see unicode.IsSpace) in a UTF-8-encoded []byte slice into a single ASCII space.
*/

package sqshuspcs

import (
	"unicode"
)

func SquashUnicodeSpaces(bytes []byte) []byte {
	var i int = 0
	var c int = 0
	var skipping bool = false
	for _, b := range bytes {
		if unicode.IsSpace(rune(b)) {
			c++
			if c >= 1 {
				if skipping {
					continue
				}
				bytes[i] = ' '
				i++
				skipping = true
			}
		} else {
			bytes[i] = b
			i++
			c = 0
			skipping = false
		}
	}
	return bytes[:i]
}
