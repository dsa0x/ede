package module

import (
	"ede/object"
	"encoding/json"
)

type JSONModule struct {
	functions   map[string]*object.Builtin
	environment *object.Environment
	evaluator   object.Evaluator
}

func (j *JSONModule) Name() string { return "json" }

func (j *JSONModule) Functions() map[string]*object.Builtin { return j.functions }

func (j *JSONModule) Init(evaluator object.Evaluator, env *object.Environment) {
	j.evaluator = evaluator
	j.environment = env
	j.functions = map[string]*object.Builtin{
		"parse":  j.Parse(),
		"string": j.String(),
	}
}

func (j *JSONModule) Parse() *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.CountArgumentError("1", len(args))
			}

			key, ok := args[0].(*object.String)
			if !ok {
				return object.NewErrorWithMsg("expected json.parse to receive argument of type 'String', got %T", args[0])
			}
			var obj map[string]any
			if err := json.Unmarshal([]byte(key.Value), &obj); err != nil {
				return object.NewErrorWithMsg("error parsing string as json: %s", err)
			}
			hash := &object.Hash{Entries: make(map[string]object.Object)}
			for k, v := range obj {
				val := object.New(v)
				if val.Type() == object.ERROR_OBJ {
					return val
				}
				hash.Entries[k] = val
			}
			return hash
		},
	}
}

func (j *JSONModule) String() *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.CountArgumentError("1", len(args))
			}

			obj, ok := args[0].(*object.Hash)
			if !ok {
				return object.NewErrorWithMsg("expected json.string to receive argument of type 'Hash', got %T", args[0])
			}
			res, err := json.Marshal(obj.Native())
			if err != nil {
				return object.NewErrorWithMsg("error parsing string as json: %s", err)
			}

			return object.NewString(string(res))
		},
	}
}
