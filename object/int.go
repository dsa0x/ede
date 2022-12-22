package object

import "fmt"

type Int struct{ Value int64 }

func (*Int) Type() Type        { return INT_OBJ }
func (v *Int) Inspect() string { return fmt.Sprint(v.Value) }
func (v *Int) Equal(obj Object) bool {
	if objInt, ok := obj.(*Int); ok {
		return objInt.Value == v.Value
	}
	return false
}

func (a *Int) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
	case "float":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) > 0 {
					return countArgumentError("0", len(args))
				}

				return &Float{Value: float64(a.Value)}
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
