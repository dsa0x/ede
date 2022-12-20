package parser

import (
	"ede/lexer"
	"ede/token"
	"fmt"
	"testing"
)

func TestParseGroupedExprs(t *testing.T) {
	tests := []struct {
		input              string
		expectedFirstToken string
		expectedLastToken  string
	}{
		{
			input: "(5 + 4)",
		},
		{
			input: "(5 * (5 + 4))",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			checkParserErrors(t, p)
			if !p.currTokenIs(token.LPAREN) {
				t.Fatalf("expected (, got %s", p.currToken.Literal)
			}
			expr := p.parseExpr(LOWEST)
			if p.prevToken.Type != token.RPAREN {
				t.Fatalf("expected (, got %s", p.currToken.Literal)
			}
			fmt.Println(expr)
		})
	}
}

func TestParseMethodExpression(t *testing.T) {
	tests := []struct {
		input              string
		expectedFirstToken string
		expectedLastToken  string
	}{
		{
			input: "(5 + 4)",
		},
		{
			input: "(5 * (5 + 4))",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			checkParserErrors(t, p)
			if !p.currTokenIs(token.LPAREN) {
				t.Fatalf("expected (, got %s", p.currToken.Literal)
			}
			expr := p.parseExpr(LOWEST)
			if p.prevToken.Type != token.RPAREN {
				t.Fatalf("expected (, got %s", p.currToken.Literal)
			}
			fmt.Println(expr)
		})
	}
}
