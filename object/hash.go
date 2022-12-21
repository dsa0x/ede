package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Hash struct{ Entries map[string]Object }

func (*Hash) Type() Type { return ARRAY_OBJ }
func (v *Hash) Inspect() string {
	buf := new(bytes.Buffer)
	buf.WriteString("{\n")
	entries := []string{}
	for key, el := range v.Entries {
		entries = append(entries, fmt.Sprintf("%s: %s", key, el.Inspect()))
	}
	buf.WriteString(strings.Join(entries, ","))
	buf.WriteString("}")
	return buf.String()
}
func (v *Hash) Equal(obj Object) bool {
	if obj, ok := obj.(*Hash); ok {
		for key, self := range v.Entries {
			entry, found := obj.Entries[key]
			if !found {
				return false
			}
			if !self.Equal(entry) {
				return false
			}
		}
		return true
	}
	return false
}

func (a *Hash) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
	case "contains":
		return a.Contains(eval)
	case "keys":
		return a.Keys()
	case "items":
		return a.Items()
	case "clear":
		return a.Clear()
	case "get":
		return a.Get()
	}
	return nil
}

func (a *Hash) Contains(evaluator Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}

			key := ToRawValue(args[0])
			if key == "" {
				return invalidKeyError(key)
			}
			if _, ok := a.Entries[key]; ok {
				return TRUE
			}
			return FALSE
		},
	}
}

func (a *Hash) Get() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}
			key := ToRawValue(args[0])
			if key == "" {
				return invalidKeyError(key)
			}
			entry, ok := a.Entries[key]
			if !ok {
				return NULL
			}
			return entry
		},
	}
}

func (a *Hash) Keys() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) > 0 {
				return countArgumentError("0", len(args))
			}
			entries := lo.Keys(a.Entries)
			objEntries := lo.Map(entries, func(item string, i int) Object {
				return &String{Value: item}
			})
			return &Array{Entries: &objEntries}
		},
	}
}

func (a *Hash) Items() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) > 0 {
				return countArgumentError("0", len(args))
			}
			entries := lo.Values(a.Entries)
			return &Array{Entries: &entries}
		},
	}
}

func (a *Hash) Clear() *Builtin {
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

func invalidKeyError(key string) *Error {
	return &Error{Message: fmt.Sprintf("invalid key '%s'", key)}
}
