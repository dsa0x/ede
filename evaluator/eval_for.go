package evaluator

import (
	"ede/ast"
	"ede/object"
	"ede/token"
)

func evalForLoopStmt(node *ast.ForLoopStmt, env *object.Environment) object.Object {
	var result object.Object
	if node == nil {
		return NULL
	}

	switch boundRange := node.Boundary.(type) {
	case *ast.ArrayLiteral:
		for i, el := range boundRange.Elements {
			env.Set(token.IndexIdentifier, &object.Int{Value: int64(i)})
			env.Set(node.Variable.Value, Eval(el, env)) // bound loop variable
			for _, stmt := range node.Statement.Statements {
				result = Eval(stmt, env)
			}
		}
	case *ast.Identifier:
		ident := Eval(boundRange, env)
		arr := ident.(*object.Array[any]) // TODO: may change when we support more
		for i, entry := range arr.Entries {
			env.Set(token.IndexIdentifier, &object.Int{Value: int64(i)})
			env.Set(node.Variable.Value, entry) // bound loop variable
			for _, stmt := range node.Statement.Statements {
				result = Eval(stmt, env)
			}
		}
	default:
		return nil
	}
	return result
}
