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
