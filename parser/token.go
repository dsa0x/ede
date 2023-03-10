package parser

import (
	"ede/token"
)

func (p *Parser) advanceToken() {
	p.prevToken = p.currToken
	p.currToken = p.nextToken
	p.tokens = append(p.tokens, p.currToken)
	p.nextToken = p.lexer.NextToken()
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

// advanceCurrTokens advances the current token if they match those orders
func (p *Parser) advanceCurrTokens(toks ...token.TokenType) bool {
	for _, tok := range toks {
		if p.currTokenIs(tok) {
			p.advanceToken()
		} else {
			return false
		}
	}
	p.eatEndToken()
	return true
}

// advanceNextTokenIs advances to the next token if it matches, else it does nothing
func (p *Parser) advanceNextTokenIs(tok token.TokenType) bool {
	found := p.nextTokenIs(tok)
	if found {
		p.advanceToken()
	}
	return found
}

func (p *Parser) eatEndToken() {
	for p.currTokenIs(token.SEMICOLON) || p.currTokenIs(token.NEWLINE) {
		p.advanceToken()
	}
}

func (p *Parser) advanceNextToEndToken() bool {
	found := p.nextTokenIs(token.SEMICOLON) || p.nextTokenIs(token.NEWLINE)
	if found {
		p.advanceToken()
	}
	return found
}
