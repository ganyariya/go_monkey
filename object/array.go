package object

import (
	"fmt"
	"strings"
)

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}
func (a *Array) AsBool() bool {
	return len(a.Elements) > 0
}
