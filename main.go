package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: keiranscript <filename.ks>")
		os.Exit(1)
	}

	filename := os.Args[1]
	input, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	lexer := NewLexer(string(input))
	tokens := lexer.Tokenize()

	parser := NewParser(tokens)
	ast := parser.Parse()

	compiler := NewCompiler(runtime.GOOS, runtime.GOARCH)
	machineCode := compiler.Compile(ast)

	outputFilename := filename[:len(filename)-3] + ".asm"
	err = os.WriteFile(outputFilename, machineCode, 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Compiled successfully. Assembly output written to %s\n", outputFilename)
	fmt.Println("To create an executable, use an assembler appropriate for your OS and architecture.")
}
