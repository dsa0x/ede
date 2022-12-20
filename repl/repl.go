package repl

import (
	"bufio"
	"ede/evaluator"
	"ede/lexer"
	"ede/object"
	"ede/parser"
	"fmt"
	"io"
)

func Start(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	env := object.NewEnvironment(nil)
	fmt.Fprintf(output, ">> ")
	for scanner.Scan() {
		lex := lexer.New(scanner.Text())
		p := parser.New(lex)
		prog := p.Parse()
		if p.Errors() != nil {
			fmt.Println(p.Errors())
			continue
		}
		if eval := evaluator.Eval(prog, env); eval != nil {
			output.Write([]byte(fmt.Sprintf("Result: %v\n", eval.Inspect())))
		}
	}
}
