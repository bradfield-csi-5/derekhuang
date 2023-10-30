package scanner

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	INVALID TokenType = iota
	L_PAREN
	R_PAREN
	NOT
	AND
	OR
	TERM
	EOF
)

var typeToString = map[TokenType]string{
	INVALID: "INVALID",
	L_PAREN: "L_PAREN",
	R_PAREN: "R_PAREN",
	NOT:     "NOT",
	AND:     "AND",
	OR:      "OR",
	TERM:    "TERM",
	EOF:     "EOF",
}

var keywordToType = map[string]TokenType{
	"NOT": NOT,
	"AND": AND,
	"OR":  OR,
}

var symbolToType = map[rune]TokenType{
	'(': L_PAREN,
	')': R_PAREN,
	'-': NOT,
}

var (
	tokenEOF     = Token{EOF, ""}
	tokenInvalid = Token{INVALID, ""}
)

type TokenType int

type Token struct {
	tokenType TokenType
	literal   string
}

func (t Token) String() string {
	return fmt.Sprintf("%s('%s')", typeToString[t.tokenType], t.literal)
}

type Scanner struct {
	cur       int
	src       []byte
	nextToken Token
}

func (s *Scanner) Scan() (Token, error) {
	if s.nextToken != tokenInvalid {
		result := s.nextToken
		s.nextToken = tokenInvalid
		return result, nil
	}
	s.skipWhitespace()
	if s.isAtEnd() {
		return tokenEOF, nil
	}
	r, _ := s.peekChar()
	if r == utf8.RuneError {
		return tokenInvalid, fmt.Errorf("Invalid utf-8 character at position %d", s.cur)
	} else if t, ok := symbolToType[r]; ok {
		start := s.cur
		s.advance()
		return Token{t, string(s.src[start:s.cur])}, nil
	} else if unicode.IsLetter(r) {
		return s.scanLetters()
	}
	return tokenInvalid, fmt.Errorf("Unsupported character %c\n", r)
}

func (s *Scanner) advance() (rune, int) {
	rune, size := s.peekChar()
	s.cur += size
	return rune, size
}

func (s *Scanner) isAtEnd() bool {
	return s.cur == len(s.src)
}

func (s *Scanner) peekChar() (rune, int) {
	return utf8.DecodeRune(s.src[s.cur:])
}

func (s *Scanner) scanLetters() (Token, error) {
	start := s.cur
	s.advance()
	for !s.isAtEnd() {
		r, _ := s.peekChar()
		if !unicode.IsLetter(r) {
			break
		}
		s.advance()
	}
	literal := string(s.src[start:s.cur])
	if keywordType, ok := keywordToType[strings.ToUpper(literal)]; ok {
		return Token{keywordType, literal}, nil
	}
	return Token{TERM, literal}, nil
}

func (s *Scanner) skipWhitespace() {
	for !s.isAtEnd() {
		r, _ := s.peekChar()
		if !unicode.IsSpace(r) {
			return
		}
		s.advance()
	}
}

func newScanner(src []byte) *Scanner {
	return &Scanner{
		cur:       0,
		src:       src,
		nextToken: tokenInvalid,
	}
}
