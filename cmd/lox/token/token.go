package token

import "fmt"

type Token struct {
	totype  TokenType
	lexeme  string
	literal any
	line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %s %v", t.totype, t.lexeme, t.literal)
}

func NewToken(totype TokenType, lexeme string, literal any, line int) Token {
	return Token{
		totype:  totype,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}
