package main

import (
	"flag"
	"fmt"
	"os"
)

var fileName string

func Run() {
	flag.Parse()

	fileName = flag.Arg(0)
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Execute(string(file))
}
