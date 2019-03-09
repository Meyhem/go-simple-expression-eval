package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type ItemType int

const (
	IERR             = 0
	INUMBER ItemType = iota
	ILPAR
	IRPAR
	IADD
	ISUB
	IMUL
	IDIV
)

const EOF = -1

const (
	numbers   = "0123456789"
	operators = "+-*/"
	white     = " \n\r\t"
	lpar      = "("
	rpar      = ")"
)

type stateFn func(*Lexer) stateFn

type LexItem struct {
	Typ ItemType
	Pos int
	Val string
}

func (li LexItem) String() string {
	strType := "UNDEFINED"

	switch li.Typ {
	case IERR:
		strType = "IERR"
	case INUMBER:
		strType = "INUMBER"
	case ILPAR:
		strType = "ILPAR"
	case IRPAR:
		strType = "IRPAR"
	case IADD:
		strType = "IADD"
	case ISUB:
		strType = "ISUB"
	case IMUL:
		strType = "IMUL"
	case IDIV:
		strType = "IDIV"
	case EOF:
		strType = "EOF"
	}

	return fmt.Sprintf("Type: %s, Val: %q, Pos: %d", strType, li.Val, li.Pos)

}

type Lexer struct {
	text  string
	start int
	Pos   int
	width int
	items chan LexItem
}

func (l *Lexer) dumpState() {
	fmt.Printf("%#v\n", l)
}

func (l *Lexer) next() rune {
	if l.Pos >= len(l.text) {
		l.width = 0
		return EOF
	}

	r, w := utf8.DecodeRuneInString(l.text[l.Pos:])

	l.width = w
	l.Pos += w

	return r
}

func (l *Lexer) backup() {
	l.Pos -= l.width
	_, w := utf8.DecodeLastRuneInString(l.text[:l.Pos])
	l.width = w
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) consume(runes string) bool {
	if strings.ContainsRune(runes, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) consumeAll(runes string) {
	for l.consume(runes) {
	}
}

func (l *Lexer) ignore() {
	l.start = l.Pos
	l.width = 0
}

func (l *Lexer) emit(typ ItemType) {
	l.items <- LexItem{
		Typ: typ,
		Pos: l.Pos,
		Val: l.text[l.start:l.Pos],
	}
	l.start = l.Pos
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- LexItem{
		Pos: l.Pos,
		Typ: IERR,
		Val: fmt.Sprintf(format, args...),
	}
	return nil
}

func lexFn(l *Lexer) stateFn {
	r := l.peek()
	switch {
	case r == EOF:
		l.emit(EOF)
		return nil

	case strings.ContainsRune(white, r):
		return lexWhite
	case strings.ContainsRune(operators, r):
		return lexOperator
	case strings.ContainsRune(numbers, r):
		return lexNumber
	case strings.ContainsRune(lpar, r):
		return lexLpar
	case strings.ContainsRune(rpar, r):
		return lexRpar
	default:
		return l.errorf("Invalid symbol: %q", r)
	}
}

func lexOperator(l *Lexer) stateFn {
	op := l.next()
	switch op {
	case '+':
		l.emit(IADD)
	case '-':
		l.emit(ISUB)
	case '*':
		l.emit(IMUL)
	case '/':
		l.emit(IDIV)
	default:
		return l.errorf("lexOperator: inValid operator: %q", op)
	}

	return lexFn
}

func lexLpar(l *Lexer) stateFn {
	l.consume(lpar)
	l.emit(ILPAR)
	return lexFn
}

func lexRpar(l *Lexer) stateFn {
	l.consume(rpar)
	l.emit(IRPAR)
	return lexFn
}

func lexNumber(l *Lexer) stateFn {
	l.consumeAll(numbers)
	l.emit(INUMBER)
	return lexFn
}

func lexWhite(l *Lexer) stateFn {
	l.consumeAll(white)
	l.ignore()
	return lexFn
}

func (l *Lexer) Items() chan LexItem {
	return l.items
}

func (l *Lexer) Run() {
	defer close(l.items)

	for fun := lexFn; fun != nil; {
		fun = fun(l)
	}

}

func Lex(text string) *Lexer {
	return &Lexer{
		items: make(chan LexItem),
		text:  text,
	}
}
