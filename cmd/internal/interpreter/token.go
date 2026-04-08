package interpreter

import "fmt"

type Token struct {
	token_type TokenType
	lexeme     string
	literal    any
	line       int
}

func (t Token) String() string {
	return fmt.Sprintf("Type[%s], Lexeme [%s], Literal [%s] and Line[%d]", tokensString[t.token_type], t.lexeme, t.literal, t.line)
}
