package parser

import (
	"ede/ast"
	"ede/lexer"
	"fmt"
	"testing"

	"github.com/hashicorp/go-multierror"
)

func TestLetStatements(t *testing.T) {
	input := `
   let x = 5;
   let y = 10;
   let foobar = 838383;
   `

	l := lexer.New(input)
	p := New(l)

	program := p.Parse()
	if program == nil {
		t.Fatalf("Parse() returned nil")
	}
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.Literal() != "let" {
		t.Errorf("s.Literal not 'let'. got=%q", s.Literal())
		return false
	}
	letStmt, ok := s.(*ast.LetStmt)
	if !ok {
		t.Errorf("s not *ast.LetStmt. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.Literal() != name {
		t.Errorf("letStmt.Name.Literal() not '%s'. got=%s",
			name, letStmt.Name.Literal())
		return false
	}
	return true
}

func TestReturnExpression(t *testing.T) {
	input := `
   <- 5;
   <- 10;
   <- 993322;
   `
	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnExpression)
		if !ok {
			t.Errorf("stmt not *ast.ReturnExpression. got=%T", stmt)
			continue
		}

		exp := fmt.Sprintf("%s%s", "<-", returnStmt.Expr.Literal())
		if returnStmt.Literal() != exp {
			t.Errorf("returnStmt.Literal not '%s', got '%q'", exp, returnStmt.Literal())
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		prog := p.Parse()
		checkParserErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("p.Statements does not contain %d statements. got=%d\n",
				1, len(prog.Statements))
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("p.Statements[0] is not ast.ExpressionStmt. got=%T",
				prog.Statements[0])
		}

		exp, ok := stmt.Expr.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expr)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingPostfixExpressions(t *testing.T) {
	postfixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"5++;", "++", 5},
		{"15--;", "--", 15},
	}

	for _, tt := range postfixTests {
		l := lexer.New(tt.input)
		p := New(l)
		prog := p.Parse()
		checkParserErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("p.Statements does not contain %d statements. got=%d\n",
				1, len(prog.Statements))
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("p.Statements[0] is not ast.ExpressionStmt. got=%T",
				prog.Statements[0])
		}

		exp, ok := stmt.Expr.(*ast.PostfixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PostfixExpression. got=%T", stmt.Expr)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Left, tt.value) {
			return
		}
	}
}

func testLiteralExpression(t *testing.T,
	exp ast.Expression,
	expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.Literal() != value {
		t.Errorf("ident.Literal not %s. got=%s", value,
			ident.Literal())
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.Literal() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.Literal not %d. got=%s", value,
			integ.Literal())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.Literal() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.Literal not %t. got=%s",
			value, bo.Literal())
		return false
	}
	return true
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		expr := p.parseExpr(LOWEST)
		if !testInfixExpression(t, expr, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if errors == nil {
		return
	}
	t.Errorf("parser has %d errors", multierror.Append(errors).Len())
	t.Errorf("parser error: %q", errors)
	t.FailNow()
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStmt. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expr.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.FunctionLiteral. got=%T",
			stmt.Expr)
	}

	if len(function.Params) != 2 {
		t.Fatalf("function literal Params wrong. want 2, got=%d\n",
			len(function.Params))
	}

	testLiteralExpression(t, function.Params[0], "x")
	testLiteralExpression(t, function.Params[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStmt. got=%T",
			function.Body.Statements[0])
	}

	if !testInfixExpression(t, bodyStmt.Expr, "x", "+", "y") {
		t.FailNow()
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "func() {};", expectedParams: []string{}},
		{input: "func(x) {};", expectedParams: []string{"x"}},
		{input: "func(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStmt)
		function := stmt.Expr.(*ast.FunctionLiteral)

		if len(function.Params) != len(tt.expectedParams) {
			t.Errorf("length Params wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Params))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Params[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStmt. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expr.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.CallExpression. got=%T",
			stmt.Expr)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Args) != 3 {
		t.Fatalf("wrong length of Args. got=%d", len(exp.Args))
	}

	testLiteralExpression(t, exp.Args[0], 1)
	testInfixExpression(t, exp.Args[1], 2, "*", 3)
	testInfixExpression(t, exp.Args[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStmt)
		exp, ok := stmt.Expr.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expr is not ast.CallExpression. got=%T",
				stmt.Expr)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Args) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of Args. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Args))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Args[i].Literal() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Args[i].Literal())
			}
		}
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStmt)
	array, ok := stmt.Expr.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expr)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStmt)
	indexExp, ok := stmt.Expr.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expr)
	}
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	t.Skipf("skipping until deciding how to handle empty hash/set")
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	hash, ok := stmt.Expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expr)
	}

	if len(hash.Pair) != 0 {
		t.Errorf("hash.Pair has wrong length. got=%d", len(hash.Pair))
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 2, "one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	hash, ok := stmt.Expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expr)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	if len(hash.Pair) != len(expected) {
		t.Errorf("hash.Pair has wrong length. got=%d", len(hash.Pair))
	}

	for key, value := range hash.Pair {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[literal.Literal()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	hash, ok := stmt.Expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expr)
	}

	expected := map[string]int64{
		"true":  1,
		"false": 2,
	}

	if len(hash.Pair) != len(expected) {
		t.Errorf("hash.Pair has wrong length. got=%d", len(hash.Pair))
	}

	for key, value := range hash.Pair {
		boolean, ok := key.(*ast.BooleanLiteral)
		if !ok {
			t.Errorf("key is not ast.BooleanLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[boolean.Literal()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 4, 2: 2, 3: 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	hash, ok := stmt.Expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expr)
	}

	expected := map[string]int64{
		"1": 4,
		"2": 2,
		"3": 3,
	}

	if len(hash.Pair) != len(expected) {
		t.Errorf("hash.Pair has wrong length. got=%d", len(hash.Pair))
	}

	for key, value := range hash.Pair {
		integer, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not ast.IntegerLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[integer.Literal()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	hash, ok := stmt.Expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expr)
	}

	if len(hash.Pair) != 3 {
		t.Errorf("hash.Pair has wrong length. got=%d", len(hash.Pair))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pair {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.Literal()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.Literal())
			continue
		}

		testFunc(value)
	}
}
