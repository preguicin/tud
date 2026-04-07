package interpreter

import (
	"encoding/binary"
	"unicode"
)

type Scanner struct {
	interpreter *Interpreter
	source      []byte
	tokens      []Token
	start       int
	current     int
	line        int
}

func NewScanner(i *Interpreter, data []byte) Scanner {
	return Scanner{
		interpreter: i,
		tokens:      make([]Token, 0),
		source:      data,
		start:       0,
		current:     0,
		line:        1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) next() uint8 {
	val := s.source[s.current]
	s.current++
	return val
}

func (s *Scanner) peek(idx int) uint8 {
	if s.isAtEnd() {
		return '\x00'
	}

	if s.current+idx >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+idx]
}

func (s *Scanner) findErrorLine() (string, int) {
	start := s.current
	data := s.source
	if start >= len(data) {
		start = len(data) - 1
	}
	for start > 0 && data[start-1] != '\n' {
		start--
	}

	end := s.current
	for end < len(data) && data[end] != '\n' {
		end++
	}

	line_text := string(data[start:end])
	col := (s.current - start) - 1

	return line_text, col
}

func (s *Scanner) match(char uint8) bool {
	if s.isAtEnd() || s.peek(0) != char {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) multOpt(normal TokenType, other TokenType, comp uint8) {

	tt := normal

	if !s.isAtEnd() && s.match(comp) {
		tt = other
	}

	s.addToken(tt, nil)
}

func (s *Scanner) skipWhiteSpace() {
	for {
		switch s.peek(0) {
		case ' ', '\t', '\r':
			s.next()
		case '\n':
			s.line++
			s.next()
		default:
			return
		}
	}
}

func (s *Scanner) skipComment() {
	for {
		ch := s.peek(0)
		if ch == '\n' {
			s.next()
			break
		} else if s.isAtEnd() {
			break
		}
		s.next()
	}
}

func (s *Scanner) addToken(tt TokenType, literal any) {
	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, Token{token_type: tt, lexeme: text, literal: literal, line: s.line})
}

func (s *Scanner) scanString() {
	for {
		ch := s.peek(0)
		if ch == '"' || s.isAtEnd() {
			break
		}
		s.next()
	}
	if s.isAtEnd() {
		line_txt, col := s.findErrorLine()
		s.interpreter.ie = &InterpreterError{
			Line:       s.line,
			SourceLine: line_txt,
			Col:        col,
			Where:      s.source[s.current-1],
			Message:    "Unterminated string.",
		}
		return
	}
	s.next()

	value := string(s.source[s.start+1 : s.current-1])
	s.addToken(STRING, value)

}

func (s *Scanner) scanNumber() {
	for {
		if !unicode.IsDigit(rune(s.peek(0))) {
			break
		}
		s.next()
	}

	if s.peek(0) == '.' && unicode.IsDigit(rune(s.peek(1))) {
		s.next()
		for {
			if !unicode.IsDigit(rune(s.peek(0))) {
				break
			}
			s.next()
		}
	}

	data := s.source[s.start:s.current]
	val := float64(binary.BigEndian.Uint64(data))
	s.addToken(NUMBER, (val))
}

func (s *Scanner) scan() {
begin:
	s.skipWhiteSpace()
	char := s.next()
	if unicode.IsDigit(rune(char)) {
		s.scanNumber()
		return
	}
	switch char {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '!':
		s.multOpt(BANG, BANG_EQUAL, '=')
	case '=':
		s.multOpt(EQUAL, EQUAL_EQUAL, '=')
	case '<':
		s.multOpt(LESS, LESS_EQUAL, '=')
	case '>':
		s.multOpt(GREATER, GREATER_EQUAL, '=')
	case '/':
		if s.peek(0) == '/' {
			s.skipComment()
			goto begin
		} else {
			s.addToken(SLASH, nil)
		}
	case '"':
		s.scanString()
	default:
		line_txt, col := s.findErrorLine()
		s.interpreter.ie = &InterpreterError{
			Line:       s.line,
			SourceLine: line_txt,
			Col:        col,
			Where:      s.source[s.current-1],
			Message:    "Unterminated string.",
		}
	}
}

func (s *Scanner) ScanTokens() []Token {
	for {
		if s.isAtEnd() || s.interpreter.ie != nil {
			break
		}
		s.scan()
	}
	return s.tokens
}
