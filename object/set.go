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

func (v *Set) Equal(obj Object) bool {
	// if obj, ok := obj.(*Set); ok {
	// 	for key, self := range v.Entries {
	// 		entry, found := obj.Entries[key]
	// 		if !found {
	// 			return false
	// 		}
	// 		if !self.Equal(key) {
	// 			return false
	// 		}
	// 	}
	// 	return true
	// }
	return false
}

func (a *Set) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
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

func (a *Set) Contains() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}

			key, ok := args[0].(Hashable)
			if !ok {
				return invalidKeyError(key.Inspect())
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
				return countArgumentError("0", len(args))
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
				return countArgumentError("0", len(args))
			}
			for key := range a.Entries {
				delete(a.Entries, key)
			}
			return a
		},
	}
}
