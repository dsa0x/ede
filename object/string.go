package object

import (
	"strings"

	"github.com/samber/lo"
)

func (*String) Type() Type        { return STRING_OBJ }
func (v *String) Inspect() string { return v.Value }
func (v *String) Equal(obj Object) bool {
	if objInt, ok := obj.(*String); ok {
		return objInt.Value == v.Value
	}
	return false
}

func NewString(val string) *String {
	return &String{Value: val}
}

func (a *String) Native() any { return a.Value }

func (v *String) HashKey() HashKey {
	return HashKey{Type: v.Type(), Value: v.Value}
}

func (a *String) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
	case "split":
		return a.Split()
	case "replace":
		return a.Replace()
	}
	return nil
}

func (a *String) Split() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return CountArgumentError("1", len(args))
			}
			sep, ok := args[0].(*String)
			if !ok {
				return methodExpectArgumentError("split", "String", string(args[0].Type()))
			}
			strs := strings.Split(a.Value, sep.Value)
			entries := lo.Map(strs, func(val string, i int) any { return val })
			return NewArray(entries)
		},
	}
}

func (a *String) Replace() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return CountArgumentError("2", len(args))
			}
			old, ok := args[0].(*String)
			if !ok {
				return methodExpectArgumentError("split", "String", string(args[0].Type()))
			}
			new, ok := args[1].(*String)
			if !ok {
				return methodExpectArgumentError("split", "String", string(args[0].Type()))
			}
			str := strings.Replace(a.Value, old.Value, new.Value, -1)
			return NewString(str)
		},
	}
}
