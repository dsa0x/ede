package object

func (*String) Type() Type        { return STRING_OBJ }
func (v *String) Inspect() string { return v.Value }
func (v *String) Equal(obj Object) bool {
	if objInt, ok := obj.(*String); ok {
		return objInt.Value == v.Value
	}
	return false
}
func (v *String) HashKey() HashKey {
	return HashKey{Type: v.Type(), Value: v.Value}
}
