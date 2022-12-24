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
	IMPORT_OBJ       Type = "IMPORT"

	NULL  = &Null{}
	TRUE  = NewBoolean(true)
	FALSE = NewBoolean(false)

	EmptyHashKey = HashKey{}
)

type Object interface {
	Type() Type
	Native() any
	Inspect() string
	Equal(obj Object) bool
}

func New(val any) Object {
	switch val := val.(type) {
	case bool:
		return NewBoolean(val)
	case string:
		return NewString(val)
	case float64:
		return NewFloat(val)
	case float32:
		return NewFloat(float64(val))
	case int64:
		return NewInt(val)
	case int:
		return NewInt(int64(val))
	case []any:
		return NewArray(val)
	case map[string]any:
		return NewHash(val)
	}
	return &Error{Message: "unsupported value"}
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
func (*Builtin) Type() Type     { return BUILTIN_OBJ }

func (v *Error) Inspect() string       { return fmt.Sprint(v.Message) }
func (v *Null) Inspect() string        { return "null" }
func (v *ReturnValue) Inspect() string { return v.Value.Inspect() }
func (*Builtin) Inspect() string       { return "builtin fn" }

func (v *Error) Equal(obj Object) bool {
	if objInt, ok := obj.(*Error); ok {
		return objInt.Message == v.Message
	}
	return false
}
func (v *Null) Equal(obj Object) bool        { return true }
func (v *ReturnValue) Equal(obj Object) bool { return false }
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
		return NewBoolean(boolVal)
	}
	return nil
}

func invalidKeyError(key string) *Error {
	return &Error{Message: fmt.Sprintf("invalid key '%s'", key)}
}
