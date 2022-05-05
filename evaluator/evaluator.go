package evaluator

import (
	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/object"
)

/*
AST Node を再帰的に評価して Object System の Object に変換する
*/
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.ExpressionValue, env)
	case *ast.IntegerLiteralExpression:
		return evalIntegerLiteralExpression(node)
	case *ast.BooleanExpression:
		return evalBooleanExpression(node)
	case *ast.IdentifierExpression:
		return evalIdentifierExpression(node, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.FunctionExpression:
		return evalFunctionExpression(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	}
	return nil
}
