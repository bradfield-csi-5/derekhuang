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
		if !compStrSet(r, tt.out) {
			t.Errorf("normalize(%v) == <%v> want <%v>", tt.in, r, tt.out)
		}
	}
}

func compStrSet(s1, s2 StrSet) bool {
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

var intersectionTests = []struct {
	s1  IntSet
	s2  IntSet
	out IntSet
}{
	{s1: IntSet{1: true, 2: true, 3: true}, s2: IntSet{}, out: IntSet{}},
	{s1: IntSet{}, s2: IntSet{1: true, 2: true, 3: true}, out: IntSet{}},
	{
		s1:  IntSet{1: true, 2: true, 3: true, 4: true, 5: true},
		s2:  IntSet{1: true, 2: true, 3: true},
		out: IntSet{1: true, 2: true, 3: true},
	},
	{
		s1:  IntSet{1: true, 2: true, 3: true},
		s2:  IntSet{1: true, 2: true, 3: true, 4: true, 5: true},
		out: IntSet{1: true, 2: true, 3: true},
	},
}

func TestSetIntersection(t *testing.T) {
	for _, tt := range intersectionTests {
		tt.s1.Intersection(tt.s2)
		if !compIntSet(tt.s1, tt.out) {
			t.Errorf("s1.Intersection(s2) == <%v> want <%v>", tt.s1, tt.out)
		}
	}
}

var unionTests = []struct {
	s1  IntSet
	s2  IntSet
	out IntSet
}{
	{s1: IntSet{1: true, 2: true, 3: true}, s2: IntSet{}, out: IntSet{1: true, 2: true, 3: true}},
	{s1: IntSet{}, s2: IntSet{1: true, 2: true, 3: true}, out: IntSet{1: true, 2: true, 3: true}},
	{
		s1:  IntSet{1: true, 2: true, 3: true, 4: true, 5: true},
		s2:  IntSet{1: true, 2: true, 3: true},
		out: IntSet{1: true, 2: true, 3: true, 4: true, 5: true},
	},
	{
		s1:  IntSet{1: true, 2: true, 3: true},
		s2:  IntSet{1: true, 2: true, 3: true, 4: true, 5: true},
		out: IntSet{1: true, 2: true, 3: true, 4: true, 5: true},
	},
}

func TestSetUnion(t *testing.T) {
	for _, tt := range unionTests {
		tt.s1.Union(tt.s2)
		if !compIntSet(tt.s1, tt.out) {
			t.Errorf("s1.Union(s2) == <%v> want <%v>", tt.s1, tt.out)
		}
	}
}

func compIntSet(s1, s2 IntSet) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if _, exists := s2[i]; !exists {
			return false
		}
	}
	return true
}
