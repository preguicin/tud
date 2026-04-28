package ast

import (
	"testing"
	"tud/cmd/internal/interpreter/lexer"
)

func TestAstPrinter(t *testing.T) {
	expression := &Binary[string]{
		Left: &Unary[string]{
			Operator: lexer.Token{TokenType: lexer.MINUS, Lexeme: "-", Literal: nil, Line: 1},
			Right:    &Literal[string]{Value: 123},
		},
		Operator: lexer.Token{TokenType: lexer.STAR, Lexeme: "*", Literal: nil, Line: 1},
		Right: &Grouping[string]{
			Expression: &Literal[string]{Value: 45.67},
		},
	}

	printer := &AstPrinter{}

	got := printer.Print(expression)
	want := "(* (- 123) (group 45.67))"

	if got != want {
		t.Errorf("AstPrinter error:\n got:  %q\n want: %q", got, want)
	}
}
