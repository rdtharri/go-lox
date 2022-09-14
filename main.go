package main

import (
	"github.com/rdtharri/go-lox/runner"
	"os"
)

func main() {
	runner := runner.LoxRunner{}
	args := os.Args[1:]
	if len(args) == 1 {
		runner.RunFile(args[0])
	} else {
		runner.RunPrompt()
	}
	if runner.HadError {
		os.Exit(65)
	}
}
