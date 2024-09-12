package lexer

import (
	"testing"
)

type lexerTest struct {
	input          string
	expectedTokens []Token
}

func runLexerTests(t *testing.T, tests []lexerTest) {
	for idx, test := range tests {
		lexer := NewLexer(test.input)
		for i, expected := range test.expectedTokens {
			token := lexer.Next()
			if token != expected {
				t.Errorf("Test case %d failed at token %d: expected %v, got %v", idx, i, expected, token)
			}
		}
	}
}

func TestLexer(t *testing.T) {
	tests := []lexerTest{
		{
			input: ``,
			expectedTokens: []Token{
				{Kind: EOF, Lexeme: ""},
			},
		},
		{
			input: `price eq 10`,
			expectedTokens: []Token{
				{Kind: IDENT, Lexeme: "price"},
				{Kind: EQ, Lexeme: "eq"},
				{Kind: NUMBER, Lexeme: "10"},
				{Kind: EOF, Lexeme: ""},
			},
		},
		{
			input: `price @ 10`,
			expectedTokens: []Token{
				{Kind: IDENT, Lexeme: "price"},
				{Kind: ILLEGAL, Lexeme: "@"},
				{Kind: NUMBER, Lexeme: "10"},
				{Kind: EOF, Lexeme: ""},
			},
		},
		// Test unclosed string
		{
			input: `price eq "10`,
			expectedTokens: []Token{
				{Kind: IDENT, Lexeme: "price"},
				{Kind: EQ, Lexeme: "eq"},
				{Kind: ILLEGAL, Lexeme: "Unclosed string"},
				{Kind: EOF, Lexeme: ""},
			},
		},
		{
			input: `
name neq "John"
name contains "J"
group and
    price lt 10
    first_name eq "John"
last_name eq "Doe"
age not between 10 and 20
products.category[0][1].nested.nested_nested eq 10
another_field eq nil
`,
			expectedTokens: []Token{
				{Kind: IDENT, Lexeme: "name"},
				{Kind: NEQ, Lexeme: "neq"},
				{Kind: STRING, Lexeme: "John"},
				{Kind: IDENT, Lexeme: "name"},
				{Kind: CONTAINS, Lexeme: "contains"},
				{Kind: STRING, Lexeme: "J"},
				{Kind: GROUP, Lexeme: "group"},
				{Kind: AND, Lexeme: "and"},
				{Kind: INDENT, Lexeme: ""},
				{Kind: IDENT, Lexeme: "price"},
				{Kind: LT, Lexeme: "lt"},
				{Kind: NUMBER, Lexeme: "10"},
				{Kind: IDENT, Lexeme: "first_name"},
				{Kind: EQ, Lexeme: "eq"},
				{Kind: STRING, Lexeme: "John"},
				{Kind: DEDENT, Lexeme: ""},
				{Kind: IDENT, Lexeme: "last_name"},
				{Kind: EQ, Lexeme: "eq"},
				{Kind: STRING, Lexeme: "Doe"},
				{Kind: IDENT, Lexeme: "age"},
				{Kind: NOT, Lexeme: "not"},
				{Kind: BETWEEN, Lexeme: "between"},
				{Kind: NUMBER, Lexeme: "10"},
				{Kind: AND, Lexeme: "and"},
				{Kind: NUMBER, Lexeme: "20"},
				{Kind: IDENT, Lexeme: "products.category[0][1].nested.nested_nested"},
				{Kind: EQ, Lexeme: "eq"},
				{Kind: NUMBER, Lexeme: "10"},
				{Kind: IDENT, Lexeme: "another_field"},
				{Kind: EQ, Lexeme: "eq"},
				{Kind: NIL, Lexeme: "nil"},
				{Kind: EOF, Lexeme: ""},
			},
		},
		{
			input: `
# This is a comment
group and
    price lt 10 # This is a comment
    price_second lt 10
    # This is a comment
    group or
        price_3_field lt 20 # This is a comment
another_field lte 5
group and # This is a comment
    price_4 gt 100
    # This is a comment
    price_5 lt 10 # This is a comment
`,
			expectedTokens: []Token{
				{Kind: GROUP, Lexeme: "group"},
				{Kind: AND, Lexeme: "and"},
				{Kind: INDENT, Lexeme: ""},
				{Kind: IDENT, Lexeme: "price"},
				{Kind: LT, Lexeme: "lt"},
				{Kind: NUMBER, Lexeme: "10"},
				{Kind: IDENT, Lexeme: "price_second"},
				{Kind: LT, Lexeme: "lt"},
				{Kind: NUMBER, Lexeme: "10"},
				{Kind: GROUP, Lexeme: "group"},
				{Kind: OR, Lexeme: "or"},
				{Kind: INDENT, Lexeme: ""},
				{Kind: IDENT, Lexeme: "price_3_field"},
				{Kind: LT, Lexeme: "lt"},
				{Kind: NUMBER, Lexeme: "20"},
				{Kind: DEDENT, Lexeme: ""},
				{Kind: DEDENT, Lexeme: ""},
				{Kind: IDENT, Lexeme: "another_field"},
				{Kind: LTE, Lexeme: "lte"},
				{Kind: NUMBER, Lexeme: "5"},
				{Kind: GROUP, Lexeme: "group"},
				{Kind: AND, Lexeme: "and"},
				{Kind: INDENT, Lexeme: ""},
				{Kind: IDENT, Lexeme: "price_4"},
				{Kind: GT, Lexeme: "gt"},
				{Kind: NUMBER, Lexeme: "100"},
				{Kind: IDENT, Lexeme: "price_5"},
				{Kind: LT, Lexeme: "lt"},
				{Kind: NUMBER, Lexeme: "10"},
				{Kind: DEDENT, Lexeme: ""},
				{Kind: EOF, Lexeme: ""},
			},
		},
	}

	runLexerTests(t, tests)
}
