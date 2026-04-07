package interpreter

import (
	"encoding/binary"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	interpreter *Interpreter
	source      []byte
	tokens      []Token
	ch          rune //Current char being proccessed
	start       int
	nextpos     int //The next char position
	line        int
}

func NewScanner(i *Interpreter, data []byte) Scanner {
	return Scanner{
		interpreter: i,
		tokens:      make([]Token, 0),
		source:      data,
		start:       0,
		nextpos:     0,
		ch:          ' ',
		line:        1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.nextpos >= len(s.source)
}

// src: https://eli.thegreenplace.net/2022/a-faster-lexer-in-go/
func (s *Scanner) next() {
	if s.nextpos < len(s.source) {
		r, w := rune(s.source[s.nextpos]), 1
		if r >= utf8.RuneSelf {
			r, w = utf8.DecodeRune(s.source[s.nextpos:])
			if r == utf8.RuneError || r == 0 {
				s.interpreter.ie = &InterpreterError{}
				return
			}
		}
		s.nextpos += w
		s.ch = r
	} else {
		s.nextpos = len(s.source)
		s.ch = -1 // EOF
	}
}

func (s *Scanner) peek(idx int) uint8 {
	if s.nextpos+idx < len(s.source) {
		return s.source[s.nextpos+idx]
	}
	return '\x00'
}

func (s *Scanner) findErrorLine() (string, int) {
	start := s.nextpos
	data := s.source
	if start >= len(data) {
		start = len(data) - 1
	}
	for start > 0 && data[start-1] != '\n' {
		start--
	}

	end := s.nextpos
	for end < len(data) && data[end] != '\n' {
		end++
	}

	line_text := string(data[start:end])
	col := (s.nextpos - start) - 1

	return line_text, col
}

func (s *Scanner) match(char uint8) bool {
	if s.isAtEnd() || s.peek(0) != char {
		return false
	}
	s.nextpos++
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
		switch s.ch {
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
	text := string(s.source[s.start:s.nextpos])
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
			Where:      s.source[s.nextpos-1],
			Message:    "Unterminated string.",
		}
		return
	}
	s.next()

	value := string(s.source[s.start+1 : s.nextpos-1])
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

	data := s.source[s.start:s.nextpos]
	val := float64(binary.BigEndian.Uint64(data))
	s.addToken(NUMBER, (val))
}

func (s *Scanner) scan() {
begin:
	s.next()
	s.skipWhiteSpace()
	char := s.ch

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
			Where:      s.source[s.nextpos-1],
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
