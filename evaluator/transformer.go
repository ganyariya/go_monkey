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

// true / false 再利用
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

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
