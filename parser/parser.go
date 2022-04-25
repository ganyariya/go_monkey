package parser

import (
	"fmt"

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

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.ExpressionValue = p.parseExpression(LOWEST)
	// セミコロンを省略可能にする（セミコロンだったら飛ばす）
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParseFns[p.curToken.Type]
	if prefixFn == nil {
		return nil
	}
	leftExp := prefixFn()
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

func (p *Parser) Errors() []string {
	return p.errors
}
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead.", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
