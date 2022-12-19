package parser

import (
	"ede/ast"
	"ede/token"
	"fmt"

	"golang.org/x/exp/slices"
)

var mathTokens = []token.TokenType{token.PLUS, token.MINUS, token.ASTERISK, token.SLASH}
var literalTokens = []token.TokenType{token.INT, token.FLOAT, token.STRING, token.IDENT}

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

	currPrecedence := p.currPrecedence()

	p.advanceToken()
	// unless the expression is an identifier, in math-like operators,
	// the left must be of same type as the right
	if slices.Contains(mathTokens, operator.Type) {
		if left.TokenType() != token.IDENT && left.TokenType() != p.currToken.Type {
			p.errors = append(p.errors, fmt.Errorf("left and right expressions do not match for operator '%s'", operator.Literal))
			return nil
		}
	}
	inf.Right = p.parseExpr(currPrecedence)
	return inf
}

func (p *Parser) parseIdent() ast.Expression {
	tok := p.currToken
	return &ast.Identifier{Value: tok.Literal, Token: tok}
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

func (p *Parser) parseReassignment(ident ast.Expression) ast.Expression {
	var expr *ast.ReassignmentStmt
	if id, ok := ident.(*ast.Identifier); ok {
		expr = &ast.ReassignmentStmt{Name: id, Token: p.currToken}
	} else {
		p.errors = append(p.errors, fmt.Errorf("unexpected token assignment: %s", ident.Literal()))
		return &ast.ErrorStmt{Value: fmt.Sprintf("unexpected token assignment: %s", ident.Literal())}
	}

	if token.IsReservedKeyword(expr.Token.Literal) {
		return &ast.ErrorStmt{Value: fmt.Sprintf("unexpected assignment of reserved keyword %s", expr.Token.Literal)}
	}

	if !p.advanceCurrTokenIs(token.ASSIGN) {
		return nil
	}
	expr.Expr = p.parseExpr(LOWEST)
	p.advanceNextTokenIs(token.SEMICOLON)
	return expr
}

func (p *Parser) parseRangeArray(start ast.Expression) ast.Expression {
	startL, ok := start.(*ast.IntegerLiteral)
	if !ok {
		return nil
	}
	expr := &ast.ArrayLiteral{Token: p.currToken, Elements: make([]ast.Expression, 0)}
	if !p.advanceNextTokenIs(token.INT) {
		return nil
	}
	end := p.parseExpr(LOWEST)
	endL, ok := end.(*ast.IntegerLiteral)
	if !ok {
		return nil
	}
	for i := startL.Value; i <= endL.Value; i++ {
		expr.Elements = append(expr.Elements, &ast.IntegerLiteral{Value: i})
	}
	if !p.nextTokenIs(token.LBRACE) && !p.nextTokenIs(token.RBRACKET) { // TODO: usage in forloop and array literal should be diff
		return nil
	}
	p.advanceToken()
	return expr
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

	// if the array uses range array e.g. [1..10]
	if p.nextTokenIs(token.RANGE_ARRAY) {
		start := p.parseInteger()
		p.advanceToken() // advance the range array
		return p.parseRangeArray(start)
	}
	// else if it is a normal array literal
	expr.Elements = p.parseArguments(token.RBRACKET)
	p.advanceNextTokenIs(token.SEMICOLON)
	return expr
}

func (p *Parser) parseHashLiteral() ast.Expression {
	expr := &ast.HashLiteral{Token: p.currToken, Pair: make(map[ast.Expression]ast.Expression)}

	keySet := map[any]ast.Expression{}
	if !p.advanceCurrTokenIs(token.LBRACE) {
		return nil
	}
	for !p.currTokenIs(token.EOF) && !p.currTokenIs(token.RBRACE) {
		key := p.parseExpr(LOWEST)
		rawValue := getRawValue(key)
		if rawValue == nil {
			p.errors = append(p.errors, fmt.Errorf("invalid type %T for hash key", key))
			return nil
		}
		// if the key exists, delete it so it is overwritten by the new token
		if key, ok := keySet[rawValue]; ok {
			delete(expr.Pair, key)
		}
		keySet[rawValue] = key
		p.advanceToken()
		if !p.advanceCurrTokenIs(token.COLON) {
			return nil
		}
		expr.Pair[key] = p.parseExpr(LOWEST)

		p.advanceNextTokenIs(token.COMMA)
		p.advanceToken()
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

func (p *Parser) parseArguments(endToken token.TokenType) []ast.Expression {
	exprs := make([]ast.Expression, 0)

	// if no args
	if p.currTokenIs(endToken) {
		return exprs
	}

	for !p.currTokenIs(token.EOF) && !p.currTokenIs(endToken) {
		exprs = append(exprs, p.parseExpr(LOWEST))

		if !p.advanceNextTokenIs(token.COMMA) && !p.nextTokenIs(endToken) {
			p.errors = append(p.errors, fmt.Errorf("unexpected end of token. expected %s, got %s", endToken, p.nextToken.Literal))
			return nil
		}
		p.advanceToken()
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

type infixPattern [3]token.TokenType

var infixPatterns = []infixPattern{
	{token.INT, token.PLUS, token.INT},
	{token.INT, token.GT, token.INT},
	{token.INT, token.LT, token.INT},
}
