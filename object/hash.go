package object

import (
	"bytes"
	"fmt"
	"strings"
)

type Hash struct{ Entries map[Object]Object }

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
	}
	return nil
}

func (a *Hash) Contains(evaluator Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return countArgumentError("1", len(args))
			}

			for key := range a.Entries {
				if key.Equal(args[0]) {
					return TRUE
				}
			}
			return FALSE
		},
	}
}
