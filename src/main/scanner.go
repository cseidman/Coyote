package main

import (
	"log"
	"strings"
)

var NULLCHAR = []byte("\000")[0]

type Scanner struct {
	Code        []byte
	Current     int
	Start       int
	Line        int
	SkipCRDepth int
	SkipCRMode  []bool
	ScanDepth   int
}

func NewScanner(source *string) Scanner {

	return Scanner{
		Code:    []byte(*source),
		Current: 0,
		Start:   0,
		Line:    0,

		SkipCRDepth: 0,
		SkipCRMode:  make([]bool, 256),

		ScanDepth: 0,
	}

}

func (s *Scanner) ScanToken() Token {

	s.SkipWhitespace()
	s.Start = s.Current

	if s.isAtEnd() {
		return s.MakeToken(TOKEN_EOF)
	}

	b := s.Advance()

	if s.isAlpha(b) {
		return s.Identifier()
	}
	if s.isDigit(b) {
		return s.Number()
	}

	switch b {
	case '(':
		s.PushCRMode(true)
		return s.MakeToken(TOKEN_LEFT_PAREN)
	case ')':
		s.PopCRMode()
		return s.MakeToken(TOKEN_RIGHT_PAREN)
	case '{':
		s.PushCRMode(false)
		return s.MakeToken(TOKEN_LEFT_BRACE)
	case '}':
		s.PopCRMode()
		return s.MakeToken(TOKEN_RIGHT_BRACE)
	case '[':
		s.PushCRMode(true)
		return s.MakeToken(TOKEN_LEFT_BRACKET)
	case ']':
		s.PopCRMode()
		return s.MakeToken(TOKEN_RIGHT_BRACKET)
	case ';':
		return s.MakeToken(TOKEN_SEMICOLON)
	case ',':
		return s.MakeToken(TOKEN_COMMA)
	case '.':
		return s.MakeToken(TOKEN_DOT)
	case '-':
		return s.MakeToken(TOKEN_MINUS)
	case '+':
		if s.Match('+') {
			return s.MakeToken(TOKEN_PLUS_PLUS)
		}
		return s.MakeToken(TOKEN_PLUS)
	case '/':
		return s.MakeToken(TOKEN_SLASH)
	case '*':
		return s.MakeToken(TOKEN_STAR)
	case '^':
		return s.MakeToken(TOKEN_HAT)
	case ':':
		return s.MakeToken(TOKEN_COLON)
	case '$':
		return s.MakeToken(TOKEN_DOLLAR)
	// -- Double tokens
	case '@':
		if s.Match('[') {
			s.PushCRMode(true)
			return s.MakeToken(TOKEN_LEFT_ARRAY)
		} else if s.Match('{') {
			s.PushCRMode(true)
			return s.MakeToken(TOKEN_LEFT_LIST)

		} else {
			return s.MakeToken(TOKEN_AT)
		}
	case '\n':
		s.Line++
		return s.MakeToken(TOKEN_CR)
	case '!':
		if s.Match('=') {
			return s.MakeToken(TOKEN_BANG_EQUAL)
		} else {
			return s.MakeToken(TOKEN_BANG)
		}
	case '=':
		if s.Match('=') {
			return s.MakeToken(TOKEN_EQUAL_EQUAL)
		} else {
			return s.MakeToken(TOKEN_EQUAL)
		}
	case '<':
		if s.Match('=') {
			return s.MakeToken(TOKEN_LESS_EQUAL)
		} else {
			return s.MakeToken(TOKEN_LESS)
		}
	case '>':
		if s.Match('=') {
			return s.MakeToken(TOKEN_GREATER_EQUAL)
		} else {
			return s.MakeToken(TOKEN_GREATER)
		}
	// -- Literals
	case '"':
		return s.String()

	default:
		log.Panicf("Unrecognized character byte:%d char: %U", b, b)
	}
	return s.MakeToken(TOKEN_IDENTIFIER)
}

func (s *Scanner) PushCRMode(CRMode bool) {
	s.SkipCRDepth++
	s.SkipCRMode[s.SkipCRDepth] = CRMode
}

func (s *Scanner) PopCRMode() {
	s.SkipCRDepth--
}

func (s *Scanner) CurrentCRMode() bool {
	return s.SkipCRMode[s.SkipCRDepth]
}

