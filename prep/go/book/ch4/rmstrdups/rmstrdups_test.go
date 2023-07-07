package rmstrdups

import (
	// "fmt"
	"testing"
)

var tests = []struct {
	in  []string
	out []string
}{
	{[]string{"one", "two", "three"}, []string{"one", "two", "three"}},
	{[]string{"one", "two", "two", "three"}, []string{"one", "two", "three"}},
	{[]string{"one", "two", "three", "three"}, []string{"one", "two", "three"}},
	{[]string{"one", "one", "one"}, []string{"one"}},
	{[]string{"one"}, []string{"one"}},
}

func TestRmstrdups(t *testing.T) {
	for _, tt := range tests {
		var cp = make([]string, len(tt.in))
		copy(cp, tt.in)
		r := Rmstrdups(tt.in)
		if !compareSlice(r, tt.out) {
			t.Errorf("Rmstrdups(%v) == <%v> want <%v>", cp, r, tt.out)
		}
	}
}

func compareSlice(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		// fmt.Printf("s1: %s s2: %s\n", s1, s2)
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			// fmt.Printf("s1[%d]: %s s2[%d]: %s\n", i, s1[i], i, s2[i])
			return false
		}
	}
	return true
}
