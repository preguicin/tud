package scanner

import (
	"badparser/cmd/lox/token"
)

type Scanner struct {
	source  string
	tokens  []token.Token
	line    int
	current int
	start   int
}

func (s *Scanner) isAtEnd() bool {

	return true
}

func (s *Scanner) ScanToken() {

}

func (s *Scanner) ScanTokens() []token.Token {
	for {
		if s.isAtEnd() {
			s.start = s.current
			s.ScanToken()
			break
		}
	}

	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens
}

func NewScanner(source string) Scanner {
	return Scanner{
		start:   0,
		current: 0,
		line:    1,
		source:  source,
		tokens:  make([]token.Token, 0, 300),
	}
}
