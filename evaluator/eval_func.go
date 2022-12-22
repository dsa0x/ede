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
	args := e.evalArgs(call.Args, env)

	methodableObj, ok := obj.(Methodable)
	if !ok {
		return nil
	}
	ident := call.Function.(*ast.Identifier)
	if ident.Value == "equal" {
		return e.evalEqualMethod(obj, args...)
	}

	method := methodableObj.GetMethod(ident.Value, e)
	if method == nil {
		return object.NewErrorWithMsg(fmt.Sprintf("unknown method '%s' for type %T", ident.Value, obj))
	}
	return method.Fn(args...)
}

// evalEqualMethod evaluates the equal method. All objects implement this, and it is false
// for functions
func (e *Evaluator) evalEqualMethod(obj object.Object, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewErrorWithMsg(fmt.Sprintf("method 'equal' requires 1 argument, got %d", len(args)))
	}
	return &object.Boolean{Value: obj.Equal(args[0])}
}
