package rotate

import "testing"

var tests = []struct {
	in  []int
	r   int
	out []int
}{
	{[]int{1, 2, 3, 4, 5}, 1, []int{5, 1, 2, 3, 4}},
	{[]int{1, 2, 3, 4, 5}, 2, []int{4, 5, 1, 2, 3}},
	{[]int{1, 2, 3, 4, 5}, 3, []int{3, 4, 5, 1, 2}},
	{[]int{1, 2, 3, 4, 5}, 4, []int{2, 3, 4, 5, 1}},
	{[]int{1, 2, 3, 4, 5}, 5, []int{1, 2, 3, 4, 5}},
	{[]int{1, 2, 3, 4, 5}, 0, []int{1, 2, 3, 4, 5}},
	{[]int{-1, 2, 3, 4, 5}, 1, []int{5, -1, 2, 3, 4}},
	{[]int{1, -1, 3, 4, 5}, 1, []int{5, 1, -1, 3, 4}},
}

func TestRotate(t *testing.T) {
	for _, tt := range tests {
		var cp = make([]int, len(tt.in))
		copy(cp, tt.in)
		Rotate(tt.in, tt.r)
		if !compareSlice(tt.in, tt.out) {
			t.Errorf("Rotate(%v, %v) == <%v> want <%v>", cp, tt.r, tt.in, tt.out)
		}
	}
}

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
