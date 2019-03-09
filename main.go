package main

import (
	"fmt"

	"./parser"
)

func main() {
	ast := parser.Parse("1/2+3/4")
	fmt.Println(ast)
}
