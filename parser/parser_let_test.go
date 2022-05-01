package parser

import (
	"testing"

	"github.com/ganyariya/go_monkey/ast"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}
	for _, tt := range tests {
		_, program := initParserProgram(t, tt.input)
		if !checkIsValidLetStatement(t, program.Statements[0], tt.expectedIdentifier, tt.expectedValue) {
			return
		}
	}
}

func checkIsValidLetStatement(t *testing.T, stmt ast.Statement, name string, value interface{}) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let', get=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt is not LetStatement. got=%T", stmt)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name is not '%s', got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}
	if !checkIsValidLiteralExpression(t, letStmt.Value, value) {
		return false
	}
	return true
}
