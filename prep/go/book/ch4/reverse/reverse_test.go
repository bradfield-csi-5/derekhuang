package reverse

import "testing"

var tests = []struct {
	in  [5]int
	out [5]int
}{
	{[5]int{1, 2, 3, 4, 5}, [5]int{5, 4, 3, 2, 1}},
	{[5]int{-2, -1, 0, 1, 2}, [5]int{2, 1, 0, -1, -2}},
	{[5]int{1, 1, 2, 3, 3}, [5]int{3, 3, 2, 1, 1}},
}

func TestReverse(t *testing.T) {
	for _, tt := range tests {
		var cp [5]int = tt.in
		Reverse(&tt.in)
		if tt.in != tt.out {
			t.Errorf("Reverse(%v) == <%v> want <%v>", cp, tt.in, tt.out)
		}
	}
}
