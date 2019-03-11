package main

import (
	"testing"
)

func runLexing(expr string) (*Lexer, []LexItem) {
	lexer := Lex(expr)
	items := make([]LexItem, 0)

	go lexer.Run()

	for item := range lexer.Items() {
		items = append(items, item)
	}

	return lexer, items
}

func TestSimpleExprLen(t *testing.T) {
	_, items := runLexing("1+2")

	// should be 4 with EOF
	itemCount := len(items)

	if itemCount != 4 {
		t.Errorf("Len = %d; want 4", itemCount)
	}
}

func TestSimpleExprItems(t *testing.T) {
	_, items := runLexing("11+2")

	tested := items[0]
	if tested.Typ != INUMBER {
		t.Errorf("Typ is %s; want %s", tested.Typ, INUMBER)
	}

	if tested.Val != "11" {
		t.Errorf("Val is %s; want %s", tested.Val, "11")
	}

	tested = items[1]
	if tested.Typ != IADD {
		t.Errorf("Typ is %s; want %s", tested.Typ, IADD)
	}

	tested = items[2]
	if tested.Typ != INUMBER {
		t.Errorf("Typ is %s; want %s", tested.Typ, INUMBER)
	}

	if tested.Val != "2" {
		t.Errorf("Val is %s; want %s", tested.Val, "2")
	}

	tested = items[3]
	if tested.Typ != EOF {
		t.Errorf("Typ is %s; want %s", tested.Typ, EOF)
	}
}

func TestEmptyExprItems(t *testing.T) {
	_, items := runLexing("")

	itemCount := len(items)

	if itemCount != 1 {
		t.Errorf("Len = %d; want 4", itemCount)
	}

	tested := items[0]

	if tested.Typ != EOF {
		t.Errorf("Typ is %s; want %s", tested.Typ, EOF)
	}
}

func TestInvalidRune(t *testing.T) {
	_, items := runLexing("1+W+3+5")

	itemCount := len(items)

	// stopped at 'W' rune
	if itemCount != 3 {
		t.Errorf("Len = %d; want 4", itemCount)
	}

	tested := items[2]

	if tested.Typ != IERR {
		t.Errorf("Typ is %s; want %s", tested.Typ, IERR)
	}
}
