package interpreter

import (
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

func (i *Interpreter) Exec() ([]lexer.Token, []string) {
	data := i.File.GetBytes()
	errReporter := error.NewErrorReporter(data)

	scanner := lexer.NewScanner(i.File, errReporter.NewError)
	tokens := scanner.ScanTokens()

	if errReporter.HasErrors() {
		return nil, errReporter.GetErrorsText()
	}

	return tokens, nil
}
