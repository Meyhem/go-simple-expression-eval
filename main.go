package main

import (
	"fmt"

	"./lexer"
)

func main() {
	lx := lexer.Lex("1+  \t\n\r  2+3*3")

	go lx.Run()

	for item := range lx.Items() {
		fmt.Printf("Item: %s\n", item)
	}

}
