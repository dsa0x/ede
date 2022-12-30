package parser

import (
	"ede/ast"
	"ede/token"
	"fmt"
	"strconv"
)

func (p *Parser) parseInteger() ast.Expression {
	literal := int64(1)
	switch p.currToken.Type {
	case token.MINUS:
		literal = -1
		p.advanceToken()
	case token.PLUS:
		p.advanceToken()
	case token.INT, token.IDENT:
		break
	default:
		err := NewParseError(fmt.Errorf("invalid prefix operator %s for integer", p.currToken.Literal), p.currPos())
		p.appendError(err)
		return nil
	}
	expr := &ast.IntegerLiteral{Token: p.currToken}

	val, err := strconv.Atoi(p.currToken.Literal)
	if err != nil {
		return nil
	}
	expr.Value = int64(val) * literal
	p.advanceToken()
	return expr
}

func (p *Parser) parseFloat() ast.Expression {
	expr := &ast.FloatLiteral{Token: p.currToken}

	val, _ := strconv.ParseFloat(p.currToken.Literal, 64)
	expr.Value = val
	p.advanceToken()
	return expr
}

func (p *Parser) parseBool() ast.Expression {
	expr := &ast.BooleanLiteral{Token: p.currToken}

	val, _ := strconv.ParseBool(p.currToken.Literal)
	expr.Value = val
	p.advanceToken()
	return expr
}

func (p *Parser) parseStringLiteral() ast.Expression {
	expr := &ast.StringLiteral{Value: p.currToken.Literal, Token: p.currToken}
	p.advanceToken()
	return expr
}
