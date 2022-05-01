package parser

import (
	"fmt"
	"testing"

	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/lexer"
)

func initParserProgram(t *testing.T, input string) (*Parser, *ast.Program) {
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	return p, program
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------

func checkIsExpressionStatements(t *testing.T, program *ast.Program, programLen int) *ast.ExpressionStatement {
	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain 1 statement. got=%d (program=%v)", len(program.Statements), program.Statements)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T", program.Statements[0])
	}
	return stmt
}

func checkIsValidLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return checkIsIntegerLiteralExpression(t, exp, int64(v))
	case int64:
		return checkIsIntegerLiteralExpression(t, exp, v)
	case bool:
		return checkIsBooleanExpression(t, exp, v)
	case string:
		return checkIsIdentifierExpression(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", expected)
	return false
}

func checkIsValidPrefixExpression(t *testing.T, exp ast.Expression, operator string, right interface{}) bool {
	preExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("exp is not ast.PrefixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if preExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, preExp.Operator)
		return false
	}
	if !checkIsValidLiteralExpression(t, preExp.Right, right) {
		return false
	}
	return true
}

func checkIsValidInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !checkIsValidLiteralExpression(t, infixExp.Left, left) {
		return false
	}
	if infixExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, infixExp.Operator)
		return false
	}
	if !checkIsValidLiteralExpression(t, infixExp.Right, right) {
		return false
	}
	return true
}

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------

// 与えられた式が IntegerLiteralExpression かテストする Helper
func checkIsIntegerLiteralExpression(t *testing.T, e ast.Expression, value int64) bool {
	ile, ok := e.(*ast.IntegerLiteralExpression)
	if !ok {
		t.Errorf("ile not *ast.IntegerLiteralExpression. got=%T", e)
		return false
	}
	if ile.Value != value {
		t.Errorf("ile.Value not %d. got=%d", value, ile.Value)
		return false
	}
	return true
}

func checkIsBooleanExpression(t *testing.T, e ast.Expression, value bool) bool {
	be, ok := e.(*ast.BooleanExpression)
	if !ok {
		t.Errorf("be not *ast.BooleanExpression. got=%T", e)
		return false
	}
	if be.Value != value {
		t.Errorf("be.Value not %v. got=%v", value, be.Value)
		return false
	}
	if be.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("be.TokenLiteral not %t. got=%s", value, be.TokenLiteral())
		return false
	}
	return true
}

// 与えられた式が IdentifierExpression かテストする Helper
func checkIsIdentifierExpression(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.IdentifierExpression)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value is not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is not %s. got=%s", value, ident.Value)
		return false
	}
	return true
}
