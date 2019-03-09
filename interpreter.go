package main

import (
	"strconv"
)

func add(a int, b int) int {
	return a + b
}

func sub(a int, b int) int {
	return a - b
}

func mul(a int, b int) int {
	return a * b
}

func div(a int, b int) int {
	return a / b
}

type arithmeticFunc func(int, int) int

type interpretMap map[AstNodeType]arithmeticFunc

func postOrderTraversal(node *AstNode, functions interpretMap) int {
	if node.Typ == ASTNODE_LEAF {
		if node.Value == nil {
			panic("Intepretation error: Expected value, got nil")
		}
		number, _ := strconv.Atoi(*node.Value)
		return number
	}

	aritFunc := functions[node.Typ]
	left := postOrderTraversal(node.Left, functions)
	right := postOrderTraversal(node.Right, functions)
	return aritFunc(left, right)
}

func Interpret(ast *AstNode) int {
	var astInterpretMap = interpretMap{
		ASTNODE_ADD: add,
		ASTNODE_SUB: sub,
		ASTNODE_MUL: mul,
		ASTNODE_DIV: div,
	}

	return postOrderTraversal(ast, astInterpretMap)
}
