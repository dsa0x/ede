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
	p.advanceNextTokenIs(token.SEMICOLON)
	return stmt
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	stmt := &ast.IfStmt{Token: p.currToken}
	stmt.Alternatives = make([]*ast.ConditionalStmt, 0)
	if !p.nextTokenIs(token.LPAREN) {
		p.errors = append(p.errors, fmt.Errorf("unexpected token %s, want %s", p.nextToken.Literal, token.LPAREN))
		return nil
	}
	p.advanceToken()
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

				if !p.currTokenIs(token.LPAREN) {
					// TODO: error handling
					return nil
				}
				p.advanceToken()

				condition := p.parseExpr(LOWEST) // parse expression inside if()

				if !p.advanceNextTokenIs(token.RPAREN) {
					// TODO: error handling
					return nil
				}
				if !p.advanceNextTokenIs(token.LBRACE) {
					// TODO: error handling
					return nil
				}
				p.advanceToken() // move past Lbrace

				elifStmt := p.parseStmt()

				condStmt := &ast.ConditionalStmt{Condition: condition, Statement: elifStmt}
				stmt.Alternatives = append(stmt.Alternatives, condStmt)

				if !p.advanceNextTokenIs(token.RBRACE) {
					// TODO: error handling
					return nil
				}
				p.advanceToken() // move past Rbrace
			} else {
				if !p.advanceCurrTokenIs(token.LBRACE) {
					// TODO: error handling
					return nil
				}
				elseStmt := p.parseStmt()
				condStmt := &ast.ConditionalStmt{Statement: elseStmt}
				stmt.Alternatives = append(stmt.Alternatives, condStmt)

				if !p.advanceNextTokenIs(token.RBRACE) {
					// TODO: error handling
					return nil
				}
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
		p.advanceToken()
		p.advanceCurrTokenIs(token.SEMICOLON)
	}
	if !p.advanceCurrTokenIs(token.RBRACE) {
		return nil
	}

	p.advanceNextTokenIs(token.SEMICOLON)
	return blockStmt
}

func (p *Parser) parseForStmt() ast.Statement {
	forLoopStmt := &ast.ForLoopStmt{Token: p.currToken}

	if !p.advanceNextTokenIs(token.IDENT) {
		return nil
	}
	forLoopStmt.Variable = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	if !p.advanceNextTokenIs(token.ASSIGN) {
		return nil
	}
	switch p.nextToken.Type {
	case token.RANGE:
		p.advanceToken()
		if !p.nextTokenIs(token.LBRACKET) && !p.nextTokenIs(token.IDENT) {
			return expectAfterTokenError("array", token.RANGE, p.nextToken.Literal)
		}
		p.advanceToken()
		startingToken := p.currToken

		expr := p.parseExpr(LOWEST)

		if startingToken.Type == token.LBRACKET { // this means the boundary is a literal
			if !p.currTokenIs(token.RBRACKET) {
				return expectAfterTokenError(token.RBRACKET, "for loop boundary", p.currToken.Literal)
			}
		}

		forLoopStmt.Boundary = expr

		if !p.advanceNextTokenIs(token.LBRACE) {
			return expectAfterTokenError(token.LBRACE, "right bracket ']'", p.currToken.Literal)
		}
		p.advanceToken() // past LBRACE

		forLoopStmt.Statement = p.parseBlockStmt()

		p.advanceNextTokenIs(token.SEMICOLON)

	default:
		return &ast.ErrorStmt{Value: ""}
	}
	return forLoopStmt
}

func (p *Parser) parseExpressionStmt() *ast.ExpressionStmt {
	if !slices.Contains(startTokens, token.LookupIdent(p.currToken.Literal)) {
		p.errors = append(p.errors, fmt.Errorf("expected expression, found '%s'", p.currToken.Literal))
		return nil
	}
	stmt := &ast.ExpressionStmt{Token: p.currToken}
	stmt.Expr = p.parseExpr(LOWEST)

	p.advanceNextTokenIs(token.SEMICOLON)
	return stmt
}
