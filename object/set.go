package object

import (
	"bytes"
	"strings"
)

type Set struct{ Entries map[HashKey]struct{} }

func (*Set) Type() Type { return SET_OBJ }
func (v *Set) Inspect() string {
	buf := new(bytes.Buffer)
	buf.WriteString("(\n")
	entries := []string{}
	for key := range v.Entries {
		entries = append(entries, key.Value)
	}
	buf.WriteString(strings.Join(entries, ","))
	buf.WriteString(")")
	return buf.String()
}

func (a *Set) Native() any {
	set := make(map[string]struct{})
	// for key, el := range a.Entries {
	// 	set[]
	// }
	return set
}

func (v *Set) Equal(obj Object) bool {
	if obj, ok := obj.(*Set); ok {
		if len(obj.Entries) != len(v.Entries) {
			return false
		}
		for el := range obj.Entries {
			if _, found := v.Entries[el]; !found {
				return false
			}
		}
		return true
	}
	return false
}

func (a *Set) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
	case "add":
		return a.Add()
	case "delete":
		return a.Delete()
	case "contains":
		return a.Contains()
	case "items":
		return a.Items()
	case "length":
		return &Builtin{
			Fn: func(args ...Object) Object {
				return &Int{Value: int64(len(a.Entries))}
			}}
	case "clear":
		return a.Clear()
	}
	return nil
}

func (a *Set) Add() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return CountArgumentError(">0", len(args))
			}
			for _, arg := range args {
				key := ToHashKey(arg)
				if key == EmptyHashKey {
					return NewErrorWithMsg("cannot add non-hashable %v item to set", arg.Inspect())
				}
				a.Entries[key] = struct{}{}
			}
			return a
		},
	}
}

func (a *Set) Delete() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return CountArgumentError(">0", len(args))
			}
			keys := []HashKey{}
			for _, arg := range args {
				key := ToHashKey(arg)
				if key == EmptyHashKey {
					return NewErrorWithMsg("cannot delete non-hashable %v item to set", arg.Inspect())
				}
				keys = append(keys, key)
			}
			for _, key := range keys {
				delete(a.Entries, key)
			}
			return a
		},
	}
}

func (a *Set) Contains() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return CountArgumentError("1", len(args))
			}

			key, ok := args[0].(Hashable)
			if !ok {
				return FALSE
			}
			if _, ok := a.Entries[key.HashKey()]; ok {
				return TRUE
			}
			return FALSE
		},
	}
}

func (a *Set) Items() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) > 0 {
				return CountArgumentError("0", len(args))
			}
			entries := make([]Object, 0)
			for el := range a.Entries {
				entries = append(entries, FromHashKey(el))
			}
			return &Array{Entries: &entries}
		},
	}
}

func (a *Set) Clear() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) > 0 {
				return CountArgumentError("0", len(args))
			}
			for key := range a.Entries {
				delete(a.Entries, key)
			}
			return a
		},
	}
}
