package parser

import (
	"testing"

	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/lexer"
)

func TestIdentifierExpressionTest(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := checkIsExpressionStatements(t, program, 1)

	ident, ok := stmt.ExpressionValue.(*ast.IdentifierExpression)
	if !ok {
		t.Fatalf("program.Statements[0].ExpressionValue is not IdentifierExpression. got=%T", stmt.ExpressionValue)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.Literal not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := checkIsExpressionStatements(t, program, 1)

	literal, ok := stmt.ExpressionValue.(*ast.IntegerLiteralExpression)
	if !ok {
		t.Fatalf("program.Statements[0].ExpressionValue is not IntegerLiteralExpression. got=%T", stmt.ExpressionValue)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() not %s. got=%s", "5", literal.TokenLiteral())
	}

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input         string
		operator      string
		intergerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := checkIsExpressionStatements(t, program, 1)

		exp, ok := stmt.ExpressionValue.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.ExpressionValue)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s', got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.intergerValue) {
			return
		}
	}
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

func testIntegerLiteral(t *testing.T, e ast.Expression, value int64) bool {
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
