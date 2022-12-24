package parser

import (
	"ede/ast"
	"ede/token"
	"fmt"

	"golang.org/x/exp/slices"
)

var (
	// token types that can be start of a statement/expression
	startTokens = []token.TokenType{token.IDENT,
		token.FUNCTION,
		token.TRUE,
		token.FALSE,
		token.INT,
		token.SINGLE_COMMENT,
	}
)

func (p *Parser) parseStmt() ast.Statement {
	defer p.eatEndToken()
	p.eatEndToken()
	switch p.currToken.Type {
	case token.LET:
		if token.IsReservedKeyword(p.nextToken.Literal) {
			p.addError(fmt.Sprintf("cannot assign to reserved keyword '%s'", p.nextToken.Literal))
			return nil
		}
		return p.parseLetStmt()
	case token.IF:
		return p.parseIfStmt()
	case token.RETURN:
		expr := p.parseReturnExpr()
		if expr, ok := expr.(*ast.ReturnExpression); ok {
			return expr
		}
	case token.FOR:
		return p.parseForStmt()
	case token.NEWLINE, token.SEMICOLON:
		p.advanceToken()
		return p.parseStmt()
	case token.SINGLE_COMMENT:
		stmt := &ast.CommentStmt{Token: p.currToken, Value: p.currToken.Literal, ValuePos: p.pos}
		p.advanceToken()
		return stmt
	case token.IMPORT:
		return p.parseImportStmt()
	}
	return p.parseExpressionStmt()
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	stmt := &ast.LetStmt{Token: p.currToken, ValuePos: p.pos}
	if !p.advanceNextTokenIs(token.IDENT) { // eat LET token
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal, ValuePos: p.pos}

	if p.nextTokenIs(token.SEMICOLON) || p.nextTokenIs(token.NEWLINE) { // e.g let foo;
		p.advanceToken()
		return stmt
	}

	if !p.advanceNextTokenIs(token.ASSIGN) { // advance to =
		return nil
	}

	p.advanceToken() // go to RHS
	stmt.Expr = p.parseExpr(LOWEST)
	return stmt
}
func (p *Parser) parseReturnExpr() ast.Expression {
	stmt := &ast.ReturnExpression{Token: p.currToken, ValuePos: p.pos}
	if !p.advanceCurrTokenIs(token.RETURN) {
		return nil
	}

	stmt.Expr = p.parseExpr(LOWEST)
	return stmt
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	stmt := &ast.IfStmt{Token: p.currToken, ValuePos: p.pos}
	stmt.Alternatives = make([]*ast.ConditionalStmt, 0)
	if !p.advanceNextTokenIs(token.LPAREN) { // eat opening token IF
		p.addError("unexpected token %s, want %s", p.nextToken.Literal, token.LPAREN)
		return nil
	}
	p.advanceToken() // eat ( opener of expr
	stmt.Condition = p.parseExpr(LOWEST)

	if !p.advanceCurrTokenIs(token.RPAREN) {
		return nil
	}
	if !p.advanceCurrTokenIs(token.LBRACE) {
		return nil
	}

	stmt.Consequence = &ast.ConditionalStmt{
		Condition: stmt.Condition,
		Statement: p.parseBlockStmt(),
		ValuePos:  p.pos,
	}

	if p.currTokenIs(token.ELSE) {
		for p.advanceCurrTokenIs(token.ELSE) {
			// else if
			if p.currTokenIs(token.IF) {
				p.advanceToken()

				if !p.advanceCurrTokenIs(token.LPAREN) {
					p.addError(unexpectedTokenError(token.LPAREN, p.currToken.Literal))
					return nil
				}

				condition := p.parseExpr(LOWEST) // parse expression inside if()

				if !p.advanceCurrTokenIs(token.RPAREN) {
					p.addError(unexpectedTokenError(token.RPAREN, p.currToken.Literal))
					return nil
				}
				if !p.advanceCurrTokenIs(token.LBRACE) {
					p.addError(unexpectedTokenError(token.LBRACE, p.currToken.Literal))
					return nil
				}

				elifStmt := p.parseStmt()

				condStmt := &ast.ConditionalStmt{Condition: condition, Statement: elifStmt, ValuePos: p.pos}
				stmt.Alternatives = append(stmt.Alternatives, condStmt)

				if !p.advanceCurrTokenIs(token.RBRACE) {
					p.addError(unexpectedTokenError(token.RBRACE, p.currToken.Literal))
					return nil
				}
			} else {
				if !p.advanceCurrTokenIs(token.LBRACE) {
					p.addError(unexpectedTokenError(token.LBRACE, p.currToken.Literal))
					return nil
				}
				elseStmt := p.parseStmt()
				condStmt := &ast.ConditionalStmt{Statement: elseStmt, ValuePos: p.pos}
				stmt.Alternatives = append(stmt.Alternatives, condStmt)

				if !p.advanceCurrTokenIs(token.RBRACE) {
					p.addError(unexpectedTokenError(token.RBRACE, p.currToken.Literal))
					return nil
				}
				return stmt // return after an else statement
			}
		}

	}

	return stmt
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	blockStmt := &ast.BlockStmt{Statements: make([]ast.Statement, 0), ValuePos: p.pos}

	for !p.currTokenIs(token.EOF) && !p.currTokenIs(token.RBRACE) {
		if stmt := p.parseStmt(); stmt != nil {
			blockStmt.Statements = append(blockStmt.Statements, stmt)
		}
		if p.Errors() != nil {
			return nil
		}

		p.eatEndToken()
	}
	if !p.advanceCurrTokenIs(token.RBRACE) {
		p.addError(unexpectedTokenError(token.RBRACE, p.currToken.Literal))
		return nil
	}

	p.eatEndToken()
	return blockStmt
}

