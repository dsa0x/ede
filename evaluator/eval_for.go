package evaluator

import (
	"ede/ast"
	"ede/object"
	"ede/token"
)

func (e *Evaluator) evalForLoopStmt(node *ast.ForLoopStmt, env *object.Environment) object.Object {
	var result object.Object
	if node == nil {
		return NULL
	}

	var arr *object.Array

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
		arr = ident.(*object.Array) // TODO: may change when we support more
	case *ast.ObjectMethodExpression:
		ident := e.Eval(boundRange, env)
		arr, _ = ident.(*object.Array)
	}

	if arr == nil {
		return object.NewErrorWithMsg("invalid for loop range boundary type: %T", node.Boundary)
	}
	for i, entry := range *arr.Entries {
		env.Set(token.IndexIdentifier, &object.Int{Value: int64(i)})
		env.Set(node.Variable.Value, entry) // bound loop variable
		for _, stmt := range node.Statement.Statements {
			result = e.Eval(stmt, env)
			if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
				return result
			}
		}
	}
	return result
}
