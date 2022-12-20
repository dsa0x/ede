package evaluator

import (
	"ede/object"
	"fmt"
)

func (e *Evaluator) evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch true {
	case left.Type() == object.INT_OBJ && right.Type() == object.INT_OBJ:
		left := left.(*object.Int)
		right := right.(*object.Int)
		return e.evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		left := left.(*object.Float)
		right := right.(*object.Float)
		return e.evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INT_OBJ:
		left := left.(*object.Float)
		right := right.(*object.Int)
		rightFloat := &object.Float{Value: float64(right.Value)}
		return e.evalFloatInfixExpression(operator, left, rightFloat)
	case left.Type() == object.INT_OBJ && right.Type() == object.FLOAT_OBJ:
		left := left.(*object.Int)
		leftFloat := &object.Float{Value: float64(left.Value)}
		right := right.(*object.Float)
		return e.evalFloatInfixExpression(operator, leftFloat, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		left := left.(*object.String)
		right := right.(*object.String)
		return e.evalStringInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		left := left.(*object.Boolean)
		right := right.(*object.Boolean)
		return e.evalBoolInfixExpression(operator, left, right)
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid infix operator %s for (%s) and (%s)", operator, left.Inspect(), right.Inspect()))
}

func (e *Evaluator) evalIntegerInfixExpression(operator string, left, right *object.Int) object.Object {
	switch operator {
	case "+":
		return &object.Int{Value: left.Value + right.Value}
	case "-":
		return &object.Int{Value: left.Value - right.Value}
	case "*":
		return &object.Int{Value: left.Value * right.Value}
	case "/":
		return &object.Int{Value: left.Value / right.Value}
	case "%":
		return &object.Int{Value: left.Value % right.Value}
	case ">":
		return e.booleanObj(left.Value > right.Value)
	case "<":
		return e.booleanObj(left.Value < right.Value)
	case "==":
		return e.booleanObj(left.Value == right.Value)
	case "!=":
		return e.booleanObj(left.Value != right.Value)
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid integer operator %s", operator))
}

func (e *Evaluator) evalFloatInfixExpression(operator string, left, right *object.Float) object.Object {
	switch operator {
	case "+":
		return &object.Float{Value: left.Value + right.Value}
	case "-":
		return &object.Float{Value: left.Value - right.Value}
	case "*":
		return &object.Float{Value: left.Value * right.Value}
	case "/":
		return &object.Float{Value: left.Value / right.Value}
	case ">":
		return e.booleanObj(left.Value > right.Value)
	case "<":
		return e.booleanObj(left.Value < right.Value)
	case "==":
		return e.booleanObj(left.Value == right.Value)
	case "!=":
		return e.booleanObj(left.Value != right.Value)
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid integer operator %s", operator))
}

func (e *Evaluator) evalStringInfixExpression(operator string, left, right *object.String) object.Object {
	switch operator {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid string operator %s", operator))
}
func (e *Evaluator) evalBoolInfixExpression(operator string, left, right *object.Boolean) object.Object {
	switch operator {
	case "&&":
		return e.booleanObj(left.Value && right.Value)
	case "||":
		return e.booleanObj(left.Value || right.Value)
	case "==":
		return e.booleanObj(left.Value == right.Value)
	case "!=":
		return e.booleanObj(left.Value != right.Value)
	}
	return object.NewErrorWithMsg(fmt.Sprintf("invalid boolean operator %s", operator))
}
