package evaluator

import (
	"fmt"

	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/object"
	"github.com/ganyariya/go_monkey/token"
)

/*
「評価」せずに ASTNode のまま返す
*/
func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquote(node, env)
	return &object.Quote{Node: node}
}

/*
	AST node ノードの子孫すべてで func を実行し Modify(変更) する
*/
func evalUnquote(node ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(node, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		// unquote (引数=1)
		call := node.(*ast.CallExpression)
		if len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToASTNode(unquoted)
	})
}

/*
unquote 関数か調べる
*/
func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExpression.Function.TokenLiteral() == "unquote"
}

/*
Modify(unquote) で得られる変換された object を ast.Node へさらに変換する
*/
func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteralExpression{Token: t, Value: obj.Value}
	default:
		return nil
	}
}
