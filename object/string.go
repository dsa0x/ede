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

func (a *String) Items() []Object {
	els := lo.Map(strings.Split(a.Value, ""), func(item string, i int) Object {
		return NewString(item)
	})
	return els
}

func (a *String) Native() any { return a.Value }

func (v *String) HashKey() HashKey {
	return HashKey{Type: v.Type(), Value: v.Value}
}

func (a *String) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
	case "split":
		return a.Split()
	case "reverse":
		return a.Reverse()
	case "replace":
		return a.Replace()
	case "length":
		return a.Length()
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

func (a *String) Reverse() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 0 {
				return CountArgumentError("0", len(args))
			}
			str := strings.Join(lo.Reverse(strings.Split(a.Value, "")), "")
			return NewString(str)
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

func (a *String) Length() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			return &Int{Value: int64(len(a.Value))}
		},
	}
}
