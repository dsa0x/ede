package evaluator

import (
	"bytes"
	"ede/object"
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Array struct{ Entries *[]object.Object }

func (*Array) Type() object.Type { return object.ARRAY_OBJ }
func (v *Array) Inspect() string {
	buf := new(bytes.Buffer)
	buf.WriteString("[")
	entries := []string{}
	for _, el := range *v.Entries {
		entries = append(entries, el.Inspect())
	}
	buf.WriteString(strings.Join(entries, ", "))
	buf.WriteString("]")
	return buf.String()
}

func (a *Array) GetMethod(name string) *object.Builtin {
	switch name {
	case "push":
		return a.Push()
	case "pop":
		return a.Pop()
	case "reverse":
		return a.Reverse()
	case "map":
		return a.Map()
	case "merge":
		return a.Merge()
	case "filter":
		return a.Filter()
	}
	return nil
}

func (a *Array) Push() *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			*a.Entries = append(*a.Entries, args...)
			return a
		},
	}
}

func (a *Array) Pop() *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(*a.Entries) == 0 {
				return a
			}

			*a.Entries = (*a.Entries)[0 : len(*a.Entries)-1]
			return a
		},
	}
}

func (a *Array) Reverse() *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(*a.Entries) == 0 {
				return a
			}
			*a.Entries = lo.Reverse(*a.Entries)
			return a
		},
	}
}

func (a *Array) Map() *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			fn, ok := args[0].(*object.Function)
			if !ok {
				return object.NewError(fmt.Errorf("method 'map' expects a function argument, got %T", fn))
			}

			if len(fn.Params) != 1 {
				object.NewError(fmt.Errorf("function should have 1 argument, got %d", len(fn.Params)))
			}

			arrs := make([]object.Object, 0)
			result := &Array{Entries: &arrs}
			for _, el := range *a.Entries {
				env := object.NewEnvironment(fn.ParentEnv)
				env.Set(fn.Params[0].Value, el)
				*result.Entries = append(*result.Entries, Eval(fn.Body, env))
			}
			*a.Entries = *result.Entries
			return a
		},
	}
}

func (a *Array) Filter() *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			fn, ok := args[0].(*object.Function)
			if !ok {
				return object.NewError(fmt.Errorf("method 'filter' expects a function argument, got %T", fn))
			}

			if len(fn.Params) != 1 {
				object.NewError(fmt.Errorf("function should have 1 argument, got %d", len(fn.Params)))
			}

			arrs := make([]object.Object, 0)
			result := &Array{Entries: &arrs}
			for _, el := range *a.Entries {
				env := object.NewEnvironment(fn.ParentEnv)
				env.Set(fn.Params[0].Value, el)
				obj := Eval(fn.Body, env)
				if boolVal := object.ToBoolean(obj); boolVal {
					*result.Entries = append(*result.Entries, el)
				}
			}
			*a.Entries = *result.Entries
			return a
		},
	}
}

func (a *Array) Merge() *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return countArgumentError(">1", len(args))
			}

			for _, arg := range args {
				fn, ok := arg.(*Array)
				if !ok {
					return object.NewError(fmt.Errorf("method 'merge' expects an array argument, got %T", fn))
				}
				*a.Entries = append(*a.Entries, *fn.Entries...)
			}
			return a
		},
	}
}

func countArgumentError(exp string, got int) *object.Error {
	return object.NewError(fmt.Errorf("expected %s arguments, got %d", exp, got))
}
