package evaluator

import (
	"testing"

	"github.com/ganyariya/go_monkey/lexer"
	"github.com/ganyariya/go_monkey/object"
	"github.com/ganyariya/go_monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{{"5", 5}, {"10", 10}}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected)
	}
}

func callEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	return Eval(program)
}
