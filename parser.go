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

func toPostfix(lx *Lexer) (*list.List, error) {

	opStack := NewStack()
	postFix := list.New()

	for item := range lx.Items() {

		// end of tok stream
		if item.Typ == EOF {
			continue
		}

		// lexing error
		if item.Typ == IERR {
			return nil, fmt.Errorf("Lexing error at %d: %s", item.pos, item.Val)
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
				return nil, fmt.Errorf("Parser error at %d: Unmatched paretheses", opStack.Top().(LexItem).pos)
			}

			// we are in rparenth so if there is no lparenh its parity error
			if opStack.Len() == 0 {
				return nil, fmt.Errorf("Parsing error at %d: Missing '('", item.pos)
			}
			// otherwise just trash it
			opStack.Pop()

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

	return postFix, nil
}

func translateLexToAstType(typ ItemType) (AstNodeType, error) {
	switch typ {
	case IADD:
		return ASTNODE_ADD, nil
	case ISUB:
		return ASTNODE_SUB, nil
	case IMUL:
		return ASTNODE_MUL, nil
	case IDIV:
		return ASTNODE_DIV, nil
	default:
		return 0, fmt.Errorf("Unexpected item type occured during parsing %q", typ)
	}
}

func constructAst(postfixList *list.List) (*AstNode, error) {
	stack := NewStack()
	for item := postfixList.Front(); item != nil; item = item.Next() {
		lexItem := item.Value.(LexItem)
		if lexItem.Typ == INUMBER {
			stack.Push(NewAstNode(ASTNODE_LEAF, &lexItem.Val))
		} else {
			nodeType, err := translateLexToAstType(lexItem.Typ)
			if err != nil {
				return nil, fmt.Errorf("Parser error at %d: Missing ')'", lexItem.pos)
			}
			node := NewAstNode(nodeType, nil)

			// order important, otherwise we switch operands
			if stack.Len() < 2 {
				return nil, fmt.Errorf("Parser error at %d: Missing operand", lexItem.pos)
			}

			node.Right = stack.Pop().(*AstNode)
			node.Left = stack.Pop().(*AstNode)
			stack.Push(node)
		}
	}

	if stack.Len() < 1 {
		return nil, fmt.Errorf("Parsing error: Expression without root")
	}

	return stack.Pop().(*AstNode), nil
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

func Parse(expr string) (*AstNode, error) {
	lx := Lex(expr)
	go lx.Run()

	postfixNotation, err := toPostfix(lx)

	if err != nil {
		return nil, err
	}

	abstractSyntaxTree, err := constructAst(postfixNotation)

	if err != nil {
		return nil, err
	}

	return abstractSyntaxTree, nil
}
