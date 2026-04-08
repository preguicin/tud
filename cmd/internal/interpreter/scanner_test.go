package interpreter

import (
	"reflect"
	"testing"
)

func newTestScanner(src string) (*Scanner, *Interpreter) {
	i := &Interpreter{}
	s := NewScanner(i, []byte(src))
	return &s, i
}

func scannedTokenTypes(tokens []Token) []TokenType {
	types := make([]TokenType, len(tokens))
	for i, t := range tokens {
		types[i] = t.token_type
	}
	return types
}

func TestScanner(t *testing.T) {
	tests := []struct {
		name          string
		src           string
		expectedTypes []TokenType
		expectError   bool
	}{
		{"Single Chars", "(){},.-+;*/", []TokenType{LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE, COMMA, DOT, MINUS, PLUS, SEMICOLON, STAR, SLASH}, false},
		{"Two Chars", "!= == <= >= < > !", []TokenType{BANG_EQUAL, EQUAL_EQUAL, LESS_EQUAL, GREATER_EQUAL, LESS, GREATER, BANG}, false},
		{"Whitespace & Comments", "  \t \n // comment \n +", []TokenType{PLUS}, false},
		{"String Literal", `"hello"`, []TokenType{STRING}, false},
		{"Keywords", "var fn return if else", []TokenType{VAR, FN, RETURN, IF, ELSE}, false},
		{"Keywords", "var dia = 1.0;", []TokenType{VAR, IDENTIFIER, EQUAL, NUMBER, SEMICOLON}, false},
		{"Unterminated String", `"no end`, nil, true},
		{"Invalid Char", "@", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, i := newTestScanner(tt.src)
			tokens := s.ScanTokens()

			if tt.expectError {
				if i.ie == nil {
					t.Errorf("expected error for %q, but got none", tt.src)
				}
				return
			}

			if i.ie != nil {
				t.Fatalf("unexpected error for %q: %s", tt.src, i.ie.Message)
			}

			actualTypes := scannedTokenTypes(tokens)
			if !reflect.DeepEqual(actualTypes, tt.expectedTypes) {
				t.Errorf("mismatch for %q\nexpected: %v\ngot:      %v", tt.src, tt.expectedTypes, actualTypes)
			}
		})
	}
}
