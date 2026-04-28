package lexer

import "fmt"

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   any
	Line      int
}

func (t Token) String() string {
	return fmt.Sprintf("Type[%s], Lexeme [%s], Literal [%s] and Line[%d]", tokensString[t.TokenType], t.Lexeme, t.Literal, t.Line)
}
