package parser

import (
	"testing"

	"github.com/ganyariya/go_monkey/ast"
)

func TestMacroLiteralParsing(t *testing.T) {
	input := "macro(x, y) {x + y;}"
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)

	macro, ok := stmt.ExpressionValue.(*ast.MacroExpression)
	if !ok {
		t.Fatalf("stmt.ExpressionValue is not ast.MacroLiteral. got=%T", stmt.ExpressionValue)
	}
	if len(macro.Parameters) != 2 {
		t.Fatalf("not 2. got=%d", len(macro.Parameters))
	}

	checkIsValidLiteralExpression(t, macro.Parameters[0], "x")
	checkIsValidLiteralExpression(t, macro.Parameters[1], "y")

	if len(macro.Body.Statements) != 1 {
		t.Fatalf("not 1. got=%d", len(macro.Body.Statements))
	}

	bodyStmt, ok := macro.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("bodyStmt is not Expression. got=%T", macro.Body.Statements[0])
	}

	checkIsValidInfixExpression(t, bodyStmt.ExpressionValue, "x", "+", "y")
}
