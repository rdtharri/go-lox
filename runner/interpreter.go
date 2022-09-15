package runner

import (
	"fmt"
)

type Interpreter struct{}

func (i *Interpreter) interpret(exp Expression) interface{} {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	return exp.Accept(i)
}

func (i *Interpreter) VisitBinaryExpression(be *BinaryExpression) interface{} {
	left := i.interpret(be.Left)
	right := i.interpret(be.Right)

	validateNum := func() (float64,float64) {
		return validateOperands[float64](
			be.Operator.Type,
			left,
			right,
		)
	}

	switch be.Operator.Type {
	case MINUS:
		leftVal, rightVal := validateNum()
		return leftVal - rightVal
	case SLASH:
		leftVal, rightVal := validateNum()
		return leftVal / rightVal
	case STAR:
		leftVal, rightVal := validateNum()
		return leftVal * rightVal
	case PLUS:
		leftString, leftOk := left.(string)
		rightString, rightOk := right.(string)
		if leftOk && rightOk {
			return leftString + rightString
		}


		leftNum, leftOk := left.(float64)
		rightNum, rightOk := right.(float64)
		if leftOk && rightOk {
			return leftNum + rightNum
		}

		panic(fmt.Errorf("invalid operands for '%v': %v, %v",be.Operator.Type,left, right))
	case GREATER:
		leftVal, rightVal := validateNum()
		return leftVal > rightVal
	case GREATER_EQUAL:
		leftVal, rightVal := validateNum()
		return leftVal >= rightVal
	case LESS:
		leftVal, rightVal := validateNum()
		return leftVal < rightVal
	case LESS_EQUAL:
		leftVal, rightVal := validateNum()
		return leftVal <= rightVal
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	}
	return nil
}

func (i *Interpreter) VisitGroupingExpression(ge *GroupingExpression) interface{} {
	return i.interpret(ge.Expression)
}

func (i *Interpreter) VisitLiteralExpression(le *LiteralExpression) interface{} {
	return le.Token.Value
}

func (i *Interpreter) VisitUnaryExpression(ue *UnaryExpression) interface{} {
	right := i.interpret(ue.Right)

	switch ue.Operator.Type {
	case MINUS:
		value, ok := right.(float64)
		if !ok {
			panic(fmt.Errorf("invalid operand for '%v': %v",MINUS,right))
		}
		return -value
	case BANG:
		return !i.isTruthy(right)
	}

	return right
}

func (i *Interpreter) isTruthy(value interface{}) bool {

	// null values false
	if value == nil {
		return false
	}

	// bools are their value
	boolVal, ok := value.(bool)
	if ok {
		return boolVal
	}

	// otherwise true
	return true
}

func (i *Interpreter) isEqual(left interface{}, right interface{}) bool {
	return left == right
}

func validateOperands[T string|float64|bool](operator TokenType, left interface{}, right interface{}) (T, T) {
	leftVal, leftOk := left.(T)
	rightVal, rightOk := right.(T)
	if !leftOk || !rightOk {
		panic(fmt.Errorf("invalid operands for '%v': %v, %v",operator,left, right))
	}
	return leftVal, rightVal
}

