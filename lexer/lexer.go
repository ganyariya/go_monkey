package lexer

import "github.com/ganyariya/go_monkey/token"

type Lexer struct {
	input        string
	position     int  // 常に、現在 ch に入っている文字の位置を指す
	readPosition int  // 常に、これから読もうとしている次の文字の位置を指す
	ch           byte // 現在検査中の文字
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = token.NewToken(token.ASSIGN, '=')
	case '+':
		tok = token.NewToken(token.PLUS, '+')
	case '(':
		tok = token.NewToken(token.LPAREN, '(')
	case ')':
		tok = token.NewToken(token.RPAREN, ')')
	case '{':
		tok = token.NewToken(token.LBRACE, '{')
	case '}':
		tok = token.NewToken(token.RBRACE, '}')
	case ',':
		tok = token.NewToken(token.COMMA, ',')
	case ';':
		tok = token.NewToken(token.SEMICOLON, ';')
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default: // 識別子・キーワード・数について

		// 変数 or Keyword
		if isIdentifierLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			// リテラルから「変数」か「Keyword」か調べる
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else {
			// 失敗したら ILLEGAL トークンを埋め込むことで、テストなどでエラーを発見しやすくする
			tok = token.NewToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// 変数の識別子を取得する
func (l *Lexer) readIdentifier() string {
	p := l.position
	for isIdentifierLetter(l.ch) {
		l.readChar()
	}
	return l.input[p:l.position]
}

// 次の一文字を読む & 現在位置を進める
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII "NUL" に対応している
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// 変数の識別子として利用できる文字
func isIdentifierLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
