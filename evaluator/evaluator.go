package evaluator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/alicavdar/logix/parser"
)

func Evaluate(p *parser.Parser, context map[string]interface{}) (bool, error) {
	for {
		parsed := p.ParseNext()
		if parsed == nil {
			break
		}

		switch item := parsed.(type) {
		case *parser.Condition:
			conditionResult, err := evaluateCondition(item, context)
			if err != nil {
				return false, err
			}

			if !conditionResult {
				return false, nil
			}
		case *parser.Group:
			groupResult, err := evaluateGroup(item, context)
			if err != nil {
				return false, err
			}

			if !groupResult {
				return false, nil
			}
		default:
			return false, fmt.Errorf("unexpected item type: %T", parsed)
		}
	}

	return true, nil
}

func evaluateCondition(cond *parser.Condition, context map[string]interface{}) (bool, error) {
	fieldValue, err := resolveFieldValue(cond.Field, context)
	if err != nil {
		return false, err
	}

	conditionValue := cond.Value[0]

	switch cond.Operator {
	case "lt", "gt", "lte", "gte":
		return compareNumeric(fieldValue, conditionValue, cond.Operator, cond.Negate)
	case "eq":
		return applyNegation(fieldValue == conditionValue, cond.Negate), nil
	case "neq":
		return applyNegation(fieldValue != conditionValue, cond.Negate), nil
	case "contains":
		strVal, ok := fieldValue.(string)
		if !ok {
			return false, fmt.Errorf("the field value is not a string for 'contains' operator")
		}

		return applyNegation(strings.Contains(strVal, conditionValue.(string)), cond.Negate), nil
	case "between":
		return evaluateBetween(fieldValue, cond.Value, cond.Negate)
	case "startsWith":
		strVal, ok := fieldValue.(string)
		if !ok {
			return false, fmt.Errorf("the field value is not a string for 'startsWith' operator")
		}

		return applyNegation(strings.HasPrefix(strVal, conditionValue.(string)), cond.Negate), nil
	case "endsWith":
		strVal, ok := fieldValue.(string)
		if !ok {
			return false, fmt.Errorf("field value is not a string for 'endsWith' operator")
		}

		return applyNegation(strings.HasSuffix(strVal, conditionValue.(string)), cond.Negate), nil
	case "in":
		return evaluateIn(fieldValue, cond.Value, cond.Negate)
	default:
		return false, fmt.Errorf("unknown operator '%s'", cond.Operator)
	}
}

func evaluateGroup(group *parser.Group, context map[string]interface{}) (bool, error) {
	var result bool

	if group.LogicalOp == "and" {
		result = true
	} else {
		result = false
	}

	for _, child := range group.Children {
		switch child := child.(type) {
		case *parser.Condition:
			condResult, err := evaluateCondition(child, context)
			if err != nil {
				return false, err
			}

			if group.LogicalOp == "and" {
				result = result && condResult
			} else {
				result = result || condResult
			}
		case *parser.Group:
			groupResult, err := evaluateGroup(child, context)
			if err != nil {
				return false, err
			}

			if group.LogicalOp == "and" {
				result = result && groupResult
			} else {
				result = result || groupResult
			}
		}
	}

	return result, nil
}

func applyNegation(result bool, negate bool) bool {
	if negate {
		return !result
	}

	return result
}

func evaluateIn(fieldValue interface{}, values parser.Value, negate bool) (bool, error) {
	for _, val := range values {
		if fieldValue == val {
			return applyNegation(true, negate), nil
		}
	}

	return applyNegation(false, negate), nil
}

func compareNumeric(fieldValue, conditionValue interface{}, operator string, negate bool) (bool, error) {
	fieldFloat, ok := fieldValue.(float64)
	conditionFloat, ok2 := conditionValue.(float64)

	if !ok || !ok2 {
		return false, fmt.Errorf("invalid types for numeric comparison: %T and %T", fieldValue, conditionValue)
	}

	var result bool
	switch operator {
	case "lt":
		result = fieldFloat < conditionFloat
	case "gt":
		result = fieldFloat > conditionFloat
	case "lte":
		result = fieldFloat <= conditionFloat
	case "gte":
		result = fieldFloat >= conditionFloat
	}

	return applyNegation(result, negate), nil
}

func evaluateBetween(fieldValue interface{}, values parser.Value, negate bool) (bool, error) {
	fieldFloat, ok := fieldValue.(float64)
	low, lowOk := values[0].(float64)
	high, highOk := values[1].(float64)

	if !ok || !lowOk || !highOk {
		return false, fmt.Errorf("invalid types for 'between' operator")
	}

	result := fieldFloat >= low && fieldFloat <= high
	return applyNegation(result, negate), nil
}

func resolveFieldValue(fieldName string, context interface{}) (interface{}, error) {
	re := regexp.MustCompile(`(\w+|\[\d+\])`)
	matches := re.FindAllString(fieldName, -1)

	var current interface{} = context
	for _, match := range matches {
		switch cur := current.(type) {
		case map[string]interface{}:
			current = cur[match]
		case []interface{}:
			if !strings.HasPrefix(match, "[") || !strings.HasSuffix(match, "]") {
				return nil, fmt.Errorf("invalid array index: %s", match)
			}

			value := match[1 : len(match)-1]
			index, err := strconv.Atoi(value)

			if index < 0 || index >= len(cur) {
				return nil, fmt.Errorf("array index out of range: %d", index)
			}

			if err != nil {
				return nil, fmt.Errorf("invalid array index: %s", match)
			}

			current = cur[index]
		default:
			return nil, fmt.Errorf("invalid field path: %s", match)
		}
	}

	return current, nil
}
