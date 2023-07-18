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

func TestLineCounter(t *testing.T) {
	var tests = []struct {
		lines []byte
		want  int
	}{
		{[]byte(""), 0},
		{[]byte("foo"), 1},
		{[]byte("\nfoo\n"), 2},
		{[]byte("\n\nfoo"), 3},
		{[]byte("\n\n\n"), 3},
		{[]byte("\nfoo\nbar\nbaz"), 4},
	}
	for _, tt := range tests {
		var c LineCounter
		c.Write(tt.lines)
		if int(c) != tt.want {
			t.Errorf("c.Write(%v) == <%d> want <%d>", tt.lines, int(c), tt.want)
		}
	}
}
