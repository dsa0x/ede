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
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.BooleanLiteral:
		return &object.Boolean{Value: node.Value}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IfStmt:
		return evalIfExpression(node, env)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PostfixExpression:
		left := Eval(node.Left, env)
		result := evalPostfixExpression(node.Operator, left)
		if _, ok := node.Left.(*ast.Identifier); ok && !isError(result) { // update identifier
			env.Set(node.Left.Literal(), result)
		}
		return result
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)
	case *ast.ReturnExpression:
		return evalReturnExpression(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
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
		return &object.Function{Body: node.Body, Params: node.Params, ParentEnv: env}
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalArgs(node.Args, env)
		return applyFunction(fn, args)
	case *ast.ArrayLiteral:
		entries := evalArgs(node.Elements, env)
		return &Array{Entries: &entries}
	case *ast.ReassignmentStmt:
		if _, found := env.Get(node.Name.Value); !found {
			return object.NewErrorWithMsg(fmt.Sprintf("cannot reassign undeclared identifier '%s'", node.Name.Value))
		}
		res := Eval(node.Expr, env)
		env.Set(node.Name.Value, res)
		return res
	case *ast.ForLoopStmt:
		return evalForLoopStmt(node, env)
	case *ast.ObjectMethodExpression:
		return evalObjectMethodExpr[any](node, env)
	}

	return nil
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
	result := make([]object.Object, 0, len(args))

	for i, arg := range args {
		result = append(result, Eval(arg, env))
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

func evalIfExpression(node *ast.IfStmt, env *object.Environment) object.Object {
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

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	switch left := node.Left.(type) {
	case *ast.Identifier:
		ident := Eval(left, env)
		if index, ok := node.Index.(*ast.IntegerLiteral); ok {
			if arr, ok := ident.(*Array); ok {
				return (*arr.Entries)[index.Value]
			}
		}

	case *ast.ArrayLiteral:
		if index, ok := node.Index.(*ast.IntegerLiteral); ok {
			if int(index.Value) >= len(left.Elements) {
				return object.NewErrorWithMsg(fmt.Sprintf("index %d out of range", index.Value))
			}
			return Eval(left.Elements[index.Value], env)
		}
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid index operator %s for %s", node.Index.Literal(), node.Left))
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
	if left.Type() == object.FLOAT_OBJ {
		left := left.(*object.Float)
		switch operator {
		case token.INC:
			return &object.Float{Value: left.Value + 1}
		case token.DEC:
			return &object.Float{Value: left.Value - 1}
		}
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid postfix operator %s for %s", operator, left.Inspect()))
}

func evalBangOperator(operator string, right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Int:
		return booleanObj(right.Value == 0)
	case *object.String:
		return booleanObj(len(right.Value) == 0)
	case *Array:
		return booleanObj(len(*right.Entries) == 0)
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

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		fnEnv := object.NewEnvironment(fn.ParentEnv)
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
