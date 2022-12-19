package evaluator

import (
	"ede/object"
	"fmt"
)

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch true {
	// bang operator for all types
	case right.Type() != object.ERROR_OBJ && operator == "!":
		return evalBangOperator(operator, right)
	case right.Type() == object.INT_OBJ:
		right := right.(*object.Int)
		return evalIntegerPrefixExpression(operator, right)
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid prefix operator %s for %s", operator, right.Inspect()))
}

func evalIntegerPrefixExpression(operator string, right *object.Int) object.Object {
	switch operator {
	case "-":
		return &object.Int{Value: -right.Value}
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid integer operator %s", operator))
}
