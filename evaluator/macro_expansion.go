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
