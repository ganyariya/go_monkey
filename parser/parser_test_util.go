package parser

import (
	"testing"

	"github.com/ganyariya/go_monkey/ast"
)

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

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------

func checkIsValidLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return checkIsIntegerLiteralExpression(t, exp, int64(v))
	case int64:
		return checkIsIntegerLiteralExpression(t, exp, v)
	case string:
		return checkIsIdentifierExpression(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", expected)
	return false
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
