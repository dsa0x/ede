package parser

import (
	"ede/ast"
	"strconv"
)

func (p *Parser) parseInteger() ast.Expression {
	expr := &ast.IntegerLiteral{Token: p.currToken}

	val, _ := strconv.Atoi(p.currToken.Literal)
	expr.Value = int64(val)
	return expr
}

func (p *Parser) parseFloat() ast.Expression {
	expr := &ast.FloatLiteral{Token: p.currToken}

	val, _ := strconv.ParseFloat(p.currToken.Literal, 64)
	expr.Value = val
	return expr
}

func (p *Parser) parseBool() ast.Expression {
	expr := &ast.BooleanLiteral{Token: p.currToken}

	val, _ := strconv.ParseBool(p.currToken.Literal)
	expr.Value = val
	return expr
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Value: p.currToken.Literal, Token: p.currToken}
}
