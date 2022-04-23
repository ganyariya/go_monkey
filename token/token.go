package token

type TokenType string

/*
大事なのは左側の「トークンタイプ」で右に何が入っているかは重要ではない
lexer などで必ず `token` package を引いて
`token.ILLEGAL` のように使うので `ILLEGAL = illEGAL` のようにしても問題はない
*/
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子 & リテラル
	IDENT = "IDENT"
	INT   = "INT"

	// 演算子
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keyword
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

/*
{ Type = "ASSIGN", Literal = "=" }
*/
type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// キーワードか調べる
// `IDENT` が帰ってきたら「変数の識別子」
func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENT
}
