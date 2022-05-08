package parser

import (
	"testing"

	"github.com/ganyariya/go_monkey/ast"
)

func TestParsingHashLiteralStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)
	exp, ok := stmt.ExpressionValue.(*ast.HashLiteralExpression)
	if !ok {
		t.Fatalf("exp not HashLiteralExpression. got=%T", stmt.ExpressionValue)
	}
	if len(exp.Pairs) != 3 {
		t.Fatalf("len(Pairs) not 3. got=%d", len(exp.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for ke, ve := range exp.Pairs {
		literal, ok := ke.(*ast.StringLiteralExpression)
		if !ok {
			t.Fatalf("hash's key not StringLiteral. got=%T", ke)
		}
		expectedValue := expected[literal.String()]
		checkIsIntegerLiteralExpression(t, ve, expectedValue)
	}

}
func TestParsingHashLiteralEmpty(t *testing.T) {
	input := "{}"
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)
	exp, ok := stmt.ExpressionValue.(*ast.HashLiteralExpression)
	if !ok {
		t.Fatalf("exp not HashLiteralExpression. got=%T", stmt.ExpressionValue)
	}
	if len(exp.Pairs) != 0 {
		t.Fatalf("len(Pairs) not 0. got=%d", len(exp.Pairs))
	}
}
func TestParsingHashLiteralIntegerKeys(t *testing.T) {
	input := `{1: "one", 2: "two", 3: "three"}`
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)
	exp, ok := stmt.ExpressionValue.(*ast.HashLiteralExpression)
	if !ok {
		t.Fatalf("exp not HashLiteralExpression. got=%T", stmt.ExpressionValue)
	}
	if len(exp.Pairs) != 3 {
		t.Fatalf("len(Pairs) not 3. got=%d", len(exp.Pairs))
	}

	expected := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}
	for ke, ve := range exp.Pairs {
		literal, ok := ke.(*ast.IntegerLiteralExpression)
		if !ok {
			t.Fatalf("hash's key not IntegerLiteral. got=%T", ke)
		}
		expectedValue := expected[int(literal.Value)]
		checkIsStringLiteralExpression(t, ve, expectedValue)
	}

}
func TestParsingHashLiteralBooleanKeys(t *testing.T) {
	input := `{true: "true", false: "false"}`
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)
	exp, ok := stmt.ExpressionValue.(*ast.HashLiteralExpression)
	if !ok {
		t.Fatalf("exp not HashLiteralExpression. got=%T", stmt.ExpressionValue)
	}
	if len(exp.Pairs) != 2 {
		t.Fatalf("len(Pairs) not 2. got=%d", len(exp.Pairs))
	}

	expected := map[bool]string{
		true:  "true",
		false: "false",
	}
	for ke, ve := range exp.Pairs {
		literal, ok := ke.(*ast.BooleanExpression)
		if !ok {
			t.Fatalf("hash's key not IntegerLiteral. got=%T", ke)
		}
		expectedValue := expected[literal.Value]
		checkIsStringLiteralExpression(t, ve, expectedValue)
	}

}
func TestParsingHashLiteralWithExpressions(t *testing.T) {
	input := `{"one": 0+1, "two": 10-8, "three": 15/5}`
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)
	exp, ok := stmt.ExpressionValue.(*ast.HashLiteralExpression)
	if !ok {
		t.Fatalf("exp not HashLiteralExpression. got=%T", stmt.ExpressionValue)
	}
	if len(exp.Pairs) != 3 {
		t.Fatalf("len(Pairs) not 3. got=%d", len(exp.Pairs))
	}

	expected := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			checkIsValidInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			checkIsValidInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			checkIsValidInfixExpression(t, e, 15, "/", 5)
		},
	}
	for ke, ve := range exp.Pairs {
		literal, ok := ke.(*ast.StringLiteralExpression)
		if !ok {
			t.Fatalf("hash's key not StringLiteral. got=%T", ke)
		}
		testFunc, ok := expected[literal.String()]
		if !ok {
			t.Fatalf("No test function for key %s", literal.String())
		}
		testFunc(ve)
	}

}
