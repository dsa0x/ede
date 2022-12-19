package object

import (
	"ede/ast"
	"strings"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type String struct{ Value string }
type Int struct{ Value int64 }
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

func NewErrorWithMsg(msgs ...string) *Error {
	return &Error{Message: strings.Join(msgs, "; ")}
}

func NewError(msgs ...error) *Error {
	messages := []string{}
	for _, msg := range msgs {
		messages = append(messages, msg.Error())
	}
	return NewErrorWithMsg(messages...)
}
