package ast

// AST のすべてのノードは Node のメソッドを実装する必要あり
type Node interface {
	TokenLiteral() string // **Token** の Literal (式ではない トークン自体のリテラル)
	String() string
}
