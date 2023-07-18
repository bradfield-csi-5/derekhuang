package diffbits

import (
	"testing"
)

var tests = []struct {
	s1  string
	s2  string
	out int
}{
	{"a", "a", 0},
	{"a", "b", 126},
	{"a", "c", 124},
	{"a", "d", 137},
}

func TestDiffBits(t *testing.T) {
	for _, tt := range tests {
		r := DiffBits(tt.s1, tt.s2)
		if r != tt.out {
			t.Errorf("GetDiffBits(%s, %s) = <%d> want <%d>", tt.s1, tt.s2, r, tt.out)
		}
	}
}
