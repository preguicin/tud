package ast

import "tud/cmd/internal/interpreter/lexer"

type Visitor[R any] interface {
	VisitBinary(*Binary[R]) R
	VisitUnary(*Unary[R]) R
	VisitLiteral(*Literal[R]) R
	VisitGrouping(*Grouping[R]) R
}

type Expr[R any] interface {
	Accept(Visitor[R]) R
}

type Binary[R any] struct {
	Left     Expr[R]
	Operator lexer.Token
	Right    Expr[R]
}

func (b *Binary[R]) Accept(v Visitor[R]) R {
	return v.VisitBinary(b)
}

type Unary[R any] struct {
	Right    Expr[R]
	Operator lexer.Token
}

func (u *Unary[R]) Accept(v Visitor[R]) R {
	return v.VisitUnary(u)
}

type Literal[R any] struct {
	Value any
}

func (l *Literal[R]) Accept(v Visitor[R]) R {
	return v.VisitLiteral(l)
}

type Grouping[R any] struct {
	Expression Expr[R]
}

func (g *Grouping[R]) Accept(v Visitor[R]) R {
	return v.VisitGrouping(g)
}
