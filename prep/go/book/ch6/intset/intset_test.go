package intset

import (
	"testing"
)

var nums = []int{1, 2, 3}

func compareSlice(s1 []int, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func create(ints []int) IntSet {
	var s IntSet
	s.AddAll(ints...)
	return s
}

func TestLen(t *testing.T) {
	var tests = []struct {
		set  IntSet
		want int
	}{
		{set: create([]int{}), want: 0},
		{set: create([]int{1}), want: 1},
		{set: create([]int{1, 2}), want: 2},
		{set: create([]int{1, 2, 3}), want: 3},
	}
	for _, test := range tests {
		if test.set.Len() != test.want {
			t.Errorf("(%s).Len() == <%d> want <%d>", test.set.String(), test.set.Len(), test.want)
		}
	}
}

func TestRemove(t *testing.T) {
	var set IntSet
	set.AddAll(nums...)
	set.Remove(2)
	if !set.Has(1) {
		t.Error("set is missing 1")
	}
	if set.Has(2) {
		t.Error("set still has 2")
	}
	if !set.Has(1) {
		t.Error("set is missing 3")
	}
}

func TestClear(t *testing.T) {
	var set IntSet
	set.AddAll(nums...)
	set.Clear()
	for _, n := range nums {
		if set.Has(n) {
			t.Errorf("set.Has(%d) after Clear()", n)
		}
	}
}

func TestCopy(t *testing.T) {
	var set IntSet
	set.AddAll(nums...)
	var copy = set.Copy()
	if &set == copy {
		t.Error("set == copy after Copy()")
	}
}

func TestElems(t *testing.T) {
	var set IntSet
	set.AddAll(nums...)
	if !compareSlice(set.Elems(), nums) {
		t.Errorf("set.Elems() == <%v> want <%v>", set.Elems(), nums)
	}
}

func TestAddAll(t *testing.T) {
	var tests = []struct {
		in []int
	}{
		{[]int{}},
		{[]int{1}},
		{[]int{1, 2}},
		{[]int{1, 2, 3}},
	}

	for _, test := range tests {
		var set IntSet
		set.AddAll(test.in...)
		for _, n := range test.in {
			if !set.Has(n) {
				t.Errorf("Has(%d) failed after AddAll(%v...)", n, test.in)
			}
		}
	}
}

func TestUnionWith(t *testing.T) {
	var s1 IntSet
	var s2 IntSet

	// Add ints to s2
	s2.AddAll(nums...)

	// Combine the sets
	s1.UnionWith(&s2)

	// Assert s1 has the ints added to s2
	for _, n := range nums {
		if !s1.Has(n) {
			t.Errorf("%s.Has(%d) failed after UnionWith %s", s1.String(), n, s2.String())
		}
	}
}

func TestIntersectWith(t *testing.T) {
	var tests = []struct {
		s1      IntSet
		s2      IntSet
		want    []int
		wantLen int
	}{
		{s1: create([]int{}), s2: create(nums), want: []int{}, wantLen: 0},
		{s1: create(nums), s2: create([]int{}), want: []int{}, wantLen: 0},
		{s1: create(nums), s2: create([]int{2, 3, 4}), want: []int{2, 3}, wantLen: 2},
	}
	for _, test := range tests {
		test.s1.IntersectWith(&test.s2)

		// This loop won't run for empty intersections so the length is asserted below
		for _, n := range test.want {
			if !test.s1.Has(n) {
				t.Errorf("%s.Has(%d) failed", test.s1.String(), n)
			}
		}
		if test.s1.Len() != test.wantLen {
			t.Errorf("%s isn't expected length %d", test.s1.String(), test.wantLen)
		}
	}
}
