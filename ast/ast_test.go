package ast

import (
	"testing"

	"github.com/ganyariya/go_monkey/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			// let myVar = anotherVar;
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
