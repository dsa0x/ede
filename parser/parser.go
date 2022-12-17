package parser

import (
	"ede/ast"
	"ede/lexer"
	"ede/token"
	"fmt"
	"reflect"
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
		currToken token.Token
		nextToken token.Token

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
	p.parseFns[token.LT] = parseFn{infix: p.parseInfixOperator}
	p.parseFns[token.DEC] = parseFn{postfix: p.parsePostfixExpression}
	p.parseFns[token.INC] = parseFn{postfix: p.parsePostfixExpression}
	p.parseFns[token.LPAREN] = parseFn{prefix: p.parseGroupedExpression, infix: p.parseCallExpression}
	p.parseFns[token.LBRACKET] = parseFn{prefix: p.parseArrayLiteral, infix: p.parseIndexExpression}
	p.parseFns[token.FUNCTION] = parseFn{prefix: p.parseFunctionLiteral}
}

func (p *Parser) Parse() *ast.Program {
	prog := &ast.Program{
		Statements:  make([]ast.Statement, 0),
		ParseErrors: make([]error, 0),
	}
	for !p.currTokenIs(token.EOF) {
		stmt := p.parseStmt()
		if !reflect.ValueOf(stmt).IsNil() {
			prog.Statements = append(prog.Statements, stmt)
		}

		if len(p.Errors()) > 0 {
			prog.ParseErrors = append(prog.ParseErrors, p.Errors()...)
			return prog
		}

		p.advanceToken()
		p.advanceCurrTokenIs(token.SEMICOLON)
	}
	return prog
}

func (p *Parser) advanceToken() {
	p.currToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) currTokenIs(tok token.TokenType) bool {
	return p.currToken.Type == tok
}
func (p *Parser) nextTokenIs(tok token.TokenType) bool {
	return p.nextToken.Type == tok
}

// advanceCurrTokenIs advances to the next token if the current token  matches, else it does nothing
func (p *Parser) advanceCurrTokenIs(tok token.TokenType) bool {
	found := p.currTokenIs(tok)
	if found {
		p.advanceToken()
	}
	return found
}

// advanceNextTokenIs advances to the next token if it matches, else it does nothing
func (p *Parser) advanceNextTokenIs(tok token.TokenType) bool {
	found := p.nextTokenIs(tok)
	if found {
		p.advanceToken()
	}
	return found
}

func (p *Parser) parseExpr(precedence int) ast.Expression {
	prefixFn := p.prefixParseFn(p.currToken.Type)
	if prefixFn == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	left := prefixFn()

	for !p.nextTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		if infixFn := p.infixParseFn(p.nextToken.Type); infixFn != nil {
			p.advanceToken()
			left = infixFn(left)
			continue
		}
		if postfixFn := p.postfixParseFn(p.nextToken.Type); postfixFn != nil {
			p.advanceToken()
			left = postfixFn(left)
			continue
		}
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
	p.errors = append(p.errors, fmt.Errorf("no prefix parse function for %s found", t))
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

func (p *Parser) Errors() []error {
	return p.errors
}
