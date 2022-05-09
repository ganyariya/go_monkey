package evaluator

import (
	"fmt"

	"github.com/ganyariya/go_monkey/object"
)

func builtinLen(args ...object.Object) object.Object {
	if ret := checkArgsLen(1, args...); ret != nil {
		return ret
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported, got=%s", arg.Type())
	}
}

func builtinFirst(args ...object.Object) object.Object {
	if ret := checkArgsLen(1, args...); ret != nil {
		return ret
	}
	errObj, arr := checkArgsIsArray("first", args...)
	if errObj != nil {
		return errObj
	}
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}
	return NULL
}

func builtinLast(args ...object.Object) object.Object {
	if ret := checkArgsLen(1, args...); ret != nil {
		return ret
	}
	errObj, arr := checkArgsIsArray("last", args...)
	if errObj != nil {
		return errObj
	}
	if len(arr.Elements) > 0 {
		return arr.Elements[len(arr.Elements)-1]
	}
	return NULL
}

func builtinRest(args ...object.Object) object.Object {
	if ret := checkArgsLen(1, args...); ret != nil {
		return ret
	}
	errObj, arr := checkArgsIsArray("rest", args...)
	if errObj != nil {
		return errObj
	}
	l := len(arr.Elements)
	if l > 0 {
		newElements := make([]object.Object, l-1)
		copy(newElements, arr.Elements[1:l])
		return &object.Array{Elements: newElements}
	}
	return NULL
}

func builtinPush(args ...object.Object) object.Object {
	if ret := checkArgsLen(2, args...); ret != nil {
		return ret
	}
	errObj, arr := checkArgsIsArray("push", args...)
	if errObj != nil {
		return errObj
	}
	l := len(arr.Elements)

	newElements := make([]object.Object, l+1)
	copy(newElements, arr.Elements)
	newElements[l] = args[1]
	return &object.Array{Elements: newElements}
}

func builtinPuts(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NULL
}

var builtins = map[string]*object.Builtin{
	"len":   {Fn: builtinLen},
	"first": {Fn: builtinFirst},
	"last":  {Fn: builtinLast},
	"rest":  {Fn: builtinRest},
	"push":  {Fn: builtinPush},
	"puts":  {Fn: builtinPuts},
}

// ------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------

func checkArgsLen(length int, args ...object.Object) *object.Error {
	if len(args) != length {
		return newError("wrong number of arguments. expected=%d, got=%d", length, len(args))
	}
	return nil
}
func checkArgsIsArray(name string, args ...object.Object) (*object.Error, *object.Array) {
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("argument to `%s` must be ARRAY, got=%T", name, args[0]), nil
	}
	return nil, arr
}
