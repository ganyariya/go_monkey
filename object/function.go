package object

import (
	"fmt"
	"strings"

	"github.com/ganyariya/go_monkey/ast"
)

/*
Env = 「Function が `定義された` 時点での環境」（親環境とも言える）
関数を`実行する`ときはその都度新しい環境がつくられ、定義された時点の環境でラップされる

定義された時点の環境を保持しているため「クロージャ」の概念がうまれる

Function Object を定義するときは「Body」と「Parameters」は Expression のまま保持するのみ
*/
type Function struct {
	Parameters []*ast.IdentifierExpression
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("fn(%s) {\n%s\n}", strings.Join(params, ","), f.Body.String())
}
func (f *Function) AsBool() bool { return true }
