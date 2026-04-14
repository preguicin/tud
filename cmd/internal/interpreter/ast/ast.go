package ast

import "tud/cmd/internal/interpreter/lexer"

type Expr interface{}

type Binary struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}
type Unary struct {
	Right    Expr
	Operator lexer.Token
}

type Literal struct {
	Value any
}

type Grouping struct {
	Expression Expr
}
