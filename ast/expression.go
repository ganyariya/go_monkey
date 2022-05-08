package ast

import (
	"bytes"
	"fmt"
	"strings"

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

type BooleanExpression struct {
	Token token.Token // token.TRUE or FALSE
	Value bool
}

func (b *BooleanExpression) expressionNode()      {}
func (b *BooleanExpression) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanExpression) String() string       { return b.Token.Literal }

type StringLiteralExpression struct {
	Token token.Token // token.STRING
	Value string
}

func (s *StringLiteralExpression) expressionNode()      {}
func (s *StringLiteralExpression) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteralExpression) String() string       { return s.Token.Literal }

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

/*
if (condition) Consequence else Alternative
Monkey で if は値を返す式だが、BlockStatement のそれぞれをIfExpression式の中に含む
（式の中に「文」が複数含まれることがあり、最後に評価された式の値を if 式は返す）
*/
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type FunctionExpression struct {
	Token      token.Token // token.FUNCTION
	Parameters []*IdentifierExpression
	Body       *BlockStatement
}

func (fe *FunctionExpression) expressionNode()      {}
func (fe *FunctionExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *FunctionExpression) String() string {
	params := []string{}
	for _, p := range fe.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("%s(%s)%s", fe.TokenLiteral(), strings.Join(params, ", "), fe.Body.String())
}

/*
`add``(2,3)`
`fn(x, y){x+y}``(2,3)`
*/
type CallExpression struct {
	Token     token.Token
	Function  Expression // IdentifierExpression or FunctionExpression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	return fmt.Sprintf("%s(%s)", ce.Function.String(), strings.Join(args, ", "))
}

type ArrayLiteralExpression struct {
	Token    token.Token // token.LBRACKET
	Elements []Expression
}

func (a *ArrayLiteralExpression) expressionNode()      {}
func (a *ArrayLiteralExpression) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteralExpression) String() string {
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

/*
Infix Type
Left に「Array」や「Arrayを指す識別子」「Arrayを返す即時関数」を含む
*/
type IndexExpression struct {
	Token token.Token // token.LBRACKET
	/* 任意の式（ただし構文解析＋評価された結果　配列となる必要がある）*/
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode()      {}
func (i *IndexExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", i.Left.String(), i.Index.String())
}

type HashLiteralExpression struct {
	Token token.Token // token.LBRACE
	Pairs map[Expression]Expression
}

func (h *HashLiteralExpression) expressionNode()      {}
func (h *HashLiteralExpression) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteralExpression) String() string {
	pairs := []string{}
	for k, v := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s:%s", k.String(), v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}
