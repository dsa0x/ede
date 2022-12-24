package evaluator

import (
	"ede/module"
	"ede/object"
)

var modules = map[string]object.Module{}

func init() {
	modules["json"] = &module.JSONModule{}

	for _, mod := range modules {
		mod.Init()
	}
}
