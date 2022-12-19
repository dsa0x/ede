package evaluator

import (
	"ede/lexer"
	"ede/object"
	"ede/parser"
	"fmt"
	"strconv"
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

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%v, want=%v",
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

func testObject(t *testing.T, obj object.Object, evaluated any) bool {
	switch evaluated := evaluated.(type) {
	case int, int32, int64, uint:
		ev, _ := strconv.ParseInt(fmt.Sprint(evaluated), 10, 64)
		return testIntegerObject(t, obj, ev)
	case float32, float64:
		ev, _ := strconv.ParseFloat(fmt.Sprint(evaluated), 64)
		return testFloatObject(t, obj, ev)
	case nil:
		return obj.Type() == object.NULL_OBJ
	case error:
		return obj.Type() == object.ERROR_OBJ && obj.Inspect() == evaluated.Error()
	}
	return false
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

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = func(x) { x; }; identity(5);", 5},
		{"let identity = func(x) { return x; }; identity(5);", 5},
		{"let double = func(x) { x * 2; }; double(5);", 10},
		{"let add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { x + x; }(5)", 10},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			testIntegerObject(t, testEval(tt.input), tt.expected)
		})
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
	let first = 10;
	let second = 10;
	let third = 10;

	let ourFunction = func(first) {
	let second = 20;

	first + second + third;
	};

	ourFunction(20) + first + second;`

	if !testIntegerObject(t, testEval(input), 70) {
		t.FailNow()
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = func(x) {
	func(y) { x + y };
	};

	let addTwo = newAdder(2);
	addTwo(2);`

	if !testObject(t, testEval(input), 4) {
		t.FailNow()
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
		{`println("hello", "world!")`, nil},
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
		case []int:
			array, ok := evaluated.(*object.Array[any])
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(*array.Entries) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(*array.Entries))
				continue
			}

			for i, expectedElem := range expected {
				testIntegerObject(t, (*array.Entries)[i], int64(expectedElem))
			}
		}
	}
}

func TestEvalStatements(t *testing.T) {

	tests := []struct {
		input  string
		result any
	}{
		{
			input: `
			let a = 10;
		let add = func(x) {
			println("a", a);
			<- x + a;
		};
		add(add(10));
		`,
			result: 30,
		},
		{
			input: `let a = 10;
			let add = func(x) {
				<- x + a;
			};
			a = add(add(10));
			a + a;
		`,
			result: 60,
		},
		{
			input: `let a = 10;
			let add = func(x) {
				<- x + a;
			};
			a = add(add(10));
			add(a + a) + add(a + a);
		`,
			result: 180,
		},
		{
			input: `let a = 10.5;
			let add = func(x) {
				<- x + a;
			};
			a = add(add(10));
			a;
		`,
			result: 31.0,
		},
		{
			input: `
			let arr = [1..10];
			arr[2];
		`,
			result: 3,
		},
		{
			input: `
			let name = "foo";
			let age = 10.5;
			for i = range [1..10] {
				age++;
			};
			age;
		`,
			result: 20.5,
		},
		{
			input: `
			let name = "foo";
			let age = 30;
			for i = range [1..3] {
				age = age + i;
			};
			age;
		`,
			result: 36,
		},
		{
			input: `let sub;
			sub;
		`,
			result: nil,
		},
		{
			input: `
			let subjects = ["english", "french"];
			for sub = range subjects {
				println(index, sub);
			};
		`,
			result: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			if !testObject(t, evaluated, tt.result) {
				t.FailNow()
			}
		})
	}
}
func TestEvalStatements_Error(t *testing.T) {

	tests := []struct {
		input  string
		result string
	}{
		{
			input: `let index = foo;
			let subjects = ["english", "french"];
		`,
			result: "cannot assign to reserved keyword 'index'",
		},
		{
			input: `let arr = [2, ( + 5];
		`,
			result: "expected closing parenthesis token ')', got '5'",
		},
		{
			input: `let arr = [2, 3 +];
		`,
			result: "invalid right expression ] for operator '+'",
		},
		{
			input: `let arr = [2, | +];
		`,
			result: "illegal token |",
		},
		{
			input: `let arr = [2;
		`,
			result: "expected closing bracket token ']', got '2'",
		},
		{
			input: `
			a = 24;
			println(a);
		`,
			result: "cannot reassign undeclared identifier 'a'",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			ev := testEval(tt.input)
			evaluated, ok := ev.(*object.Error)
			if !ok {
				t.Fatalf("expected result of type *object.Error, got %T", evaluated)
			}
			if !strings.Contains(evaluated.Message, tt.result) {
				t.Fatalf("expected \"%s\" to contain error \"%s\"", evaluated.Message, tt.result)
			}
		})
	}
}

func TestEval(t *testing.T) {
	input := `
	let arr = [1..10];
	let double = func(x) {
		<- x + x;
	};
	arr.push(5);
	// let arrx = arr.map(double);
	// arrx[1]
	println(arr);
	`

	input = `
	let arr = [1,2];
	arr.push(5);
	println(arr);
	`

	evaluated := testEval(input)
	if !testObject(t, evaluated, 2) {
		t.FailNow()
	}
}
