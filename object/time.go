package object

import (
	"fmt"
	"time"
)

type Time struct {
	Value  time.Time
	format string
}

var _ Object = (*Time)(nil)
var DefaultFormat = "2006-01-02 15:04:05"

func (*Time) Type() Type        { return TIME_OBJ }
func (v *Time) Inspect() string { return v.Value.Format(v.format) }
func (v *Time) Equal(obj Object) bool {
	if objTime, ok := obj.(*Time); ok {
		return objTime.Value == v.Value
	}
	return false
}

func NewTime(val time.Time, format string) *Time {
	if format == "" {
		format = DefaultFormat
	}
	return &Time{Value: val, format: format}
}

func (a *Time) Native() any {
	return a.Value
}

func (a *Time) GetMethod(name string, eval Evaluator) *Builtin {
	switch name {
	case "string":
		return a.String(eval)
	case "sub":
		return a.Sub(eval)
	}
	return nil
}

func (v *Time) HashKey() HashKey {
	return HashKey{Type: v.Type(), Value: fmt.Sprint(v.Value)}
}

func (a *Time) String(eval Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			format := "2006-01-02 15:04:05"
			if len(args) > 1 {
				if args[0].Type() != STRING_OBJ {
					return TypeError(STRING_OBJ, args[0].Type())
				}
				format = args[0].(*String).Value
			}

			return NewString(a.Value.Format(format))
		},
	}
}

func (a *Time) Sub(eval Evaluator) *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return CountArgumentError(">=1", len(args))
			}

			if args[0].Type() != TIME_OBJ {
				return TypeError(TIME_OBJ, args[0].Type())
			}
			sub := a.Value.Sub(args[0].(*Time).Value)

			if len(args) > 1 {
				if args[1].Type() != HASH_OBJ {
					return TypeError(HASH_OBJ, args[1].Type())
				}

				opts := args[1].(*Hash).Entries
				if opts["unit"] != nil {
					if opts["unit"].Type() != STRING_OBJ {
						return TypeError(STRING_OBJ, opts["unit"].Type())
					}
					unit := opts["unit"].(*String).Value
					return formatDuration(sub, unit)

				}
			}
			return NewFloat(sub.Seconds())
		},
	}
}

func formatDuration(dur time.Duration, unit string) Object {

	switch unit {
	case "ns":
		return NewFloat(float64(dur.Nanoseconds()))
	case "us":
		return NewFloat(float64(dur.Microseconds()))
	case "ms":
		return NewFloat(float64(dur.Milliseconds()))
	case "s":
		return NewFloat(dur.Seconds())
	case "m":
		return NewFloat(dur.Minutes())
	case "h":
		return NewFloat(dur.Hours())
	case "d":
		return NewFloat(dur.Hours() / 24)
	case "w":
		return NewFloat(dur.Hours() / 24 / 7)
	case "M":
		return NewFloat(dur.Hours() / 24 / 30)
	case "y":
		return NewFloat(dur.Hours() / 24 / 365)
	}
	return NewFloat(dur.Seconds())
}
