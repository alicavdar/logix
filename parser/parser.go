package parser

import (
	"fmt"
	"strconv"

	"github.com/alicavdar/logix/lexer"
)

type SingleValue interface{}
type Value []SingleValue

type Condition struct {
	Field    string
	Operator string
	Value    Value
	Negate   bool
}

type Group struct {
	LogicalOp string        // "and" or "or"
	Children  []interface{} // can be either Condition or Group (for nested groups)
}

type Parser struct {
	lexer     *lexer.Lexer
	currToken lexer.Token
	peekToken lexer.Token
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	return p
}

func (p *Parser) parseCondition() *Condition {
	field := p.currToken.Lexeme
	p.nextToken()

	negate := false
	if p.currToken.Kind == lexer.NOT {
		negate = true
		p.nextToken()
	}

	operator := p.currToken.Lexeme
	if negate && !allowedNegateSuffix(operator) {
		panic(fmt.Sprintf("Negation is not supported for operator: %s", operator))
	}
	p.nextToken()

	var value Value
	if operator == "between" {
		value = p.parseRange()
	} else if p.currToken.Kind == lexer.LSQUARE {
		value = p.parseArray()
	} else {
		value = append(value, parseValue(p.currToken))
	}
	return &Condition{Field: field, Operator: operator, Value: value, Negate: negate}
}

func (p *Parser) parseGroup() *Group {
	p.nextToken()

	if p.currToken.Kind != lexer.AND && p.currToken.Kind != lexer.OR {
		panic(fmt.Sprintf("Expected logical operator 'and' or 'or', got: %s", p.currToken.Lexeme))
	}

	group := &Group{LogicalOp: p.currToken.Lexeme}
	p.nextToken()

	for p.currToken.Kind != lexer.DEDENT && p.currToken.Kind != lexer.EOF {
		switch p.currToken.Kind {
		case lexer.GROUP:
			group.Children = append(group.Children, p.parseGroup())
		case lexer.IDENT:
			group.Children = append(group.Children, p.parseCondition())
		}

		p.nextToken()
	}

	return group
}

func (p *Parser) parseRange() Value {
	var value Value
	value = append(value, parseValue(p.currToken))
	p.nextToken()

	if p.currToken.Kind != lexer.AND {
		panic("Expected 'and' in range condition")
	}
	p.nextToken()

	value = append(value, parseValue(p.currToken))
	return value
}

func (p *Parser) nextToken() {
	p.currToken = p.lexer.Next()

	for p.currToken.Kind == lexer.INDENT {
		p.currToken = p.lexer.Next()
	}
}

func (p *Parser) ParseNext() interface{} {
	if p.currToken.Kind == lexer.EOF {
		return nil
	}

	var result interface{}

	switch p.currToken.Kind {
	case lexer.GROUP:
		result = p.parseGroup()
	case lexer.IDENT:
		result = p.parseCondition()
	default:
		panic(fmt.Sprintf("Unexpected token: '%s' of kind '%s'", p.currToken.Lexeme, p.currToken.Kind))
	}

	p.nextToken()

	return result
}

func (p *Parser) parseArray() Value {
	p.nextToken()

	var arrayValues Value
	for p.currToken.Kind != lexer.RSQUARE {
		arrayValues = append(arrayValues, parseValue(p.currToken))

		p.nextToken()

		if p.currToken.Kind == lexer.COMMA {
			p.nextToken()
		}
	}

	if p.currToken.Kind != lexer.RSQUARE {
		panic("Expected ']' to close array")
	}

	return arrayValues
}

func parseValue(token lexer.Token) SingleValue {
	var value SingleValue

	switch token.Kind {
	case lexer.TRUE:
		value = true
	case lexer.FALSE:
		value = false
	case lexer.NIL:
		value = nil
	case lexer.NUMBER:
		number, err := strconv.ParseFloat(token.Lexeme, 64)
		if err != nil {
			panic(fmt.Sprintf("Invalid number literal: %v", err))
		}

		value = number
	case lexer.STRING:
		value = token.Lexeme
	}

	return value
}

func allowedNegateSuffix(op string) bool {
	switch op {
	case "in", "contains", "between", "startsWith", "endsWith":
		return true
	default:
		return false
	}
}
