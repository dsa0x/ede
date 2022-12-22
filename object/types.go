package object

import (
	"ede/ast"
	"fmt"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type String struct{ Value string }
type Float struct{ Value float64 }
type Boolean struct{ Value bool }
type Null struct{}
type Error struct{ Message string }
type ReturnValue struct{ Value Object }
type Function struct {
	Params    []*ast.Identifier
	Body      *ast.BlockStmt
	ParentEnv *Environment
}
type BuiltinFn func(args ...Object) Object
type Builtin struct{ Fn BuiltinFn }

func NewErrorWithMsg(msg string, format ...any) *Error {
	return &Error{Message: fmt.Sprintf(msg, format...)}
}

func NewError(msg error) *Error {
	return NewErrorWithMsg(msg.Error())
}
