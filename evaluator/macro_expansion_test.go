package evaluator

import (
	"testing"

	"github.com/ganyariya/go_monkey/lexer"
	"github.com/ganyariya/go_monkey/object"
	"github.com/ganyariya/go_monkey/parser"
)

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// マクロは必ず *object.Quote を返すようにする（macro() 内では quote を使う）
		{
			`
			let infix = macro() {quote(1 + 2);};
			infix();		
			`,
			`(1 + 2)`,
		},
		// unquote がないと `b - a` が帰ってくる
		{
			`
			let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
			reverse(2+2, 10-5);
			`,
			`(10 - 5) - (2 + 2)`,
		},
	}

	for _, tt := range tests {
		// AST (評価はしてない)
		expected := parser.NewParser(lexer.NewLexer(tt.expected)).ParseProgram()
		program := parser.NewParser(lexer.NewLexer(tt.input)).ParseProgram()

		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)

		if expected.String() != expanded.String() {
			t.Errorf("not equal. want=%s, got=%s", expected.String(), expanded.String())
		}
	}
}
