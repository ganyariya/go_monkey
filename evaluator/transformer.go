package evaluator

import (
	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/object"
)

func evalStatements(stmts []ast.Statement) object.Object {
	var ret object.Object
	for _, stmt := range stmts {
		ret = Eval(stmt)
	}
	return ret
}

// -----------------------------------------------------------
// -----------------------------------------------------------

func evalBooleanExpression(exp *ast.BooleanExpression) object.Object {
	if exp.Value {
		return TRUE
	} else {
		return FALSE
	}
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

func evalBangPrefixOperator(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		switch right := right.(type) {
		case *object.Integer:
			if right.Value == 0 {
				return TRUE
			} else {
				return FALSE
			}
		default:
			return FALSE
		}
	}
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
