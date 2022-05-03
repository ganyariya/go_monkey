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
