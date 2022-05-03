package evaluator

import (
	"fmt"
	"testing"

	"github.com/ganyariya/go_monkey/lexer"
	"github.com/ganyariya/go_monkey/object"
	"github.com/ganyariya/go_monkey/parser"
	"github.com/stretchr/testify/require"
)

// Source Code -> 字句解析 -> 構文解析 -> 評価 -> Object
func callEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	return Eval(program)
}

func checkIntegerObject(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	require.True(t, ok, fmt.Sprintf("object is not Integer. got=%T (%+v)", obj, obj))
	require.Equal(t, expected, result.Value)
}

func checkBooleanObject(t *testing.T, obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)
	require.True(t, ok, fmt.Sprintf("object is not Boolean. got=%T (%+v)", obj, obj))
	require.Equal(t, expected, result.Value)
}
