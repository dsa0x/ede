package object

import (
	"bytes"
	"ede/ast"
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Array struct{ Entries *[]Object }
type Evaluator interface {
	Eval(node ast.Node, env *Environment) Object
}

func (*Array) Type() Type { return ARRAY_OBJ }
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
func (v *Array) Equal(obj Object) bool {
	if obj, ok := obj.(*Array); ok {
		for idx, o := range *v.Entries {
			if !o.Equal((*obj.Entries)[idx]) {
				return false
			}
		}
		return true
	}
	return false
}

func (a *Array) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
	case "push":
		return a.Push()
	case "pop":
		return a.Pop()
	case "reverse":
		return a.Reverse()
	case "map":
		return a.Map(eval)
	case "merge":
		return a.Merge()
	case "filter":
		return a.Filter(eval)
	case "contains":
		return a.Contains(eval)
	}
	return nil
}

func (a *Array) Push() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			*a.Entries = append(*a.Entries, args...)
			return a
		},
	}
}

func (a *Array) Pop() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(*a.Entries) == 0 {
				return a
			}

			*a.Entries = (*a.Entries)[0 : len(*a.Entries)-1]
			return a
		},
	}
}

func (a *Array) Reverse() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(*a.Entries) == 0 {
				return a
			}
			*a.Entries = lo.Reverse(*a.Entries)
			return a
		},
	}
}
func (a *Array) Contains(eval Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			var found bool
			for _, entry := range *a.Entries {
				if entry.Equal(args[0]) {
					found = true
					break
				}
			}
			return &Boolean{Value: found}
		},
	}
}

func (a *Array) Map(eval Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			fn, ok := args[0].(*Function)
			if !ok {
				return NewError(fmt.Errorf("method 'map' expects a function argument, got %T", fn))
			}

			if len(fn.Params) != 1 {
				NewError(fmt.Errorf("function should have 1 argument, got %d", len(fn.Params)))
			}

			arrs := make([]Object, 0)
			result := &Array{Entries: &arrs}
			for _, el := range *a.Entries {
				env := NewEnvironment(fn.ParentEnv)
				env.Set(fn.Params[0].Value, el)
				*result.Entries = append(*result.Entries, eval.Eval(fn.Body, env))
			}
			*a.Entries = *result.Entries
			return a
		},
	}
}

func (a *Array) Filter(eval Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			fn, ok := args[0].(*Function)
			if !ok {
				return NewError(fmt.Errorf("method 'filter' expects a function argument, got %T", fn))
			}

			if len(fn.Params) != 1 {
				NewError(fmt.Errorf("function should have 1 argument, got %d", len(fn.Params)))
			}

			arrs := make([]Object, 0)
			result := &Array{Entries: &arrs}
			for _, el := range *a.Entries {
				env := NewEnvironment(fn.ParentEnv)
				env.Set(fn.Params[0].Value, el)
				obj := eval.Eval(fn.Body, env)
				if boolVal := ToBoolean(obj); boolVal {
					*result.Entries = append(*result.Entries, el)
				}
			}
			*a.Entries = *result.Entries
			return a
		},
	}
}

func (a *Array) Merge() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return countArgumentError(">1", len(args))
			}

			for _, arg := range args {
				fn, ok := arg.(*Array)
				if !ok {
					return NewError(fmt.Errorf("method 'merge' expects an array argument, got %T", fn))
				}
				*a.Entries = append(*a.Entries, *fn.Entries...)
			}
			return a
		},
	}
}

func countArgumentError(exp string, got int) *Error {
	return NewError(fmt.Errorf("expected %s argument(s), got %d", exp, got))
}
