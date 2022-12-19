package parser

import "ede/token"

func (p *Parser) advanceToken() {
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

// advanceNextTokenIs advances to the next token if it matches, else it does nothing
func (p *Parser) advanceNextTokenIs(tok token.TokenType) bool {
	found := p.nextTokenIs(tok)
	if found {
		p.advanceToken()
	}
	return found
}
