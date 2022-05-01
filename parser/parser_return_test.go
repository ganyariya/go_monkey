package parser

import (
	"testing"

	"github.com/ganyariya/go_monkey/ast"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}
	for _, tt := range tests {
		_, program := initParserProgram(t, tt.input)
		returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", program.Statements[0])
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if checkIsValidLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}
