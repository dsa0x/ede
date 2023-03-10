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
	p.eatEndToken()
	return expr
}

func (p *Parser) parseSetExpression(startPos token.Pos) ast.Expression {
	expr := &ast.SetLiteral{Token: p.prevToken}
	set := make(map[ast.Expression]struct{})
	elements := p.parseArguments(token.RBRACE)
	for _, el := range elements {
		if val := getRawValue(el); val == nil {
			p.addError("invalid set entry. entry is of type %T", el)
			return nil
		}
		set[el] = struct{}{}
	}
	if !p.advanceCurrTokenIs(token.RBRACE) {
		p.addError("expected closing parenthesis token '}', got '%s'", p.currToken.Literal)
		return nil
	}
	expr.Elements = set
	return expr
}

func (p *Parser) parseObjectMethodExpression(obj ast.Expression) ast.Expression {
	expr := &ast.ObjectMethodExpression{Token: p.currToken, Object: obj}
	p.advanceToken()
	if !p.currTokenIs(token.IDENT) {
		p.addError(unexpectedTokenError(token.IDENT, p.currToken.Literal))
		return nil
	}

	// if it is a method call
	if p.nextTokenIs(token.LPAREN) {
		// subtract 1, so that we can parse the closing parenthesis, and we stop there.
		expr.Method = p.parseExpr(CALL - 1)
	} else {
		expr.Method = p.parseIdent()
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
	if left == nil {
		// if it is nil, then a parse error should have been added to the internal list
		return nil
	}
	operator := p.currToken
	inf := &ast.InfixExpression{
		Operator: operator.Literal,
		Left:     left,
		Token:    operator,
	}

	operatorPrecedence := p.currPrecedence()

	p.advanceToken()
	if slices.Contains([]token.TokenType{token.LBRACE, token.LBRACKET}, left.TokenType()) {
		p.addError("invalid left expression %s for operator '%s'", left.TokenType(), operator.Literal)
		return nil
	}

	if inf.Right = p.parseExpr(operatorPrecedence); inf.Right == nil {
		// if we couldn't parse the right
		p.addError("invalid right expression %s for operator '%s'", p.currToken.Literal, operator.Literal)
		p.advanceToken()
		return nil
	}
	return inf
}

func (p *Parser) parseIdent() ast.Expression {
	expr := &ast.Identifier{Value: p.currToken.Literal, Token: p.currToken}
	p.advanceToken()
	return expr
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
	switch ident := ident.(type) {
	case *ast.Identifier:
		expr = &ast.ReassignmentStmt{Name: ident, Token: p.currToken}
		if token.IsReservedKeyword(expr.Token.Literal) {
			p.addError("unexpected assignment of reserved keyword %s", expr.Token.Literal)
			return nil
		}
	case *ast.IndexExpression:
		expr = &ast.ReassignmentStmt{Name: ident, Token: p.currToken}
	default:
		p.addError("unexpected token assignment: %s", ident.Literal())
		return nil
	}

	if !p.advanceCurrTokenIs(token.ASSIGN) {
		return nil
	}
	expr.Expr = p.parseExpr(LOWEST)
	return expr
}

func (p *Parser) parsePlusEqual(ident ast.Expression) ast.Expression {
	var expr *ast.ReassignmentStmt
	if id, ok := ident.(*ast.Identifier); ok {
		expr = &ast.ReassignmentStmt{Name: id, Token: p.currToken}
	} else {
		p.addError("unexpected token assignment: %s", ident.Literal())
		return nil
	}

	if token.IsReservedKeyword(expr.Token.Literal) {
		p.addError("unexpected assignment of reserved keyword %s", expr.Token.Literal)
		return nil
	}

	if !p.advanceCurrTokenIs(token.PLUS_EQUAL) {
		p.addError(unexpectedTokenError(token.PLUS_EQUAL, p.currToken.Literal))
		return nil
	}

	expr.Expr = &ast.InfixExpression{Left: ident,
		Right:    p.parseExpr(LOWEST),
		Operator: token.PLUS,
	}
	return expr
}

func (p *Parser) parseMinusEqual(ident ast.Expression) ast.Expression {
	var expr *ast.ReassignmentStmt
	if id, ok := ident.(*ast.Identifier); ok {
		expr = &ast.ReassignmentStmt{Name: id, Token: p.currToken}
	} else {
		p.addError("unexpected token assignment: %s", ident.Literal())
		return nil
	}

	if token.IsReservedKeyword(expr.Token.Literal) {
		p.addError("unexpected assignment of reserved keyword %s", expr.Token.Literal)
		return nil
	}

	if !p.advanceCurrTokenIs(token.MINUS_EQUAL) {
		p.addError(unexpectedTokenError(token.MINUS_EQUAL, p.currToken.Literal))
		return nil
	}

	expr.Expr = &ast.InfixExpression{Left: ident,
		Right:    p.parseExpr(LOWEST),
		Operator: token.MINUS,
	}
	return expr
}

// getInteger is a helper function to retrieve an integer
func getInteger(expr ast.Expression) (int64, bool) {
	switch expr := expr.(type) {
	case *ast.IntegerLiteral:
		return int64(expr.Value), true
	case *ast.PrefixExpression:
		if integer, ok := expr.Right.(*ast.IntegerLiteral); ok {
			if expr.Operator == token.PLUS {
				return int64(integer.Value), true
			} else if expr.Operator == token.MINUS {
				return int64(-integer.Value), true
			}
		}
	}
	return 0, false
}

func (p *Parser) parseRangeArray(start ast.Expression) ast.Expression {
	expr := &ast.RangeArrayLiteral{Token: p.currToken, Start: start}
	p.advanceToken()
	expr.End = p.parseExpr(LOWEST)
	p.advanceToken()
	return expr
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Function: fn, Token: p.currToken}

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
	expr := &ast.HashLiteral{Token: p.currToken, Pair: make(map[ast.Expression]ast.Expression)}

	keySet := map[any]ast.Expression{}
	if !p.advanceCurrTokenIs(token.LBRACE) { // eat opening token
		return nil
	}

	if !p.nextTokenIs(token.COLON) { // if it is a set (e.g. {1,2,3})
		return p.parseSetExpression(expr.Token.Pos)
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

// getRawValue returns the raw value of an expression
// typically, these are the literals that can be comparable and also easily stringified (i.e. string, int, bool)
func getRawValue(expr ast.Expression) any {
	switch expr := expr.(type) {
	case *ast.StringLiteral:
		return expr.Value
	case *ast.IntegerLiteral:
		return expr.Value
	case *ast.BooleanLiteral:
		return expr.Value
	case *ast.Identifier:
		return expr
	}
	return nil
}
