package evaluator

import (
	"ede/ast"
	"ede/object"
	"ede/token"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}
	case *ast.BooleanLiteral:
		return &object.Boolean{Value: node.Value}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PostfixExpression:
		left := Eval(node.Left, env)
		return evalPostfixExpression(node.Operator, left)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)
	case *ast.ReturnExpression:
		return evalReturnExpression(node, env)
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.LetStmt:
		ident := Eval(node.Expr, env)
		if !isError(ident) {
			env.Set(node.Name.Value, ident)
		}
		return ident
	case *ast.BlockStmt:
		return evalBlockStmt(node, env)
	case *ast.ExpressionStmt:
		return Eval(node.Expr, env)
	case *ast.ConditionalStmt:
		return Eval(node.Statement, env)
	case *ast.FunctionLiteral:
		return &object.Function{Body: node.Body, Params: node.Params}
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalArgs(node.Args, env)
		return applyFunction(fn, args, env)
	case *ast.ArrayLiteral:
		entries := make([]object.Object, len(node.Elements))
		for i, el := range node.Elements {
			entries[i] = Eval(el, env)
		}
		return &object.Array{Entries: entries}
	}

	return nil
}

func evalProgram(node *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	if len(node.ParseErrors) > 0 {
		return object.NewError(node.ParseErrors...)
	}

	for _, stmt := range node.Statements {
		result = Eval(stmt, env)
		// return internal value
		if result.Type() == object.RETURN_VALUE_OBJ {
			return result.(*object.ReturnValue).Value
		}

		// terminate after encountering an error
		if result.Type() == object.ERROR_OBJ {
			return result
		}
	}
	return result
}

func evalReturnExpression(node *ast.ReturnExpression, env *object.Environment) object.Object {
	returnVal := Eval(node.Expr, env)
	if returnVal.Type() == object.ERROR_OBJ {
		return returnVal
	}
	// wrap the value so that block statements can terminate early if they encounter a return
	return &object.ReturnValue{Value: returnVal}
}

// evalArgs evaluates arguments
func evalArgs(args []ast.Expression, env *object.Environment) []object.Object {
	result := make([]object.Object, len(args))

	for i, arg := range args {
		result[i] = Eval(arg, env)
		if isError(result[i]) {
			return []object.Object{result[i]}
		}
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if obj, ok := env.Get(node.Value); ok {
		return obj
	}

	if b, ok := builtins[node.Value]; ok {
		return b
	}

	return nil
}

func evalBlockStmt(node *ast.BlockStmt, env *object.Environment) object.Object {
	var result object.Object
	if node == nil {
		return NULL
	}

	for _, stmt := range node.Statements {
		result = Eval(stmt, env)
		if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}
	}
	return result
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(node.Consequence.Condition, env)
	if isTruthy(cond) {
		return Eval(node.Consequence, env)
	} else {
		for _, alt := range node.Alternatives {
			if alt.Condition == nil { // normal else branch (else)
				return Eval(alt, env)
			}
			cond := Eval(alt.Condition, env) // (else if branch)
			if isTruthy(cond) {
				return Eval(alt, env)
			}
		}
	}
	return NULL
}
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch true {
	// bang operator for all types
	case right.Type() != object.ERROR_OBJ && operator == "!":
		return evalBangOperator(operator, right)
	case right.Type() == object.INT_OBJ:
		right := right.(*object.Int)
		return evalIntegerPrefixExpression(operator, right)
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid prefix operator %s for %s", operator, right.Inspect()))
}

func evalPostfixExpression(operator string, left object.Object) object.Object {
	if left.Type() == object.INT_OBJ {
		left := left.(*object.Int)
		switch operator {
		case token.INC:
			return &object.Int{Value: left.Value + 1}
		case token.DEC:
			return &object.Int{Value: left.Value - 1}
		}
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid postfix operator %s for %s", operator, left.Inspect()))
}

func evalIntegerPrefixExpression(operator string, right *object.Int) object.Object {
	switch operator {
	case "-":
		return &object.Int{Value: -right.Value}
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid integer operator %s", operator))
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch true {
	case left.Type() == object.INT_OBJ && right.Type() == object.INT_OBJ:
		left := left.(*object.Int)
		right := right.(*object.Int)
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		left := left.(*object.String)
		right := right.(*object.String)
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		left := left.(*object.Boolean)
		right := right.(*object.Boolean)
		return evalBoolInfixExpression(operator, left, right)
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid infix operator %s for %s and %s", operator, left.Inspect(), right.Inspect()))
}

func evalIntegerInfixExpression(operator string, left, right *object.Int) object.Object {
	switch operator {
	case "+":
		return &object.Int{Value: left.Value + right.Value}
	case "-":
		return &object.Int{Value: left.Value - right.Value}
	case "*":
		return &object.Int{Value: left.Value * right.Value}
	case "/":
		return &object.Int{Value: left.Value / right.Value}
	case ">":
		return booleanObj(left.Value > right.Value)
	case "<":
		return booleanObj(left.Value < right.Value)
	case "==":
		return booleanObj(left.Value == right.Value)
	case "!=":
		return booleanObj(left.Value != right.Value)
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid integer operator %s", operator))
}

func evalStringInfixExpression(operator string, left, right *object.String) object.Object {
	switch operator {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid string operator %s", operator))
}
func evalBoolInfixExpression(operator string, left, right *object.Boolean) object.Object {
	switch operator {
	case "&&":
		return booleanObj(left.Value && right.Value)
	case "||":
		return booleanObj(left.Value || right.Value)
	case "==":
		return booleanObj(left.Value == right.Value)
	case "!=":
		return booleanObj(left.Value != right.Value)
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid boolean operator %s", operator))
}

func evalBangOperator(operator string, right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Int:
		return booleanObj(right.Value == 0)
	case *object.String:
		return booleanObj(len(right.Value) == 0)
	case *object.Array:
		return booleanObj(len(right.Entries) == 0)
	case *object.Boolean:
		return booleanObj(!right.Value)
	case *object.Null:
		return TRUE
	}
	return object.NewErrorWithMsg(fmt.Sprintf("bang operator not valid for type %T", right))
}

func booleanObj(boolean bool) *object.Boolean {
	if boolean {
		return TRUE
	}
	return FALSE
}

func applyFunction(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		fnEnv := object.NewEnvironment(env)
		for i, p := range fn.Params {
			fnEnv.Set(p.Value, args[i])
		}
		result := Eval(fn.Body, fnEnv)
		return unwrapReturnValue(result)
	case *object.Builtin:
		return fn.Fn(args...)
	}
	return nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func isTruthy(obj object.Object) bool {
	switch obj.Type() {
	case object.BOOLEAN_OBJ:
		obj := obj.(*object.Boolean)
		return obj.Value
	case object.INT_OBJ:
		obj := obj.(*object.Int)
		return obj.Value != 0
	case object.ERROR_OBJ:
		obj := obj.(*object.Error)
		return obj.Message == ""
	}
	return true
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
