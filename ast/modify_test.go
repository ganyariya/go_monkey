package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression { return &IntegerLiteralExpression{Value: 1} }
	two := func() Expression { return &IntegerLiteralExpression{Value: 2} }
	// 条件の一致するノードだけ動的に変更する
	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteralExpression)
		if !ok {
			return node
		}
		if integer.Value != 1 {
			return node
		}
		integer.Value = 2
		return integer
	}
	tests := []struct {
		input    Node
		expected Node
	}{
		{one(), two()},
		{
			&Program{Statements: []Statement{&ExpressionStatement{ExpressionValue: one()}}},
			&Program{Statements: []Statement{&ExpressionStatement{ExpressionValue: two()}}},
		},
		{
			&InfixExpression{Left: one(), Operator: "+", Right: two()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&PrefixExpression{Operator: "-", Right: one()},
			&PrefixExpression{Operator: "-", Right: two()},
		},
		{
			&IndexExpression{Index: one(), Left: one()},
			&IndexExpression{Index: two(), Left: two()},
		},
		{
			&IfExpression{
				Condition: one(),
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{ExpressionValue: one()},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{ExpressionValue: one()},
					},
				},
			},
			&IfExpression{
				Condition: two(),
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{ExpressionValue: two()},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{ExpressionValue: two()},
					},
				},
			},
		},
		{
			&ReturnStatement{ReturnValue: one()},
			&ReturnStatement{ReturnValue: two()},
		},
		{
			&LetStatement{Value: one()},
			&LetStatement{Value: two()},
		},
		{
			&FunctionExpression{
				Parameters: []*IdentifierExpression{},
				Body:       &BlockStatement{Statements: []Statement{&ExpressionStatement{ExpressionValue: one()}}},
			},
			&FunctionExpression{
				Parameters: []*IdentifierExpression{},
				Body:       &BlockStatement{Statements: []Statement{&ExpressionStatement{ExpressionValue: two()}}},
			},
		},
		{
			&ArrayLiteralExpression{Elements: []Expression{one(), one()}},
			&ArrayLiteralExpression{Elements: []Expression{two(), two()}},
		},
	}
	for _, tt := range tests {
		modified := Modify(tt.input, turnOneIntoTwo)
		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal. got=%#v, want=%#v", modified, tt.expected)
		}
	}

	hashLiteral := &HashLiteralExpression{
		Pairs: map[Expression]Expression{
			one(): one(),
			one(): one(),
		},
	}
	Modify(hashLiteral, turnOneIntoTwo)

	for key, val := range hashLiteral.Pairs {
		key, _ := key.(*IntegerLiteralExpression)
		if key.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, key.Value)
		}
		val, _ := val.(*IntegerLiteralExpression)
		if val.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, val.Value)
		}
	}
}