func (s *Scanner) Advance() byte {
	s.Current++
	return s.Code[s.Current-1]
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) Identifier() Token {
	// This is for the name of the variable or other identifier
	for s.isAlpha(s.Peek()) || s.isDigit(s.Peek()) {
		s.Advance()
	}
	return s.MakeToken(s.IdentifierType())
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) Number() Token {
	thisToken := TOKEN_INTEGER

	for s.isDigit(s.Peek()) {
		s.Advance()
	}
	// See if this is a float or decimal
	if s.Peek() == '.' && s.isDigit(s.PeekNext()) {
		thisToken = TOKEN_DECIMAL
		s.Advance()
		for s.isDigit(s.Peek()) {
			s.Advance()
		}
	}
	return s.MakeToken(thisToken)
}

func (s *Scanner) String() Token {
	for s.Peek() != '"' && !s.isAtEnd() {
		if s.Peek() == '\n' {
			s.Line++
		}
		s.Advance()
	}
	if s.isAtEnd() {
		return s.ErrorToken("Unterminated string")
	}
	s.Advance()
	return s.MakeToken(TOKEN_STRING)
}

func (s *Scanner) Match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.Code[s.Current] != expected {
		return false
	}
	s.Current++
	return true
}

func (s *Scanner) isAtEnd() bool {
	if s.Current >= len(s.Code) {
		return true
	}
	return s.Code[s.Current] == NULLCHAR
}

func (s *Scanner) MakeToken(t_type TokenType) Token {
	var token = Token{}
	token.Type = t_type
	token.Length = s.Current
	token.Line = s.Line
	token.Value = s.Code[s.Start:s.Current]

	return token
}

func (s *Scanner) ErrorToken(message string) Token {
	var token = Token{}
	token.Type = TOKEN_ERROR
	token.Value = []byte(message)
	token.Length = len(message)
	token.Line = s.Line

	return token
}

func (s *Scanner) IsTokenCaseInsensitive(tokenString string) bool {
	return TokenLabels[tokenString].IsCaseSensitive
}

func (s *Scanner) IdentifierType() TokenType {

	tokenString := string(s.GetTokenValue())
	// This is what we use to see if this keyword is case insensitive
	testString := strings.ToLower(tokenString)
	// Check to see if this allows mixed case
	if s.IsTokenCaseInsensitive(testString) {
		// If so, then change it to lowercase to match the token list
		tokenString = testString
	}
	// Does this value exist as a keyword? If so, return it
	if val, ok := TokenLabels[tokenString]; ok {
		return val.Type
	} else {
		return TOKEN_IDENTIFIER
	}

}

func (s *Scanner) GetTokenValue() []byte {
	return s.Code[s.Start:s.Current]
}

func (s *Scanner) Peek() byte {
	if s.Current >= len(s.Code) {
		return NULLCHAR
	}
	return s.Code[s.Current]
}

func (s *Scanner) PeekNext() byte {
	if s.isAtEnd() {
		return NULLCHAR
	}
	return s.Code[s.Current+1]
}

func (s *Scanner) SkipWhitespace() {
	for {
		c := s.Peek()
		switch c {
		case '\r', ' ', '\v', '\f', '\t':
			s.Advance()
		case '\n':
			if s.CurrentCRMode() {
				s.Line++
				s.Advance()
			}
			return
		case '/':
			// If the next character is a *
			if s.PeekNext() == '*' {
				// We are in the beginning of comment mode
				commentHeaderDepth := 0
				s.PushCRMode(true)

				// Move to the *
				s.Advance()
				// Keep reading characters
				for {
					// Move to the next character
					s.Advance()
					// If the current character is / and the next is *
					if s.Peek() == '/' && s.PeekNext() == '*' {
						// Opens a new comment block
						commentHeaderDepth++
						// Move to the *
						s.Advance()
						// Move to the next character
						s.Advance()
					}

					// If the current character is * followed by /
					if s.Peek() == '*' && s.PeekNext() == '/' {

						// This closes the comment block
						// Move to the /
						s.Advance()
						s.Advance()
						// If there's a newline ..
						if s.Peek() == '\n' {
							// Move to the newline
							s.Advance()
						}
						if commentHeaderDepth == 0 {
							s.Advance()
							s.PopCRMode()
							break
						}
						commentHeaderDepth--
					}

				}
			}

			if s.PeekNext() == '/' {
				// A comment goes until the end of the line.
				for s.Peek() != '\n' && !s.isAtEnd() {
					s.Advance()
				}
			} else {
				return
			}

		default:

			return
		}
	}
}
