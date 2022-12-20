package parser

import "fmt"

type parseError struct {
	line   int
	column int
	err    error
}

func NewParseError(err error, line, column int) *parseError {
	return &parseError{
		err:    err,
		line:   line,
		column: column,
	}
}

func (p *parseError) Error() string {
	return fmt.Sprintf(`
	Error: %s
	Line: %d
	Column: %d
	`, p.err, p.line, p.column)
}
