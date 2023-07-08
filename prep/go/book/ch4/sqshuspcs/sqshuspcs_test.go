package sqshuspcs

import (
	"testing"
)

var tests = []struct {
	in  []byte
	out []byte
}{
	{[]byte{}, []byte{}},
	{[]byte{'a', 'b', 'c'}, []byte{'a', 'b', 'c'}},
	{[]byte{'a', 'b', ' ', '\t', '\f', 'c'}, []byte{'a', 'b', ' ', 'c'}},
	{[]byte{'\n', 'a', '\t', '\t', '\v', 'b', 'c'}, []byte{' ', 'a', ' ', 'b', 'c'}},
	{[]byte{'\r', ' ', '\r', ' ', ' ', 'a', 'b', 'c'}, []byte{' ', 'a', 'b', 'c'}},
	{[]byte{'\r', ' ', '\r', ' ', '\f', ' '}, []byte{' '}},
}

func TestSquashUnicodeSpaces(t *testing.T) {
	for _, tt := range tests {
		var cp = make([]byte, len(tt.in))
		copy(cp, tt.in)
		r := SquashUnicodeSpaces(tt.in)
		if !compareSlice(r, tt.out) {
			t.Errorf("SquashUnicodeSpaces(%q) == <%q> want <%q>", cp, r, tt.out)
		}
	}
}

func compareSlice(s1 []byte, s2 []byte) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
