package evaluator

import (
	"ede/ast"
	"ede/object"
	"fmt"
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
		if method == nil {
			return object.NewErrorWithMsg(fmt.Sprintf("unknown method '%s' for type %T", ident.Value, obj))
		}
		args := e.evalArgs(call.Args, env)
		return method.Fn(args...)
	case *object.Hash:
		ident := call.Function.(*ast.Identifier)
		method := obj.GetMethod(ident.Value, e)
		if method == nil {
			return object.NewErrorWithMsg(fmt.Sprintf("unknown method '%s' for type %T", ident.Value, obj))
		}
		args := e.evalArgs(call.Args, env)
		return method.Fn(args...)
	}

	return obj
}
