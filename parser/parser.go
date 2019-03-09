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
				// just trash it
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

func Parse(expr string) {
	lx := lexer.Lex(expr)
	go lx.Run()

	postfixNotation := toPostfix(lx)

	for item := postfixNotation.Front(); item != nil; item = item.Next() {
		fmt.Println("Postfix item: ", item.Value.(lexer.LexItem).Val)
	}

	fmt.Println(postfixNotation)

}
