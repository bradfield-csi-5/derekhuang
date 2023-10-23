package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strconv"
)

// Given an expression containing only int types, evaluate
// the expression and return the result.
func Evaluate(expr ast.Expr) (int, error) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		switch v.Kind {
		case token.INT:
			i, err := strconv.Atoi(v.Value)
			return i, err
		}
	case *ast.BinaryExpr:
		x, err := Evaluate(v.X)
		if err != nil {
			return -1, err
		}
		y, err := Evaluate(v.Y)
		if err != nil {
			return -1, err
		}
		switch v.Op {
		case token.ADD:
			return x + y, nil
		case token.SUB:
			return x - y, nil
		case token.MUL:
			return x * y, nil
		case token.QUO:
			return x / y, nil
		}
	case *ast.ParenExpr:
		return Evaluate(v.X)
	}
	return -1, fmt.Errorf("Missing case for ast.Expr: %v\n", expr)
}

func main() {
	expr, err := parser.ParseExpr("1 + 2 - 3 * 4")
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	err = ast.Print(fset, expr)
	if err != nil {
		log.Fatal(err)
	}
}
