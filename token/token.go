package token

type TokenType string

/*
「タイプ」を識別するための値（右側の値）
大事なのは左側の「トークンタイプ」で右に何が入っているかは重要ではない
lexer などで必ず `token` package を引いて
`token.ILLEGAL` のように使うので `ILLEGAL = illEGAL` のようにしても問題はない（なんなら ILLEGAL = 43 でもいい）
*/
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子 & リテラル
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"
	STRING     = "STRING"

	// 演算子
	ASSIGN   = "ASSIGN"
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	BANG     = "BANG"
	ASTERISK = "ASTERISK"
	SLASH    = "SLASH"

	LT = "LT"
	GT = "GT"

	EQ     = "EQ"
	NOT_EQ = "NOT_EQ"

	// デリミタ
	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"
	COLON     = "COLON"

	LPAREN   = "LPAREN"
	RPAREN   = "RPAREN"
	LBRACE   = "LBRACE"
	RBRACE   = "RBRACE"
	LBRACKET = "LBRACKET"
	RBRACKET = "RBRACKET"

	// Keyword
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	MACRO = "MACRO"
)

/*
{ Type = "ASSIGN", Literal = "=" }
{ Type = "IDENT", Literal = "xyZ" } // Literal はソースコードに書かれている実際の値が入る
*/
type Token struct {
	Type    TokenType // const の右辺値が入る
	Literal string    // ソースコードにおける実際の値が入る
}

func NewToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"macro":  MACRO,
}

// リテラルの値からその値がキーワードか調べて「タイプ」を返す
// `IDENT` が帰ってきたら「変数の識別子」
func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
