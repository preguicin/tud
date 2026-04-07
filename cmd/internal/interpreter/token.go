package interpreter

type Token struct {
	token_type TokenType
	lexeme     string
	literal    any
	line       int
}
