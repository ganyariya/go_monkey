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

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statement. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T", program.Statements[0])
	}

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
