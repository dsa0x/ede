package object

type Module interface {
	Name() string
	Init()
	Functions() map[string]*Builtin
}

type Import struct {
	Module Module
	Name   string
}

func NewImport(mod Module) *Import {
	return &Import{Module: mod, Name: mod.Name()}
}

func (a *Import) GetMethod(name string, eval Evaluator) *Builtin {
	return a.Module.Functions()[name]
}

func (*Import) Type() Type        { return IMPORT_OBJ }
func (v *Import) Inspect() string { return v.Module.Name() }
func (v *Import) Equal(obj Object) bool {
	if obj, ok := obj.(*Import); ok {
		return v.Module.Name() == obj.Module.Name()
	}
	return false
}

func (a *Import) Native() any {
	return a.Module.Name()
}

// func NewModule(name string, environment *Environment) *Module {
// 	return &Module{Name: name, Environment: *environment}
// }
