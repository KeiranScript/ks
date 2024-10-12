package main

import (
	"fmt"
	"strings"
)

type Compiler struct {
	variables  map[string]int
	code       []string
	labelCount int
	os         string
	arch       string
}

func NewCompiler(os, arch string) *Compiler {
	return &Compiler{
		variables:  make(map[string]int),
		code:       []string{},
		labelCount: 0,
		os:         os,
		arch:       arch,
	}
}

func (c *Compiler) Compile(program *Program) []byte {
	c.generatePrologue()

	for _, stmt := range program.Statements {
		c.compileStatement(stmt)
	}

	c.generateEpilogue()

	return []byte(strings.Join(c.code, "\n"))
}

func (c *Compiler) generatePrologue() {
	if c.os == "linux" {
		c.code = append(c.code, "global _start")
		c.code = append(c.code, "section .text")
		c.code = append(c.code, "_start:")
	} else if c.os == "darwin" {
		c.code = append(c.code, "global start")
		c.code = append(c.code, "section .text")
		c.code = append(c.code, "start:")
	} else if c.os == "windows" {
		c.code = append(c.code, "global main")
		c.code = append(c.code, "extern ExitProcess")
		c.code = append(c.code, "section .text")
		c.code = append(c.code, "main:")
	}
}

func (c *Compiler) generateEpilogue() {
	if c.os == "linux" {
		c.code = append(c.code, "    mov eax, 60")
		c.code = append(c.code, "    xor edi, edi")
		c.code = append(c.code, "    syscall")
	} else if c.os == "darwin" {
		c.code = append(c.code, "    mov rax, 0x2000001")
		c.code = append(c.code, "    xor rdi, rdi")
		c.code = append(c.code, "    syscall")
	} else if c.os == "windows" {
		c.code = append(c.code, "    push 0")
		c.code = append(c.code, "    call ExitProcess")
	}
}

func (c *Compiler) compileStatement(stmt Statement) {
	switch s := stmt.(type) {
	case *AssignmentStatement:
		c.compileAssignment(s)
	case *PrintStatement:
		c.compilePrint(s)
	case *IfStatement:
		c.compileIf(s)
	case *WhileStatement:
		c.compileWhile(s)
	case *Block:
		for _, blockStmt := range s.Statements {
			c.compileStatement(blockStmt)
		}
	}
}

func (c *Compiler) compileAssignment(stmt *AssignmentStatement) {
	value := c.compileExpression(stmt.Value)
	if _, exists := c.variables[stmt.Identifier]; !exists {
		c.variables[stmt.Identifier] = len(c.variables)
	}
	c.code = append(c.code, fmt.Sprintf("    mov [var_%s], %s", stmt.Identifier, value))
}

func (c *Compiler) compilePrint(stmt *PrintStatement) {
	value := c.compileExpression(stmt.Expression)
	c.code = append(c.code, "    push "+value)
	c.code = append(c.code, "    call print_int")
	c.code = append(c.code, "    add esp, 4")
}

func (c *Compiler) compileIf(stmt *IfStatement) {
	condition := c.compileExpression(stmt.Condition)
	elseLabel := c.nextLabel()
	endLabel := c.nextLabel()

	c.code = append(c.code, fmt.Sprintf("    cmp %s, 0", condition))
	c.code = append(c.code, fmt.Sprintf("    je %s", elseLabel))

	for _, thenStmt := range stmt.ThenBranch.Statements {
		c.compileStatement(thenStmt)
	}

	c.code = append(c.code, fmt.Sprintf("    jmp %s", endLabel))
	c.code = append(c.code, fmt.Sprintf("%s:", elseLabel))

	if stmt.ElseBranch != nil {
		for _, elseStmt := range stmt.ElseBranch.Statements {
			c.compileStatement(elseStmt)
		}
	}

	c.code = append(c.code, fmt.Sprintf("%s:", endLabel))
}

