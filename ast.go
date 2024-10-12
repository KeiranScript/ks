package main

type Node interface{}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type AssignmentStatement struct {
	Identifier string
	Value      Expression
}

func (as *AssignmentStatement) statementNode() {}

type PrintStatement struct {
	Expression Expression
}

func (ps *PrintStatement) statementNode() {}

type IfStatement struct {
	Condition  Expression
	ThenBranch *Block
	ElseBranch *Block
}

func (is *IfStatement) statementNode() {}

type WhileStatement struct {
	Condition Expression
	Body      *Block
}

func (ws *WhileStatement) statementNode() {}

type Block struct {
	Statements []Statement
}

func (b *Block) statementNode() {}

type NumberLiteral struct {
	Value int
}

func (nl *NumberLiteral) expressionNode() {}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}

type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}

type Identifier struct {
	Name string
}

func (i *Identifier) expressionNode() {}

type BinaryExpression struct {
	Left     Expression
	Operator TokenType
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}
