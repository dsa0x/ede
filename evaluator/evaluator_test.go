package evaluator

import (
	"ede/ast"
	"ede/lexer"
	"ede/object"
	"ede/parser"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
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
		{"6 % 2", 0},
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
	ev := &Evaluator{}
	resp := ev.Eval(program, env)
	return resp
}
func testEval2(input string) (*Evaluator, object.Object) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Parse()
	env := object.NewEnvironment(nil)
	ev := &Evaluator{}
	return ev, ev.Eval(program, env)
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

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
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
		})
	}

	t.Run("invalid order", func(t *testing.T) {
		evaluated := testEval("if (1 > 2) { 10 } else { 20 } else if (true) { 5 }")
		result, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
		}
		if !strings.Contains(result.Message, "expected start of expression") {
			t.Fatalf("Error message %s does not contain 'expected start of expression'", result.Message)
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
	if obj, ok := obj.(*object.Hash); ok {
		evaluated := evaluated.([]string)
		exp := []string{}
		for entry := range obj.Entries {
			exp = append(exp, entry)
		}
		return slices.Equal(exp, evaluated)
	}
	switch evaluated := evaluated.(type) {
	case int, int32, int64, uint:
		ev, _ := strconv.ParseInt(fmt.Sprint(evaluated), 10, 64)
		return testIntegerObject(t, obj, ev)
	case float32, float64:
		ev, _ := strconv.ParseFloat(fmt.Sprint(evaluated), 64)
		return testFloatObject(t, obj, ev)
	case nil:
		return obj == nil || obj.Type() == object.NIL_OBJ
	case error:
		return obj.Type() == object.ERROR_OBJ && obj.Inspect() == evaluated.Error()
	case []string:
		obj, ok := obj.(*object.Array)
		if !ok {
			return false
		}
		if len(*obj.Entries) != len(evaluated) {
			return false
		}
		entries := []string{}
		for _, el := range *obj.Entries {
			entries = append(entries, el.Inspect())
		}
		sort.Strings(entries)
		sort.Strings(evaluated)
		return slices.Equal(entries, evaluated)
	case []any:
		obj, ok := obj.(*object.Array)
		if !ok {
			return false
		}
		if len(*obj.Entries) != len(evaluated) {
			return false
		}
		entries, evEntries := []string{}, []string{}
		for _, el := range *obj.Entries {
			entries = append(entries, el.Inspect())
		}
		for _, el := range evaluated {
			evEntries = append(evEntries, fmt.Sprint(el))
		}
		sort.Strings(entries)
		sort.Strings(evEntries)
		return slices.Equal(entries, evEntries)
	case bool:
		obj, ok := obj.(*object.Boolean)
		if !ok {
			return false
		}
		return obj.Value
	case string:
		obj, ok := obj.(*object.String)
		if !ok {
			return false
		}
		return obj.Value == evaluated
	}
	return false
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; }", 10},
		{
			`
		if (10 > 1) {
		  if (10 > 2) {
		    return 10;
		  }

		  return 1;
		}
		`,
			10,
		},
		{
			`
		let f = func(x) {
		  return x;
		  x + 10;
		};
		f(10);`,
			10,
		},
		{
			`
		let f = func(x) {
		   let result = x + 10;
		   return result;
		   return 10;
		};
		f(10);`,
			20,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			if !testIntegerObject(t, evaluated, tt.expected) {
				t.FailNow()
			}
		})
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
		{`len(1)`, "error: argument to `len` not supported, got INT"},
		{`len("one", "two")`, "error: builtin function 'len' requires exactly one argument, got 2"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`println("hello", "world!")`, nil},
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
			array, ok := evaluated.(*object.Array)
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
			return x + a;
		};
		add(add(10));
		`,
			result: 30,
		},
		{
			input: `let a = 10;
			let add = func(x) {
				return x + a;
			};
			a = add(add(10));
			a + a;
		`,
			result: 60,
		},
		{
			input: `let a = 10;
			let add = func(x) {
				return x + a;
			};
			a = add(add(10));
			add(a + a) + add(a + a);
		`,
			result: 180,
		},
		{
			input: `let a = 10.5;
			let add = func(x) {
				return x + a;
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
			let age = 10.5;
			for i = range [1..10] {
				let age = 4
				age++;
			};
			age;
		`,
			result: 10.5,
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
			input: `let sub
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
		{
			input: `let name = "foo";
			let age = 10.5;
			age += 10
			age`,
			result: 20.5,
		},
		{
			input: `let name = "foo";
			let age = 20.5;
			age -= 10
			age`,
			result: 10.5,
		},
		{
			input: `let arr = [1..10];
			arr[2] = 100;
			arr
		`,
			result: []any{1, 2, 100, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			input: `
			let arrk = [-10..-5].reverse()
			let arr = arrk[arrk.length()-1]
			arr
		`,
			result: -10,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			if !testObject(t, evaluated, tt.result) {
				t.Fatalf("expected %v, got %v", tt.result, evaluated.Inspect())
			}
		})
	}
}

