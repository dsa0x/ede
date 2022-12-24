package object

import "fmt"

func (a *Error) Native() any {
	return fmt.Errorf(a.Message)
}
