package main

import (
	"container/list"
	"fmt"
)

func precedence(typ ItemType) int {
	switch typ {
	case IADD:
		fallthrough
	case ISUB:
		return 1

	case IMUL:
		fallthrough
	case IDIV:
		return 2

	default:
		return -1
	}
}

func toPostfix(lx *Lexer) *list.List {

	opStack := NewStack()
	postFix := list.New()

	for item := range lx.Items() {

		// end of tok stream
		if item.Typ == EOF {
			continue
		}

		// lexing error
		if item.Typ == IERR {
			panic("Lexing error")
		}

		// if its number put to output
		if item.Typ == INUMBER {
			postFix.PushBack(item)
			continue
		}

		// if left parenth put to output
		if item.Typ == ILPAR {
			opStack.Push(item)
			continue
		}

		// if right parenth
		if item.Typ == IRPAR {
			// pop stack to output until we find left parenth in stack
			for opStack.Len() > 0 && opStack.Top().(LexItem).Typ != ILPAR {
				postFix.PushBack(opStack.Pop())
			}

			// if there is none then there is error in parity
			if opStack.Len() > 0 && opStack.Top().(LexItem).Typ != ILPAR {
				panic("Invalid expr")
			} else {
				// otherwise just trash it
				opStack.Pop()
			}
		} else {
			// is any other operator
			// check precedence
			for opStack.Len() > 0 && precedence(item.Typ) <= precedence(opStack.Top().(LexItem).Typ) {
				// just put it to output
				postFix.PushBack(opStack.Pop())
			}
			// put it to stack
			opStack.Push(item)
		}
	}

	// empty stack to output
	for opStack.Len() > 0 {
		postFix.PushBack(opStack.Pop())
	}

	return postFix
}

func translateLexToAstType(typ ItemType) AstNodeType {
	switch typ {
	case IADD:
		return ASTNODE_ADD
	case ISUB:
		return ASTNODE_SUB
	case IMUL:
		return ASTNODE_MUL
	case IDIV:
		return ASTNODE_DIV
	default:
		panic(fmt.Sprintf("Unexpected item type occured during parsing %q", typ))
	}
}

func constructAst(postfixList *list.List) *AstNode {
	stack := NewStack()
	for item := postfixList.Front(); item != nil; item = item.Next() {
		lexItem := item.Value.(LexItem)
		if lexItem.Typ == INUMBER {
			stack.Push(NewAstNode(ASTNODE_LEAF, &lexItem.Val))
		} else {
			nodeType := translateLexToAstType(lexItem.Typ)
			node := NewAstNode(nodeType, nil)

			// order important, otherwise we switch operands
			node.Right = stack.Pop().(*AstNode)
			node.Left = stack.Pop().(*AstNode)
			stack.Push(node)
		}
	}

	return stack.Pop().(*AstNode)
}

func traversePreorder(root *AstNode) {
	if root == nil {
		return
	}

	fmt.Println(root)

	traversePreorder(root.Left)
	traversePreorder(root.Right)
}

func traverseInorder(root *AstNode) {
	if root == nil {
		return
	}

	traversePreorder(root.Left)
	fmt.Println(root)
	traversePreorder(root.Right)
}

func traversePostorder(root *AstNode) {
	if root == nil {
		return
	}

	traversePreorder(root.Left)
	traversePreorder(root.Right)
	fmt.Println(root)
}

func Parse(expr string) *AstNode {
	lx := Lex(expr)
	go lx.Run()

	postfixNotation := toPostfix(lx)
	abstractSyntaxTree := constructAst(postfixNotation)

	return abstractSyntaxTree
}
