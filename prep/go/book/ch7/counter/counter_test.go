package counter

import (
	"testing"
)

func TestWordCounter(t *testing.T) {
	var tests = []struct {
		words []byte
		want  int
	}{
		{[]byte(""), 0},
		{[]byte("foo"), 1},
		{[]byte("foo bar"), 2},
		{[]byte("foo bar baz"), 3},
		{[]byte("foo\nbar baz"), 3},
	}
	for _, tt := range tests {
		var c WordCounter
		c.Write(tt.words)
		if int(c) != tt.want {
			t.Errorf("c.Write(%v) == <%d> want <%d>", tt.words, int(c), tt.want)
		}
	}
}
