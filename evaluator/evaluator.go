package evaluator

import (
	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/object"
)

/*
AST Node を再帰的に評価して Object System の Object に変換する
*/
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.ExpressionValue)
	case *ast.IntegerLiteralExpression:
		return evalIntegerLiteralExpression(node)
	case *ast.BooleanExpression:
		return evalBooleanExpression(node)
	}
	return nil
}
