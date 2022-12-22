package evaluator

import (
	"ede/ast"
	"ede/object"
	"fmt"
)

type Methodable interface {
	GetMethod(name string, eval object.Evaluator) *object.Builtin
}

func (e *Evaluator) evalObjectMethodExpr(node *ast.ObjectMethodExpression, env *object.Environment) object.Object {
	if node == nil {
		return NULL
	}

	obj := e.Eval(node.Object, env)
	call, ok := node.Method.(*ast.CallExpression)
	if !ok {
		return nil
	}

	methodableObj, ok := obj.(Methodable)
	if !ok {
		return nil
	}
	ident := call.Function.(*ast.Identifier)
	method := methodableObj.GetMethod(ident.Value, e)
	if method == nil {
		return object.NewErrorWithMsg(fmt.Sprintf("unknown method '%s' for type %T", ident.Value, obj))
	}
	args := e.evalArgs(call.Args, env)
	return method.Fn(args...)
}
