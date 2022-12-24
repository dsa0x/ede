package parser

import (
	"ede/ast"
	"ede/lexer"
	"ede/token"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-multierror"
)

type (
	parseFn struct {
		prefix  func() ast.Expression
		infix   func(ast.Expression) ast.Expression
		postfix func(ast.Expression) ast.Expression
	}

	Parser struct {
		lexer *lexer.Lexer

		pos       token.Pos
		prevToken token.Token
		// line      int
		// column    int
		currToken token.Token
		nextToken token.Token

		tokens []token.Token

		parseFns map[token.TokenType]parseFn

		errors []error
	}
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.registerParseFns()

	p.advanceToken()
	p.advanceToken()
	return p
}

func (p *Parser) registerParseFns() {
	p.parseFns = make(map[token.TokenType]parseFn)
	p.parseFns[token.INT] = parseFn{prefix: p.parseInteger}
	p.parseFns[token.FLOAT] = parseFn{prefix: p.parseFloat}
	p.parseFns[token.TRUE] = parseFn{prefix: p.parseBool}
	p.parseFns[token.FALSE] = parseFn{prefix: p.parseBool}
	p.parseFns[token.IDENT] = parseFn{prefix: p.parseIdent}
	p.parseFns[token.STRING] = parseFn{prefix: p.parseStringLiteral}
	p.parseFns[token.BANG] = parseFn{prefix: p.parsePrefixExpression}
	p.parseFns[token.PLUS] = parseFn{prefix: p.parsePrefixExpression, infix: p.parseInfixOperator}
	p.parseFns[token.MINUS] = parseFn{prefix: p.parsePrefixExpression, infix: p.parseInfixOperator}
	p.parseFns[token.ASTERISK] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.SLASH] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.EQ] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.NEQ] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.GT] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.GTE] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.LT] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.LTE] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.OR_OR] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.AND_AND] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.MODULO] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.DEC] = parseFn{postfix: p.parsePostfixExpression}
	p.parseFns[token.INC] = parseFn{postfix: p.parsePostfixExpression}
	p.parseFns[token.LPAREN] = parseFn{prefix: p.parseGroupedExpression, infix: p.parseCallExpression}
	p.parseFns[token.LBRACKET] = parseFn{prefix: p.parseArrayLiteral, infix: p.parseIndexExpression}
	p.parseFns[token.LBRACE] = parseFn{prefix: p.parseHashLiteral}
	p.parseFns[token.FUNCTION] = parseFn{prefix: p.parseFunctionLiteral}
	p.parseFns[token.ASSIGN] = parseFn{infix: p.parseReassignment}
	p.parseFns[token.PLUS_EQUAL] = parseFn{infix: p.parsePlusEqual}
	p.parseFns[token.MINUS_EQUAL] = parseFn{infix: p.parseMinusEqual}
	p.parseFns[token.RANGE_ARRAY] = parseFn{infix: p.parseRangeArray}
	p.parseFns[token.DOT] = parseFn{infix: p.parseObjectMethodExpression}
	p.parseFns[token.RETURN] = parseFn{prefix: p.parseReturnExpr}
	p.parseFns[token.MATCH] = parseFn{prefix: p.parseMatchExpression}
	p.registerIllegalFns()
}

func (p *Parser) registerIllegalFns() {
	ilFn := func() ast.Expression {
		p.addError("illegal token '%s'", p.currToken.Literal)
		return nil
	}
	ilFn2 := func(ast.Expression) ast.Expression { return ilFn() }
	p.parseFns[token.ILLEGAL] = parseFn{prefix: ilFn, infix: ilFn2, postfix: ilFn2}
}

func (p *Parser) Parse() *ast.Program {
	prog := &ast.Program{
		Statements: make([]ast.Statement, 0),
		ValuePos:   p.pos,
	}
	for !p.currTokenIs(token.EOF) {
		if p.advanceCurrTokenIs(token.NEWLINE) { // ignore new line
			continue
		}
		stmt := p.parseStmt()
		if stmt != nil && !reflect.ValueOf(stmt).IsNil() {
			prog.Statements = append(prog.Statements, stmt)
		}

		if p.Errors() != nil {
			prog.ParseErrors = multierror.Append(prog.ParseErrors, p.Errors())
			return prog
		}

		p.eatEndToken() // advance all end tokens
	}
	return prog
}

func (p *Parser) parseExpr(precedence int) ast.Expression {
	prefixFn := p.prefixParseFn(p.currToken.Type)
	if prefixFn == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	left := prefixFn()

	for !p.currTokenIs(token.SEMICOLON) && precedence < p.currPrecedence() {
		infixFn := p.infixParseFn(p.currToken.Type)
		if infixFn != nil {
			left = infixFn(left)
			continue
		}
		postfixFn := p.postfixParseFn(p.currToken.Type)
		if postfixFn != nil {
			left = postfixFn(left)
			continue
		}
		// if it get here, no infix or posfix function was found
		// we need to return to prevent infinite loop
		p.addError("no infix or postfix parse function for token '%s'", p.currToken.Literal)
		return nil
	}

	return left
}

func (p *Parser) prefixParseFn(tok token.TokenType) func() ast.Expression {
	if parseFn, ok := p.parseFns[tok]; ok && parseFn.prefix != nil {
		return parseFn.prefix
	}
	return nil
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.addError(fmt.Sprintf("no prefix parse function for '%s' found", t))
}

func (p *Parser) infixParseFn(tok token.TokenType) func(ast.Expression) ast.Expression {
	if parseFn, ok := p.parseFns[tok]; ok && parseFn.infix != nil {
		return parseFn.infix
	}
	return nil
}
func (p *Parser) postfixParseFn(tok token.TokenType) func(ast.Expression) ast.Expression {
	if parseFn, ok := p.parseFns[tok]; ok && parseFn.postfix != nil {
		return parseFn.postfix
	}
	return nil
}
