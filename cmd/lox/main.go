package main

import (
	"badparser/cmd/lox/scanner"
	"bufio"
	"fmt"
	"log"
	"os"
)

type Lox struct {
	hadError bool
}

func (l *Lox) Run(source string) {
	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()

	for i := range tokens {
		token := tokens[i]
		print(token)
	}

}

func (l *Lox) RunFile(file_path string) {
	bytes, err := os.ReadFile(file_path)
	if err != nil {
		log.Fatal(err)
	}
	l.Run(string(bytes))

	if l.hadError {
		os.Exit(65)
	}
}

func (l *Lox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if line == "q\n" || line == "exit\n" {
			break
		}
		l.Run(line)
		l.hadError = false
	}
}

func main() {
	args := os.Args[1:]
	lox := &Lox{}

	if len(args) > 1 {
		fmt.Println("Usage: pass a badp script [script.badp or script]")
		os.Exit(64)
	} else if len(args) == 1 {
		fmt.Println("Running file..")
		lox.RunFile(args[0])
	} else {
		fmt.Println("Starting interactive mode..")
		lox.RunPrompt()
	}
}
