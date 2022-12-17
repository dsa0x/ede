package evaluator

import (
	"ede/lexer"
	"ede/object"
	"ede/parser"
	"strings"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"10--", 9},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testIntegerObject(t, evaluated, tt.expected) {
			t.FailNow()
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Parse()
	env := object.NewEnvironment(nil)
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Int)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			t.FailNow()
		}
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{`!"bro"`, false},
		{`!!"bro"`, true},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			t.FailNow()
		}
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else if (true) { 15 } else { 20 }", 15},
		{"if (1 > 2) { 10 } else if (2) { 15 } else { 20 }", 15},
		{"if (1 > 2) { 10 } else if (false) { 15 } else { 20 }", 20},
		{"if (1 > 2) { 10 } else if (5 < 2) { 15 } else { 20 }", 20},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			if !testIntegerObject(t, evaluated, int64(integer)) {
				t.FailNow()
			}
		} else {
			if !testNullObject(t, evaluated) {
				t.FailNow()
			}
		}
	}

	t.Run("invalid order", func(t *testing.T) {
		evaluated := testEval("if (1 > 2) { 10 } else { 20 } else if (true) { 5 }")
		result, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
		}
		if !strings.Contains(result.Message, "expected expression") {
			t.Fatalf("Error message %s does not contain 'expected expression'", result.Message)
		}
	})
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"<- 10;", 10},
		{"<- 10; 9;", 10},
		{"<- 2 * 5; 9;", 10},
		{"9; <- 2 * 5; 9;", 10},
		{"if (10 > 1) { <- 10; }", 10},
		{
			`
		if (10 > 1) {
		  if (10 > 2) {
		    <- 10;
		  }

		  <- 1;
		}
		`,
			10,
		},
		{
			`
		let f = func(x) {
		  <- x;
		  x + 10;
		};
		f(10);`,
			10,
		},
		{
			`
		let f = func(x) {
		   let result = x + 10;
		   <- result;
		   <- 10;
		};
		f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testIntegerObject(t, evaluated, tt.expected) {
			t.FailNow()
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INT"},
		{`len("one", "two")`, "builtin function 'len' requires exactly one argument, got 2"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		// {`puts("hello", "world!")`, nil},
		// {`first([1, 2, 3])`, 1},
		// {`first([])`, nil},
		// {`first(1)`, "argument to `first` must be ARRAY, got INTEGER"},
		// {`last([1, 2, 3])`, 3},
		// {`last([])`, nil},
		// {`last(1)`, "argument to `last` must be ARRAY, got INTEGER"},
		// {`rest([1, 2, 3])`, []int{2, 3}},
		// {`rest([])`, nil},
		// {`push([], 1)`, []int{1}},
		// {`push(1, 1)`, "argument to `push` must be ARRAY, got INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
			// case []int:
			// 	array, ok := evaluated.(*object.Array)
			// 	if !ok {
			// 		t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
			// 		continue
			// 	}

			// 	if len(array.Elements) != len(expected) {
			// 		t.Errorf("wrong num of elements. want=%d, got=%d",
			// 			len(expected), len(array.Elements))
			// 		continue
			// 	}

			// 	for i, expectedElem := range expected {
			// 		testIntegerObject(t, array.Elements[i], int64(expectedElem))
			// 	}
		}
	}
}
