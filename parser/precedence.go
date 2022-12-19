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

var precedences = map[token.TokenType]int{
	token.ASSIGN: ASSIGN,
	token.EQ:     EQ,
	token.NEQ:    EQ,
}

func (p *Parser) precedence(tokenType token.TokenType) int {
	switch tokenType {
	case token.ASSIGN:
		return ASSIGN
	case token.EQ, token.NEQ:
		return EQ
	case token.GT, token.LT:
		return LESSGREATER
	case token.PLUS, token.MINUS, token.DEC, token.INC:
		return SUM
	case token.ASTERISK, token.SLASH:
		return PRODUCT
	case token.LPAREN:
		return CALL
	case token.LBRACKET:
		return INDEX
	case token.RANGE_ARRAY:
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
