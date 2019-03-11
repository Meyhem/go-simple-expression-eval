package main

import (
	"testing"
)

func TestInterpreter(t *testing.T) {
	ast, _ := Parse("1+2")

	res, err := Interpret(ast)

	if err != nil {
		t.Errorf("err = %s; want nil", err)
	}

	if res != 3 {
		t.Errorf("res = %d; want %d", res, 3)
	}
}

func TestSingleValue(t *testing.T) {
	ast, _ := Parse("1")

	res, err := Interpret(ast)

	if err != nil {
		t.Errorf("err = %s; want nil", err)
	}

	if res != 1 {
		t.Errorf("res = %d; want %d", res, 1)
	}
}

func TestOperator(t *testing.T) {

	ast := NewAstNode(ASTNODE_ADD, nil)

	_, err := Interpret(ast)

	if err == nil {
		t.Error("err = nil; want 'Expected evaluatable node, got nil'")
	}
}
