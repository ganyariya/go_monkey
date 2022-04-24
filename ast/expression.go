package ast

import "github.com/ganyariya/go_monkey/token"

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

func (i *IdentifierExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IdentifierExpression) expressionNode()      {}
func (i *IdentifierExpression) String() string       { return i.Value }
