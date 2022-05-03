package evaluator

import (
	"testing"

	"github.com/ganyariya/go_monkey/lexer"
	"github.com/ganyariya/go_monkey/object"
	"github.com/ganyariya/go_monkey/parser"
	"github.com/stretchr/testify/assert"
)

// Source Code -> 字句解析 -> 構文解析 -> 評価 -> Object
func callEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	return Eval(program)
}

func checkIntegerObject(t *testing.T, obj object.Object, expected int64, text string) {
	result, ok := obj.(*object.Integer)
	assert.True(t, ok, text)
	assert.Equal(t, expected, result.Value)
}

func checkBooleanObject(t *testing.T, obj object.Object, expected bool, text string) {
	result, ok := obj.(*object.Boolean)
	assert.True(t, ok, text)
	assert.Equal(t, expected, result.Value)
}

func checkNullObject(t *testing.T, obj object.Object, text string) {
	assert.Equal(t, NULL, obj, text)
}
