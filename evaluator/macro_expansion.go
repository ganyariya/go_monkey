package evaluator

import (
	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	// マクロ定義を探す
	for i, stmt := range program.Statements {
		if isMacroDefinition(stmt) {
			addMacro(stmt, env)
			definitions = append(definitions, i)
		}
	}

	// マクロをASTから取り除く
	for i := len(definitions) - 1; i >= 0; i-- {
		defIdx := definitions[i]
		program.Statements = append(program.Statements[:defIdx], program.Statements[defIdx+1:]...)
	}
}

/*
マクロ Statement かチェックする
- let x = `macro` の構文である
- MacroExpression である
*/
func isMacroDefinition(node ast.Statement) bool {
	letStmt, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}
	_, ok = letStmt.Value.(*ast.MacroExpression)
	return ok
}

/*
マクロ固有の環境に閉じ込める
*/
func addMacro(stmt ast.Statement, env *object.Environment) {
	letStmt := stmt.(*ast.LetStatement)
	macroExp := letStmt.Value.(*ast.MacroExpression)
	macro := &object.Macro{
		Parameters: macroExp.Parameters,
		Body:       macroExp.Body,
		Env:        env,
	}
	env.Set(letStmt.Name.Value, macro)
}

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExp, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}
		macro, ok := isMacroCall(callExp, env)
		if !ok {
			return node
		}

		args := quoteArgs(callExp)
		enclosedEnv := encloseMacroEnv(macro, args)

		/*
			let a = macro(x, y) {quote(unquote(y) - macro(x))}; macro(10 + 4, 2 - 3) のとき
			enclosedEnv{x = object.Quote(2 - 3), y=object.Quote(10 + 4)} のようになる
		*/
		evaluated := Eval(macro.Body, enclosedEnv)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macro")
		}

		return quote.Node
	})
}

/*
関数呼び出しについて `macro()` かどうかチェックする
*/
func isMacroCall(exp *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	ident, ok := exp.Function.(*ast.IdentifierExpression)
	if !ok {
		return nil, false
	}
	obj, ok := env.Get(ident.Value)
	if !ok {
		return nil, false
	}
	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}
	return macro, true
}

/*
macro(10, 40 + 2) のような引数([]ast.Expression)を Quote へそれぞれ入れて object.Quote にする
*/
func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	args := []*object.Quote{}
	for _, a := range exp.Arguments {
		args = append(args, &object.Quote{Node: a})
	}
	return args
}

/*
仮引数（変数）と実引数（object.Quote(未評価の「macro に渡された実引数」)）を紐付けた Environment を返す
*/
func encloseMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	enclosedEnv := object.NewEnclosedEnvironment(macro.Env)
	// Parameters = 仮引数[x, y]  args = &object.Macros[4+10, 2]
	for idx, param := range macro.Parameters {
		// x = 4 + 10
		enclosedEnv.Set(param.Value, args[idx])
	}
	return enclosedEnv
}
