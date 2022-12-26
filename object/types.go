package object

import (
	"ede/ast"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}
type (
	String      struct{ Value string }
	Float       struct{ Value float64 }
	Boolean     struct{ Value bool }
	Nil         struct{}
	Error       struct{ Message string }
	ReturnValue struct{ Value Object }
	Hash        struct{ Entries map[string]Object }

	Function struct {
		Params    []*ast.Identifier
		Body      *ast.BlockStmt
		ParentEnv *Environment
	}

	BuiltinFn func(args ...Object) Object
	Builtin   struct{ Fn BuiltinFn }
)

func (a *Nil) Native() any {
	return nil
}
