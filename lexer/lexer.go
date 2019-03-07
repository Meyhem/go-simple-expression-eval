package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type itemType int

const (
	iErr             = 0
	iNUMBER itemType = iota
	iWHITE
	iADD
	iSUB
	iMUL
	iDIV
)

const eof = -1

const (
	numbers   = "0123456789"
	operators = "+-*/"
	white     = " \n\r\t"
)

type stateFn func(*lexer) stateFn

type lexItem struct {
	typ itemType
	pos int
	val string
}

func (li lexItem) String() string {
	strType := "UNDEFINED"

	switch li.typ {
	case iErr:
		strType = "iErr"
	case iNUMBER:
		strType = "iNUMBER"
	case iWHITE:
		strType = "iWHITE"
	case iADD:
		strType = "iADD"
	case iSUB:
		strType = "iSUB"
	case iMUL:
		strType = "iMUL"
	case iDIV:
		strType = "iDIV"
	case eof:
		strType = "EOF"
	}

	return fmt.Sprintf("type: %s, val: %q, start: %d", strType, li.val, li.pos)

}

type lexer struct {
	text  string
	start int
	pos   int
	width int
	items chan lexItem
}

func (l *lexer) dumpState() {
	fmt.Printf("%#v\n", l)
}

func (l *lexer) next() rune {
	if l.pos >= len(l.text) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.text[l.pos:])

	l.width = w
	l.pos += w

	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
	_, w := utf8.DecodeLastRuneInString(l.text[:l.pos])
	l.width = w
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) consume(runes string) bool {
	if strings.ContainsRune(runes, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) consumeAll(runes string) {
	for l.consume(runes) {
	}
}

func (l *lexer) emit(typ itemType) {
	l.items <- lexItem{
		typ: typ,
		pos: l.pos,
		val: l.text[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- lexItem{
		pos: l.pos,
		typ: iErr,
		val: fmt.Sprintf(format, args...),
	}
	return nil
}

func lexFn(l *lexer) stateFn {
	r := l.peek()
	switch {
	case r == eof:
		l.emit(eof)
		return nil
	case strings.ContainsRune(operators, r):
		return lexOperator
	case strings.ContainsRune(numbers, r):
		return lexNumber
	case strings.ContainsRune(white, r):
		return lexWhite

	default:
		return l.errorf("Invalid symbol: %q", r)
	}
}

func lexOperator(l *lexer) stateFn {
	op := l.next()
	switch op {
	case '+':
		l.emit(iADD)
	case '-':
		l.emit(iSUB)
	case '*':
		l.emit(iMUL)
	case '/':
		l.emit(iDIV)
	default:
		return l.errorf("lexOperator: invalid operator: %q", op)
	}

	return lexFn
}

func lexNumber(l *lexer) stateFn {
	l.consumeAll(numbers)
	l.emit(iNUMBER)
	return lexFn
}

func lexWhite(l *lexer) stateFn {
	l.consumeAll(white)
	l.emit(iWHITE)
	return lexFn
}

func (l *lexer) Items() chan lexItem {
	return l.items
}

func (l *lexer) Run() {
	defer close(l.items)

	for fun := lexFn; fun != nil; {
		fun = fun(l)
	}

}

func Lex(text string) *lexer {
	return &lexer{
		items: make(chan lexItem),
		text:  text,
	}
}
