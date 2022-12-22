package object

import (
	"fmt"
	"strconv"
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
	HASH_OBJ         Type = "HASH"
	SET_OBJ          Type = "SET"

	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}

	EmptyHashKey = HashKey{}
)

type Object interface {
	Type() Type
	Inspect() string
	Equal(obj Object) bool
}

type HashKey struct {
	Type  Type
	Value string
}
type Hashable interface {
	HashKey() HashKey
	Inspect() string
}

func (*Error) Type() Type       { return ERROR_OBJ }
func (*Null) Type() Type        { return NULL_OBJ }
func (*ReturnValue) Type() Type { return RETURN_VALUE_OBJ }
func (*Function) Type() Type    { return FUNCTION_OBJ }
func (*Builtin) Type() Type     { return BUILTIN_OBJ }

func (v *Error) Inspect() string       { return fmt.Sprint(v.Message) }
func (v *Null) Inspect() string        { return "null" }
func (v *ReturnValue) Inspect() string { return v.Value.Inspect() }
func (v *Function) Inspect() string    { return "fn" }
func (*Builtin) Inspect() string       { return "builtin fn" }

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
	if obj, ok := obj.(Hashable); ok {
		return obj.HashKey().Value
	}
	return ""
}

func ToHashKey(obj Object) HashKey {
	if obj, ok := obj.(Hashable); ok {
		return obj.HashKey()
	}
	return EmptyHashKey
}

func FromHashKey(key HashKey) Object {
	switch key.Type {
	case STRING_OBJ:
		return &String{Value: key.Value}
	case INT_OBJ:
		intVal, _ := strconv.ParseInt(key.Value, 10, 64)
		return &Int{Value: intVal}
	case BOOLEAN_OBJ:
		boolVal, _ := strconv.ParseBool(key.Value)
		return &Boolean{Value: boolVal}
	}
	return nil
}

func (v *Int) HashKey() HashKey {
	return HashKey{Type: v.Type(), Value: fmt.Sprint(v.Value)}
}
