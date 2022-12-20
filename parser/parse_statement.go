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
	for p.advanceCurrToEndToken() {
	}
	switch p.currToken.Type {
	case token.LET:
		if token.IsReservedKeyword(p.nextToken.Literal) {
			return &ast.ErrorStmt{Value: fmt.Sprintf("cannot assign to reserved keyword '%s'", p.nextToken.Literal)}
		}
		return p.parseLetStmt()
	case token.IF:
		return p.parseIfStmt()
	case token.RETURN:
		return p.parseReturnExpr()
	case token.FOR:
		return p.parseForStmt()
	case token.NEWLINE, token.SEMICOLON:
		p.advanceToken()
		return p.parseStmt()
	case token.SINGLE_COMMENT:
		return &ast.CommentStmt{Token: p.currToken, Value: p.currToken.Literal}
	}
	return p.parseExpressionStmt()
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	stmt := &ast.LetStmt{Token: p.currToken}
	if !p.advanceNextTokenIs(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if p.advanceNextTokenIs(token.SEMICOLON) { // e.g let foo;
		return stmt
	}

	if !p.advanceNextTokenIs(token.ASSIGN) { // advance to =
		return nil
	}

	p.advanceToken() // go to RHS
	stmt.Expr = p.parseExpr(LOWEST)
	return stmt
}
func (p *Parser) parseReturnExpr() *ast.ReturnExpression {
	stmt := &ast.ReturnExpression{Token: p.currToken}
	if !p.advanceCurrTokenIs(token.RETURN) {
		return nil
	}

	stmt.Expr = p.parseExpr(LOWEST)
	p.advanceNextToEndToken()
	return stmt
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	stmt := &ast.IfStmt{Token: p.currToken}
	stmt.Alternatives = make([]*ast.ConditionalStmt, 0)
	if !p.advanceNextTokenIs(token.LPAREN) { // eat opening token IF
		p.errors = append(p.errors, fmt.Errorf("unexpected token %s, want %s", p.nextToken.Literal, token.LPAREN))
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
	}

	if p.currTokenIs(token.ELSE) {
		for p.advanceCurrTokenIs(token.ELSE) {
			// else if
			if p.currTokenIs(token.IF) {
				p.advanceToken()

				if !p.advanceCurrTokenIs(token.LPAREN) {
					unexpectedTokenError(token.LPAREN, p.currToken.Literal)
					return nil
				}

				condition := p.parseExpr(LOWEST) // parse expression inside if()

				if !p.advanceCurrTokenIs(token.RPAREN) {
					unexpectedTokenError(token.RPAREN, p.currToken.Literal)
					return nil
				}
				if !p.advanceCurrTokenIs(token.LBRACE) {
					unexpectedTokenError(token.LBRACE, p.currToken.Literal)
					return nil
				}

				elifStmt := p.parseStmt()

				condStmt := &ast.ConditionalStmt{Condition: condition, Statement: elifStmt}
				stmt.Alternatives = append(stmt.Alternatives, condStmt)

				if !p.advanceCurrTokenIs(token.RBRACE) {
					unexpectedTokenError(token.RBRACE, p.currToken.Literal)
					return nil
				}
			} else {
				if !p.advanceCurrTokenIs(token.LBRACE) {
					unexpectedTokenError(token.LBRACE, p.currToken.Literal)
					return nil
				}
				elseStmt := p.parseStmt()
				condStmt := &ast.ConditionalStmt{Statement: elseStmt}
				stmt.Alternatives = append(stmt.Alternatives, condStmt)

				if !p.advanceCurrTokenIs(token.RBRACE) {
					unexpectedTokenError(token.RBRACE, p.currToken.Literal)
					return nil
				}
				return stmt // return after an else statement
			}
		}

	}

	return stmt
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	blockStmt := &ast.BlockStmt{Statements: make([]ast.Statement, 0)}

	for !p.currTokenIs(token.EOF) && !p.currTokenIs(token.RBRACE) {
		if stmt := p.parseStmt(); stmt != nil {
			blockStmt.Statements = append(blockStmt.Statements, stmt)
		}
		for p.advanceCurrToEndToken() {
		}
	}
	if !p.advanceCurrTokenIs(token.RBRACE) {
		return nil
	}

	for p.advanceCurrToEndToken() {
	}
	return blockStmt
}

func (p *Parser) parseForStmt() ast.Statement {
	forLoopStmt := &ast.ForLoopStmt{Token: p.currToken}

	if !p.advanceCurrTokenIs(token.FOR) {
		return nil
	}
	if !p.currTokenIs(token.IDENT) {
		return nil
	}
	forLoopStmt.Variable = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	p.advanceToken()
	if !p.advanceCurrTokenIs(token.ASSIGN) {
		return nil
	}
	switch p.currToken.Type {
	case token.RANGE:
		p.advanceToken()
		if !p.currTokenIs(token.LBRACKET) && !p.currTokenIs(token.IDENT) {
			err := expectAfterTokenErrorStr("array", token.RANGE, p.nextToken.Literal)
			p.errors = append(p.errors, NewParseError(err, p.lexer.Position(), p.column))
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
			p.errors = append(p.errors, NewParseError(err, p.lexer.Position(), p.column))
			return nil
		}

		forLoopStmt.Statement = p.parseBlockStmt()

	default:
		return &ast.ErrorStmt{Value: ""}
	}
	return forLoopStmt
}

func (p *Parser) parseExpressionStmt() *ast.ExpressionStmt {
	if !slices.Contains(startTokens, token.LookupIdent(p.currToken.Literal)) {
		p.errors = append(p.errors, fmt.Errorf("expected start of expression, found '%s'", p.currToken.Literal))
		return nil
	}
	stmt := &ast.ExpressionStmt{Token: p.currToken}
	stmt.Expr = p.parseExpr(LOWEST)

	return stmt
}
