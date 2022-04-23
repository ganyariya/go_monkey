package lexer

import (
	"testing"

	"github.com/ganyariya/go_monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		nToken := l.NextToken()

		if nToken.Type != tt.expectedType {
			t.Fatalf("tests[%d] - TokenType Wrong. expected=%q, got=%q", i, tt.expectedType, nToken.Type)
		}

		if nToken.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal Wrong. expected=%q, got=%q", i, tt.expectedLiteral, nToken.Literal)
		}
	}

}
