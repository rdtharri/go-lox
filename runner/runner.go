package runner

import (
	"bufio"
	"fmt"
	"os"
)

type LoxRunner struct {
	HadError bool
	Scanner  *Scanner
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

	// TODO: remove
	printer := &PrinterVistor{}
	exp := &BinaryExpression{
		Operator: Token{
			Lexeme: "*",
		},
		Left: &LiteralExpression{
			Token: Token{
				Lexeme: "500",
			},
		},
		Right: &LiteralExpression{
			Token: Token{
				Lexeme: "200",
			},
		},
	}
	printer.print(exp)
}

func (r *LoxRunner) error(line int, message string) {
	r.report(line, "", message)
}

func (r *LoxRunner) report(line int, where string, message string) {
	fmt.Printf("[line %v] Error %v: %v\n", line, where, message)
	r.HadError = true
}
