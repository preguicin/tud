package main

import "log"

func Error(line int, message string, l *Lox) {
	report(line, "", message, l)
}

func report(line int, where string, message string, l *Lox) {
	log.Printf("[line + %d] Error %s: %s", line, where, message)
	l.hadError = true
}
