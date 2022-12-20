package evaluator

import (
	"ede/ast"
	"ede/object"
)

func evalProgram(node *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	if node.ParseErrors != nil {
		return object.NewError(node.ParseErrors)
	}

	for _, stmt := range node.Statements {
		if _, isComment := stmt.(*ast.CommentStmt); isComment {
			continue
		}
		result = (&Evaluator{}).Eval(stmt, env)
		if result == nil {
			return NULL
		}
		// return internal value
		if result.Type() == object.RETURN_VALUE_OBJ {
			return result.(*object.ReturnValue).Value
		}

		// terminate after encountering an error
		if result.Type() == object.ERROR_OBJ {
			return result
		}
	}
	return result
}
