package parser

import (
	"testing"

	"github.com/alicavdar/logix/lexer"
)

func newTestParser(input string) *Parser {
	l := lexer.NewLexer(input)
	return NewParser(l)
}

func TestParseCondition(t *testing.T) {
	input := `
field1 eq 10
field2 eq "Hello"
field3 neq 20
field4 gt 10
field5 lt 100
field6 gte 500
field7 lte 10
field8 contains "hello"
field9 between 10 and 50
field10 not contains "hello"
field11 eq nil
field12 startsWith "hel"
field13 endsWith "lo"
field14 in [1, 2, 3]
field14 in ["lorem", 1, 2]
field15 eq true
field16 eq false
`
	p := newTestParser(input)

	assertCondition(t, p.ParseNext(), "field1", "eq", Value{10.0}, false)
	assertCondition(t, p.ParseNext(), "field2", "eq", Value{"Hello"}, false)
	assertCondition(t, p.ParseNext(), "field3", "neq", Value{20.0}, false)
	assertCondition(t, p.ParseNext(), "field4", "gt", Value{10.0}, false)
	assertCondition(t, p.ParseNext(), "field5", "lt", Value{100.0}, false)
	assertCondition(t, p.ParseNext(), "field6", "gte", Value{500.0}, false)
	assertCondition(t, p.ParseNext(), "field7", "lte", Value{10.0}, false)
	assertCondition(t, p.ParseNext(), "field8", "contains", Value{"hello"}, false)
	assertCondition(t, p.ParseNext(), "field9", "between", Value{10.0, 50.0}, false)
	assertCondition(t, p.ParseNext(), "field10", "contains", Value{"hello"}, true)
	assertCondition(t, p.ParseNext(), "field11", "eq", Value{nil}, false)
	assertCondition(t, p.ParseNext(), "field12", "startsWith", Value{"hel"}, false)
	assertCondition(t, p.ParseNext(), "field13", "endsWith", Value{"lo"}, false)
	assertCondition(t, p.ParseNext(), "field14", "in", Value{1.0, 2.0, 3.0}, false)
	assertCondition(t, p.ParseNext(), "field14", "in", Value{"lorem", 1.0, 2.0}, false)
	assertCondition(t, p.ParseNext(), "field15", "eq", Value{true}, false)
	assertCondition(t, p.ParseNext(), "field16", "eq", Value{false}, false)
}

func TestParseGroup(t *testing.T) {
	input := `
group and
    field1 eq 10
    field2 gt 20
`
	p := newTestParser(input)

	result := p.ParseNext()

	group := assertGroup(t, result, "and", 2)
	assertCondition(t, group.Children[0], "field1", "eq", Value{10.0}, false)
	assertCondition(t, group.Children[1], "field2", "gt", Value{20.0}, false)
}

func TestParseNestedGroups(t *testing.T) {
	input := `
group or
    group and
        field1 eq 10
        field2 lt 20
    field3 eq 30
    group or
        field4 lt 10
        field5 neq 10
field_6 eq "10"
`
	p := newTestParser(input)

	groupResult := p.ParseNext()

	group := assertGroup(t, groupResult, "or", 3)

	nestedGroup := assertGroup(t, group.Children[0], "and", 2)
	assertCondition(t, nestedGroup.Children[0], "field1", "eq", Value{10.0}, false)
	assertCondition(t, nestedGroup.Children[1], "field2", "lt", Value{20.0}, false)

	assertCondition(t, group.Children[1], "field3", "eq", Value{30.0}, false)

	nestedGroup2 := assertGroup(t, group.Children[2], "or", 2)
	assertCondition(t, nestedGroup2.Children[0], "field4", "lt", Value{10.0}, false)
	assertCondition(t, nestedGroup2.Children[1], "field5", "neq", Value{10.0}, false)

	singleCondition := p.ParseNext()
	assertCondition(t, singleCondition, "field_6", "eq", Value{"10"}, false)
}

func TestNegateCondition(t *testing.T) {
	input := `
title not contains "Berlin"
`
	p := newTestParser(input)
	result := p.ParseNext()

	assertCondition(t, result, "title", "contains", Value{"Berlin"}, true)
}

func TestBetween(t *testing.T) {
	input := `
age between 10 and 30
`
	p := newTestParser(input)
	result := p.ParseNext()

	assertCondition(t, result, "age", "between", Value{10.0, 30.0}, false)
}

func slicesEqual(a, b Value) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		switch v := a[i].(type) {
		case float64:
			if bv, ok := b[i].(float64); !ok || v != bv {
				return false
			}
		case string:
			if bv, ok := b[i].(string); !ok || v != bv {
				return false
			}
		case bool:
			if bv, ok := b[i].(bool); !ok || v != bv {
				return false
			}
		case nil:
			if b[i] != nil {
				return false
			}
		default:
			return false
		}
	}

	return true
}

func assertCondition(t *testing.T, result interface{}, expectedField, expectedOperator interface{}, expectedValue Value, expectedNegate bool) *Condition {
	condition, ok := result.(*Condition)
	if !ok {
		t.Fatalf("Expected *Condition, got %T", result)
	}

	if condition.Field != expectedField {
		t.Errorf("Expected field %s, got %s", expectedField, condition.Field)
	}
	if condition.Operator != expectedOperator {
		t.Errorf("Expected operator %s, got %s", expectedOperator, condition.Operator)
	}
	if !slicesEqual(condition.Value, expectedValue) {
		t.Errorf("Expected value %v, got %v", expectedValue, condition.Value)
	}
	if condition.Negate != expectedNegate {
		t.Errorf("Expected value %t, got %t", expectedNegate, condition.Negate)
	}

	return condition
}

func assertGroup(t *testing.T, result interface{}, expectedLogicalOp string, expectedChildrenCount int) *Group {
	group, ok := result.(*Group)
	if !ok {
		t.Fatalf("Expected *Group, got %T", result)
	}

	if group.LogicalOp != expectedLogicalOp {
		t.Errorf("Expected logical operator %s, got %s", expectedLogicalOp, group.LogicalOp)
	}

	if len(group.Children) != expectedChildrenCount {
		t.Fatalf("Expected %d children, got %d", expectedChildrenCount, len(group.Children))
	}

	return group
}
