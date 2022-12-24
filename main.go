package main

import (
	"bytes"
	"ede/evaluator"
	"ede/lexer"
	"ede/object"
	"ede/parser"
	"fmt"
	"os"
)

func main() {
	// repl.Start(os.Stdin, os.Stdout)
	input := `
	let a = 10;
	let add = func(x) {
		println("a", a, "\n");
		return x + a;
	};
	a = a * 2;
	a = add(add(10));
	`
	buf := new(bytes.Buffer)
	file, err := os.ReadFile("./hello.ede")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, err := buf.Write(file); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	input = buf.String()
	if err := Execute(input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Execute(input string) error {
	env := object.NewEnvironment(nil)
	lex := lexer.New(input)
	p := parser.New(lex)
	if p.Errors() != nil {
		return p.Errors()
	}
	prog := p.Parse()
	if prog.ParseErrors != nil {
		return prog.ParseErrors
	}
	eval := (&evaluator.Evaluator{}).Eval(prog, env)
	if eval != nil {
		fmt.Println(eval.Inspect())
	}
	return nil
}

func unwrappedError(errs []error) error {
	var err error
	for _, e := range errs {
		err = fmt.Errorf("%s, %w", err, e)
	}
	return err
}
