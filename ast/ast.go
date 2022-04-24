package ast

// AST のすべてのノードは Node のメソッドを実装する必要あり
type Node interface {
	TokenLiteral() string
}
