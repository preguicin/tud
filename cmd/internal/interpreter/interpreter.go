package interpreter

import (
	"fmt"
	"os"
	"strings"
)

type InterpreterError struct {
	Line       int
	Col        int
	Message    string
	Where      uint8
	SourceLine string
}

func (ie *InterpreterError) PrintError() {
	fmt.Printf("Error: %s [%d:%d] at %c\n", ie.Message, ie.Line, ie.Col+1, ie.Where)
	fmt.Println(ie.SourceLine)
	fmt.Printf("%s^\n", strings.Repeat(" ", ie.Col))
}

type Interpreter struct {
	ie *InterpreterError
}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i *Interpreter) Exec(data []byte) error {
	scanner := NewScanner(i, data)
	tokens := scanner.ScanTokens()

	if i.ie != nil {
		i.ie.PrintError()
		os.Exit(65)
	}

	for _, t := range tokens {
		fmt.Println(t)
	}

	return nil
}
