package object

import (
	"fmt"
	"strings"

	"github.com/ganyariya/go_monkey/ast"
)

type Macro struct {
	Parameters []*ast.IdentifierExpression
	Body       *ast.BlockStatement
	Env        *Environment
}

func (m *Macro) Type() ObjectType { return MACRO_OBJ }
func (m *Macro) Inspect() string {
	params := []string{}
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("macro(%s){\n%s\n}", strings.Join(params, ","), m.Body.String())
}
func (m *Macro) AsBool() bool { return true }
