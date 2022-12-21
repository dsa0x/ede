package parser

import (
	"ede/token"
)

// Precedence order for operators
const (
	_ int = iota
	LOWEST
	COND        // OR or AND
	ASSIGN      // =
	EQ          // == or !=
	LESSGREATER // > or <
	SUM         // + or -
	PRODUCT     // * or /
	POWER       // **
	MOD         // %
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index], map[key]
	HIGHEST
)

func (p *Parser) precedence(tokenType token.TokenType) int {
	switch tokenType {
	case token.ASSIGN, token.PLUS_EQUAL, token.MINUS_EQUAL:
		return ASSIGN
	case token.EQ, token.NEQ:
		return EQ
	case token.GT, token.LT, token.LTE, token.GTE:
		return LESSGREATER
	case token.PLUS, token.MINUS, token.DEC, token.INC, token.MODULO:
		return SUM
	case token.ASTERISK, token.SLASH:
		return PRODUCT
	case token.LPAREN:
		return CALL
	case token.LBRACKET:
		return INDEX
	case token.RANGE_ARRAY, token.DOT:
		return HIGHEST
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	return p.precedence(p.currToken.Type)
}
func (p *Parser) peekPrecedence() int {
	return p.precedence(p.nextToken.Type)
}
