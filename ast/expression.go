package ast

import (
	"fmt"

	"github.com/ganyariya/go_monkey/token"
)

// 式
type Expression interface {
	Node
	expressionNode()
}

/*
識別子   Expression の種類の一つ

本来 `let <identifier> = <expression>` において 左の identifier は「式」ではないが
他のプログラムの箇所では `let x = valueProducingIdentifier` のように 値を生成する Identifier 識別子がある（valueProducingIdentifier）
よって簡易化のために Let の Identifier についても Expression の一部としている
*/
type IdentifierExpression struct {
	Token token.Token // token.IDENTIFIER
	Value string      // xyZ などの変数名
}

func (i *IdentifierExpression) expressionNode()      {}
func (i *IdentifierExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IdentifierExpression) String() string       { return i.Value } // for Debug

// **Token 以外の値である** Value が構文解析とそのあとで「実際に使う」値っぽい（整数に変換しているため）
// Token はレキサーの時点で使うもの
type IntegerLiteralExpression struct {
	Token token.Token // token.INT
	Value int64       // 5 (Token.Literal を変換する)
}

func (i *IntegerLiteralExpression) expressionNode()      {}
func (i *IntegerLiteralExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteralExpression) String() string       { return i.Token.Literal } // for Debug

type PrefixExpression struct {
	Token    token.Token // token.MINUS, BANG
	Operator string      // 前置演算子
	Right    Expression  // 演算子の右にくる式
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

type InfixExpression struct {
	Token    token.Token // token.PLUS, ...
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}
