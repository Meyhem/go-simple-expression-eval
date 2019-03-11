package main

import (
	"testing"
)

func TestParserAst(t *testing.T) {
	ast, err := Parse("1+2")

	if err != nil {
		t.Errorf("got %s; want nil", err)
	}

	plus := ast
	one := ast.Left
	two := ast.Right

	if plus.Typ != ASTNODE_ADD {
		t.Errorf("Typ = %s; want %s", ast.Typ, ASTNODE_ADD)
	}

	if plus.Value != nil {
		t.Errorf("Value = %s; want %v", *ast.Value, nil)
	}

	if one.Typ != ASTNODE_LEAF {
		t.Errorf("Typ = %s; want %s", one.Typ, ASTNODE_LEAF)
	}

	if one.Value == nil || *one.Value != "1" {
		t.Errorf("Value = %v; want %s", one.Value, "1")
	}

	if two.Typ != ASTNODE_LEAF {
		t.Errorf("Typ = %s; want %s", two.Typ, ASTNODE_LEAF)
	}

	if two.Value == nil || *two.Value != "2" {
		t.Errorf("Value = %v; want %s", two.Value, "2")
	}
}

func TestParserEmptyExpr(t *testing.T) {
	ast, err := Parse("")

	if ast != nil {
		t.Error("ast not nil; want nil")
	}

	if err.code != ErrParser {
		t.Errorf("code = %s; want %s", err.code, ErrParser)
	}
}

func TestParserEmptyParenExpr(t *testing.T) {
	ast, err := Parse("(())")

	if ast != nil {
		t.Error("ast not nil; want nil")
	}

	if err.code != ErrParser {
		t.Errorf("code = %s; want %s", err.code, ErrParser)
	}
}

func TestParserPrecedence(t *testing.T) {
	ast, err := Parse("3*(1-2)")

	if ast == nil {
		t.Error("ast nil; want not nil")
	}

	if err != nil {
		t.Errorf("err = %s; want nil", err)
	}

	mul := ast
	three := ast.Left
	minus := ast.Right
	one := ast.Right.Left
	two := ast.Right.Right

	if mul.Typ != ASTNODE_MUL {
		t.Errorf("Typ = %s; want %s", mul.Typ, ASTNODE_MUL)
	}

	if three.Typ != ASTNODE_LEAF {
		t.Errorf("Typ = %s; want %s", three.Typ, ASTNODE_LEAF)
	}

	if minus.Typ != ASTNODE_SUB {
		t.Errorf("Typ = %s; want %s", minus.Typ, ASTNODE_SUB)
	}

	if one.Typ != ASTNODE_LEAF {
		t.Errorf("Typ = %s; want %s", one.Typ, ASTNODE_LEAF)
	}

	if two.Typ != ASTNODE_LEAF {
		t.Errorf("Typ = %s; want %s", two.Typ, ASTNODE_LEAF)
	}
}
