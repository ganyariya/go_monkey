package evaluator

import (
	"fmt"

	"github.com/ganyariya/go_monkey/object"
)

func newError(format string, x ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, x...)}
}
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
