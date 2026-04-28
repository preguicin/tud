package ast

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (ap *AstPrinter) Print(expr Expr[string]) string {
	return expr.Accept(ap)
}

func (ap *AstPrinter) VisitBinary(b *Binary[string]) string {
	return ap.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (ap *AstPrinter) VisitGrouping(g *Grouping[string]) string {
	return ap.parenthesize("group", g.Expression)
}

func (ap *AstPrinter) VisitLiteral(l *Literal[string]) string {
	if l.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", l.Value)
}

func (ap *AstPrinter) VisitUnary(u *Unary[string]) string {
	return ap.parenthesize(u.Operator.Lexeme, u.Right)
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr[string]) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(ap))
	}

	builder.WriteString(")")
	return builder.String()
}