func (p *Parser) parseImportStmt() *ast.ImportStmt {
	stmt := &ast.ImportStmt{ValuePos: p.pos, Token: p.currToken}
	p.advanceToken()
	if !p.currTokenIs(token.IDENT) {
		p.addError("invalid import statement")
		return nil
	}
	stmt.Value = p.currToken.Literal
	p.advanceToken()
	return stmt
}

func (p *Parser) parseMatchExpression() ast.Expression {
	stmt := &ast.MatchExpression{ValuePos: p.pos, Token: p.currToken, Cases: make([]ast.MatchCase, 0)}
	if !p.advanceCurrTokens(token.MATCH, token.LPAREN) {
		p.addError(unexpectedTokenError(token.LPAREN, string(p.currToken.Type))) // TODO: may be incorrect
		return nil
	}

	stmt.Expression = p.parseExpr(LOWEST)

	if !p.advanceCurrTokens(token.RPAREN, token.LBRACE) {
		p.addError(unexpectedTokenError(token.RBRACE, p.currToken.Literal))
		return nil
	}

	// function to parse each case of the match block
	parseMatchCase := func() ast.Expression {
		if !p.advanceCurrTokenIs(token.COLON) {
			p.addError(unexpectedTokenError(token.COLON, string(p.currToken.Type)))
			return nil
		}
		expr := p.parseExpr(LOWEST)
		p.eatEndToken()
		return expr
	}

outerloop:
	for {
		switch p.currToken.Type {
		case token.RBRACE: // single case
			p.advanceToken()
			return stmt
		case token.CASE:
			p.advanceToken()
			pattern := p.parseExpr(LOWEST)
			matchCase := ast.MatchCase{Pattern: pattern, Output: parseMatchCase()}
			stmt.Cases = append(stmt.Cases, matchCase)
			if matchCase.Pattern == nil || matchCase.Output == nil {
				return nil
			}
		case token.DEFAULT:
			p.advanceToken()
			if stmt.Default = parseMatchCase(); stmt.Default == nil {
				return nil
			}
		default:
			break outerloop
		}
	}

	p.advanceToken()
	return stmt
}

func (p *Parser) parseForStmt() ast.Statement {
	forLoopStmt := &ast.ForLoopStmt{Token: p.currToken, ValuePos: p.pos}

	if !p.advanceCurrTokenIs(token.FOR) {
		p.addError(unexpectedTokenError(token.FOR, p.currToken.Literal))
		return nil
	}
	if !p.currTokenIs(token.IDENT) {
		p.addError(unexpectedTokenError(token.IDENT, p.currToken.Literal))
		return nil
	}
	forLoopStmt.Variable = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal, ValuePos: p.pos}
	p.advanceToken()
	if !p.advanceCurrTokenIs(token.ASSIGN) {
		p.addError(unexpectedTokenError(token.ASSIGN, p.currToken.Literal))
		return nil
	}
	switch p.currToken.Type {
	case token.RANGE:
		p.advanceToken()
		if !p.currTokenIs(token.LBRACKET) && !p.currTokenIs(token.IDENT) {
			p.addError(expectAfterTokenErrorStr("array", token.RANGE, p.nextToken.Literal))
			return nil
		}
		startingToken := p.currToken
		if startingToken.Type != token.IDENT { // ident does not need to advance
			p.advanceToken()
		}

		expr := p.parseExpr(LOWEST)

		if startingToken.Type == token.LBRACKET { // this means the boundary is a literal
			if expr == nil {
				return nil
			}
		}

		forLoopStmt.Boundary = expr

		if !p.advanceCurrTokenIs(token.LBRACE) {
			err := expectAfterTokenErrorStr(token.LBRACE, "right bracket ']'", p.currToken.Literal)
			p.addError(err)
			return nil
		}

		forLoopStmt.Statement = p.parseBlockStmt()

	default:
		p.addError("error")
		return nil
	}
	return forLoopStmt
}

func (p *Parser) parseExpressionStmt() *ast.ExpressionStmt {
	if !slices.Contains(startTokens, token.LookupIdent(p.currToken.Literal)) {
		p.addError("expected start of expression, found '%s'", p.currToken.Literal)
		return nil
	}
	stmt := &ast.ExpressionStmt{Token: p.currToken, ValuePos: p.pos}
	stmt.Expr = p.parseExpr(LOWEST)

	return stmt
}
