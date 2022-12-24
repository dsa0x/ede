package main

import (
	"ede/evaluator"
	"ede/lexer"
	"ede/object"
	"ede/parser"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	fileName := flag.Arg(0)
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Execute(string(file))
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
		fmt.Println(prog.ParseErrors)
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
