package object

import (
	"bytes"
	"ede/ast"
	"ede/token"
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
	case "first":
		return a.First()
	case "last":
		return a.Last()
	case "length":
		return a.Length()
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
	case "find":
		return a.Find(eval)
	case "join":
		return a.Join(eval)
	case "clear":
		return a.Clear()
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

func (a *Array) First() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(*a.Entries) == 0 {
				return NULL
			}
			return (*a.Entries)[0]
		},
	}
}

func (a *Array) Last() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(*a.Entries) == 0 {
				return NULL
			}
			return (*a.Entries)[len(*a.Entries)-1]
		},
	}
}

func (a *Array) Length() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			return &Int{Value: int64(len(*a.Entries))}
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

func (a *Array) Contains(evaluator Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			_, found := lo.Find(*a.Entries, func(entry Object) bool {
				return entry.Equal(args[0])
			})
			return &Boolean{Value: found}
		},
	}
}

func (a *Array) Join(evaluator Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			strArg, ok := args[0].(*String)
			if !ok {
				return methodExpectArgumentError("join", "string", string(args[0].Type()))
			}
			entriesStr := make([]string, 0, len(*a.Entries))
			for idx, entry := range *a.Entries {
				str := ToRawValue(entry)
				if str == "" {
					return NewErrorWithMsg("cannot join non string-like item at index %d", idx)
				}
				entriesStr = append(entriesStr, str)
			}
			return &String{Value: strings.Join(entriesStr, strArg.Value)}
		},
	}
}

func (a *Array) Find(evaluator Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			fn, ok := args[0].(*Function)
			if !ok {
				return methodExpectArgumentError("find", "function", string(args[0].Type()))
			}
			if len(fn.Params) != 1 {
				NewError(fmt.Errorf("function should have 1 argument, got %d", len(fn.Params)))
			}

			arrs := make([]Object, 0)
			result := &Array{Entries: &arrs}
			for _, el := range *a.Entries {
				env := NewEnvironment(fn.ParentEnv)
				env.Set(fn.Params[0].Value, el)
				obj := evaluator.Eval(fn.Body, env)
				if boolVal := ToBoolean(obj); boolVal {
					*result.Entries = append(*result.Entries, el)
					return el
				}
			}
			*a.Entries = *result.Entries
			return a
		},
	}
}

func (a *Array) Map(evaluator Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			fn, ok := args[0].(*Function)
			if !ok {
				return methodExpectArgumentError("map", "function", string(args[0].Type()))
			}

			if len(fn.Params) != 1 {
				NewError(fmt.Errorf("function should have 1 argument, got %d", len(fn.Params)))
			}

			arrs := make([]Object, 0)
			result := &Array{Entries: &arrs}
			for _, el := range *a.Entries {
				env := NewEnvironment(fn.ParentEnv)
				env.Set(fn.Params[0].Value, el)
				*result.Entries = append(*result.Entries, evaluator.Eval(fn.Body, env))
			}
			*a.Entries = *result.Entries
			return a
		},
	}
}

func (a *Array) Filter(evaluator Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			fn, ok := args[0].(*Function)
			if !ok {
				return methodExpectArgumentError("filter", "function", string(args[0].Type()))
			}

			if len(fn.Params) != 1 {
				NewError(fmt.Errorf("function should have 1 argument, got %d", len(fn.Params)))
			}

			arrs := make([]Object, 0)
			result := &Array{Entries: &arrs}
			for idx, el := range *a.Entries {
				env := NewEnvironment(fn.ParentEnv)
				env.Set(token.IndexIdentifier, &Int{Value: int64(idx)})
				env.Set(fn.Params[0].Value, el)
				obj := evaluator.Eval(fn.Body, env)
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

func (a *Array) Clear() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) > 0 {
				return countArgumentError("0", len(args))
			}
			*a.Entries = (*a.Entries)[:0]
			return a
		},
	}
}

func countArgumentError(exp string, got int) *Error {
	return NewError(fmt.Errorf("expected %s argument(s), got %d", exp, got))
}

func methodExpectArgumentError(methodName, argType, gotType string) *Error {
	return NewError(fmt.Errorf("method '%s' expects a %s argument, got %s", methodName, argType, gotType))
}
