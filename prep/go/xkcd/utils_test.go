package main

import (
	"testing"
)

var normalizeTests = []struct {
	in  []string
	out StrSet
}{
	{
		in:  []string{"Good Cop, Dadaist Cop"},
		out: StrSet{"good": true, "cop": true, "dadaist": true},
	},
	{
		in:  []string{"Pi Equals", "Pi = 3.141592653589793helpimtrappedinauniversefactory7108914..."},
		out: StrSet{"equals": true, "helpimtrappedinauniversefactory": true},
	},
	{
		in:  []string{"foo\nbar\nbaz"},
		out: StrSet{"foo": true, "bar": true, "baz": true},
	},
	{
		in:  []string{"10th Annual Symposium on Formal Languages\n\u003c\u003cCRASH\u003e\u003e\n\n"},
		out: StrSet{"annual": true, "symposium": true, "formal": true, "languages": true, "crash": true},
	},
}

func TestNormalize(t *testing.T) {
	for _, tt := range normalizeTests {
		r := normalize(tt.in...)
		if !equal(r, tt.out) {
			t.Errorf("normalize(%v) == <%v> want <%v>", tt.in, r, tt.out)
		}
	}
}

func equal(s1 StrSet, s2 StrSet) bool {
	if len(s1) != len(s2) {
		return false
	}
	for word := range s1 {
		if _, exists := s2[word]; !exists {
			return false
		}
	}
	return true
}
