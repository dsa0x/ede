package parser

import (
	"ede/ast"
	"ede/token"
)

func (p *Parser) parseArrayLiteral() ast.Expression {
	expr := &ast.ArrayLiteral{Token: p.currToken}

	if !p.advanceCurrTokenIs(token.LBRACKET) {
		p.addError(unexpectedTokenError(token.LBRACKET, p.currToken.Literal))
		return nil
	}

	// if the first element of the array has a unary operator,
	// we save it, and wrap the parse element in a prefix expression
	// This is so that the look-ahead can still see the RANGE_ARRAY token
	var unary token.Token
	if p.currTokenIs(token.MINUS) || p.currTokenIs(token.PLUS) {
		unary = p.currToken
		p.advanceToken()
	}

	// if the array is from a range array e.g. let arr = [1..10]
	if p.nextTokenIs(token.RANGE_ARRAY) {
		start := p.parseExpr(p.precedence(token.RANGE_ARRAY) + 1) // so it parses only the left side of RANGE_ARRAY
		if unary != (token.Token{}) {
			start = &ast.PrefixExpression{Operator: unary.Literal, Right: start, Token: unary}
		}
		return p.parseRangeArray(start)
	}
	// else if it is a normal array literal
	expr.Elements = p.parseArguments(token.RBRACKET)
	if expr.Elements == nil || !p.currTokenIs(token.RBRACKET) {
		p.addError("expected closing bracket token ']', got '%s'", p.currToken.Literal)
		return nil
	}
	p.advanceToken() // eat closing token

	if unary != (token.Token{}) {
		expr.Elements[0] = &ast.PrefixExpression{Operator: unary.Literal, Right: expr.Elements[0], Token: unary}
	}

	return expr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Left: left, Token: p.currToken}

	if !p.advanceCurrTokenIs(token.LBRACKET) { // eat starting token
		p.addError(unexpectedTokenError(token.LBRACKET, p.currToken.Literal))
		return nil
	}
	expr.Index = p.parseExpr(LOWEST)
	if !p.advanceCurrTokenIs(token.RBRACKET) { // eat closing bracket
		p.addError(unexpectedTokenError(token.RBRACKET, p.currToken.Literal))
		return nil
	}
	return expr
}
