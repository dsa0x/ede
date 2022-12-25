package evaluator

import (
	"ede/ast"
	"ede/object"

	"github.com/hashicorp/go-multierror"
)

func (e *Evaluator) evalProgram(node *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	if node.ParseErrors != nil {
		return object.NewError(node.ParseErrors)
	}

	for _, stmt := range node.Statements {
		if _, isComment := stmt.(*ast.CommentStmt); isComment {
			continue
		}
		result = e.Eval(stmt, env)
		if result == nil {
			continue
		}
		// return internal value
		if result.Type() == object.RETURN_VALUE_OBJ {
			return result.(*object.ReturnValue).Value
		}

		// terminate after encountering an error
		if result.Type() == object.ERROR_OBJ {
			e.errStack = multierror.Append(e.errStack, result.Native().(error))
			return result
		}
	}
	return result
}
