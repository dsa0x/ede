package evaluator

import (
	"ede/ast"
	"ede/object"
)

func (e *Evaluator) evalObjectMethodExpr(node *ast.ObjectMethodExpression, env *object.Environment) object.Object {
	if node == nil {
		return NULL
	}

	obj := e.Eval(node.Object, env)
	call, ok := node.Method.(*ast.CallExpression)
	if !ok {
		return nil
	}

	switch obj := obj.(type) {
	case *object.Array:
		ident := call.Function.(*ast.Identifier)
		method := obj.GetMethod(ident.Value, e)
		args := e.evalArgs(call.Args, env)
		return method.Fn(args...)
	}

	return obj
}
