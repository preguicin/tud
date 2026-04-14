package interpreter

import (
	"fmt"
	"tud/cmd/internal/file"
	"tud/cmd/internal/interpreter/error"
	"tud/cmd/internal/interpreter/lexer"
)

type Interpreter struct {
	file.File
}

func NewInterpreter(file file.File) Interpreter {
	return Interpreter{
		File: file,
	}
}

func (i *Interpreter) Exec() {
	data := i.File.GetBytes()
	errReporter := error.NewErrorReporter(data)

	scanner := lexer.NewScanner(i.File, errReporter.NewError)
	tokens := scanner.ScanTokens()

	if errReporter.HasErrors() {

	}

	for i, t := range tokens {
		fmt.Printf("Idx[%d]: %s\n", i, t)
	}
}
