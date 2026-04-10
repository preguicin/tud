// package ast

// import "tud/cmd/internal/lexer"

// type Expr interface {
// 	Accept(visitor Visitor) any
// }

// type Binary struct {
// 	Left     Expr
// 	Operator lexer.Token
// 	Right    Expr
// }

// func (b *Binary) Accept(v Visitor) any { return v.VisitBinaryExpr(b) }

// type Literal struct {
// 	Value any
// }

// func (l *Literal) Accept(v Visitor) any { return v.VisitLiteralExpr(l) }

// type Grouping struct {
// 	Expression Expr
// }

// func (g *Grouping) Accept(v Visitor) any { return v.VisitGroupingExpr(g) }
