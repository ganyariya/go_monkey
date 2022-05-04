package evaluator

import (
	"testing"

	"github.com/ganyariya/go_monkey/object"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"--10", 10},
		{"5 + 5", 10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2", 16},
		{"-10 + 20 - 10", 0},
		{"5 + 2 * -10", -15},
		{"5 * 2 - 10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (3 + 5)", 16},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"3 * (3 * 3 + 10)", 57},
	}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{{"true", true}, {"false", false}}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkBooleanObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!0", true},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!!0", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"true == true", true},
		{"true != true", false},
		{"false == false", true},
		{"false != false", false},
		{"(1 < 2) == true", true},
		{"(1 > 2) == true", false},
		{"3 == true", false},
		{"0 == true", false},
		{"4 != true", true},
	}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkBooleanObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (0) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := callEval(tt.input)
		integer, ok := tt.expected.(int) // int64 にキャストできない（tests で int64(x) にしてないため）
		if ok {
			checkIntegerObject(t, evaluated, int64(integer), tt.input)
		} else {
			checkNullObject(t, evaluated, tt.input)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9", 10},
		{"4; return 2*5; 9", 10},
		{"return 1; return 2;", 1},
		{
			/*
				if(10 > 1){statements} で 2 (`if(2>1)`でReturnValue(2)から取り出した2) を受け取るが
				その後の return 10 が ReturnValue(10) として 10 が返されてしまう (p147)
			*/
			`if (10 > 1) {
				if (2 > 1) {
					return 2;
				}
				return 10;
			}`, 2,
		},
	}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true;", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 4;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { return true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { if (2 > 1) { return true + false; } return true; }", "unknown operator: BOOLEAN + BOOLEAN"},
	}

	for _, tt := range tests {
		evaluated := callEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("Error object is not returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("Wrong error message. expected=%s, got=%s", tt.expectedMessage, errObj.Message)
		}
	}
}
