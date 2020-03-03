package parser

import (
	."../token"
	."../scanner"
)

type Parser struct {
	Prev2 Token
	Previous Token
	Current Token
	TokenScanner Scanner

	PanicMode bool
	HadError bool
}

func NewParser(code *string) Parser {
	return Parser{
		Prev2: Token{},
		Previous: Token{},
		Current: Token{},
		TokenScanner: NewScanner(code),
	}
}

func (p *Parser) NextToken() Token {
	return p.TokenScanner.ScanToken()
}


