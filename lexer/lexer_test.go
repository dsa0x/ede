package lexer

import (
	"ede/token"
	"testing"
)

func TestNextTokenLet(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let name = "sam";
	`

	tests := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "name"},
		{token.ASSIGN, "="},
		{token.STRING, "sam"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expType, tok.Type)
		}
		if tok.Literal != tt.expLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expLiteral, tok.Literal)
		}
	}
}
func TestNextTokenStruct(t *testing.T) {
	input := `
	struct User {name, age};
	let profile = User{"joe", 10};
	`

	tests := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{token.STRUCT, "struct"},
		{token.IDENT, "User"},
		{token.LBRACE, "{"},
		{token.IDENT, "name"},
		{token.COMMA, ","},
		{token.IDENT, "age"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "profile"},
		{token.ASSIGN, "="},
		{token.IDENT, "User"},
		{token.LBRACE, "{"},
		{token.STRING, "joe"},
		{token.COMMA, ","},
		{token.INT, "10"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expType, tok.Type)
		}
		if tok.Literal != tt.expLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expLiteral, tok.Literal)
		}
	}
}

func TestNextTokenConditionals(t *testing.T) {
	input := `
	if name > 1 {
		return Error();
	} else {
		return 10;
	};
	`

	tests := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{token.IF, "if"},
		{token.IDENT, "name"},
		{token.GT, ">"},
		{token.INT, "1"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "Error"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expType, tok.Type)
		}
		if tok.Literal != tt.expLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expLiteral, tok.Literal)
		}
	}
}

func TestNextTokenFunctions(t *testing.T) {
	input := `
	let age = func() { return 10 };
	// let age = func() { 
	// 	if name > 1 {
	// 		return Error()
	// 	}
	// 	return 10
	//  }
	`

	tests := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "age"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.INT, "10"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expType, tok.Type)
		}
		if tok.Literal != tt.expLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expLiteral, tok.Literal)
		}
	}
}

func TestNextTokenOperators(t *testing.T) {
	input := `
	let sum = 10 + 10;
	let mult = 4 * 100;
	let div = 2 / 1;
	let gt = 100 > 4;
	`
	tests := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "sum"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.PLUS, "+"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "mult"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
		{token.ASTERISK, "*"},
		{token.INT, "100"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "div"},
		{token.ASSIGN, "="},
		{token.INT, "2"},
		{token.SLASH, "/"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "gt"},
		{token.ASSIGN, "="},
		{token.INT, "100"},
		{token.GT, ">"},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expType, tok.Type)
		}
		if tok.Literal != tt.expLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expLiteral, tok.Literal)
		}
	}
}

// func TestNextTokens(t *testing.T) {
// 	input := `let five = 5;
// 	let ten = 10;
// 	let add = func(x, y) {
// 		x + y;
// 	};

// 	let result = add(five, ten);
// 	!-/*5;
// 	5 < 10 > 5;

// 	if (5 < 10) {
// 		return true;
// 	} else {
// 		return false;
//  	}

// 	10 == 10;
// 	10 != 9;
// 	"foobar"
// 	"foo bar"
// 	[1, 2];
// 	`

// 	tests := []struct {
// 		expType    token.TokenType
// 		expLiteral string
// 	}{
// 		{token.LET, "let"},
// 		{token.IDENT, "five"},
// 		{token.ASSIGN, "="},
// 		{token.INT, "5"},
// 		{token.SEMICOLON, ";"},
// 		{token.LET, "let"},
// 		{token.IDENT, "ten"},
// 		{token.ASSIGN, "="},
// 		{token.INT, "10"},
// 		{token.SEMICOLON, ";"},
// 		{token.LET, "let"},
// 		{token.IDENT, "add"},
// 		{token.ASSIGN, "="},
// 		{token.FUNCTION, "func"},
// 		{token.LPAREN, "("},
// 		{token.IDENT, "x"},
// 		{token.COMMA, ","},
// 		{token.IDENT, "y"},
// 		{token.RPAREN, ")"},
// 		{token.LBRACE, "{"},
// 		{token.IDENT, "x"},
// 		{token.PLUS, "+"},
// 		{token.IDENT, "y"},
// 		{token.SEMICOLON, ";"},
// 		{token.RBRACE, "}"},
// 		{token.SEMICOLON, ";"},
// 		{token.LET, "let"},
// 		{token.IDENT, "result"},
// 		{token.ASSIGN, "="},
// 		{token.IDENT, "add"},
// 		{token.LPAREN, "("},
// 		{token.IDENT, "five"},
// 		{token.COMMA, ","},
// 		{token.IDENT, "ten"},
// 		{token.RPAREN, ")"},
// 		{token.SEMICOLON, ";"},
// 		{token.BANG, "!"},
// 		{token.MINUS, "-"},
// 		{token.SLASH, "/"},
// 		{token.ASTERISK, "*"},
// 		{token.INT, "5"},
// 		{token.SEMICOLON, ";"},
// 		{token.INT, "5"},
// 		{token.LT, "<"},
// 		{token.INT, "10"},
// 		{token.GT, ">"},
// 		{token.INT, "5"},
// 		{token.SEMICOLON, ";"},
// 		{token.IF, "if"},
// 		{token.LPAREN, "("},
// 		{token.INT, "5"},
// 		{token.LT, "<"},
// 		{token.INT, "10"},
// 		{token.RPAREN, ")"},
// 		{token.LBRACE, "{"},
// 		{token.RETURN, "return"},
// 		{token.TRUE, "true"},
// 		{token.SEMICOLON, ";"},
// 		{token.RBRACE, "}"},
// 		{token.ELSE, "else"},
// 		{token.LBRACE, "{"},
// 		{token.RETURN, "return"},
// 		{token.FALSE, "false"},
// 		{token.SEMICOLON, ";"},
// 		{token.RBRACE, "}"},
// 		{token.INT, "10"},
// 		{token.EQ, "=="},
// 		{token.INT, "10"},
// 		{token.SEMICOLON, ";"},
// 		{token.INT, "10"},
// 		{token.NOT_EQ, "!="},
// 		{token.INT, "9"},
// 		{token.SEMICOLON, ";"},
// 		{token.STRING, "foobar"},
// 		{token.STRING, "foo bar"},
// 		{token.LBRACKET, "["},
// 		{token.INT, "1"},
// 		{token.COMMA, ","},
// 		{token.INT, "2"},
// 		{token.RBRACKET, "]"},
// 		{token.SEMICOLON, ";"},
// 		{token.EOF, ""},
// 	}

// 	l := New(input)

// 	for i, tt := range tests {
// 		tok := l.NextToken()
// 		if tok.Type != tt.expType {
// 			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
// 				i, tt.expType, tok.Type)
// 		}
// 		if tok.Literal != tt.expLiteral {
// 			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
// 				i, tt.expLiteral, tok.Literal)
// 		}
// 	}
// }
