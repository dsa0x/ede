package evaluator

import (
	"ede/ast"
	"ede/object"
	"fmt"
)

type Methodable interface {
	GetMethod(name string, eval object.Evaluator) *object.Builtin
}

func (e *Evaluator) evalObjectDotExpr(node *ast.ObjectMethodExpression, env *object.Environment) object.Object {
	if node == nil {
		return NULL
	}

	obj := e.Eval(node.Object, env)
	if isError(obj) {
		return obj
	}
	if obj == nil {
		return object.NewErrorWithMsg("identifier not found '%s'", node.Object.Literal())
	}

	if method, ok := node.Method.(*ast.CallExpression); ok {
		return e.evalObjectMethodExpr(obj, method, env)
	}
	// if the right side of the dot is not a method call
	if ident, ok := node.Method.(*ast.Identifier); ok {
		return e.evalObjectAttrExpr(obj, ident, env)
	}
	return object.NewErrorWithMsg("expected method call or identifier, got %s", node.Method.TokenType())
}

func (e *Evaluator) evalObjectMethodExpr(obj object.Object, call *ast.CallExpression, env *object.Environment) object.Object {

	args := e.evalArgs(call.Args, env)

	methodableObj, ok := obj.(Methodable)
	if !ok {
		return nil
	}
	ident := call.Function.(*ast.Identifier)
	if ident.Value == "equal" {
		return e.evalEqualMethod(obj, args...)
	} else if ident.Value == "type" {
		return e.evalTypeMethod(obj, args...)
	}

	method := methodableObj.GetMethod(ident.Value, e)
	if method == nil {
		return object.NewErrorWithMsg(fmt.Sprintf("unknown method '%s' for module '%s'", ident.Value, obj.Inspect()))
	}
	return method.Fn(args...)
}

func (e *Evaluator) evalObjectAttrExpr(obj object.Object, attr *ast.Identifier, env *object.Environment) object.Object {
	switch obj := obj.(type) {
	case *object.Hash:
		return obj.Entries[attr.Value]
	}
	return nil
}

// evalEqualMethod evaluates the equal method. All objects implement this, and it is false
// for functions
func (e *Evaluator) evalEqualMethod(obj object.Object, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewErrorWithMsg(fmt.Sprintf("method 'equal' requires 1 argument, got %d", len(args)))
	}
	return object.NewBoolean(obj.Equal(args[0]))
}

func (e *Evaluator) evalTypeMethod(obj object.Object, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewErrorWithMsg(fmt.Sprintf("method 'type' requires no argument, got %d", len(args)))
	}
	return object.NewString(string(obj.Type()))
}
