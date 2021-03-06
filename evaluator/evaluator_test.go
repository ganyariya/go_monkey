package evaluator

import (
	"fmt"
	"testing"

	"github.com/ganyariya/go_monkey/lexer"
	"github.com/ganyariya/go_monkey/object"
	"github.com/ganyariya/go_monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"--10", 10},
		{"5 + 5", 10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2", 16},
		{"-10 + 20 - 10", 0},
		{"5 + 2 * -10", -15},
		{"5 * 2 - 10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (3 + 5)", 16},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"3 * (3 * 3 + 10)", 57},
	}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{{"true", true}, {"false", false}}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkBooleanObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!0", true},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!!0", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"true == true", true},
		{"true != true", false},
		{"false == false", true},
		{"false != false", false},
		{"(1 < 2) == true", true},
		{"(1 > 2) == true", false},
		{"3 == true", false},
		{"0 == true", false},
		{"4 != true", true},
		{`"Hello" == "Hello"`, true},
		{`"Hello" == "World"`, false},
		{`"Hello" != "World"`, true},
		{`!"World"`, false},
		{`!""`, true},
	}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkBooleanObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (0) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := callEval(tt.input)
		integer, ok := tt.expected.(int) // int64 ??????????????????????????????tests ??? int64(x) ????????????????????????
		if ok {
			checkIntegerObject(t, evaluated, int64(integer), tt.input)
		} else {
			checkNullObject(t, evaluated, tt.input)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9", 10},
		{"4; return 2*5; 9", 10},
		{"return 1; return 2;", 1},
		{
			/*
				if(10 > 1){statements} ??? 2 (`if(2>1)`???ReturnValue(2)?????????????????????2) ??????????????????
				???????????? return 10 ??? ReturnValue(10) ????????? 10 ???????????????????????? (p147)
			*/
			`if (10 > 1) {
				if (2 > 1) {
					return 2;
				}
				return 10;
			}`, 2,
		},
	}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true;", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 4;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { return true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { if (2 > 1) { return true + false; } return true; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar;", "identifier not found: foobar"},
		{"let x = 10 + foobar;", "identifier not found: foobar"},
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
	}

	for _, tt := range tests {
		evaluated := callEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("Error object is not returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("Wrong error message. expected=%s, got=%s", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5*5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = 2 * a; let c = a + b + 5;", 20},
		{"let a = 5; let a = 2 * a; a;", 10},
	}

	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected, tt.input)
	}

}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) {x + 2;}"
	evaluated := callEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T(%v)", evaluated, evaluated)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not x. got=%s", fn.Parameters[0].String())
	}
	if fn.Body.String() != "(x + 2)" {
		t.Fatalf("Body is wrong. got=%s", fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) {x;} identity(5);", 5},
		{"let identity = fn(x) {return x;} identity(5);", 5},
		{"let double = fn(x) {2 * x;} double(5);", 10},
		{"let add = fn(x, y) {x + y;} add(5 + 4, add(2, 3));", 14},
		{"fn(x){x;}(5)", 5},
		{"fn(x, y){if (x > y) {return x;} else {return y;}}(1, 2)", 2},
		{"let x = 10; let add = fn(y) {x + y;}; add(2)", 12},
		{"let x = 10; let identity = fn(x) {x;} identity(5);", 5},
		{"let x = 10; let identity = fn(x) {x;} identity(5); x;", 10},
		{
			`
			let adder = fn(x) { return fn(y) {return x + y;}; }	
			let addTwo = adder(2); addTwo(4);
			`, 6,
		}, // Closure
		{
			`
			let z = 10;
			let adder = fn(x) { return fn(y) {return x + y;}; }	
			let addTen = adder(z); let z = 20; addTen(4);
			`, 14,
		}, // Closure adder(z) ???????????? z = 10 ??????????????????????????? fn(y) ??? Env ?????????????????????
		{
			`
				let add = fn(x, y) {x + y;};
				let apply = fn(x, y, func) {func(x, y);};
				apply(1, 2, add);
			`, 3,
		},
		{"let fact = fn(x) {if (x == 1) {return 1;} else {x * fact(x - 1);} }; fact(5);", 120},
	}

	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\"Hello, World!\";", "Hello, World!"},
		{"\"Sei\" + \"Kin\";", "SeiKin"},
	}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		checkStringObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestArrayLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{"[1, 2 * 2, 3 + 3]", []int64{1, 4, 6}},
		{"[1, 2 * 2, fn(x, y){x+y}(1, 2)]", []int64{1, 4, 3}},
	}

	for _, tt := range tests {
		evaluated := callEval(tt.input)
		arr, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf("object is not Array. got=%T", evaluated)
		}
		if len(arr.Elements) != 3 {
			t.Fatalf("len(elements) not 3. got=%d", len(arr.Elements))
		}
		for i := 0; i < 3; i++ {
			checkIntegerObject(t, arr.Elements[i], tt.expected[i], tt.input)
		}
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 1; [10, 2][i]", 2},
		{"let i = 0; [10, 2][i]", 10},
		{"[1, 2, 3][1+1]", 3},
		{"let a = [1, 2, 3]; a[2]", 3},
		{"let a = [1, 2, 3]; a[0] + a[1] + a[2]", 6},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-10]", nil},
	}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			checkIntegerObject(t, evaluated, int64(integer), tt.input)
		} else {
			checkNullObject(t, evaluated, tt.input)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `
		let two = "two";
		{
			"one": 10 - 9,
			two: 1 + 1,
			"thr" + "ee": 6 / 2,
			4: 4,
			true: 5,
			false: 6
		}
	`
	evaluated := callEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("not object.Hash, got=%T", evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		(&object.Boolean{Value: true}).HashKey():   5,
		(&object.Boolean{Value: false}).HashKey():  6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	i := 0
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Fatal("no pair for given key in Pairs")
		}
		checkIntegerObject(t, pair.Value, expectedValue, fmt.Sprintf("%d", i))
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "foo"; {key: 5}["foo"]`, 5},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5:5}[5]`, 5},
		{`{true:5}[true]`, 5},
	}
	for _, tt := range tests {
		evaluated := callEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			checkIntegerObject(t, evaluated, int64(integer), tt.input)
		} else {
			checkNullObject(t, evaluated, tt.input)
		}
	}
}

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"quote(5)", "5"},
		{"quote(5 + 8)", "(5 + 8)"},
		{"quote(foobar)", "foobar"},
		{"quote(foobar + barfoo)", "(foobar + barfoo)"},
	}

	for _, tt := range tests {
		evaluated := callEval(tt.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected object.Quote. got=%T (%+v)", evaluated, evaluated)
		}
		if quote.Node == nil {
			t.Fatal("quote.Node is nil.")
		}

		if quote.Node.String() != tt.expected {
			t.Errorf("not equal. got=%s, expected=%s", quote.Node.String(), tt.expected)
		}
	}

}

func TestQuoteUnQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"quote(unquote(5))", "5"},
		{"quote(unquote(5 + 8))", "13"},
		{"quote(8 + unquote(4 + 4))", "(8 + 8)"},
		{"quote(unquote(4 + 4) + unquote(4 + 4))", "(8 + 8)"},
		{"let x = 8; quote(x)", "x"},
		{"let x = 8; quote(unquote(x))", "8"},
		{"quote(unquote(true))", "true"},
		{"quote(unquote(true == false))", "false"},
		{"quote(unquote(quote(4 + 4)))", "(4 + 4)"},
		// ??????????????????????????????????????????????????????quote ????????????????????? unquote ???????????? ast.Node ????????????????????? ast.Node ????????????????????????
		{"let x = quote(4 + 4); quote(unquote(4 + 4) + unquote(x))", "(8 + (4 + 4))"},
	}

	for _, tt := range tests {
		evaluated := callEval(tt.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected object.Quote. got=%T (%+v)", evaluated, evaluated)
		}
		if quote.Node == nil {
			t.Fatal("quote.Node is nil.")
		}

		if quote.Node.String() != tt.expected {
			t.Errorf("not equal. got=%s, expected=%s", quote.Node.String(), tt.expected)
		}
	}

}

func TestDefineMacros(t *testing.T) {
	input := `
	let number = 1;
	let function = fn(x, y) { x + y; };
	let mymacro = macro(x, y) { x + y; };
	`

	env := object.NewEnvironment()
	program := parser.NewParser(lexer.NewLexer(input)).ParseProgram()

	// ????????????????????? AST ?????? macro ??? ?????????env ???????????????AST ??????????????????
	DefineMacros(program, env)

	if len(program.Statements) != 2 {
		t.Fatalf("not 2. got=%d", len(program.Statements))
	}

	_, ok := env.Get("number")
	if ok {
		t.Fatal("number should not be defined")
	}
	_, ok = env.Get("function")
	if ok {
		t.Fatal("number should not be defined")
	}
	obj, ok := env.Get("mymacro")
	if !ok {
		t.Fatal("mymacro should be defined")
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("object is not Macro. got=%T", obj)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("not 2. got=%d", len(macro.Parameters))
	}

	if macro.Parameters[0].String() != "x" {
		t.Fatalf("not x")
	}
	if macro.Parameters[1].String() != "y" {
		t.Fatalf("not y")
	}

	expectedBody := "(x + y)"
	if macro.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, macro.Body.String())
	}
}
