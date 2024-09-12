package logix

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alicavdar/logix/evaluator"
	"github.com/alicavdar/logix/lexer"
	"github.com/alicavdar/logix/parser"
)

func EvaluateLogix(logixContent string, context map[string]interface{}) (bool, error) {
	lex := lexer.NewLexer(logixContent)
	pr := parser.NewParser(lex)

	return evaluator.Evaluate(pr, context)
}

func LoadContextFromFile(filepath string) (map[string]interface{}, error) {
	contextContent, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error occurred while reading the context JSON file: %v", err)
	}

	var context map[string]interface{}
	err = json.Unmarshal(contextContent, &context)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON file for context: %v", err)
	}

	return context, nil
}
