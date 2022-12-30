package evaluator

import (
	"ede/ast"
	"ede/object"
	"ede/token"
)

// Iterable is any object that can be iterated over, e.g. arrays, strings
type Iterable interface {
	Items() []object.Object
}

func (e *Evaluator) evalForLoopStmt(node *ast.ForLoopStmt, env *object.Environment) object.Object {
	var result object.Object
	if node == nil {
		return NULL
	}

	var arr Iterable
	var iter object.Object
	var ok bool

	switch boundRange := node.Boundary.(type) {
	case *ast.ArrayLiteral:
		// TODO: make range array zero index
		for i, el := range boundRange.Elements {
			// create an environment for the block statemet
			blockEnv := object.NewEnvironment(env)
			blockEnv.Set(token.IndexIdentifier, &object.Int{Value: int64(i)})
			blockEnv.Set(node.Variable.Value, e.Eval(el, blockEnv)) // bound loop variable
			result = e.evalBlockStmt(node.Statement, blockEnv)
		}
		return result
	case *ast.Identifier:
		ident := e.Eval(boundRange, env)
		if ident == nil {
			return object.NewErrorWithMsg("invalid identifier '%s'", boundRange.Value)
		}
		iter = ident
	case *ast.ObjectMethodExpression, *ast.RangeArrayLiteral:
		iter = e.Eval(boundRange, env)
	}

	if arr, ok = iter.(Iterable); !ok {
		return object.NewErrorWithMsg("for loop boundary type is not iterable, got %T", iter)
	}

	for i, entry := range arr.Items() {
		stmtEnv := object.NewEnvironment(env)
		stmtEnv.Set(token.IndexIdentifier, &object.Int{Value: int64(i)})
		stmtEnv.Set(node.Variable.Value, entry) // bound loop variable
		for _, stmt := range node.Statement.Statements {
			result = e.Eval(stmt, stmtEnv)
			if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
				return result
			}
		}
	}
	// if the returned value is a return object or an error,
	// then we would have returned
	// else we return nil, because it's a statement
	return NULL
}
