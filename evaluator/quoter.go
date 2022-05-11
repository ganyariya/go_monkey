package evaluator

import (
	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/object"
)

/*
「評価」せずに ASTNode をそのまま返す
*/
func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
