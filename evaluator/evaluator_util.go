package evaluator

import (
	"fmt"
	"testing"

	"github.com/ganyariya/go_monkey/object"
	"github.com/stretchr/testify/require"
)

func checkIntegerObject(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	require.True(t, ok, fmt.Sprintf("object is not Integer. got=%T (%+v)", obj, obj))
	require.Equal(t, expected, result.Value)
}
