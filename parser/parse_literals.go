package parser

import (
	"ede/ast"
	"strconv"
)

func (p *Parser) parseInteger() ast.Expression {
	expr := &ast.IntegerLiteral{Token: p.currToken}

	val, _ := strconv.Atoi(p.currToken.Literal)
	expr.Value = int64(val)
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
