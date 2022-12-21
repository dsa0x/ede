package parser

import (
	"ede/ast"
	"ede/token"

	"golang.org/x/exp/slices"
)

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceToken() // eat opening token
	expr := p.parseExpr(LOWEST)
	if !p.advanceCurrTokenIs(token.RPAREN) {
		p.addError("expected closing parenthesis token ')', got '%s'", p.currToken.Literal)
		return nil
	}
	for {
		switch p.currToken.Type {
		case token.SEMICOLON, token.NEWLINE:
			p.advanceToken()
		default:
			return expr
		}
	}
}

func (p *Parser) parseMethodExpression(obj ast.Expression) ast.Expression {
	expr := &ast.ObjectMethodExpression{Token: p.currToken, Object: obj, ValuePos: p.pos}
	p.advanceToken()
	if !p.currTokenIs(token.IDENT) {
		return nil
	}
	expr.Method = p.parseExpr(LOWEST)
	return expr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Operator: p.currToken.Literal,
		Token:    p.currToken,
		ValuePos: p.pos,
	}

	p.advanceToken()
	expr.Right = p.parseExpr(PREFIX)
	return expr
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.PostfixExpression{
		Operator: p.currToken.Literal,
		Token:    p.currToken,
		ValuePos: p.pos,
	}

	p.advanceToken()

	expr.Left = left
	return expr
}

func (p *Parser) parseInfixOperator(left ast.Expression) ast.Expression {
	if left == nil {
		// if it is nil, then a parse error should have been added to the internal list
		return nil
	}
	operator := p.currToken
	inf := &ast.InfixExpression{
		Operator: operator.Literal,
		Left:     left,
		Token:    operator,
		ValuePos: p.pos,
	}

	currPrecedence := p.currPrecedence()

	p.advanceToken()
	if slices.Contains([]token.TokenType{token.LBRACE, token.LBRACKET}, left.TokenType()) {
		p.addError("invalid left expression %s for operator '%s'", left.TokenType(), operator.Literal)
		return nil
	}

	if inf.Right = p.parseExpr(currPrecedence); inf.Right == nil {
		// if we couldn't parse the right
		p.addError("invalid right expression %s for operator '%s'", p.currToken.Literal, operator.Literal)
		return nil
	}
	return inf
}

func (p *Parser) parseIdent() ast.Expression {
	expr := &ast.Identifier{Value: p.currToken.Literal, Token: p.currToken, ValuePos: p.pos}
	p.advanceToken()
	return expr
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	stmt := &ast.FunctionLiteral{Token: p.currToken, ValuePos: p.pos}
	if !p.advanceNextTokenIs(token.LPAREN) {
		return nil
	}
	p.advanceToken()
	stmt.Params = p.parseFunctionParams()

	if !p.advanceCurrTokenIs(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStmt()

	// if the literal is called immediately
	if p.currTokenIs(token.LPAREN) {
		return p.parseCallExpression(stmt)
	}

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

func (p *Parser) parseReassignment(ident ast.Expression) ast.Expression {
	var expr *ast.ReassignmentStmt
	if id, ok := ident.(*ast.Identifier); ok {
		expr = &ast.ReassignmentStmt{Name: id, Token: p.currToken, ValuePos: p.pos}
	} else {
		p.addError("unexpected token assignment: %s", ident.Literal())
		return nil
	}

	if token.IsReservedKeyword(expr.Token.Literal) {
		p.addError("unexpected assignment of reserved keyword %s", expr.Token.Literal)
		return nil
	}

	if !p.advanceCurrTokenIs(token.ASSIGN) {
		return nil
	}
	expr.Expr = p.parseExpr(LOWEST)
	p.advanceNextToEndToken()
	return expr
}

func (p *Parser) parseRangeArray(start ast.Expression) ast.Expression {
	startL, ok := start.(*ast.IntegerLiteral)
	if !ok {
		return nil
	}
	expr := &ast.ArrayLiteral{Token: p.currToken, Elements: make([]ast.Expression, 0), ValuePos: p.pos}
	if !p.advanceNextTokenIs(token.INT) {
		return nil
	}
	end := p.parseExpr(LOWEST)
	endL, ok := end.(*ast.IntegerLiteral)
	if !ok {
		return nil
	}
	for i := startL.Value; i <= endL.Value; i++ {
		expr.Elements = append(expr.Elements, &ast.IntegerLiteral{Value: i, ValuePos: p.pos})
	}
	if !p.currTokenIs(token.LBRACE) && !p.currTokenIs(token.RBRACKET) { // TODO: usage in forloop and array literal should be diff
		return nil
	}
	p.advanceToken()
	return expr
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Function: fn, Token: p.currToken, ValuePos: p.pos}

	if !p.advanceCurrTokenIs(token.LPAREN) { // eat opening token
		return nil
	}
	expr.Args = p.parseArguments(token.RPAREN)
	if expr.Args == nil {
		p.addError("expected closing parenthesis token ')', got '%s'", p.currToken.Literal)
	}
	if !p.advanceCurrTokenIs(token.RPAREN) { // eat closing token
		return nil
	}
	return expr
}

func (p *Parser) parseHashLiteral() ast.Expression {
	expr := &ast.HashLiteral{Token: p.currToken, ValuePos: p.pos, Pair: make(map[ast.Expression]ast.Expression)}

	keySet := map[any]ast.Expression{}
	if !p.advanceCurrTokenIs(token.LBRACE) { // eat opening token
		return nil
	}
	for !p.currTokenIs(token.EOF) && !p.currTokenIs(token.RBRACE) {
		key := p.parseExpr(LOWEST)
		rawValue := getRawValue(key)
		if rawValue == nil {
			p.addError("invalid type %T for hash key", key)
			return nil
		}
		// if the key exists, delete it so it is overwritten by the new token
		if key, ok := keySet[rawValue]; ok {
			delete(expr.Pair, key)
		}
		keySet[rawValue] = key
		if !p.advanceCurrTokenIs(token.COLON) {
			return nil
		}
		expr.Pair[key] = p.parseExpr(LOWEST)

		p.advanceCurrTokenIs(token.COMMA)
	}
	if !p.advanceCurrTokenIs(token.RBRACE) {
		p.addError("unexpected end of token. expected }, got %s", p.nextToken.Literal)
		return nil
	}
	return expr
}

func (p *Parser) parseArguments(endToken token.TokenType) []ast.Expression {
	// parseArguments will not meet the opening brace token and should not advance the closing brace
	exprs := make([]ast.Expression, 0)

	// if no args
	if p.currTokenIs(endToken) {
		return exprs
	}

	for !p.currTokenIs(token.EOF) && !p.currTokenIs(endToken) {
		exprs = append(exprs, p.parseExpr(LOWEST))

		if !p.currTokenIs(token.COMMA) && !p.currTokenIs(endToken) {
			p.addError("unexpected end of token. expected %s, got %s", endToken, p.nextToken.Literal)
			return nil
		}
		p.advanceCurrTokenIs(token.COMMA) // advance to next expr if comma
	}

	return exprs
}

func getRawValue(expr ast.Expression) any {
	switch expr := expr.(type) {
	case *ast.StringLiteral:
		return expr.Value
	case *ast.IntegerLiteral:
		return expr.Value
	case *ast.BooleanLiteral:
		return expr.Value
	}
	return nil
}
