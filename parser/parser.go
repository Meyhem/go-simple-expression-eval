package parser

import (
	"container/list"
	"fmt"

	"../lexer"
)

type parser struct {
	x int
}

func precedence(typ lexer.ItemType) int {
	switch typ {
	case lexer.IADD:
		fallthrough
	case lexer.ISUB:
		return 1

	case lexer.IMUL:
		fallthrough
	case lexer.IDIV:
		return 2

	default:
		return -1
	}
}

func toPostfix(lx *lexer.Lexer) *list.List {

	opStack := NewStack()
	postFix := list.New()

	for item := range lx.Items() {

		// end of tok stream
		if item.Typ == lexer.EOF {
			continue
		}

		// lexing error
		if item.Typ == lexer.IERR {
			panic("Lexing error")
		}

		// if its number put to output
		if item.Typ == lexer.INUMBER {
			postFix.PushBack(item)
			continue
		}

		// if left parenth put to output
		if item.Typ == lexer.ILPAR {
			opStack.Push(item)
			continue
		}

		// if right parenth
		if item.Typ == lexer.IRPAR {
			// pop stack to output until we find left parenth in stack
			for opStack.Len() > 0 && opStack.Top().(lexer.LexItem).Typ != lexer.ILPAR {
				postFix.PushBack(opStack.Pop())
			}

			// if there is none then there is error in parity
			if opStack.Len() > 0 && opStack.Top().(lexer.LexItem).Typ != lexer.ILPAR {
				panic("Invalid expr")
			} else {
				// otherwise just trash it
				opStack.Pop()
			}
		} else {
			// is any other operator
			// check precedence
			for opStack.Len() > 0 && precedence(item.Typ) <= precedence(opStack.Top().(lexer.LexItem).Typ) {
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

func translateLexToAstType(typ lexer.ItemType) AstNodeType {
	switch typ {
	case lexer.IADD:
		return ASTNODE_ADD
	case lexer.ISUB:
		return ASTNODE_SUB
	case lexer.IMUL:
		return ASTNODE_MUL
	case lexer.IDIV:
		return ASTNODE_DIV
	default:
		panic(fmt.Sprintf("Unexpected lexer item type occured during parsing %q", typ))
	}
}

func constructAst(postfixList *list.List) *AstNode {
	stack := NewStack()
	for item := postfixList.Front(); item != nil; item = item.Next() {
		lexItem := item.Value.(lexer.LexItem)
		if lexItem.Typ == lexer.INUMBER {
			stack.Push(NewAstNode(ASTNODE_LEAF, &lexItem.Val))
		} else {
			nodeType := translateLexToAstType(lexItem.Typ)
			node := NewAstNode(nodeType, nil)
			node.Left = stack.Pop().(*AstNode)
			node.Right = stack.Pop().(*AstNode)
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
	lx := lexer.Lex(expr)
	go lx.Run()

	postfixNotation := toPostfix(lx)
	abstractSyntaxTree := constructAst(postfixNotation)

	return abstractSyntaxTree
}
