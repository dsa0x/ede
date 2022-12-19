package object

import (
	"fmt"
)

type Type string

var (
	STRING_OBJ       Type = "STRING"
	BOOLEAN_OBJ      Type = "BOOLEAN"
	INT_OBJ          Type = "INT"
	FLOAT_OBJ        Type = "FLOAT"
	FUNCTION_OBJ     Type = "FUNCTION"
	ERROR_OBJ        Type = "ERROR"
	NULL_OBJ         Type = "NULL"
	RETURN_VALUE_OBJ Type = "RETURN_VALUE"
	BUILTIN_OBJ      Type = "BUILTIN"
	ARRAY_OBJ        Type = "ARRAY"
)

type Object interface {
	Type() Type
	Inspect() string
}

func (*String) Type() Type      { return STRING_OBJ }
func (*Int) Type() Type         { return INT_OBJ }
func (*Float) Type() Type       { return FLOAT_OBJ }
func (*Boolean) Type() Type     { return BOOLEAN_OBJ }
func (*Error) Type() Type       { return ERROR_OBJ }
func (*Null) Type() Type        { return NULL_OBJ }
func (*ReturnValue) Type() Type { return RETURN_VALUE_OBJ }
func (*Function) Type() Type    { return FUNCTION_OBJ }
func (*Builtin) Type() Type     { return BUILTIN_OBJ }

func (v *String) Inspect() string      { return v.Value }
func (v *Int) Inspect() string         { return fmt.Sprint(v.Value) }
func (v *Float) Inspect() string       { return fmt.Sprint(v.Value) }
func (v *Boolean) Inspect() string     { return fmt.Sprint(v.Value) }
func (v *Error) Inspect() string       { return fmt.Sprint(v.Message) }
func (v *Null) Inspect() string        { return "null" }
func (v *ReturnValue) Inspect() string { return v.Value.Inspect() }
func (v *Function) Inspect() string    { return "fn" }
func (*Builtin) Inspect() string       { return "builtin fn" }
