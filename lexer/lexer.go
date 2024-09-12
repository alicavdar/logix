package lexer

import (
	"strings"
	"unicode"
)

type TokenKind string

const (
	EOF         TokenKind = "EOF"
	IDENT       TokenKind = "IDENT"
	EQ          TokenKind = "EQ"
	NEQ         TokenKind = "NEQ"
	GT          TokenKind = "GT"
	GTE         TokenKind = "GTE"
	LT          TokenKind = "LT"
	LTE         TokenKind = "LTE"
	CONTAINS    TokenKind = "CONTAINS"
	BETWEEN     TokenKind = "BETWEEN"
	IN          TokenKind = "IN"
	NOT         TokenKind = "NOT"
	STRING      TokenKind = "STRING"
	NUMBER      TokenKind = "NUMBER"
	NIL         TokenKind = "NIL"
	STARTS_WITH TokenKind = "STARTS_WITH"
	ENDS_WITH   TokenKind = "ENDS_WITH"
	LSQUARE     TokenKind = "LSQUARE"
	RSQUARE     TokenKind = "RSQUARE"
	COMMA       TokenKind = "COMMA"
	GROUP       TokenKind = "GROUP"
	AND         TokenKind = "AND"
	OR          TokenKind = "OR"
	INDENT      TokenKind = "INDENT"
	DEDENT      TokenKind = "DEDENT"
	ILLEGAL     TokenKind = "ILLEGAL"
	TRUE        TokenKind = "TRUE"
	FALSE       TokenKind = "FALSE"
)

var keywords = map[string]TokenKind{
	"eq":         EQ,
	"neq":        NEQ,
	"gt":         GT,
	"lt":         LT,
	"gte":        GTE,
	"lte":        LTE,
	"contains":   CONTAINS,
	"between":    BETWEEN,
	"not":        NOT,
	"nil":        NIL,
	"startsWith": STARTS_WITH,
	"endsWith":   ENDS_WITH,
	"in":         IN,
	"true":       TRUE,
	"false":      FALSE,
	"group":      GROUP,
	"and":        AND,
	"or":         OR,
}

type Token struct {
	Kind   TokenKind
	Lexeme string
}

type Lexer struct {
	input        string // the entire input string being lexed
	position     int    // current position (points to the current char)
	readPosition int    // the next position (used for lookahead)
	ch           rune   // the current char being processed
	indentWidth  int    // the number of spaces or tabs that represent one level of indentation
	useSpaces    bool   // whether the input is using spaces for indentation (spaces: true, tabs: false)
	indentStack  []int  // stack to track the current indentation levels (used for handling nested blocks)
	dedentCount  int    // count of DEDENT tokens pending to be emitted (after reducing indentation levels)
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:       input,
		indentWidth: -1,
		indentStack: []int{0},
		dedentCount: 0,
	}

	l.readRune()

	return l
}

func (l *Lexer) Next() Token {
	l.setIndentationMode()

	if l.dedentCount > 0 {
		l.dedentCount--
		return l.newToken(DEDENT, "")
	}

	if l.ch == '\n' {
		indentLevel := l.readIndentLevel()

		if indentLevel > l.indentStack[len(l.indentStack)-1] {
			l.indentStack = append(l.indentStack, indentLevel)
			return l.newToken(INDENT, "")
		} else if indentLevel < l.indentStack[len(l.indentStack)-1] {
			for indentLevel < l.indentStack[len(l.indentStack)-1] {
				l.indentStack = l.indentStack[:len(l.indentStack)-1]
				l.dedentCount++
			}
			l.dedentCount--
			return l.newToken(DEDENT, "")
		}
	}

	l.skipWhitespace()

	// Logix only supports comments with # and we ignore all comments here
	if l.ch == '#' {
		for l.ch != '\n' && l.ch != 0 {
			l.readRune()
		}

		return l.Next()
	}

	if l.ch == '"' {
		return l.readStringToken()
	} else if l.ch == '[' {
		l.readRune()
		return l.newToken(LSQUARE, "[")
	} else if l.ch == ']' {
		l.readRune()
		return l.newToken(RSQUARE, "]")
	} else if l.ch == ',' {
		l.readRune()
		return l.newToken(COMMA, ",")
	} else if l.isAlpha(l.ch) {
		var lexeme = l.readLexeme()
		return l.newToken(l.lookupKeyword(lexeme), lexeme)
	} else if l.isDigit(l.ch) {
		return l.newToken(NUMBER, l.readNumber())
	} else if l.ch == 0 {
		return l.newToken(EOF, "")
	} else {
		tok := l.newToken(ILLEGAL, string(l.ch))
		l.readRune()
		return tok
	}
}

func (l *Lexer) readStringToken() Token {
	var builder strings.Builder

	l.readRune() // Skip the opening quote

	for l.ch != '"' && l.ch != 0 {
		builder.WriteRune(l.ch)
		l.readRune()
	}

	if l.ch == '"' {
		l.readRune() // Skip the closing quote
		return l.newToken(STRING, builder.String())
	}

	return l.newToken(ILLEGAL, "Unclosed string")
}

func (l *Lexer) setIndentationMode() {
	var isIndent bool = l.ch == '\n' &&
		l.peek() != 0 &&
		(l.peek() == ' ' || l.peek() == '\t')

	if l.indentWidth != -1 || !isIndent {
		return
	}

	if l.input[l.readPosition] == '\t' {
		l.useSpaces = false
		l.indentWidth = 1
	} else {
		var i int
		var pos int = l.readPosition
		for {
			if unicode.IsSpace(rune(l.input[pos])) {
				i++
				pos++
			} else {
				break
			}
		}
		l.useSpaces = true
		l.indentWidth = i
	}
}

func (l *Lexer) peek() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return rune(l.input[l.readPosition])
}

func (l *Lexer) readIndentLevel() int {
	indentCount := 0

	var ch rune
	if l.useSpaces {
		ch = ' '
	} else {
		ch = '\t'
	}

	for l.peek() == ch {
		indentCount += 1
		l.readRune()
	}

	if l.useSpaces {
		return indentCount / l.indentWidth
	}

	return indentCount
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' ||
		l.ch == '\t' ||
		l.ch == '\n' ||
		l.ch == '\r' {
		l.readRune()
	}
}

func (l *Lexer) newToken(tokenKind TokenKind, lexeme string) Token {
	return Token{Kind: tokenKind, Lexeme: lexeme}
}

func (l *Lexer) lookupKeyword(lexeme string) TokenKind {
	if tok, ok := keywords[lexeme]; ok {
		return tok
	}

	return IDENT
}

func (l *Lexer) readRune() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = rune(l.input[l.readPosition])
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readNumber() string {
	position := l.position
	for l.isDigit(l.ch) {
		l.readRune()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readLexeme() string {
	position := l.position

	for l.isAlphaNumeric(l.ch) || l.ch == '.' || l.ch == '[' || l.ch == ']' {
		l.readRune()
	}

	return l.input[position:l.position]
}

func (l *Lexer) isAlphaNumeric(ch rune) bool {
	return l.isAlpha(ch) || l.isDigit(ch)
}

func (l *Lexer) isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) isAlpha(ch rune) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_'
}
