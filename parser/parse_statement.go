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
		return p.parseReturnExpr()
	case token.FOR:
		return p.parseForStmt()
	case token.NEWLINE, token.SEMICOLON:
		p.advanceToken()
		return p.parseStmt()
	case token.SINGLE_COMMENT:
		return &ast.CommentStmt{Token: p.currToken, Value: p.currToken.Literal, ValuePos: p.pos}
	}
	return p.parseExpressionStmt()
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	stmt := &ast.LetStmt{Token: p.currToken}
	if !p.advanceNextTokenIs(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal, ValuePos: p.pos}

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
	stmt := &ast.ReturnExpression{Token: p.currToken, ValuePos: p.pos}
	if !p.advanceCurrTokenIs(token.RETURN) {
		return nil
	}

	stmt.Expr = p.parseExpr(LOWEST)
	p.advanceNextToEndToken()
	return stmt
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	stmt := &ast.IfStmt{Token: p.currToken, ValuePos: p.pos}
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
		p.eatEndToken()
	}
	if !p.advanceCurrTokenIs(token.RBRACE) {
		p.addError(unexpectedTokenError(token.RBRACE, p.currToken.Literal))
		return nil
	}

	p.eatEndToken()
	return blockStmt
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
		return &ast.ErrorStmt{Value: ""}
	}
	return forLoopStmt
}

func (p *Parser) parseExpressionStmt() *ast.ExpressionStmt {
	if !slices.Contains(startTokens, token.LookupIdent(p.currToken.Literal)) {
		p.errors = append(p.errors, fmt.Errorf("expected start of expression, found '%s'", p.currToken.Literal))
		return nil
	}
	stmt := &ast.ExpressionStmt{Token: p.currToken, ValuePos: p.pos}
	stmt.Expr = p.parseExpr(LOWEST)

	return stmt
}
