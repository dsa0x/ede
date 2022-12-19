package evaluator

import (
	"ede/object"
	"fmt"
)

var builtins = map[string]*object.Builtin{
	"len":   {Fn: applyBuiltinLen},
	"print": {Fn: applyBuiltinPrint},
	"println": {Fn: func(args ...object.Object) object.Object {
		applyBuiltinPrint(args...)
		fmt.Println()
		return NULL
	}},
}

func applyBuiltinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewErrorWithMsg(fmt.Sprintf("builtin function 'len' requires exactly one argument, got %d", len(args)))
	}
	arg := args[0]
	switch arg := arg.(type) {
	case *object.String:
		return &object.Int{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Int{Value: int64(len(arg.Entries))}
	}
	return object.NewErrorWithMsg(fmt.Sprintf("argument to `len` not supported, got %s", arg.Type()))
}

func applyBuiltinPrint(args ...object.Object) object.Object {
	for i, arg := range args {
		if arg == nil {
			fmt.Println()
		} else if arg.Inspect() == "\\n" {
			fmt.Println()
		} else if i == len(args)-1 {
			fmt.Print(arg.Inspect())
		} else {
			fmt.Printf("%s ", arg.Inspect())
		}
	}
	return NULL
}
