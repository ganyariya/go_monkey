package evaluator

import (
	"testing"

	"github.com/ganyariya/go_monkey/object"
	"github.com/stretchr/testify/assert"
)

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("");`, 0},
		{`len("four");`, 4},
		{`len("hello, world");`, 12},
		{`len(1);`, "argument to `len` not supported, got=INTEGER"},
		{`len("one", "two");`, "wrong number of arguments. expected=1, got=2"},
		{`len([]);`, 0},
		{`len([1, "hello"]);`, 2},
		{`first([10, 20]);`, 10},
		{`last([10, 20]);`, 20},
	}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			checkIntegerObject(t, evaluated, int64(expected), tt.input)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("object is not error object. got=%T", errObj)
			}
			assert.Equal(t, errObj.Message, tt.expected)
		}
	}
}
