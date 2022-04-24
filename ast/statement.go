package ast

import "github.com/ganyariya/go_monkey/token"

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

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) statementNode()       {}
