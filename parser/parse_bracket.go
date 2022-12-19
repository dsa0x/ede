package parser

import (
	"ede/ast"
	"ede/token"
)

func (p *Parser) parseArrayLiteral() ast.Expression {
	expr := &ast.ArrayLiteral{Token: p.currToken}

	if !p.advanceCurrTokenIs(token.LBRACKET) {
		return nil
	}

	// if the array is from a range array e.g. let arr = [1..10]
	if p.nextTokenIs(token.RANGE_ARRAY) {
		start := p.parseInteger()
		p.advanceToken() // advance the range array
		return p.parseRangeArray(start)
	}
	// else if it is a normal array literal
	expr.Elements = p.parseArguments(token.RBRACKET)
	if expr.Elements == nil {
		p.addError("expected closing bracket token ']', got '%s'", p.currToken.Literal)
	}
	p.advanceNextTokenIs(token.SEMICOLON)
	return expr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Left: left, Token: p.currToken}

	if !p.advanceCurrTokenIs(token.LBRACKET) {
		return nil
	}
	expr.Index = p.parseExpr(LOWEST)
	if !p.advanceNextTokenIs(token.RBRACKET) {
		return nil
	}
	p.advanceNextTokenIs(token.SEMICOLON)
	return expr
}
