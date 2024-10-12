package main

import (
	"strconv"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) Parse() *Program {
	program := &Program{}
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != TOKEN_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.tokens[p.pos].Type {
	case TOKEN_IDENT:
		return p.parseAssignment()
	case TOKEN_PRINT:
		return p.parsePrintStatement()
	case TOKEN_IF:
		return p.parseIfStatement()
	case TOKEN_WHILE:
		return p.parseWhileStatement()
	default:
		return nil
	}
}

func (p *Parser) parseAssignment() *AssignmentStatement {
	ident := p.tokens[p.pos].Value
	p.pos++
	if p.tokens[p.pos].Type != TOKEN_ASSIGN {
		panic("Expected '=' after identifier in assignment")
	}
	p.pos++
	expr := p.parseExpression()
	if p.tokens[p.pos].Type != TOKEN_SEMICOLON {
		panic("Expected ';' at end of assignment")
	}
	p.pos++
	return &AssignmentStatement{Identifier: ident, Value: expr}
}

func (p *Parser) parsePrintStatement() *PrintStatement {
	p.pos++ // Skip 'print' token
	expr := p.parseExpression()
	if p.tokens[p.pos].Type != TOKEN_SEMICOLON {
		panic("Expected ';' at end of print statement")
	}
	p.pos++
	return &PrintStatement{Expression: expr}
}

func (p *Parser) parseIfStatement() *IfStatement {
	p.pos++ // Skip 'if' token
	if p.tokens[p.pos].Type != TOKEN_LPAREN {
		panic("Expected '(' after 'if'")
	}
	p.pos++
	condition := p.parseExpression()
	if p.tokens[p.pos].Type != TOKEN_RPAREN {
		panic("Expected ')' after if condition")
	}
	p.pos++
	thenBranch := p.parseBlock()
	var elseBranch *Block
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == TOKEN_ELSE {
		p.pos++
		elseBranch = p.parseBlock()
	}
	return &IfStatement{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) parseWhileStatement() *WhileStatement {
	p.pos++ // Skip 'while' token
	if p.tokens[p.pos].Type != TOKEN_LPAREN {
		panic("Expected '(' after 'while'")
	}
	p.pos++
	condition := p.parseExpression()
	if p.tokens[p.pos].Type != TOKEN_RPAREN {
		panic("Expected ')' after while condition")
	}
	p.pos++
	body := p.parseBlock()
	return &WhileStatement{Condition: condition, Body: body}
}

func (p *Parser) parseBlock() *Block {
	if p.tokens[p.pos].Type != TOKEN_LBRACE {
		panic("Expected '{' at start of block")
	}
	p.pos++
	block := &Block{}
	for p.tokens[p.pos].Type != TOKEN_RBRACE {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}
	p.pos++ // Skip closing '}'
	return block
}

func (p *Parser) parseExpression() Expression {
	return p.parseLogicalOr()
}

func (p *Parser) parseLogicalOr() Expression {
	expr := p.parseLogicalAnd()
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == TOKEN_OR {
		op := p.tokens[p.pos].Type
		p.pos++
		right := p.parseLogicalAnd()
		expr = &BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *Parser) parseLogicalAnd() Expression {
	expr := p.parseEquality()
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == TOKEN_AND {
		op := p.tokens[p.pos].Type
		p.pos++
		right := p.parseEquality()
		expr = &BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *Parser) parseEquality() Expression {
	expr := p.parseComparison()
	for p.pos < len(p.tokens) && (p.tokens[p.pos].Type == TOKEN_EQ || p.tokens[p.pos].Type == TOKEN_NEQ) {
		op := p.tokens[p.pos].Type
		p.pos++
		right := p.parseComparison()
		expr = &BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *Parser) parseComparison() Expression {
	expr := p.parseAdditive()
	for p.pos < len(p.tokens) && (p.tokens[p.pos].Type == TOKEN_LT || p.tokens[p.pos].Type == TOKEN_GT) {
		op := p.tokens[p.pos].Type
		p.pos++
		right := p.parseAdditive()
		expr = &BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *Parser) parseAdditive() Expression {
	expr := p.parseMultiplicative()
	for p.pos < len(p.tokens) && (p.tokens[p.pos].Type == TOKEN_PLUS || p.tokens[p.pos].Type == TOKEN_MINUS) {
		op := p.tokens[p.pos].Type
		p.pos++
		right := p.parseMultiplicative()
		expr = &BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *Parser) parseMultiplicative() Expression {
	expr := p.parsePrimary()
	for p.pos < len(p.tokens) && (p.tokens[p.pos].Type == TOKEN_MULTIPLY || p.tokens[p.pos].Type == TOKEN_DIVIDE) {
		op := p.tokens[p.pos].Type
		p.pos++
		right := p.parsePrimary()
		expr = &BinaryExpression{Left: expr,

			Operator: op, Right: right}
	}
	return expr
}

func (p *Parser) parsePrimary() Expression {
	switch p.tokens[p.pos].Type {
	case TOKEN_NUMBER:
		value, _ := strconv.Atoi(p.tokens[p.pos].Value)
		p.pos++
		return &NumberLiteral{Value: value}
	case TOKEN_STRING:
		value := p.tokens[p.pos].Value
		p.pos++
		return &StringLiteral{Value: value}
	case TOKEN_TRUE:
		p.pos++
		return &BooleanLiteral{Value: true}
	case TOKEN_FALSE:
		p.pos++
		return &BooleanLiteral{Value: false}
	case TOKEN_IDENT:
		ident := p.tokens[p.pos].Value
		p.pos++
		return &Identifier{Name: ident}
	case TOKEN_LPAREN:
		p.pos++
		expr := p.parseExpression()
		if p.tokens[p.pos].Type != TOKEN_RPAREN {
			panic("Expected ')' after expression")
		}
		p.pos++
		return expr
	default:
		panic("Unexpected token in expression")
	}
}
