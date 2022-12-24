package object

import "fmt"

type Int struct{ Value int64 }

var _ Object = (*Int)(nil)

func (*Int) Type() Type        { return INT_OBJ }
func (v *Int) Inspect() string { return fmt.Sprint(v.Value) }
func (v *Int) Equal(obj Object) bool {
	if objInt, ok := obj.(*Int); ok {
		return objInt.Value == v.Value
	}
	// 1.0 = 1
	if objInt, ok := obj.(*Float); ok {
		return objInt.Value == float64(v.Value)
	}
	return false
}

func NewInt(val int64) *Int {
	return &Int{Value: val}
}

func (a *Int) Native() any {
	return a.Value
}

func (a *Int) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
	case "float":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) > 0 {
					return CountArgumentError("0", len(args))
				}

				return &Float{Value: float64(a.Value)}
			},
		}
	case "string":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) > 0 {
					return CountArgumentError("0", len(args))
				}

				return &String{Value: fmt.Sprint(a.Value)}
			},
		}
	}
	return nil
}

func (v *Int) HashKey() HashKey {
	return HashKey{Type: v.Type(), Value: fmt.Sprint(v.Value)}
}
