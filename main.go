package main

import (
	"fmt"
)

func main() {
	ast := Parse("-2+3")
	result := Interpret(ast)
	fmt.Println("Result=", result)
}
