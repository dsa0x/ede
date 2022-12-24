package object

func (*Function) Type() Type              { return FUNCTION_OBJ }
func (v *Function) Inspect() string       { return "func" }
func (v *Function) Equal(obj Object) bool { return false }
func (a *Function) Native() any           { return "func" }
