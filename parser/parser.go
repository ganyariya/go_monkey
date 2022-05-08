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
	INDEX       // array[index]
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
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
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
	p.registerPrefixFn(token.TRUE, p.parseBooleanExpression)
	p.registerPrefixFn(token.FALSE, p.parseBooleanExpression)
	p.registerPrefixFn(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixFn(token.IF, p.parseIfExpression)
	p.registerPrefixFn(token.FUNCTION, p.parseFunctionExpression)
	p.registerPrefixFn(token.STRING, p.parseStringLiteralExpression)
	p.registerPrefixFn(token.LBRACKET, p.parseArrayLiteralExpression)
	p.registerPrefixFn(token.LBRACE, p.parseHashLiteralExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)
	p.registerInfixFn(token.LPAREN, p.parseCallExpression)
	p.registerInfixFn(token.LBRACKET, p.parseIndexExpression)

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

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
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

// curToken が `{` で開始し `}` で終了する
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken} // Token = token.LBRACE
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
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
		// 次の中間演算子の優先順位の方が高いなら次の中置演算子に紐付いた関数でパースする
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

// token.TRUE か FALSE がトークンとして与えられるので「token.TRUE」かを判定することで true/false 式にする
func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.BooleanExpression{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseStringLiteralExpression() ast.Expression {
	return &ast.StringLiteralExpression{Token: p.curToken, Value: p.curToken.Literal}
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

/*
`(` で呼び出される。 expectPeek で `)` を飛ばす。
つまり、`(` は構文解析時に除去され 専用の Expression はない。
*/
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// TODO: elif を追加する
func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		exp.Alternative = p.parseBlockStatement()
	}
	return exp
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	fe := &ast.FunctionExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	fe.Parameters = []*ast.IdentifierExpression{}
	for _, p := range p.parseExpressionList(token.RPAREN) {
		fe.Parameters = append(fe.Parameters, p.(*ast.IdentifierExpression))
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fe.Body = p.parseBlockStatement()
	// 最後は curToken = } を指した状態で終了する
	return fe
}

/*
`(` の **infix** として登録される。
これは add(x, y) や <fn(x, y){x+y}>(x, y) のように「式」のあとの中置演算 `(`（関数呼び出し）だからである。
function には `add` のような IdentifierExpression か `fn(x, y){x+y}` のような FunctionExpression が入る。
*/
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	call := &ast.CallExpression{Token: p.curToken, Function: function}
	call.Arguments = p.parseExpressionList(token.RPAREN)
	return call
}

func (p *Parser) parseArrayLiteralExpression() ast.Expression {
	array := &ast.ArrayLiteralExpression{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) parseHashLiteralExpression() ast.Expression {
	exp := &ast.HashLiteralExpression{Token: p.curToken, Pairs: make(map[ast.Expression]ast.Expression)}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) {
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)
		exp.Pairs[key] = value
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
		p.nextToken()
	}
	return exp
}

func (p *Parser) parseExpressionList(endToken token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	p.nextToken()
	for !p.curTokenIs(endToken) {
		list = append(list, p.parseExpression(LOWEST))
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		} else if !p.peekTokenIs(endToken) {
			return nil
		}
		p.nextToken()
	}
	return list
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

/*
次のトークン（p.peekToken）が token.TokenType と一致しているか調べる
一致しているならそれを読み出したいので curToken <- peekToken へ更新する

アサーション関数と呼ばれ多くの構文解析器に存在する（IfExpression のパースなどで上手に利用できる）
次に来るトークンとして「期待する型」をチェックし正しい場合のみトークンを次に進める
失敗したら false を返しエラーとする
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

// `)` は優先順位が最低=LOWESTとなる
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
