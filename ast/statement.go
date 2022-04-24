package ast

import (
	"bytes"
	"fmt"

	"github.com/ganyariya/go_monkey/token"
)

// 文
type Statement interface {
	Node
	statementNode()
}

// (let x = 5;) Statement
type LetStatement struct {
	Token token.Token           // token.LET (for トークン)
	Name  *IdentifierExpression // x (for 識別子（式）)
	Value Expression            // 5 (for 式)
}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("%s %s = ", ls.TokenLiteral(), ls.Name.String()))
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("%s ", rs.TokenLiteral()))
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token           token.Token // 式に含まれる最初のトークン
	ExpressionValue Expression
}

func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) String() string {
	if es.ExpressionValue != nil {
		return es.ExpressionValue.String()
	}
	return ""
}
