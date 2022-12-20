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
		return p.parseRangeArray(start)
	}
	// else if it is a normal array literal
	expr.Elements = p.parseArguments(token.RBRACKET)
	if expr.Elements == nil || !p.currTokenIs(token.RBRACKET) {
		p.addError("expected closing bracket token ']', got '%s'", p.currToken.Literal)
		return nil
	}
	p.advanceToken() // eat closing token
	return expr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Left: left, Token: p.currToken}

	if !p.advanceCurrTokenIs(token.LBRACKET) { // eat starting token
		return nil
	}
	expr.Index = p.parseExpr(LOWEST)
	if !p.advanceCurrTokenIs(token.RBRACKET) { // eat closing bracket
		return nil
	}
	return expr
}