package interpreter

type TokenType uint

var tokensString = [...]string{
	// Single char
	"LEFT_PAREN",
	"RIGHT_PAREN",
	"LEFT_BRACE",
	"RIGHT_BRACE",
	"COMMA",
	"DOT",
	"MINUS",
	"PLUS",
	"SEMICOLON",
	"SLASH",
	"STAR",
	// One or two char
	"BANG",
	"BANG_EQUAL",
	"EQUAL",
	"EQUAL_EQUAL",
	"GREATER",
	"GREATER_EQUAL",
	"LESS",
	"LESS_EQUAL",
	// Literals
	"IDENTIFIER",
	"STRING",
	"NUMBER",
	// Keywords
	"AND",
	"CLASS",
	"ELSE",
	"FALSE",
	"FN",
	"FOR",
	"IF",
	"NIL",
	"OR",
	"PRINT",
	"RETURN",
	"SUPER",
	"THIS",
	"TRUE",
	"VAR",
	"WHILE",
	"EOF",
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fn":     FN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

const (
	//Single char
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	//One or two char
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	//Literals
	IDENTIFIER
	STRING
	NUMBER
	//Keywords
	AND
	CLASS
	ELSE
	FALSE
	FN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	EOF
)
