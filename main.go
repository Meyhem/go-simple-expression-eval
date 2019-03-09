package main

import (
	"fmt"
	"os"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Specify expressions to evaluate...\ne.g.: 1+2*(6-8)")
		return
	}

	ast, err := Parse(args[1])

	if err != nil {
		fmt.Println(err)
		return
	}

	result := Interpret(ast)
	fmt.Println(result)
}
