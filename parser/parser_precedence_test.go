package parser

import (
	"testing"

	"github.com/ganyariya/go_monkey/lexer"
)

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a", "a"},
		{"-a", "(-a)"},
		{"-a * b", "((-a) * b)"},
		/*
			* {
				Left: - {
					Right: a
				}
				Right: ! {
					Right: b
				}
			}
		*/
		{"-a * !b", "((-a) * (!b))"},
		{"-a * !!b", "((-a) * (!(!b)))"},
		/*
			! {
				Right: - {
					Right: a
				}
			}
		*/
		{"!-a", "(!(-a))"},
		{"!!!!a", "(!(!(!(!a))))"},
		{"a + b + c", "((a + b) + c)"},
		/*
			// 数字はその記号が元の式で左から何番目に出てくるか
			+2 {
				Left: +1 {
					Left: a
					Right: b
				}
				Right: -1 {
					Right: c
				}
			}
		*/
		{"a + b + -c", "((a + b) + (-c))"},
		{"a * b * -c", "((a * b) * (-c))"},
		{"a * b / -c", "((a * b) / (-c))"},
		/*
			+1 {
				Left: a
				Right: *1 {
					Left: b
					Right: c
				}
			}
		*/
		{"a + b * c", "(a + (b * c))"},
		/*
			-1 {
				Left: +2 {
					Left: +1 {
						Left: a
						Right: *1 {
							Left: b
							Right: c
						}
					}
					Right: /1 {
						Left: d
						Right: e
					}
				}
				Right: f
			}
		*/
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5;", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 > 4 != 3 < 4", "((5 > 4) != (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}

}
