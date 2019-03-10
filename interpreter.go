package main

import (
	"strconv"
)

// Implementation of ADD operation
func add(a int, b int) int {
	return a + b
}

// Implementation of SUBTRACT operation
func sub(a int, b int) int {
	return a - b
}

// Implementation of MULTIPLY operation
func mul(a int, b int) int {
	return a * b
}

// Implementation of DIVIDE operation
func div(a int, b int) int {
	return a / b
}

// Type of arithmetic function taking two ints and returning ints
type arithmeticFunc func(int, int) int

// type for mapping AST node types to corresponding arithmetic operation
type interpretMap map[AstNodeType]arithmeticFunc

// Recursive post order traversal that evaluates AST
// 1. Visit left
// 2. Visit right
// 3. Visit self
func postOrderTraversal(node *AstNode, functions interpretMap) int {
	// If we are on number node
	if node.Typ == ASTNODE_LEAF {
		if node.Value == nil {
			panic("Intepretation error: Expected value, got nil")
		}

		// Parse string val to integer
		number, _ := strconv.Atoi(*node.Value)

		// return it to higher stack frame (numbers should occur only in leaf nodes)
		return number
	}

	// pick correct computation function
	aritFunc := functions[node.Typ]

	// recursively evaluate left subtree
	left := postOrderTraversal(node.Left, functions)

	// recursively evaluate rightsubtree
	right := postOrderTraversal(node.Right, functions)

	// use its value to do computation
	return aritFunc(left, right)
}

// Interpret is function that evaluates AST and returns corresponding result
func Interpret(ast *AstNode) int {
	var astInterpretMap = interpretMap{
		ASTNODE_ADD: add,
		ASTNODE_SUB: sub,
		ASTNODE_MUL: mul,
		ASTNODE_DIV: div,
	}

	return postOrderTraversal(ast, astInterpretMap)
}
