package evaluator

import (
	"ede/ast"
	"ede/object"
	"ede/token"
	"fmt"
)

var (
	NULL  = object.NULL
	TRUE  = object.TRUE
	FALSE = object.FALSE
)

type Evaluator struct {
	pos token.Pos
}

func (e *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	if node == nil {
		return nil
	}
	e.pos = node.Pos()
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
		return e.evalIfExpression(node, env)
	case *ast.InfixExpression:
		left := e.Eval(node.Left, env)
		right := e.Eval(node.Right, env)
		return e.evalInfixExpression(node.Operator, left, right)
	case *ast.PostfixExpression:
		left := e.Eval(node.Left, env)
		result := e.evalPostfixExpression(node.Operator, left)
		if _, ok := node.Left.(*ast.Identifier); ok && !isError(result) { // update identifier
			env.Set(node.Left.Literal(), result)
		}
		return result
	case *ast.PrefixExpression:
		right := e.Eval(node.Right, env)
		return e.evalPrefixExpression(node.Operator, right)
	case *ast.ReturnExpression:
		return e.evalReturnExpression(node, env)
	case *ast.IndexExpression:
		return e.evalIndexExpression(node, env)
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.LetStmt:
		ident := e.Eval(node.Expr, env)
		if !isError(ident) {
			env.Set(node.Name.Value, ident)
		}
		return ident
	case *ast.BlockStmt:
		return e.evalBlockStmt(node, env)
	case *ast.ExpressionStmt:
		return e.Eval(node.Expr, env)
	case *ast.ConditionalStmt:
		return e.Eval(node.Statement, env)
	case *ast.FunctionLiteral:
		return &object.Function{Body: node.Body, Params: node.Params, ParentEnv: env}
	case *ast.CallExpression:
		fn := e.Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := e.evalArgs(node.Args, env)
		return e.applyFunction(fn, args)
	case *ast.ArrayLiteral:
		entries := e.evalArgs(node.Elements, env)
		return &object.Array{Entries: &entries}
	case *ast.HashLiteral:
		entries := e.evalPairs(node.Pair, env)
		for _, val := range entries {
			if isError(val) {
				return val
			}
		}
		return &object.Hash{Entries: entries}
	case *ast.ReassignmentStmt:
		if _, found := env.Get(node.Name.Value); !found {
			return e.EvalError(fmt.Sprintf("cannot reassign undeclared identifier '%s'", node.Name.Value), node.Pos())
		}
		res := e.Eval(node.Expr, env)
		env.Set(node.Name.Value, res)
		return res
	case *ast.ForLoopStmt:
		return e.evalForLoopStmt(node, env)
	case *ast.ObjectMethodExpression:
		return e.evalObjectMethodExpr(node, env)
	}

	return nil
}

func (e *Evaluator) evalReturnExpression(node *ast.ReturnExpression, env *object.Environment) object.Object {
	returnVal := e.Eval(node.Expr, env)
	if returnVal.Type() == object.ERROR_OBJ {
		return returnVal
	}
	// wrap the value so that block statements can terminate early if they encounter a return
	return &object.ReturnValue{Value: returnVal}
}

// evalArgs evaluates arguments
func (e *Evaluator) evalArgs(args []ast.Expression, env *object.Environment) []object.Object {
	result := make([]object.Object, 0, len(args))

	for i, arg := range args {
		result = append(result, e.Eval(arg, env))
		if isError(result[i]) {
			return []object.Object{result[i]}
		}
	}

	return result
}

func (e *Evaluator) evalPairs(args map[ast.Expression]ast.Expression, env *object.Environment) map[string]object.Object {
	result := make(map[string]object.Object)

	for key, value := range args {
		keyObj := e.Eval(key, env)
		keyString := object.ToRawValue(keyObj)
		if keyString == "" {
			result[keyString] = object.NewErrorWithMsg(fmt.Sprintf("invalid key '%s'", keyString))
			return result
		}
		valueObj := e.Eval(value, env)
		if isError(keyObj) || isError(valueObj) {
			result[keyString] = valueObj
			return result
		}
		result[keyString] = valueObj
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

func (e *Evaluator) evalBlockStmt(node *ast.BlockStmt, env *object.Environment) object.Object {
	var result object.Object
	if node == nil {
		return NULL
	}

	for _, stmt := range node.Statements {
		result = e.Eval(stmt, env)
		if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}
	}
	return result
}

func (e *Evaluator) evalIfExpression(node *ast.IfStmt, env *object.Environment) object.Object {
	cond := e.Eval(node.Consequence.Condition, env)
	if isTruthy(cond) {
		return e.Eval(node.Consequence, env)
	} else {
		for _, alt := range node.Alternatives {
			if alt.Condition == nil { // normal else branch (else)
				return e.Eval(alt, env)
			}
			cond := e.Eval(alt.Condition, env) // (else if branch)
			if isTruthy(cond) {
				return e.Eval(alt, env)
			}
		}
	}
	return NULL
}

func (e *Evaluator) evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	switch left := node.Left.(type) {
	case *ast.Identifier:
		ident := e.Eval(left, env)
		if index, ok := node.Index.(*ast.IntegerLiteral); ok {
			if arr, ok := ident.(*object.Array); ok {
				if int(index.Value) >= len(*arr.Entries) {
					return e.EvalError(fmt.Sprintf("index %d out of range with length %d", index.Value, len(*arr.Entries)), node.Pos())
				}
				return (*arr.Entries)[index.Value]
			}
		}

	case *ast.ArrayLiteral:
		if index, ok := node.Index.(*ast.IntegerLiteral); ok {
			if int(index.Value) >= len(left.Elements) {
				return e.EvalError(fmt.Sprintf("index %d out of range with length %d", index.Value, len(left.Elements)), node.Pos())
			}
			return e.Eval(left.Elements[index.Value], env)
		}
	}
	return e.EvalError(fmt.Sprintf("invalid index operator %s for %s", node.Index.Literal(), node.Left), node.Pos())
}

func (e *Evaluator) evalPostfixExpression(operator string, left object.Object) object.Object {
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
	return e.EvalError(fmt.Sprintf("invalid postfix operator %s for %s", operator, left.Inspect()), e.pos)
}

func (e *Evaluator) evalBangOperator(operator string, right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Int:
		return e.booleanObj(right.Value == 0)
	case *object.String:
		return e.booleanObj(len(right.Value) == 0)
	case *object.Array:
		return e.booleanObj(len(*right.Entries) == 0)
	case *object.Boolean:
		return e.booleanObj(!right.Value)
	case *object.Null:
		return TRUE
	}
	return e.EvalError(fmt.Sprintf("bang operator not valid for type %T", right), e.pos)
}

func (e *Evaluator) booleanObj(boolean bool) *object.Boolean {
	if boolean {
		return TRUE
	}
	return FALSE
}

func (e *Evaluator) applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		fnEnv := object.NewEnvironment(fn.ParentEnv)
		for i, p := range fn.Params {
			fnEnv.Set(p.Value, args[i])
		}
		result := e.Eval(fn.Body, fnEnv)
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

func (e *Evaluator) EvalError(err string, pos token.Pos) *object.Error {
	msg := fmt.Sprintf(`
	Error: %s
	Line: %d
	Column: %d
	`, err, pos.Line, pos.Column)
	return &object.Error{Message: msg}
}
