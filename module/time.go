package module

import (
	"ede/object"
	"time"
)

type TimeModule struct {
	functions   map[string]*object.Builtin
	environment *object.Environment
	evaluator   object.Evaluator
}

func (j *TimeModule) Name() string { return "time" }

func (j *TimeModule) Functions() map[string]*object.Builtin { return j.functions }

func (j *TimeModule) Init(evaluator object.Evaluator, env *object.Environment) {
	j.evaluator = evaluator
	j.environment = env
	j.functions = map[string]*object.Builtin{
		"parse": j.Parse(),
		"now":   j.Now(),
	}
}

func (j *TimeModule) Parse() *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return object.CountArgumentError("2", len(args))
			}

			str, ok := args[0].(*object.String)
			if !ok {
				return object.NewErrorWithMsg("expected time.parse to receive argument of type 'String', got %T", args[0])
			}

			format, ok := args[1].(*object.String)
			if !ok {
				return object.NewErrorWithMsg("expected time.parse to receive argument of type 'String', got %T", args[0])
			}

			t, err := time.Parse(format.Value, str.Value)
			if err != nil {
				return object.NewErrorWithMsg("error parsing time: got %s", err)
			}
			return object.NewTime(t, format.Value)
		},
	}
}

func (j *TimeModule) Now() *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 1 {
				return object.CountArgumentError(">1", len(args))
			}

			now := time.Now()
			format := "2006-01-02 15:04:05"

			if len(args) == 1 {
				obj, ok := args[0].(*object.Hash)
				if !ok {
					return object.NewErrorWithMsg("expected time.now to receive argument of type 'Hash', got %T", args[0])
				}
				formatObj := obj.Entries["format"]
				if formatObj != nil {
					formatObj, ok := formatObj.(*object.String)
					if !ok {
						return object.NewErrorWithMsg("expected time format to be of type 'String', got %T", format)
					}
					format = formatObj.Inspect()
				}
			}

			return object.NewTime(now, format)
		},
	}
}
