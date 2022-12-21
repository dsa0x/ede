package parser

import (
	"ede/token"
)

func (p *Parser) advanceToken() {
	p.prevToken = p.currToken
	p.currToken = p.nextToken
	p.tokens = append(p.tokens, p.currToken)
	p.nextToken = p.lexer.NextToken()
	p.pos = token.Pos{Column: p.lexer.Column(), Line: p.lexer.Line()}
}

func (p *Parser) currTokenIs(tok token.TokenType) bool {
	return p.currToken.Type == tok
}
func (p *Parser) nextTokenIs(tok token.TokenType) bool {
	return p.nextToken.Type == tok
}

// advanceCurrTokenIs advances to the next token if the current token  matches, else it does nothing
func (p *Parser) advanceCurrTokenIs(tok token.TokenType) bool {
	found := p.currTokenIs(tok)
	if found {
		p.advanceToken()
	}
	return found
}

// advanceNextTokenIs advances to the next token if it matches, else it does nothing
func (p *Parser) advanceNextTokenIs(tok token.TokenType) bool {
	found := p.nextTokenIs(tok)
	if found {
		p.advanceToken()
	}
	return found
}

func (p *Parser) advanceCurrToEndToken() bool {
	found := p.currTokenIs(token.SEMICOLON) || p.currTokenIs(token.NEWLINE)
	if found {
		p.advanceToken()
	}
	return found
}

func (p *Parser) eatEndToken() {
	for p.advanceCurrToEndToken() {
	}
}

func (p *Parser) advanceNextToEndToken() bool {
	found := p.nextTokenIs(token.SEMICOLON) || p.nextTokenIs(token.NEWLINE)
	if found {
		p.advanceToken()
	}
	return found
}
