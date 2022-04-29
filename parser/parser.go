package parser

import (
	"fmt"
	"strconv"

	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/lexer"
	"github.com/ganyariya/go_monkey/token"
)

// 順序が重要（PRODUCT は EQUALS よりも高い優先順位）
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // - or !
	CALL        // func()
)

// 中置演算子の優先順位
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

/*
2.6.6 p60

すべての構文解析関数（式を解析する関数が構文解析関数）は以下の規約に従う。
- 「構文解析関数に関連付けられたトークンが p.curToken にセットされた状態で」関数の動作を開始する。
- 「構文解析関数が処理対象とする式の一番最後のトークン」が curToken にセットされた状態で関数は終了する。
*/
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression // 中置演算子の左側の式が () に入る
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token // 今見ているトークン
	peekToken token.Token // 先読みトークン

	// トークンに対応する構文解析関数 map
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// 式を構文解析する prefixParseExpression をトークンタイプごとに登録する
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixFn(token.IDENTIFIER, p.parseIdentifierExpression)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteralExpression)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parser は与えられたソースコードをトークンごとに読み込んでパースする
// パースした結果の Statement 列を ast.Program として返す
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// 様々な Statement をパースする
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		// let return 以外は Expression のみからなる Statement
		return p.parseExpressionStatement()
	}
}

// Let Statement をパースする
func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: セミコロンにたどり着くまで読み飛ばしている（本当はここで x = Expression の「式」をパースする必要がある）
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: 式を呼び飛ばして curToken = SEMICOLON になるまで進める
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// セミコロンが来るまで「一つの大きな式」として ExpressionStatement をパースする
func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.ExpressionValue = p.parseExpression(LOWEST) // 最も低い優先順位でパースする
	// セミコロンを省略可能にする（セミコロンだったら飛ばす）
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

/*
Note: p82
parseExpression の precedence の値は右結合力を表す。
一方、今見ている演算子に対して、「次の演算子の左結合力」は peekPrecedence の値になる。

precedence の値が大きいほど（右結合力が高いほど）「文の直接1だけ右にあるトークンを Right に「直接」吸収しやすい」
かつ「自身が `他の親ノードがもつ Left 子ノード`に配置されづらい」。
`-1 + 2` で最初のマイナスを解析するとき PREFIX を parseExpression にわたす。
このとき、 `1` は「直接」マイナスの　Right に吸収される。
また、別の例として、 `*`の演算子における parseInfixExpression 内部で呼び出される parseExpression では
precedence = PRODUCT の優先順位が渡され `*` は右結合力が大きい。

precedence の値が小さいほど「これまで構文解析したものを自身の Left におき」かつ
「自分は `別のノード（将来的に親となるノード）の Left（つまり自分は`トークン列で右に出てくる`Expressionノードの子供になる）` になりやすい」
`2 + 4 + 3` の場合、 1 つ目のプラス（+1）を解釈したとき、この +1 までに解析したものは  +1 の Left に配置されやすい。
また、 +1 は他のノードの Left に配置されやすい。（今回の場合 2 つ目のプラス(+2) の Left ノードになる）
*/
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Statement に含まれる最も左側にある式を処理する
	prefixFn := p.prefixParseFns[p.curToken.Type]
	if prefixFn == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefixFn()

	// セミコロンが来る もしくは 優先順位が上がらなくなったら
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// 中間演算子の優先順位が高いなら中置演算子に紐付いた関数でパースする
		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil {
			return leftExp
		}
		p.nextToken()
		// 「これまで見ていた"中置演算子の左側にある"式」を「これから見る中間演算子式のLeft」として埋め込む
		leftExp = infixFn(leftExp)
	}

	return leftExp
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return &ast.IdentifierExpression{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteralExpression() ast.Expression {
	ile := &ast.IntegerLiteralExpression{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	ile.Value = value
	return ile
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()                        // トークンを進めて式を読む
	pe.Right = p.parseExpression(PREFIX) // PrefixExpression は強制的に前置演算子の優先順位を渡す
	return pe
}
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	ie := &ast.InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left}
	precedence := p.curPrecedence()
	p.nextToken()
	ie.Right = p.parseExpression(precedence)
	return ie
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

/*
次のトークン（p.peekToken）が token.TokenType と一致しているか調べる
一致しているならそれを読み出したいので curToken <- peekToken へ更新する

アサーション関数と呼ばれ多くの構文解析器に存在する
次に来るトークンとして「期待する型」をチェックし正しい場合のみトークンを次に進める
*/
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func getPrecedence(t token.TokenType) int {
	if p, ok := precedences[t]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) curPrecedence() int  { return getPrecedence(p.curToken.Type) }
func (p *Parser) peekPrecedence() int { return getPrecedence(p.peekToken.Type) }

func (p *Parser) Errors() []string {
	return p.errors
}
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead.", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s", t)
	p.errors = append(p.errors, msg)
}
