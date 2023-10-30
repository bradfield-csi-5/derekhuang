package scanner

import (
	"testing"
)

func TestScanner(t *testing.T) {
	for _, testCase := range []struct {
		input    []byte
		expected []Token
	}{
		{
			[]byte("alice AND bob"),
			[]Token{
				Token{TERM, "alice"},
				Token{AND, "AND"},
				Token{TERM, "bob"},
			},
		},
		{
			[]byte("alice and bob"),
			[]Token{
				Token{TERM, "alice"},
				Token{AND, "and"},
				Token{TERM, "bob"},
			},
		},
		{
			[]byte("alice AND -bob"),
			[]Token{
				Token{TERM, "alice"},
				Token{AND, "AND"},
				Token{NOT, "-"},
				Token{TERM, "bob"},
			},
		},
		{
			[]byte("alice AND NOT bob"),
			[]Token{
				Token{TERM, "alice"},
				Token{AND, "AND"},
				Token{NOT, "NOT"},
				Token{TERM, "bob"},
			},
		},
		{
			[]byte("(alice OR bob) and (carol AND (dave or eve))"),
			[]Token{
				Token{L_PAREN, "("},
				Token{TERM, "alice"},
				Token{OR, "OR"},
				Token{TERM, "bob"},
				Token{R_PAREN, ")"},
				Token{AND, "and"},
				Token{L_PAREN, "("},
				Token{TERM, "carol"},
				Token{AND, "AND"},
				Token{L_PAREN, "("},
				Token{TERM, "dave"},
				Token{OR, "or"},
				Token{TERM, "eve"},
				Token{R_PAREN, ")"},
				Token{R_PAREN, ")"},
			},
		},
	} {
		var actual []Token
		s := newScanner(testCase.input)
		for {
			token, err := s.Scan()
			if err != nil {
				t.Fatal(err)
			}
			if token == tokenEOF {
				break
			}
			actual = append(actual, token)
		}
		if len(testCase.expected) != len(actual) {
			t.Fatalf("expected %d tokens, got %d:\n%v\n%v", len(testCase.expected), len(actual), testCase.expected, actual)
		}
		for i := 0; i < len(actual); i++ {
			if testCase.expected[i] != actual[i] {
				t.Fatalf("expected token %d to be %v, got %v", i, testCase.expected[i], actual[i])
			}
		}
	}
}
