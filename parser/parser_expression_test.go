package parser

import (
	"testing"

	"github.com/ganyariya/go_monkey/ast"
)

func TestIdentifierExpressionTest(t *testing.T) {
	input := "foobar;"

	_, program := initParserProgram(t, input)

	stmt := checkIsExpressionStatements(t, program, 1)

	checkIsValidLiteralExpression(t, stmt.ExpressionValue, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)
	checkIsValidLiteralExpression(t, stmt.ExpressionValue, 5)
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		_, program := initParserProgram(t, tt.input)
		stmt := checkIsExpressionStatements(t, program, 1)
		checkIsValidLiteralExpression(t, stmt.ExpressionValue, tt.expected)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\"Hello, World\";", "Hello, World"},
		{"\"Love...\";", "Love..."},
	}

	for _, tt := range tests {
		_, program := initParserProgram(t, tt.input)
		stmt := checkIsExpressionStatements(t, program, 1)
		literal, ok := stmt.ExpressionValue.(*ast.StringLiteralExpression)
		if !ok {
			t.Fatalf("exp not *ast.StringLiteralExpression. got=%T", stmt.ExpressionValue)
		}
		if literal.Value != tt.expected {
			t.Fatalf("expected=%s, got=%s", tt.expected, literal.Value)
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		_, program := initParserProgram(t, tt.input)
		stmt := checkIsExpressionStatements(t, program, 1)
		checkIsValidPrefixExpression(t, stmt.ExpressionValue, tt.operator, tt.value)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	inputTexts := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range inputTexts {
		_, program := initParserProgram(t, tt.input)
		stmt := checkIsExpressionStatements(t, program, 1)
		if !checkIsValidInfixExpression(t, stmt.ExpressionValue, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)

	exp, ok := stmt.ExpressionValue.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.ExpressionValue is not ast.IfExpression. got=%T", stmt.ExpressionValue)
	}

	if !checkIsValidInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Consequence length is not 1. got=%d", len(exp.Consequence.Statements))
	}

	// ここでは ExpressionStatement とテストケースで仮定しているが本来は任意の Statement
	conExp, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}
	if !checkIsIdentifierExpression(t, conExp.ExpressionValue, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Fatalf("exp.Alternative will be nil, got=%v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)

	exp, ok := stmt.ExpressionValue.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.ExpressionValue)
	}

	if !checkIsValidInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}
	if !checkIsIdentifierExpression(t, consequence.ExpressionValue, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}
	if !checkIsIdentifierExpression(t, alternative.ExpressionValue, "y") {
		return
	}
}

func TestFunctionExpressionParsing(t *testing.T) {
	input := `fn(x, y) {x + y;}`
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)

	function, ok := stmt.ExpressionValue.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("stmt.ExpressionValue is not FunctionExpression. got=%T", function)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("Parameter's length is not 2. got=%d", len(function.Parameters))
	}

	checkIsIdentifierExpression(t, function.Parameters[0], "x")
	checkIsIdentifierExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("Body.Statements's length is not 1. got=%d", len(function.Body.Statements))
	}
	checkIsValidInfixExpression(t, function.Body.Statements[0].(*ast.ExpressionStatement).ExpressionValue, "x", "+", "y")
}

func TestFunctionExpressionParametersParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn(){};", []string{}},
		{"fn(x){};", []string{"x"}},
		{"fn(x,){};", []string{"x"}},
		{"fn(x,y,z){};", []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		_, program := initParserProgram(t, tt.input)
		stmt := checkIsExpressionStatements(t, program, 1)
		function := stmt.ExpressionValue.(*ast.FunctionExpression)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Fatalf("len(parameters) wrong. expected=%d, got=%d", len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			checkIsValidLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)

	exp, ok := stmt.ExpressionValue.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Not ast.CallExpression. got=%T", stmt.ExpressionValue)
	}

	if !checkIsIdentifierExpression(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("exp.Arguments not 3. got=%d", len(exp.Arguments))
	}

	checkIsIntegerLiteralExpression(t, exp.Arguments[0], 1)
	checkIsValidInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	checkIsValidInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestArrayLiteralExpression(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3];"
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)

	exp, ok := stmt.ExpressionValue.(*ast.ArrayLiteralExpression)
	if !ok {
		t.Fatalf("exp not ArrayLiteralExpression. got=%T", stmt.ExpressionValue)
	}

	if len(exp.Elements) != 3 {
		t.Fatalf("len(exp.Elements) not 3, got=%d", len(exp.Elements))
	}

	checkIsIntegerLiteralExpression(t, exp.Elements[0], 1)
	checkIsValidInfixExpression(t, exp.Elements[1], 2, "*", 2)
	checkIsValidInfixExpression(t, exp.Elements[2], 3, "+", 3)

	input = "[]"
	_, program = initParserProgram(t, input)
	stmt = checkIsExpressionStatements(t, program, 1)
	exp = stmt.ExpressionValue.(*ast.ArrayLiteralExpression)
	if len(exp.Elements) != 0 {
		t.Fatalf("len(exp.Elements) not 0, got=%d", len(exp.Elements))
	}
}

func TestParsingArrayIndexExpressions(t *testing.T) {
	input := "myArray[1+1];"
	_, program := initParserProgram(t, input)
	stmt := checkIsExpressionStatements(t, program, 1)
	exp, ok := stmt.ExpressionValue.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not IndexExpression. got=%T", stmt.ExpressionValue)
	}
	if !checkIsIdentifierExpression(t, exp.Left, "myArray") {
		return
	}
	if !checkIsValidInfixExpression(t, exp.Index, 1, "+", 1) {
		return
	}
}
