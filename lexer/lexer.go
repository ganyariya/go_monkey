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
	}

	l.readChar()
	return tok
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
