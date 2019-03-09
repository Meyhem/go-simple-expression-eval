package main

import (
	"./parser"
)

func main() {
	parser.Parse("1+2*3/4-5/6+7*8-9")
}
