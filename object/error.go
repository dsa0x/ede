package object

import "fmt"

func (a *Error) Native() any {
	return fmt.Errorf(a.Message)
}

func NewErrorWithMsg(msg string, format ...any) *Error {
	return &Error{Message: "error: " + fmt.Sprintf(msg, format...)}
}

func NewError(msg error) *Error {
	return NewErrorWithMsg(msg.Error())
}
