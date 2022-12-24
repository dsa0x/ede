package object

import "fmt"

func (*Boolean) Type() Type        { return BOOLEAN_OBJ }
func (v *Boolean) Inspect() string { return fmt.Sprint(v.Value) }

func (v *Boolean) Equal(obj Object) bool {
	if objInt, ok := obj.(*Boolean); ok {
		return objInt.Value == v.Value
	}
	return false
}

func (v *Boolean) HashKey() HashKey {
	return HashKey{Type: v.Type(), Value: fmt.Sprint(v.Value)}
}

func NewBoolean(val bool) *Boolean {
	return &Boolean{Value: val}
}

func (a *Boolean) Native() any {
	return a.Value
}