func TestEvalLet(t *testing.T) {
	input := `let sub
	sub;
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.Parse()
	if len(program.Statements) != 2 {
		t.Fatalf("expected %d statements, got %d", 2, len(program.Statements))
	}

	if stmt, ok := program.Statements[0].(*ast.LetStmt); !ok {
		t.Fatalf("expected type *ast.LetStmt, got %T", stmt)
	}

}
func TestEvalStatements_Error(t *testing.T) {

	tests := []struct {
		input  string
		result []string
	}{
		{
			input: `let subjects = ["english", "french"];
			let index = foo;
		`,
			result: []string{"cannot assign to reserved keyword 'index'", "Line: 2"},
		},
		{
			input: `let arr = [2, ( + 5];
		`,
			result: []string{"expected closing parenthesis token ')', got", "Line: 1"},
		},
		{
			input: `let arr = [2, 3 +];
		`,
			result: []string{"invalid right expression ] for operator '+'", "Line: 1"},
		},
		{
			input: `let arr = [2, | +];
		`,
			result: []string{"illegal token '|'", "Line: 1"},
		},
		{
			input: `let arr = [2;
		`,
			result: []string{"expected closing bracket token ']', got ';'"},
		},
		{
			input: `a = 24;
			println(a);
		`,
			result: []string{"cannot reassign undeclared identifier 'a'", "Line: 1"},
		},
		{
			input: `let name = "foo;
		`,
			result: []string{"illegal token 'foo", "Line: 1"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			ev := testEval(tt.input)
			evaluated, ok := ev.(*object.Error)
			if !ok {
				t.Fatalf("expected result of type *object.Error, got %T", evaluated)
			}
			for _, str := range tt.result {
				if !strings.Contains(evaluated.Message, str) {
					t.Fatalf("expected \"%s\" to contain error \"%s\"", evaluated.Message, str)
				}
			}
		})
	}
}

func TestEvalStatements_ArrayOperations(t *testing.T) {

	tests := []struct {
		input  string
		result any
	}{
		{
			input: `
			let arr = [1,2];
			arr.push(5);
			println(arr);
			arr;
			`,
			result: []string{"1", "2", "5"},
		},
		{
			input: `
			let arr = [1,2,4];
			arr.pop();
			arr;
			`,
			result: []string{"1", "2"},
		},
		{
			input: `
			let arr = [1,2,4];
			let first = arr.first();
			first;
			`,
			result: 1,
		},
		{
			input: `
			let arr = [1,2,4];
			let last = arr.last();
			last;
			`,
			result: 4,
		},
		{
			input: `
			let arr = [1,2,4];
			arr.reverse(double);
			arr;
			`,
			result: []string{"4", "2", "1"},
		},
		{
			input: `
			let arr = [1,2,4];
			let double = func(x) { x + x; };
			arr.map(double);
			arr;
			`,
			result: []string{"2", "4", "8"},
		},
		{
			input: `
			let arr = [1,2,4];
			arr.merge([5,6,7]);
			arr;
			`,
			result: []string{"1", "2", "4", "5", "6", "7"},
		},
		{
			input: `
			let arr = [1,2,3,4,5,6];
			let even = func(x) { x % 2 == 0 };
			arr.filter(even);
			arr;
			`,
			result: []string{"2", "4", "6"},
		},
		{
			input: `
			let arr = [1,2,3,4,5,6];
			let even = func(x) { index >= 3 };
			arr.filter(even);
			arr;
			`,
			result: []string{"4", "5", "6"},
		},
		{
			input: `
			let arr = [1,2,3,4,5,6];
			arr.filter(func(x) { x % 2 == 0});
			arr;
			`,
			result: []string{"2", "4", "6"},
		},
		{
			input: `
			let arr = [1,2,3,4,5,6];
			let found = arr.contains(2);
			found;
			`,
			result: true,
		},
		{
			input: `
			let arr = [1,2,3,4,5,6];
			let found = arr.contains(10);
			found;
			`,
			result: false,
		},
		{
			input: `
			let arr = [1,2,3,4,5,6];
			let first_even = arr.find(func(x) { x % 2 == 0});
			first_even;
			`,
			result: 2,
		},
		{
			input: `
			let arr = [1,2,3,4,5,6];
			let arr_str = arr.join(" ");
			arr_str;
			`,
			result: "1 2 3 4 5 6",
		},
		{
			input: `
			let foo = [1, 2, 3];
			let len = foo.length();
			len
			`,
			result: 3,
		},
		{
			input: `
			let foo = [1, 2, 3];
			foo.clear();
			let len = foo.length();
			len
			`,
			result: 0,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			if tt.result == false {
				if evaluated, ok := evaluated.(*object.Boolean); ok {
					if evaluated.Value == false {
						return
					}
				}
			}
			if !testObject(t, evaluated, tt.result) {
				t.Fatalf("expected %v, got %v", tt.result, evaluated.Inspect())
			}
		})
	}
}

func TestEvalStatements_RangeArray(t *testing.T) {

	tests := []struct {
		input  string
		result any
	}{
		{
			input: `
			let arr = [1..5]
			arr.length()
			`,
			result: 5,
		},
		{
			input: `
			let arr = [-5..-1];
			arr.pop();
			arr;
			`,
			result: []string{"-5", "-4", "-3", "-2"},
		},
		{
			input: `
			let arr = [-3..10];
			let last = arr.last();
			last;
			`,
			result: 10,
		},
		{
			input: `
			let start = 5
			let arr = [start..10];
			arr.length();
			`,
			result: 6,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			if tt.result == false {
				if evaluated, ok := evaluated.(*object.Boolean); ok {
					if evaluated.Value == false {
						return
					}
				}
			}
			if !testObject(t, evaluated, tt.result) {
				t.Fatalf("expected %v, got %v", tt.result, evaluated.Inspect())
			}
		})
	}
}

func TestEvalStatements_HashOperations(t *testing.T) {

	tests := []struct {
		input  string
		result any
	}{
		{
			input: `
			let foo = {"a":"b"};
			let age = foo.contains("a");
			age
			`,
			result: true,
		},
		{
			input: `
			let foo = {"a":"b"};
			let age = foo.contains("c");
			age
			`,
			result: false,
		},
		{
			input: `
			let hash = {"a":10, 5.5:2,"bar":3};
			hash
			`,
			result: errors.New("invalid type *ast.FloatLiteral for hash key"),
		},
		{
			input: `
			let hash = {"a":"b","foo":2,"bar":3};
			let keys = hash.keys();
			keys
			`,
			result: []string{"a", "foo", "bar"},
		},
		{
			input: `
			let hash = {"a":"b","foo":2,"bar":3};
			let keys = hash.items();
			keys
			`,
			result: []any{"b", 2, 3},
		},
		{
			input: `
			let hash = {"a":"b","foo":2,"bar":3};
			let foo = hash.get("foo");
			foo
			`,
			result: 2,
		},
		{
			input: `
			let hash = {"a":"b","foo":2,"bar":3};
			hash.set("foo", 3);
			let foo = hash.get("foo");
			foo
			`,
			result: 3,
		},
		{
			input: `
			let hash = {"a":"b","foo":2,"bar":3};
			hash["foo"] = 100;
			let foo = hash.get("foo");
			foo
			`,
			result: 100,
		},
		{
			input: `
			let hash = {"a":"b","bar":3};
			hash["foo"] = "hello";
			let foo = hash.get("foo");
			foo
			`,
			result: "hello",
		},
		{
			input: `
			let hash = {"a":"b","foo":2,"bar":3};
			hash.clear();
			hash
			`,
			result: []string{},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			if err, ok := tt.result.(error); ok {
				if errObj, ok := evaluated.(*object.Error); ok {
					if !strings.Contains(errObj.Message, err.Error()) {
						t.Fatalf("expected \"%s\" to contain error \"%s\"", errObj.Message, err.Error())
					}
				} else {
					t.Fatalf("expected object to be of type *object.Error, got %T", evaluated)
				}
				return
			}
			if tt.result == false {
				if evaluated, ok := evaluated.(*object.Boolean); ok {
					if evaluated.Value == false {
						return
					}
				}
			}
			if !testObject(t, evaluated, tt.result) {
				t.Fatalf("expected %v, got %v", tt.result, evaluated.Inspect())
			}
		})
	}
}

func TestEvalStatements_SetOperations(t *testing.T) {

	tests := []struct {
		input  string
		result any
	}{
		{
			input: `
			let foo = {1, 2, 3};
			let found = foo.contains(1);
			found
			`,
			result: true,
		},
		{
			input: `
			let one = 1;
			let two = one + one;
			let three = 3;
			let foo = {one, two, three};
			let found = foo.contains(1);
			found
			`,
			result: true,
		},
		{
			input: `
			let three = 3;
			let foo = {1, 2, three};
			let found = foo.contains(4);
			found
			`,
			result: false,
		},
		{
			input: `
			let foo = {1, 2, 3};
			let len = foo.length();
			len
			`,
			result: 3,
		},
		{
			input: `
			let foo = {1, 2, 3};
			foo.clear();
			let len = foo.length();
			len
			`,
			result: 0,
		},
		{
			input: `
			let foo = {1, 2, 2, 3, 3, 3, 3, "3", "3"};
			foo.length();
			`,
			result: 4,
		},
		{
			input: `
			let foo = {1, 2};
			foo.add(4)
			foo.add(1)
			foo.length()
			`,
			result: 3,
		},
		{
			input: `
			let foo = {1, 2, 3};
			foo.delete(3)
			foo.length()
			`,
			result: 2,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			if err, ok := tt.result.(error); ok {
				if errObj, ok := evaluated.(*object.Error); ok {
					if !strings.Contains(errObj.Message, err.Error()) {
						t.Fatalf("expected \"%s\" to contain error \"%s\"", errObj.Message, err.Error())
					}
				} else {
					t.Fatalf("expected object to be of type *object.Error, got %T", evaluated)
				}
				return
			}
			if tt.result == false {
				if evaluated, ok := evaluated.(*object.Boolean); ok {
					if evaluated.Value == false {
						return
					}
				}
			}
			if !testObject(t, evaluated, tt.result) {
				t.Fatalf("expected %v, got %v", tt.result, evaluated.Inspect())
			}
		})
	}
}

func TestEvalStatements_SetOperations_Error(t *testing.T) {

	tests := []struct {
		input  string
		result any
	}{
		{
			input: `
			let arr = [1,2]
			let foo = {1, 2, arr};
			let found = match (foo.contains(1)) {
				case error: return error
			};
			`,
			result: errors.New("invalid set entry"),
		},
		{
			input: `
			let foo = {1, 2, 3};
			let found = foo.contains([1,2]);
			found
			`,
			result: false,
		},
		{
			input: `
			let foo = {1, 2, 3};
			foo.delete([1,2]);
			`,
			result: errors.New("cannot delete non-hashable"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			ev, evaluated := testEval2(tt.input)
			_ = ev
			if err, ok := tt.result.(error); ok {
				if ev.errStack == nil {
					t.Fatalf("expected an error, got nil")
				}
				if !strings.Contains(ev.errStack.Error(), err.Error()) {
					t.Fatalf("expected \"%s\" to contain error \"%s\"", ev.errStack.Error(), err.Error())
				}

				return
			}
			if tt.result == false {
				if evaluated, ok := evaluated.(*object.Boolean); ok {
					if evaluated.Value == false {
						return
					}
				}
			}
			if !testObject(t, evaluated, tt.result) {
				t.Fatalf("expected %v, got %v", tt.result, evaluated.Inspect())
			}
		})
	}
}

func TestEval_Files(t *testing.T) {
	data, err := os.ReadFile("../examples/primitives/string.ede")
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	t.Run("true", func(t *testing.T) {
		evaluated := testEval(string(data))
		_evaluated, ok := evaluated.(*object.Boolean)
		if !ok {
			t.Fatalf("expected a boolean to be returned, got %T", evaluated)
		}
		if !_evaluated.Value {
			t.Fatalf("expected value to be true, got %v", _evaluated.Value)
		}
	})

	t.Run("false", func(t *testing.T) {
		str := strings.Replace(string(data), "let input = \"()[]{}\"", "let input = \"([]{}\"", -1)
		evaluated := testEval(str)
		_evaluated, ok := evaluated.(*object.Boolean)
		if !ok {
			t.Fatalf("expected a boolean to be returned, got %T", evaluated)
		}
		if _evaluated.Value {
			t.Fatalf("expected value to be false, got %v", _evaluated.Value)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		data, err := os.ReadFile("../examples/modules/invalid_json.ede")
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
		evaluated := testEval(string(data))
		_evaluated, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("expected an error to be returned, got %T", evaluated)
		}
		exp := "error parsing string as json"
		if !strings.Contains(_evaluated.Message, exp) {
			t.Fatalf("expected value to be %s, got %s", _evaluated.Message, exp)
		}
	})
}

func TestEval_Match(t *testing.T) {
	t.Run("errored", func(t *testing.T) {
		input := `
	println("starting")
	let obj = match 10*"a" {
	case error: return error
	default: println(obj)
	}
	println("should not get here")
	`

		evaluated := testEval(input)
		_evaluated, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("expected an error to be returned, got %T", evaluated)
		}
		if !strings.Contains(_evaluated.Message, "invalid infix operator * for (10) and (a)") {
			t.Fatalf("expected error message to be 'invalid infix operator * for (10) and (a)', got '%s'", _evaluated.Message)
		}
	})

	t.Run("no error", func(t *testing.T) {
		input := `
	println("starting")
	let obj = match (10*10) {
	case error: return error
	default: println(obj)
	}
	println("should get here")
	`

		evaluated := testEval(input)
		if evaluated, ok := evaluated.(*object.Nil); !ok {
			t.Fatalf("expected *object.Nil to be returned, got %T", evaluated)
		}
	})

	t.Run("match statement", func(t *testing.T) {
		input := `
	println("starting")
	let age = 20
	match age < 10*10 {
	case true: age = age + (10*10)
	default: println("not true", age, "is not less than", 10 * 10)
	}
	println("should get here")

	age 
	`

		evaluated := testEval(input)
		if !testObject(t, evaluated, 120) {
			t.Fatalf("expected integer to be returned, got %T", evaluated)
		}
	})
}

func TestEval_TwoSum(t *testing.T) {
	input := `
	let two_sum = func(nums, target) {
		let comp;
		let map = {}
		for num = range nums {
			let comp = target - num
			if (map.contains(num)) {
				return [num, comp]
			}
			map.add(comp)
		}
	}

	let nums = [2,7,11,15]
	let target = 9
	let res = two_sum(nums, target)
	let exp = {2, 7}
	println(res)
	println("(res == [7, 2]) = ", res == [7, 2])
	println("(res.set() == exp) = ",res.set() == exp)
	res
`

	evaluated := testEval(input)
	exp := []string{"7", "2"}
	if !testObject(t, evaluated, exp) {
		t.Fatalf("expected %v, got %v", exp, evaluated.Inspect())
	}
}

func TestEval_BlockStmt(t *testing.T) {
	input := `
	let val = 30
	for i = range [1..2] {
		let val = 10
	}
	return val
`

	evaluated := testEval(input)
	exp := 30
	if !testObject(t, evaluated, exp) {
		t.Fatalf("expected %v, got %v", exp, evaluated.Inspect())
	}
}

func TestEval_Module(t *testing.T) {
	t.Run("json.parse", func(t *testing.T) {
		input := "import json;" +
			"let obj = json.parse(`{\"numbers\":[1,2],\"subjects\":{\"foo\":\"bar\"}}`);" +
			"obj;"

		evaluated := testEval(input)
		ev, ok := evaluated.(*object.Hash)
		if !ok {
			t.Fatalf("expected an hash to be returned, got %T", ev)
		}
		if len(ev.Entries) != 2 {
			t.Fatalf("expected hash of length 2 , got %v", len(ev.Entries))
		}
	})

	t.Run("json.string", func(t *testing.T) {
		input := "import json;" +
			"let hash = {\"numbers\":[1,2],\"subjects\":{\"foo\":\"bar\"}};" +
			"let obj = json.string(hash);" +
			"let reparsed = json.parse(obj);" +
			"println(hash);" +
			"println(reparsed);" +
			"reparsed == hash;"

		evaluated := testEval(input)
		ev, ok := evaluated.(*object.Boolean)
		if !ok {
			t.Fatalf("expected a boolean to be returned, got %T", ev)
		}
		if !ev.Value {
			t.Fatalf("expected reparsed object to be same as original object")
		}
	})
}

func TestEval_Method_Error(t *testing.T) {
	t.Run("unhandled(identifier not found)", func(t *testing.T) {
		input := "let obj = json.parse(`{\"numbers\":[1,2],\"subjects\":{\"foo\":\"bar\"}}`);" +
			"obj;"

		evaluated := testEval(input)
		ev, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("expected an error to be returned, got %T", ev)
		}
		exp := "identifier not found 'json'"
		if !strings.Contains(ev.Message, exp) {
			t.Fatalf("Error message '%s' does not contain '%s'", ev.Message, exp)
		}
	})

	t.Run("handled error", func(t *testing.T) {
		input := "let obj = match (json.parse(`{\"numbers\":[1,2],\"subjects\":{\"foo\":\"bar\"}}`)) {\n" +
			"case error: return \"fallback\"" +
			"};" +
			"obj;"

		evaluated := testEval(input)
		ev, ok := evaluated.(*object.String)
		if !ok {
			t.Fatalf("expected an error to be returned, got %T", ev)
		}
		exp := "fallback"
		if ev.Value != exp {
			t.Fatalf("expected '%s', gor '%s'", ev.Value, exp)
		}
	})

	t.Run("attr not found", func(t *testing.T) {
		input :=
			"let obj ={\"numbers\":[1,2],\"subjects\":{\"foo\":\"bar\"}};" +
				"obj[\"foo\"];"

		evaluated := testEval(input)
		ev, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("expected an error to be returned, got %T", ev)
		}
		exp := "invalid index entry 'foo' for 'obj'"
		if !strings.Contains(ev.Message, exp) {
			t.Fatalf("Error message '%s' does not contain '%s'", ev.Message, exp)
		}
	})
}

func TestEval(t *testing.T) {
	t.Skip()
	input := `
	let arr = [1..10];
	`

	input = `let arrk = [-10..-5].reverse()
	println(arrk)
	let arr = arrk[arrk.length()-1]
	println(arr)
	`

	evaluated := testEval(input)
	if !testObject(t, evaluated, []string{"2", "4", "foofoo"}) {
		t.Fatalf("%v", evaluated)
	}
}
