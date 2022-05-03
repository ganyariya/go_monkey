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
		return evalProgram(node)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements)
	case *ast.ReturnStatement:
		return evalReturnStatement(node)
	case *ast.ExpressionStatement:
		return Eval(node.ExpressionValue)
	case *ast.IntegerLiteralExpression:
		return evalIntegerLiteralExpression(node)
	case *ast.BooleanExpression:
		return evalBooleanExpression(node)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)
	case *ast.InfixExpression:
		return evalInfixExpression(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}
	return nil
}
