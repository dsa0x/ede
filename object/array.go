package object

// import (
// 	"bytes"
// 	"fmt"
// 	"strings"
// )

// type Array[T any] struct{ Entries *[]Object }

// func (*Array[any]) Type() Type { return ARRAY_OBJ }
// func (v *Array[any]) Inspect() string {
// 	buf := new(bytes.Buffer)
// 	buf.WriteString("[")
// 	entries := []string{}
// 	for _, el := range *v.Entries {
// 		entries = append(entries, el.Inspect())
// 	}
// 	buf.WriteString(strings.Join(entries, ", "))
// 	buf.WriteString("]")
// 	return buf.String()
// }

// func (a Array[T]) GetMethod(name string) *Builtin {
// 	switch name {
// 	case "map":
// 		return a.Map()
// 	case "push":
// 		return a.Push()
// 	}
// 	return nil
// }

// func (a Array[T]) Map() *Builtin {
// 	return &Builtin{
// 		Fn: func(args ...Object) Object {
// 			if len(args) != 1 {
// 				return countArgumentError(1, len(args))
// 			}
// 			fn, ok := args[0].(*Function)
// 			if !ok {
// 				return NewError(fmt.Errorf("method 'map' expects a function argument, got %T", fn))
// 			}

// 			// for _, el := range *a.Entries {
// 			// 	eval(fn, el)
// 			// }
// 			return fn
// 		},
// 	}
// }

// func (a *Array[T]) Push() *Builtin {
// 	return &Builtin{
// 		Fn: func(args ...Object) Object {
// 			*a.Entries = append(*a.Entries, args...)
// 			return a
// 		},
// 	}
// }

// func countArgumentError(exp, got int) *Error {
// 	return NewError(fmt.Errorf("expected %d arguments, got %d", exp, got))
// }
