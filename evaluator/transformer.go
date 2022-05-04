package evaluator

import (
	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/object"
)

func evalProgram(program *ast.Program) object.Object {
	var ret object.Object
	for _, stmt := range program.Statements {
		ret = Eval(stmt)
		switch ret := ret.(type) {
		case *object.ReturnValue:
			return ret.Value
		case *object.Error:
			return ret
		}
	}
	return ret
}

func evalBlockStatements(stmts []ast.Statement) object.Object {
	var ret object.Object
	for _, stmt := range stmts {
		ret = Eval(stmt)
		// BlockStatement では ReturnValue.Value にアンラップしない（ブロック文ネストでバグる)
		if ret != nil {
			if ret.Type() == object.RETURN_VALUE_OBJ || ret.Type() == object.ERROR_OBJ {
				return ret
			}
		}
	}
	return ret
}

func evalReturnStatement(stmt *ast.ReturnStatement) object.Object {
	obj := Eval(stmt.ReturnValue)
	if isError(obj) {
		return obj
	}
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
	if isError(rightObj) {
		return rightObj
	}
	switch exp.Operator {
	case "!":
		return evalBangPrefixOperator(rightObj)
	case "-":
		return evalMinusPrefixOperator(rightObj)
	default:
		return newError("unknown operator: %s%s", exp.Operator, rightObj.Type())
	}
}

func evalInfixExpression(exp *ast.InfixExpression) object.Object {
	leftObj := Eval(exp.Left)
	if isError(leftObj) {
		return leftObj
	}
	rightObj := Eval(exp.Right)
	if isError(rightObj) {
		return rightObj
	}
	switch {
	// 整数は「値」で処理する
	case leftObj.Type() == object.INTEGER_OBJ && rightObj.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(exp.Operator, leftObj, rightObj)
	// reference (pointer) （異なる型 -> false）
	case exp.Operator == "==":
		return nativeBoolToBooleanObject(leftObj == rightObj)
	case exp.Operator == "!=":
		return nativeBoolToBooleanObject(leftObj != rightObj)
	case leftObj.Type() != rightObj.Type():
		return newError("type mismatch: %s %s %s", leftObj.Type(), exp.Operator, rightObj.Type())
	default:
		return newError("unknown operator: %s %s %s", leftObj.Type(), exp.Operator, rightObj.Type())
	}
}

func evalIfExpression(exp *ast.IfExpression) object.Object {
	condition := Eval(exp.Condition)
	if isError(condition) {
		return condition
	}
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
		return newError("unknown operator: -%s", right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
