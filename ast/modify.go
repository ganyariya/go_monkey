package ast

type ModifierFunc func(Node) Node

/*
すべての葉 Node に対して ModifierFunc を適応する
*/
func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.ExpressionValue, _ = Modify(node.ExpressionValue, modifier).(Expression)
	case *BlockStatement:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)
	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
	case *FunctionExpression:
		for i, param := range node.Parameters {
			node.Parameters[i], _ = Modify(param, modifier).(*IdentifierExpression)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteralExpression:
		for i, el := range node.Elements {
			node.Elements[i], _ = Modify(el, modifier).(Expression)
		}
	case *HashLiteralExpression:
		pairs := make(map[Expression]Expression)
		for k, v := range node.Pairs {
			k, _ := Modify(k, modifier).(Expression)
			v, _ := Modify(v, modifier).(Expression)
			pairs[k] = v
		}
		node.Pairs = pairs
	}
	return modifier(node)
}
