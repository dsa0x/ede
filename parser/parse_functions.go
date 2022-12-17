package parser

import (
	"ede/ast"
	"ede/token"
	"strconv"
)

func (p *Parser) parseInteger() ast.Expression {
	expr := &ast.IntegerLiteral{}

	val, _ := strconv.Atoi(p.currToken.Literal)
	expr.Value = int64(val)
	return expr
}

func (p *Parser) parseBool() ast.Expression {
	expr := &ast.BooleanLiteral{}

	val, _ := strconv.ParseBool(p.currToken.Literal)
	expr.Value = val
	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceToken()
	expr := p.parseExpr(LOWEST)
	if !p.advanceNextTokenIs(token.RPAREN) {
		return nil
	}
	return expr
}
func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Operator: p.currToken.Literal,
		Token:    p.currToken,
	}

	p.advanceToken()
	expr.Right = p.parseExpr(PREFIX)
	return expr
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.PostfixExpression{
		Operator: p.currToken.Literal,
		Token:    p.currToken,
	}

	p.advanceToken()

	expr.Left = left
	return expr
}

func (p *Parser) parseInfixOperator(left ast.Expression) ast.Expression {
	inf := &ast.InfixExpression{
		Operator: p.currToken.Literal,
		Left:     left,
		Token:    p.currToken,
	}

	currPrecedence := p.currPrecedence()

	p.advanceToken()
	inf.Right = p.parseExpr(currPrecedence)
	return inf
}

func (p *Parser) parseIdent() ast.Expression {
	tok := p.currToken
	return &ast.Identifier{Value: tok.Literal, Token: tok}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	tok := p.currToken
	return &ast.StringLiteral{Value: tok.Literal, Token: tok}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	stmt := &ast.FunctionLiteral{Token: p.currToken}
	if !p.advanceNextTokenIs(token.LPAREN) {
		return nil
	}
	p.advanceToken()
	stmt.Params = p.parseFunctionParams()

	if !p.advanceCurrTokenIs(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStmt()

	p.advanceNextTokenIs(token.SEMICOLON)
	return stmt
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	identifiers := make([]*ast.Identifier, 0)

	// if no function params
	if p.advanceCurrTokenIs(token.RPAREN) {
		return identifiers
	}

	for p.currTokenIs(token.IDENT) {
		identifiers = append(identifiers, &ast.Identifier{Value: p.currToken.Literal, Token: p.currToken})

		if !p.advanceNextTokenIs(token.COMMA) && !p.advanceNextTokenIs(token.RPAREN) {
			return nil
		}
		p.advanceToken()
	}

	return identifiers
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Function: fn, Token: p.currToken}

	if !p.advanceCurrTokenIs(token.LPAREN) {
		return nil
	}
	expr.Args = p.parseArguments(token.RPAREN)
	p.advanceNextTokenIs(token.SEMICOLON)
	return expr
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	expr := &ast.ArrayLiteral{Token: p.currToken}

	if !p.advanceCurrTokenIs(token.LBRACKET) {
		return nil
	}
	expr.Elements = p.parseArguments(token.RBRACKET)
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

func (p *Parser) parseArguments(endToken token.TokenType) []ast.Expression {
	exprs := make([]ast.Expression, 0)

	// if no function args
	if p.advanceCurrTokenIs(endToken) {
		return exprs
	}

	for !p.currTokenIs(token.EOF) && !p.currTokenIs(endToken) {
		exprs = append(exprs, p.parseExpr(LOWEST))

		p.advanceNextTokenIs(token.COMMA)
		p.advanceToken()
	}

	return exprs
}
