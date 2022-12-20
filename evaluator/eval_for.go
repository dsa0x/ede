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

	switch boundRange := node.Boundary.(type) {
	case *ast.ArrayLiteral:
		for i, el := range boundRange.Elements {
			env.Set(token.IndexIdentifier, &object.Int{Value: int64(i)})
			env.Set(node.Variable.Value, e.Eval(el, env)) // bound loop variable
			for _, stmt := range node.Statement.Statements {
				result = e.Eval(stmt, env)
			}
		}
	case *ast.Identifier:
		ident := e.Eval(boundRange, env)
		arr := ident.(*object.Array) // TODO: may change when we support more
		for i, entry := range *arr.Entries {
			env.Set(token.IndexIdentifier, &object.Int{Value: int64(i)})
			env.Set(node.Variable.Value, entry) // bound loop variable
			for _, stmt := range node.Statement.Statements {
				result = e.Eval(stmt, env)
			}
		}
	default:
		return nil
	}
	return result
}
