package evaluator

import "github.com/ganyariya/go_monkey/object"

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)
