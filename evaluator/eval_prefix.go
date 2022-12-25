package evaluator

import (
	"ede/object"
	"fmt"
)

func (e *Evaluator) evalPrefixExpression(operator string, right object.Object) object.Object {
	if e.isError(right) {
		return right
	}
	switch true {
	// bang operator for all types
	case right.Type() != object.ERROR_OBJ && operator == "!":
		return e.evalBangOperator(operator, right)
	case right.Type() == object.INT_OBJ:
		right := right.(*object.Int)
		return e.evalIntegerPrefixExpression(operator, right)
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid prefix operator %s for %s", operator, right.Inspect()))
}

func (e *Evaluator) evalIntegerPrefixExpression(operator string, right *object.Int) object.Object {
	switch operator {
	case "-":
		return &object.Int{Value: -right.Value}
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid integer operator %s", operator))
}
