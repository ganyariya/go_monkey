package lexer

import "github.com/ganyariya/go_monkey/token"

// レキサー
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
		if l.isTwoCharToken('=', '=') {
			tok = token.Token{Type: token.EQ, Literal: l.readTwoCharToken()}
		} else {
			tok = token.NewToken(token.ASSIGN, '=')
		}
	case '!':
		if l.isTwoCharToken('!', '=') {
			tok = token.Token{Type: token.NOT_EQ, Literal: l.readTwoCharToken()}
		} else {
			tok = token.NewToken(token.BANG, '!')
		}
	case '+':
		tok = token.NewToken(token.PLUS, '+')
	case '-':
		tok = token.NewToken(token.MINUS, '-')
	case '*':
		tok = token.NewToken(token.ASTERISK, '*')
	case '/':
		tok = token.NewToken(token.SLASH, '/')
	case '<':
		tok = token.NewToken(token.LT, '<')
	case '>':
		tok = token.NewToken(token.GT, '>')
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
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			// 失敗したら ILLEGAL トークンを埋め込むことで、テストなどでエラーを発見しやすくする
			tok = token.NewToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// 次の一文字を読む & 現在位置を進める
// => l.ch が更新される
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII "NUL" に対応している
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// 次の一文字を先読みする
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// 変数の識別子を取得する
func (l *Lexer) readIdentifier() string {
	p := l.position
	for isIdentifierLetter(l.ch) {
		l.readChar()
	}
	return l.input[p:l.position]
}

func (l *Lexer) readNumber() string {
	p := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[p:l.position]
}

func (l *Lexer) isTwoCharToken(c1, c2 byte) bool {
	return l.ch == c1 && l.peekChar() == c2
}

// 2文字からなるトークンを読み込む
// 先に isTwoCharToken を呼び出してチェックする
func (l *Lexer) readTwoCharToken() string {
	ch := l.ch
	l.readChar()
	return string(ch) + string(l.ch)
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

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
