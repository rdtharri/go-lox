package runner

import (
	"strconv"
)

type Scanner struct {
	Source string
	Tokens []Token
	Runner *LoxRunner

	start   int
	current int
	line    int
}

func NewScanner(source string, runner *LoxRunner) *Scanner {
	scanner := new(Scanner)
	scanner.line = 1
	scanner.Source = source
	scanner.Runner = runner
	return scanner
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.start = s.current
	s.addNullToken(EOF)
	return s.Tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '(':
		s.addNullToken(LEFT_PAREN)
	case ')':
		s.addNullToken(RIGHT_PAREN)
	case '{':
		s.addNullToken(LEFT_BRACE)
	case '}':
		s.addNullToken(RIGHT_BRACE)
	case ',':
		s.addNullToken(COMMA)
	case '.':
		s.addNullToken(DOT)
	case '-':
		s.addNullToken(MINUS)
	case '+':
		s.addNullToken(PLUS)
	case ';':
		s.addNullToken(SEMICOLON)
	case '*':
		s.addNullToken(STAR)
	case '!':
		if s.match('=') {
			s.addNullToken(BANG_EQUAL)
		} else {
			s.addNullToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addNullToken(EQUAL_EQUAL)
		} else {
			s.addNullToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addNullToken(LESS_EQUAL)
		} else {
			s.addNullToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addNullToken(GREATER_EQUAL)
		} else {
			s.addNullToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addNullToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(char) {
			s.number()
		} else if s.isAlpha(char) {
			s.identifier()
		} else {
			s.Runner.error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) advance() rune {
	retVal := []rune(s.Source)[s.current]
	s.current++
	return retVal
}

func (s *Scanner) match(check rune) bool {
	if s.isAtEnd() {
		return false
	}

	if []rune(s.Source)[s.current] != check {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\n'
	}
	return []rune(s.Source)[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.Source) {
		return '\n'
	}
	return []rune(s.Source)[s.current+1]
}

func (s *Scanner) string() {

	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.Runner.error(s.line, "Unterminated string.")
	}
	s.advance()
	s.addValueToken(STRING, string(s.Source[(s.start+1):(s.current-1)]))

}

func (s *Scanner) isDigit(char rune) bool {
	return char >= '0' && char <= '9'

}

func (s *Scanner) isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char == '_')
}

func (s *Scanner) isAlphaNumeric(char rune) bool {
	return s.isDigit(char) || s.isAlpha(char)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	s.addValueToken(NUMBER, string(s.Source[(s.start):(s.current)]))

}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	lexeme := string(s.Source[(s.start):(s.current)])
	ttype, found := KeyMap[lexeme]
	if !found {
		ttype = IDENTIFIER
	}

	s.addNullToken(ttype)
}

func (s *Scanner) addNullToken(ttype TokenType) {
	s.appendToken(
		Token{
			Type:   ttype,
			Lexeme: string(s.Source[(s.start):(s.current)]),
			Value:  nil,
			Line:   s.line,
		},
	)
}

func (s *Scanner) addValueToken(ttype TokenType, value string) {
	newToken := Token{
		Type:   ttype,
		Lexeme: string(s.Source[(s.start):(s.current)]),
		Line:   s.line,
	}

	switch ttype {
	case NUMBER:
		numVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			// Not sure how to hit this or what should happen
			// TODO: add better panic-recover to scanner
			panic(err)
		}
		newToken.Value = numVal
	case STRING:
		newToken.Value = value
	}

	s.appendToken(newToken)
}

func (s *Scanner) appendToken(token Token) {
	s.Tokens = append(s.Tokens, token)
}
