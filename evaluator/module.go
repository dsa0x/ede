package evaluator

import (
	"ede/module"
	"ede/object"
)

var modules = map[string]object.Module{}

func InitModules(eval object.Evaluator, env *object.Environment) {
	modules["json"] = &module.JSONModule{}
	modules["time"] = &module.TimeModule{}

	for _, mod := range modules {
		mod.Init(eval, env)
	}
}
