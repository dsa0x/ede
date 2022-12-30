package parser

import (
	"ede/token"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type parseError struct {
	line   int
	column int
	err    error
}

func NewParseError(err error, pos token.Pos) *parseError {
	return &parseError{
		err:    err,
		line:   pos.Line,
		column: pos.Column,
	}
}

func (p *parseError) Error() string {
	return fmt.Sprintf(`
	Error: %s
	Line: %d
	Column: %d
	`, p.err, p.line, p.column)
}

func (p *Parser) addError(msg string, format ...interface{}) {
	p.errors = append(p.errors, NewParseError(fmt.Errorf(msg, format...), p.currPos()))
}

func (p *Parser) appendError(err error) {
	p.errors = append(p.errors, err)
}

func unexpectedTokenError(exp, got string) string {
	return fmt.Sprintf("expected token %s, got %s", exp, got)
}

func expectAfterTokenErrorStr(exp, prev, got string) string {
	return fmt.Sprintf("expected %s after %s, got %s", exp, prev, got)
}

func (p *Parser) Errors() error {
	return multierror.Append(nil, p.errors...).ErrorOrNil()
}
