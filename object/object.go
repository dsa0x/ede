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

	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Object interface {
	Type() Type
	Inspect() string
	Equal(obj Object) bool
}

func (*String) Type() Type      { return STRING_OBJ }
func (*Float) Type() Type       { return FLOAT_OBJ }
func (*Boolean) Type() Type     { return BOOLEAN_OBJ }
func (*Error) Type() Type       { return ERROR_OBJ }
func (*Null) Type() Type        { return NULL_OBJ }
func (*ReturnValue) Type() Type { return RETURN_VALUE_OBJ }
func (*Function) Type() Type    { return FUNCTION_OBJ }
func (*Builtin) Type() Type     { return BUILTIN_OBJ }

func (v *String) Inspect() string      { return v.Value }
func (v *Float) Inspect() string       { return fmt.Sprint(v.Value) }
func (v *Boolean) Inspect() string     { return fmt.Sprint(v.Value) }
func (v *Error) Inspect() string       { return fmt.Sprint(v.Message) }
func (v *Null) Inspect() string        { return "null" }
func (v *ReturnValue) Inspect() string { return v.Value.Inspect() }
func (v *Function) Inspect() string    { return "fn" }
func (*Builtin) Inspect() string       { return "builtin fn" }

func (v *String) Equal(obj Object) bool {
	if objInt, ok := obj.(*String); ok {
		return objInt.Value == v.Value
	}
	return false
}
func (v *Float) Equal(obj Object) bool {
	if objInt, ok := obj.(*Float); ok {
		return objInt.Value == v.Value
	}
	return false
}
func (v *Boolean) Equal(obj Object) bool {
	if objInt, ok := obj.(*Boolean); ok {
		return objInt.Value == v.Value
	}
	return false
}
func (v *Error) Equal(obj Object) bool {
	if objInt, ok := obj.(*Error); ok {
		return objInt.Message == v.Message
	}
	return false
}
func (v *Null) Equal(obj Object) bool        { return true }
func (v *ReturnValue) Equal(obj Object) bool { return false }
func (v *Function) Equal(obj Object) bool    { return false }
func (*Builtin) Equal(obj Object) bool       { return false }

func ToBoolean(obj Object) bool {
	switch obj := obj.(type) {
	case *String:
		return len(obj.Value) > 0
	case *Boolean:
		return obj.Value
	case *Int:
		return obj.Value > 0
	case *Float:
		return obj.Value > 0
	case *Null:
		return false
	}
	return false
}

func ToRawValue(obj Object) string {
	switch obj := obj.(type) {
	case *String:
		return obj.Value
	case *Boolean:
		return fmt.Sprint(obj.Value)
	case *Int:
		return fmt.Sprint(obj.Value)
	case *Float:
		return fmt.Sprint(obj.Value)
	}
	return ""
}
