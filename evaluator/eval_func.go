package evaluator

import (
	"ede/ast"
	"ede/object"
)

func evalObjectMethodExpr[T any](node *ast.ObjectMethodExpression, env *object.Environment) object.Object {
	if node == nil {
		return NULL
	}

	obj := Eval(node.Object, env)
	call, ok := node.Method.(*ast.CallExpression)
	if !ok {
		return nil
	}

	switch obj := obj.(type) {
	case *Array:
		ident := call.Function.(*ast.Identifier)
		method := obj.GetMethod(ident.Value)
		args := evalArgs(call.Args, env)
		return method.Fn(args...)
	}

	return obj
}
