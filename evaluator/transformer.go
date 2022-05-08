package evaluator

import (
	"fmt"

	"github.com/ganyariya/go_monkey/ast"
	"github.com/ganyariya/go_monkey/object"
)

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var ret object.Object
	for _, stmt := range program.Statements {
		ret = Eval(stmt, env)
		switch ret := ret.(type) {
		case *object.ReturnValue:
			return ret.Value
		case *object.Error:
			return ret
		}
	}
	return ret
}

func evalBlockStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var ret object.Object
	for _, stmt := range stmts {
		ret = Eval(stmt, env)
		// BlockStatement では ReturnValue.Value にアンラップしない（ブロック文ネストでバグる)
		if ret != nil {
			if ret.Type() == object.RETURN_VALUE_OBJ || ret.Type() == object.ERROR_OBJ {
				return ret
			}
		}
	}
	return ret
}

func evalReturnStatement(stmt *ast.ReturnStatement, env *object.Environment) object.Object {
	obj := Eval(stmt.ReturnValue, env)
	if isError(obj) {
		return obj
	}
	return &object.ReturnValue{Value: obj}
}

func evalLetStatement(stmt *ast.LetStatement, env *object.Environment) object.Object {
	expObj := Eval(stmt.Value, env)
	if isError(expObj) {
		return expObj
	}
	env.Set(stmt.Name.Value, expObj)
	return expObj
}

// -----------------------------------------------------------
// -----------------------------------------------------------

func nativeBoolToBooleanObject(b bool) object.Object {
	if b {
		return TRUE
	} else {
		return FALSE
	}
}
func evalBooleanExpression(exp *ast.BooleanExpression) object.Object {
	return nativeBoolToBooleanObject(exp.Value)
}

func evalIntegerLiteralExpression(exp *ast.IntegerLiteralExpression) object.Object {
	return &object.Integer{Value: exp.Value}
}

func evalStringLiteralExpression(exp *ast.StringLiteralExpression) object.Object {
	return &object.String{Value: exp.Value}
}

func evalIdentifierExpression(exp *ast.IdentifierExpression, env *object.Environment) object.Object {
	if obj, ok := env.Get(exp.Value); ok {
		return obj
	}
	/* 組み込み関数は言語側ではじめから定義されており「識別子（）」で呼び出される */
	if builtin, ok := builtins[exp.Value]; ok {
		return builtin
	}
	return newError(fmt.Sprintf("identifier not found: %s", exp.Value))
}

func evalPrefixExpression(exp *ast.PrefixExpression, env *object.Environment) object.Object {
	rightObj := Eval(exp.Right, env)
	if isError(rightObj) {
		return rightObj
	}
	switch exp.Operator {
	case "!":
		return evalBangPrefixOperator(rightObj)
	case "-":
		return evalMinusPrefixOperator(rightObj)
	default:
		return newError("unknown operator: %s%s", exp.Operator, rightObj.Type())
	}
}

func evalInfixExpression(exp *ast.InfixExpression, env *object.Environment) object.Object {
	leftObj := Eval(exp.Left, env)
	if isError(leftObj) {
		return leftObj
	}
	rightObj := Eval(exp.Right, env)
	if isError(rightObj) {
		return rightObj
	}
	switch {
	// 整数は「値」で処理する
	case leftObj.Type() == object.INTEGER_OBJ && rightObj.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(exp.Operator, leftObj, rightObj)
	case leftObj.Type() == object.STRING_OBJ && rightObj.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(exp.Operator, leftObj, rightObj)
	// reference (pointer) （異なる型 -> false）
	case exp.Operator == "==":
		return nativeBoolToBooleanObject(leftObj == rightObj)
	case exp.Operator == "!=":
		return nativeBoolToBooleanObject(leftObj != rightObj)
	case leftObj.Type() != rightObj.Type():
		return newError("type mismatch: %s %s %s", leftObj.Type(), exp.Operator, rightObj.Type())
	default:
		return newError("unknown operator: %s %s %s", leftObj.Type(), exp.Operator, rightObj.Type())
	}
}

func evalIfExpression(exp *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(exp.Condition, env)
	if isError(condition) {
		return condition
	}
	if condition.AsBool() {
		return Eval(exp.Consequence, env)
	} else if exp.Alternative != nil {
		return Eval(exp.Alternative, env)
	} else {
		return NULL
	}
}

/*
Function が「定義された」時点における Env を保持する（関数を実行するときに新しい EnclosedEnv をつくる）
*/
func evalFunctionExpression(exp *ast.FunctionExpression, env *object.Environment) object.Object {
	return &object.Function{Parameters: exp.Parameters, Body: exp.Body, Env: env}
}

/*
引数にある env = 定義時点での Env
*/
func evalCallExpression(exp *ast.CallExpression, env *object.Environment) object.Object {
	/*
		Identifier -> 識別子に対応する関数を取り出す
		Function -> 関数を直接得る
		Function は`定義時点`における env を保持する
	*/
	fnObj := Eval(exp.Function, env)
	if isError(fnObj) {
		return fnObj
	}
	args := evalExpressions(exp.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}
	return applyCallFunction(fnObj, args)
}

func evalArrayLiteralExpression(exp *ast.ArrayLiteralExpression, env *object.Environment) object.Object {
	elements := evalExpressions(exp.Elements, env)
	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}
	return &object.Array{Elements: elements}
}

func evalIndexExpression(exp *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(exp.Left, env)
	if isError(left) {
		return left
	}
	index := Eval(exp.Index, env)
	if isError(index) {
		return index
	}
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return extractArrayByIndex(left, index)
	default:
		return newError("index operator not supported: %s", index.Type())
	}
}

// ------------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------------

func evalBangPrefixOperator(right object.Object) object.Object {
	return nativeBoolToBooleanObject(!right.AsBool())
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func extractArrayByIndex(arr, index object.Object) object.Object {
	arrObj := arr.(*object.Array)
	idx := index.(*object.Integer).Value
	if idx < 0 || int(idx) >= len(arrObj.Elements) {
		return NULL
	}
	return arrObj.Elements[idx]
}

// ------------------------------------------------------------------------------------------------------------
// Call Function
// ------------------------------------------------------------------------------------------------------------

/*
引数や Array などで出現する 「式の列」を「Object の配列」に評価する
*/
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var ret []object.Object
	for _, exp := range exps {
		evaluated := Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		ret = append(ret, evaluated)
	}
	return ret
}

/*
評価済みの arg objects を function object に与えて関数式を評価する。
*/
func applyCallFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		registeredEnv := registerEnclosedCallEnv(fn, args)
		evaluated := Eval(fn.Body, registeredEnv)
		/* Unwrap しないと return 効果が関数をまたいで浮上して実行が途中で停止してしまう */
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

/*
仮引数（変数）と実引数（実値）を紐付けた 新たな記憶容量 Environment を返す
**Function Object が持つ親環境に 新しい環境はラップされる**
*/
func registerEnclosedCallEnv(fnObj *object.Function, args []object.Object) *object.Environment {
	enclosedEnv := object.NewEnclosedEnvironment(fnObj.Env)
	// Parameters = 仮引数[x, y, z]  args = 評価済実引数[10, 1, 4]
	for i := 0; i < len(fnObj.Parameters); i++ {
		// 変数に値を登録する (x = 10)
		enclosedEnv.Set(fnObj.Parameters[i].Value, args[i])
	}
	return enclosedEnv
}

// ------------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------------

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
