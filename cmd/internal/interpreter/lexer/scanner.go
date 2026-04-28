package lexer

import (
	"strconv"
	"tud/cmd/internal/file"
	"unicode"
	"unicode/utf8"
)

const eof = -1

type ErrorHandler func(pos file.Position, msg string)

type Scanner struct {
	file        file.File
	source      []byte
	tokens      []Token
	ch          rune //Current char being proccessed
	createError ErrorHandler
}

func NewScanner(file file.File, err ErrorHandler) Scanner {
	return Scanner{
		file:        file,
		createError: err,
		tokens:      make([]Token, 0),
		source:      file.GetBytes(),
		ch:          ' ',
	}
}
func (s *Scanner) pos() *file.Position {
	return s.file.Pos()
}

func (s *Scanner) isAtEnd() bool {
	return s.pos().Offset >= len(s.source)
}

// src: https://eli.thegreenplace.net/2022/a-faster-lexer-in-go/
func (s *Scanner) next() {
	if s.pos().Offset < len(s.source) {
		r, w := rune(s.source[s.pos().Offset]), 1
		if r >= utf8.RuneSelf {
			r, w = utf8.DecodeRune(s.source[s.pos().Offset:])
			if r == utf8.RuneError || r == 0 {
				s.createError(*s.file.Pos(), "Invalid UTF-8 charachter")
				return
			}
		}
		s.pos().Offset += w
		s.ch = r
	} else {
		s.pos().Offset = len(s.source)
		s.ch = eof
	}
}

func (s *Scanner) peek(idx int) uint8 {

	nextChar := s.pos().Offset + idx
	if nextChar < len(s.source) {
		return s.source[nextChar]
	}

	return '\x00'
}

func (s *Scanner) match(char uint8) bool {
	if s.isAtEnd() || s.peek(0) != char {
		return false
	}
	s.pos().Offset += 1
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
			s.pos().Line += 1
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

func (s *Scanner) skipcComment() {
	for {
		ch := s.peek(0)
		if ch == '*' && s.peek(1) == '/' {
			s.next()
			s.next()
			break
		} else if s.isAtEnd() {
			break
		}
		s.next()
	}
}
func (s *Scanner) addToken(tt TokenType, literal any) {
	text := ""
	s.tokens = append(s.tokens, Token{TokenType: tt, Lexeme: text, Literal: literal, Line: s.pos().Line})
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
		s.createError(*s.file.Pos(), "Unterminated string.")
		return
	}
	s.next()

	value := string(s.source[s.pos().Line+1 : s.pos().Offset-1])
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

	data := s.source[s.pos().Start:s.pos().Offset]
	// TODO: Make add toke use slice of bytes from origin and convert the numbers on the tokenizer
	val, err := strconv.ParseFloat(string(data), 64)

	if err != nil {
		s.createError(*s.file.Pos(), "Failed to convert number.")
	}

	s.addToken(NUMBER, (val))
}

func (s *Scanner) scan() {
begin:
	s.skipWhiteSpace()
	s.pos().Start = s.pos().Offset

	s.next()
	char := s.ch

	if char == eof {
		return
	}

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
		} else if s.peek(0) == '*' {
			s.skipcComment()
			goto begin
		} else {
			s.addToken(SLASH, nil)
		}
	case '"':
		s.scanString()
	default:
		if unicode.IsLetter(s.ch) || unicode.IsNumber(s.ch) {
			s.indetifier()
		} else {
			s.createError(*s.file.Pos(), "Unsupported Type.")
		}
	}
}

func (s *Scanner) indetifier() {
	for {
		ch := rune(s.peek(0))
		if unicode.IsLetter(ch) || unicode.IsNumber(ch) {
			s.next()
		} else {
			break
		}
	}

	text := string(s.source[s.pos().Start:s.pos().Offset])

	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType, nil)
}

func (s *Scanner) ScanTokens() []Token {
	for {
		if s.isAtEnd() {
			break
		}
		s.scan()
	}
	return s.tokens
}
