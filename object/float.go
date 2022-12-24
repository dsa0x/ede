package object

import "fmt"

func (*Float) Type() Type        { return FLOAT_OBJ }
func (v *Float) Inspect() string { return fmt.Sprint(v.Value) }
func (v *Float) Equal(obj Object) bool {
	if obj, ok := obj.(*Float); ok {
		return obj.Value == v.Value
	}
	// 1.0 = 1. Doing this because go automatically parses numbers as float
	if obj, ok := obj.(*Int); ok {
		return float64(obj.Value) == v.Value
	}
	return false
}

func NewFloat(val float64) *Float {
	return &Float{Value: val}
}

func (a *Float) Native() any {
	return a.Value
}

func (a *Float) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
	case "int":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) > 0 {
					return CountArgumentError("0", len(args))
				}

				return &Int{Value: int64(a.Value)}
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
