package runner

import (
	"bufio"
	"fmt"
	"os"
)

type LoxRunner struct {
	HadError bool
	Scanner  *Scanner
	Parser   *Parser
}

func (r *LoxRunner) RunFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	r.run(string(data))
}

func (r *LoxRunner) RunPrompt() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("")
			break
		}
		r.run(text)
	}

}

func (r *LoxRunner) run(program string) {
	r.Scanner = NewScanner(program, r)
	r.Scanner.ScanTokens()
	for _, token := range r.Scanner.Tokens {
		fmt.Println(token.ToString())
	}

	r.Parser = NewParser(r.Scanner.Tokens, r)
	stmts := r.Parser.parse()


	interpreter := &Interpreter{}
	interpreter.interpret(stmts)
}

func (r *LoxRunner) error(line int, message string) {
	r.report(line, "", message)
}

func (r *LoxRunner) report(line int, where string, message string) {
	fmt.Printf("[line %v] Error %v: %v\n", line, where, message)
	r.HadError = true
}

func (r *LoxRunner) tokenError(token Token, message string) {
	if token.Type == EOF {
		r.report(token.Line, " at end", message)
	} else {
		r.report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}
