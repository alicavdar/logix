package evaluator

import (
	"testing"

	"github.com/alicavdar/logix/lexer"
	"github.com/alicavdar/logix/parser"
)

func newTestParser(input string) *parser.Parser {
	l := lexer.NewLexer(input)
	return parser.NewParser(l)
}

func TestEvaluator(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		context     map[string]interface{}
		expected    bool
		expectError bool
		errorMsg    string
	}{
		{
			name:  "Simple eq condition",
			input: "age eq 25",
			context: map[string]interface{}{
				"age": 25.0,
			},
			expected: true,
		},
		{
			name:  "Invalid operator error",
			input: "age xyz 30",
			context: map[string]interface{}{
				"age": 30.0,
			},
			expectError: true,
			errorMsg:    "unknown operator 'xyz'",
		},
		{
			name:  "Invalid numeric comparison types",
			input: `age lt "twenty"`,
			context: map[string]interface{}{
				"age": 30.0,
			},
			expectError: true,
			errorMsg:    "invalid types for numeric comparison: float64 and string",
		},
		{
			name: "Group with 'and' logic",
			input: `
group and
	age gt 18
	title contains "Hello"
`,
			context: map[string]interface{}{
				"age":   25.0,
				"title": "Hello World",
			},
			expected: true,
		},
		{
			name: "Group with 'or' logic failing",
			input: `
group or
	age gt 30
	title contains "Hi"
`,
			context: map[string]interface{}{
				"age":   25.0,
				"title": "Hello World",
			},
			expected: false,
		},
		{
			name: "Field resolution with out of range index",
			input: `
products[100].category.name eq "Electronics"
`,
			context: map[string]interface{}{
				"products": []interface{}{
					map[string]interface{}{
						"category": map[string]interface{}{
							"name": "Electronics",
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "array index out of range: 100",
		},
		{
			name: "Field resolution with valid index",
			input: `
products[0].category.name eq "Electronics"
`,
			context: map[string]interface{}{
				"products": []interface{}{
					map[string]interface{}{
						"category": map[string]interface{}{
							"name": "Electronics",
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestParser(tt.input)
			result, err := Evaluate(p, tt.context)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message: %s, but got: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
