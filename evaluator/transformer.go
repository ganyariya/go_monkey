package evaluator

import (
	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/object"
)

func evalProgram(program *ast.Program) object.Object {
	var ret object.Object
	for _, stmt := range program.Statements {
		ret = Eval(stmt)
		if returnValue, ok := ret.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return ret
}

func evalBlockStatements(stmts []ast.Statement) object.Object {
	var ret object.Object
	for _, stmt := range stmts {
		ret = Eval(stmt)
		// BlockStatement では ReturnValue.Value にアンラップしない（ブロック文ネストでバグる)
		if ret != nil && ret.Type() == object.RETURN_VALUE_OBJ {
			return ret
		}
	}
	return ret
}

func evalReturnStatement(stmt *ast.ReturnStatement) object.Object {
	obj := Eval(stmt.ReturnValue)
	return &object.ReturnValue{Value: obj}
}

// -----------------------------------------------------------
// -----------------------------------------------------------

func nativeBoolToBooleanObject(b bool) object.Object {
	if b {
		return TRUE
	} else {
		return FALSE
	}
}
func evalBooleanExpression(exp *ast.BooleanExpression) object.Object {
	return nativeBoolToBooleanObject(exp.Value)
}

func evalIntegerLiteralExpression(exp *ast.IntegerLiteralExpression) object.Object {
	return &object.Integer{Value: exp.Value}
}

func evalPrefixExpression(exp *ast.PrefixExpression) object.Object {
	rightObj := Eval(exp.Right)
	switch exp.Operator {
	case "!":
		return evalBangPrefixOperator(rightObj)
	case "-":
		return evalMinusPrefixOperator(rightObj)
	default:
		return NULL
	}
}

func evalInfixExpression(exp *ast.InfixExpression) object.Object {
	leftObj := Eval(exp.Left)
	rightObj := Eval(exp.Right)
	switch {
	// 整数は「値」で処理する
	case leftObj.Type() == object.INTEGER_OBJ && rightObj.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(exp.Operator, leftObj, rightObj)
	// reference (pointer) （異なる型 -> false）
	case exp.Operator == "==":
		return nativeBoolToBooleanObject(leftObj == rightObj)
	case exp.Operator == "!=":
		return nativeBoolToBooleanObject(leftObj != rightObj)
	default:
		return NULL
	}
}

func evalIfExpression(exp *ast.IfExpression) object.Object {
	condition := Eval(exp.Condition)
	if condition.AsBool() {
		return Eval(exp.Consequence)
	} else if exp.Alternative != nil {
		return Eval(exp.Alternative)
	} else {
		return NULL
	}
}

// ------------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------------

func evalBangPrefixOperator(right object.Object) object.Object {
	return nativeBoolToBooleanObject(!right.AsBool())
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	default:
		return NULL
	}
}
