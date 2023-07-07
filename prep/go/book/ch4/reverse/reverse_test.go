package reverse

import (
	"testing"
)

var tests = []struct {
	in  []int
	out []int
}{
	{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
	{[]int{1, 1, 2, 3, 3}, []int{3, 3, 2, 1, 1}},
	{[]int{-2, -1, 0, 1, 2}, []int{2, 1, 0, -1, -2}},
}

func TestReverse(t *testing.T) {
	for _, tt := range tests {
		in := make([]int, len(tt.in))
		copy(in, tt.in)
		Reverse(tt.in)
		if !compSlice(tt.in, tt.out) {
			t.Errorf("Reverse(%v) == <%v> want <%v>", in, tt.in, tt.out)
		}
	}
}

func compSlice(s1 []int, s2 []int) bool {
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
