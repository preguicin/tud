package lexer

import (
	"reflect"
	"testing"
	"tud/cmd/internal/file"
	"tud/cmd/internal/interpreter/error"
)

func newTestScanner(src string) (*Scanner, *error.Reporter) {
	data := []byte(src)
	f := file.NewInMemFile(data)
	errReporter := error.NewErrorReporter(data)

	s := NewScanner(f, errReporter.NewError)
	return &s, &errReporter
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
			s, errRep := newTestScanner(tt.src)
			tokens := s.ScanTokens()

			if tt.expectError {
				println(errRep.Errors[0].Message)
				if len(errRep.Errors) == 0 {
					t.Errorf("expected error for %q, but got none", tt.src)
				}
				return
			}

			if len(errRep.Errors) != 0 {
				for _, err := range errRep.Errors {
					t.Fatalf("unexpected error for %q: %s", tt.src, err.Message)
				}
			}

			actualTypes := scannedTokenTypes(tokens)
			if !reflect.DeepEqual(actualTypes, tt.expectedTypes) {
				t.Errorf("mismatch for %q\nexpected: %v\ngot:      %v", tt.src, tt.expectedTypes, actualTypes)
			}
		})
	}
}
