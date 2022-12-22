package object

import "fmt"

func (*Float) Type() Type        { return FLOAT_OBJ }
func (v *Float) Inspect() string { return fmt.Sprint(v.Value) }
func (v *Float) Equal(obj Object) bool {
	if objInt, ok := obj.(*Float); ok {
		return objInt.Value == v.Value
	}
	return false
}

func (a *Float) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
	case "int":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) > 0 {
					return countArgumentError("0", len(args))
				}

				return &Int{Value: int64(a.Value)}
			},
		}
	case "string":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) > 0 {
					return countArgumentError("0", len(args))
				}

				return &String{Value: fmt.Sprint(a.Value)}
			},
		}
	}
	return nil
}
