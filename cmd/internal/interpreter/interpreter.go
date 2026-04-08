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
	Where      string
	SourceLine string
}

func (ie *InterpreterError) PrintError() {
	fmt.Printf("Error: %s [%d:%d] at %s\n", ie.Message, ie.Line, ie.Col+1, ie.Where)
	fmt.Println(ie.SourceLine)
	fmt.Printf("%s^\n", strings.Repeat(" ", ie.Col))
}

type Interpreter struct {
	ie           *InterpreterError
	previousData []byte
}

func NewInterpreter() Interpreter {
	return Interpreter{
		previousData: make([]byte, 0),
	}
}

func (i *Interpreter) Exec(data []byte) error {
	scanner := NewScanner(i, data)
	tokens := scanner.ScanTokens()

	if i.ie != nil {
		i.ie.PrintError()
		os.Exit(65)
	}

	for i, t := range tokens {
		fmt.Printf("Idx[%d]: %s\n", i, t)
	}

	return nil
}

func (i *Interpreter) ExecInteractive(data []byte) {
	oldLen := len(i.previousData)
	i.previousData = append(i.previousData, data...)

	scanner := NewScanner(i, i.previousData)
	tokens := scanner.ScanTokens()

	if i.ie != nil {
		i.ie.PrintError()
		i.previousData = i.previousData[:oldLen]
	}

	for _, t := range tokens {
		fmt.Println(t)
	}
}