func (c *Compiler) compileWhile(stmt *WhileStatement) {
	startLabel := c.nextLabel()
	endLabel := c.nextLabel()

	c.code = append(c.code, fmt.Sprintf("%s:", startLabel))
	condition := c.compileExpression(stmt.Condition)
	c.code = append(c.code, fmt.Sprintf("    cmp %s, 0", condition))
	c.code = append(c.code, fmt.Sprintf("    je %s", endLabel))

	for _, bodyStmt := range stmt.Body.Statements {
		c.compileStatement(bodyStmt)
	}

	c.code = append(c.code, fmt.Sprintf("    jmp %s", startLabel))
	c.code = append(c.code, fmt.Sprintf("%s:", endLabel))
}

func (c *Compiler) compileExpression(expr Expression) string {
	switch e := expr.(type) {
	case *NumberLiteral:
		return fmt.Sprintf("%d", e.Value)
	case *StringLiteral:
		// For simplicity, we'll just store the string's address
		// In a real compiler, you'd need to handle string storage and manipulation
		return fmt.Sprintf("str_%d", len(c.variables))
	case *BooleanLiteral:
		if e.Value {
			return "1"
		}
		return "0"
	case *Identifier:
		return fmt.Sprintf("[var_%s]", e.Name)
	case *BinaryExpression:
		left := c.compileExpression(e.Left)
		right := c.compileExpression(e.Right)
		c.code = append(c.code, fmt.Sprintf("    mov eax, %s", left))
		switch e.Operator {
		case TOKEN_PLUS:
			c.code = append(c.code, fmt.Sprintf("    add eax, %s", right))
		case TOKEN_MINUS:
			c.code = append(c.code, fmt.Sprintf("    sub eax, %s", right))
		case TOKEN_MULTIPLY:
			c.code = append(c.code, fmt.Sprintf("    imul eax, %s", right))
		case TOKEN_DIVIDE:
			c.code = append(c.code, fmt.Sprintf("    xor edx, edx"))
			c.code = append(c.code, fmt.Sprintf("    div %s", right))
		case TOKEN_EQ:
			c.code = append(c.code, fmt.Sprintf("    cmp eax, %s", right))
			c.code = append(c.code, "    sete al")
			c.code = append(c.code, "    movzx eax, al")
		case TOKEN_NEQ:
			c.code = append(c.code, fmt.Sprintf("    cmp eax, %s", right))
			c.code = append(c.code, "    setne al")
			c.code = append(c.code, "    movzx eax, al")
		case TOKEN_LT:
			c.code = append(c.code, fmt.Sprintf("    cmp eax, %s", right))
			c.code = append(c.code, "    setl al")
			c.code = append(c.code, "    movzx eax, al")
		case TOKEN_GT:
			c.code = append(c.code, fmt.Sprintf("    cmp eax, %s", right))
			c.code = append(c.code, "    setg al")
			c.code = append(c.code, "    movzx eax, al")
		case TOKEN_AND:
			endLabel := c.nextLabel()
			c.code = append(c.code, "    cmp eax, 0")
			c.code = append(c.code, fmt.Sprintf("    je %s", endLabel))
			c.code = append(c.code, fmt.Sprintf("    cmp %s, 0", right))
			c.code = append(c.code, fmt.Sprintf("    setne al"))
			c.code = append(c.code, "    movzx eax, al")
			c.code = append(c.code, fmt.Sprintf("%s:", endLabel))
		case TOKEN_OR:
			trueLabel := c.nextLabel()
			endLabel := c.nextLabel()
			c.code = append(c.code, "    cmp eax, 0")
			c.code = append(c.code, fmt.Sprintf("    jne %s", trueLabel))
			c.code = append(c.code, fmt.Sprintf("    cmp %s, 0", right))
			c.code = append(c.code, fmt.Sprintf("    setne al"))
			c.code = append(c.code, "    movzx eax, al")
			c.code = append(c.code, fmt.Sprintf("    jmp %s", endLabel))
			c.code = append(c.code, fmt.Sprintf("%s:", trueLabel))
			c.code = append(c.code, "    mov eax, 1")
			c.code = append(c.code, fmt.Sprintf("%s:", endLabel))
		default:
			panic("ill add smth here later")
		}
		return "eax"
	default:
		panic("Unknown expression type")
	}
}

func (c *Compiler) nextLabel() string {
	c.labelCount++
	return fmt.Sprintf("label_%d", c.labelCount)
}
