package interpreter

import (
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

func TestSingleCharTokens(t *testing.T) {
	cases := []struct {
		src      string
		expected TokenType
	}{
		{"(", LEFT_PAREN},
		{")", RIGHT_PAREN},
		{"{", LEFT_BRACE},
		{"}", RIGHT_BRACE},
		{",", COMMA},
		{".", DOT},
		{"-", MINUS},
		{"+", PLUS},
		{";", SEMICOLON},
		{"*", STAR},
		{"/", SLASH},
	}

	for _, tc := range cases {
		t.Run(tc.src, func(t *testing.T) {
			s, _ := newTestScanner(tc.src)
			tokens := s.ScanTokens()
			found := false
			for _, tok := range tokens {
				if tok.token_type == tc.expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected token %v for input %q, got %v", tc.expected, tc.src, scannedTokenTypes(tokens))
			}
		})
	}
}

func TestTwoCharTokens(t *testing.T) {
	cases := []struct {
		src      string
		expected TokenType
	}{
		{"!=", BANG_EQUAL},
		{"!", BANG},
		{"==", EQUAL_EQUAL},
		{"<=", LESS_EQUAL},
		{"<", LESS},
		{">=", GREATER_EQUAL},
		{">", GREATER},
	}

	for _, tc := range cases {
		t.Run(tc.src, func(t *testing.T) {
			s, _ := newTestScanner(tc.src)
			tokens := s.ScanTokens()
			found := false
			for _, tok := range tokens {
				if tok.token_type == tc.expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected token %v for input %q, got %v", tc.expected, tc.src, scannedTokenTypes(tokens))
			}
		})
	}
}

func TestCommentIsSkipped(t *testing.T) {
	s, _ := newTestScanner("//t\n+")
	tokens := s.ScanTokens()
	for _, tok := range tokens {
		if tok.token_type == SLASH {
			t.Errorf("comment should have been skipped, got SLASH token")
		}
	}
	found := false
	for _, tok := range tokens {
		if tok.token_type == PLUS {
			found = true
		}
	}
	if !found {
		t.Errorf("expected PLUS after comment, got %v", scannedTokenTypes(tokens))
	}
}

func TestWhitespaceIsSkipped(t *testing.T) {
	s, _ := newTestScanner("   \t+")
	tokens := s.ScanTokens()
	found := false
	for _, tok := range tokens {
		if tok.token_type == PLUS {
			found = true
		}
	}
	if !found {
		t.Errorf("expected PLUS after whitespace, got %v", scannedTokenTypes(tokens))
	}
}

func TestNewlineIncrementsLine(t *testing.T) {
	s, _ := newTestScanner("\n+")
	tokens := s.ScanTokens()
	for _, tok := range tokens {
		if tok.token_type == PLUS && tok.line != 2 {
			t.Errorf("expected PLUS on line 2, got line %d", tok.line)
		}
	}
}

func TestValidString(t *testing.T) {
	s, i := newTestScanner(`"hello"`)
	tokens := s.ScanTokens()
	if i.ie != nil {
		t.Fatalf("unexpected error: %s", i.ie.Message)
	}
	for _, tok := range tokens {
		if tok.token_type == STRING {
			if tok.literal != "hello" {
				t.Errorf("expected literal 'hello', got %v", tok.literal)
			}
			return
		}
	}
	t.Errorf("expected STRING token, got %v", scannedTokenTypes(tokens))
}

func TestUnterminatedString(t *testing.T) {
	s, i := newTestScanner(`"hello`)
	s.ScanTokens()
	if i.ie == nil {
		t.Error("expected an InterpreterError for unterminated string, got nil")
	}
}

func TestUnknownCharacterSetsError(t *testing.T) {
	s, i := newTestScanner("@")
	s.ScanTokens()
	if i.ie == nil {
		t.Error("expected an InterpreterError for unknown character '@', got nil")
	}
}

func TestMultipleTokenSequence(t *testing.T) {
	s, _ := newTestScanner("(+)")
	tokens := s.ScanTokens()
	expected := []TokenType{LEFT_PAREN, PLUS, RIGHT_PAREN}
	found := 0
	for _, tok := range tokens {
		for _, e := range expected {
			if tok.token_type == e {
				found++
			}
		}
	}
	if found != len(expected) {
		t.Errorf("expected tokens %v, got %v", expected, scannedTokenTypes(tokens))
	}
}
