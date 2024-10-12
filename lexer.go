package main

import (
	"fmt"
	"unicode"
)

type TokenType int

const (
	TOKEN_EOF TokenType = iota
	TOKEN_IDENT
	TOKEN_NUMBER
	TOKEN_STRING
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_MULTIPLY
	TOKEN_DIVIDE
	TOKEN_ASSIGN
	TOKEN_SEMICOLON
	TOKEN_PRINT
	TOKEN_IF
	TOKEN_ELSE
	TOKEN_WHILE
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_EQ
	TOKEN_NEQ
	TOKEN_LT
	TOKEN_GT
	TOKEN_AND
	TOKEN_OR
	TOKEN_TRUE
	TOKEN_FALSE
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

func (l *Lexer) Tokenize() []Token {
	var tokens []Token

	for l.pos < len(l.input) {
		switch c := l.input[l.pos]; {
		case unicode.IsSpace(rune(c)):
			l.pos++
		case unicode.IsLetter(rune(c)):
			tokens = append(tokens, l.readIdentifier())
		case unicode.IsDigit(rune(c)):
			tokens = append(tokens, l.readNumber())
		case c == '"':
			tokens = append(tokens, l.readString())
		case c == '+':
			tokens = append(tokens, Token{Type: TOKEN_PLUS, Value: "+"})
			l.pos++
		case c == '-':
			tokens = append(tokens, Token{Type: TOKEN_MINUS, Value: "-"})
			l.pos++
		case c == '*':
			tokens = append(tokens, Token{Type: TOKEN_MULTIPLY, Value: "*"})
			l.pos++
		case c == '/':
			tokens = append(tokens, Token{Type: TOKEN_DIVIDE, Value: "/"})
			l.pos++
		case c == '=':
			if l.peekNext() == '=' {
				tokens = append(tokens, Token{Type: TOKEN_EQ, Value: "=="})
				l.pos += 2
			} else {
				tokens = append(tokens, Token{Type: TOKEN_ASSIGN, Value: "="})
				l.pos++
			}
		case c == '!':
			if l.peekNext() == '=' {
				tokens = append(tokens, Token{Type: TOKEN_NEQ, Value: "!="})
				l.pos += 2
			} else {
				panic(fmt.Sprintf("Unexpected character: %c", c))
			}
		case c == '<':
			tokens = append(tokens, Token{Type: TOKEN_LT, Value: "<"})
			l.pos++
		case c == '>':
			tokens = append(tokens, Token{Type: TOKEN_GT, Value: ">"})
			l.pos++
		case c == '&':
			if l.peekNext() == '&' {
				tokens = append(tokens, Token{Type: TOKEN_AND, Value: "&&"})
				l.pos += 2
			} else {
				panic(fmt.Sprintf("Unexpected character: %c", c))
			}
		case c == '|':
			if l.peekNext() == '|' {
				tokens = append(tokens, Token{Type: TOKEN_OR, Value: "||"})
				l.pos += 2
			} else {
				panic(fmt.Sprintf("Unexpected character: %c", c))
			}
		case c == ';':
			tokens = append(tokens, Token{Type: TOKEN_SEMICOLON, Value: ";"})
			l.pos++
		case c == '(':
			tokens = append(tokens, Token{Type: TOKEN_LPAREN, Value: "("})
			l.pos++
		case c == ')':
			tokens = append(tokens, Token{Type: TOKEN_RPAREN, Value: ")"})
			l.pos++
		case c == '{':
			tokens = append(tokens, Token{Type: TOKEN_LBRACE, Value: "{"})
			l.pos++
		case c == '}':
			tokens = append(tokens, Token{Type: TOKEN_RBRACE, Value: "}"})
			l.pos++
		default:
			panic(fmt.Sprintf("Unknown character: %c", c))
		}
	}

	tokens = append(tokens, Token{Type: TOKEN_EOF})
	return tokens
}

func (l *Lexer) readIdentifier() Token {
	start := l.pos
	for l.pos < len(l.input) && (unicode.IsLetter(rune(l.input[l.pos])) || unicode.IsDigit(rune(l.input[l.pos]))) {
		l.pos++
	}
	value := l.input[start:l.pos]
	switch value {
	case "print":
		return Token{Type: TOKEN_PRINT, Value: value}
	case "if":
		return Token{Type: TOKEN_IF, Value: value}
	case "else":
		return Token{Type: TOKEN_ELSE, Value: value}
	case "while":
		return Token{Type: TOKEN_WHILE, Value: value}
	case "true":
		return Token{Type: TOKEN_TRUE, Value: value}
	case "false":
		return Token{Type: TOKEN_FALSE, Value: value}
	default:
		return Token{Type: TOKEN_IDENT, Value: value}
	}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
		l.pos++
	}
	return Token{Type: TOKEN_NUMBER, Value: l.input[start:l.pos]}
}

func (l *Lexer) readString() Token {
	l.pos++ // Skip opening quote
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		l.pos++
	}
	if l.pos == len(l.input) {
		panic("Unterminated string")
	}
	value := l.input[start:l.pos]
	l.pos++ // Skip closing quote
	return Token{Type: TOKEN_STRING, Value: value}
}

func (l *Lexer) peekNext() byte {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}
