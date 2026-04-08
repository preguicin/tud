package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tud/cmd/internal/interpreter"
)

type InteractiveManager struct {
	ShouldClose bool
	Message     string
	Err         error
}

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: tud [script]")
	} else if len(args) == 1 {
		bytes, err := os.ReadFile(args[0])

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		i := interpreter.NewInterpreter()
		i.Exec(bytes)

	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Starting TUD session...")

		i := interpreter.NewInterpreter()
		for {
			fmt.Print("> ")
			data, err := reader.ReadBytes('\n')

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if handleInteractiveExitInputs(data) {
				break
			}

			i.ExecInteractive(data)
		}
	}
}

func handleInteractiveExitInputs(data []byte) bool {
	if len(data) > 6 {
		return false
	}

	res := string(data)
	res = strings.TrimSpace(res)

	if res == "q" || res == "quit" || res == "exit" {
		return true
	}

	return false
}
