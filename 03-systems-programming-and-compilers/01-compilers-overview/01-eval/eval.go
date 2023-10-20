package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strconv"
)

// Given an expression containing only int types, evaluate
// the expressiron and return the result.
func Evaluate(expr ast.Expr) (int, error) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		switch v.Kind {
		case token.INT:
			i, err := strconv.Atoi(v.Value)
			check(err)
			return i, nil
		}
	case *ast.BinaryExpr:
		x, err := Evaluate(v.X)
		check(err)
		y, err := Evaluate(v.Y)
		check(err)
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
	return 0, errors.New("Missing case for ast.Expr")
}

func main() {
	expr, err := parser.ParseExpr("1 + 2 - 3 * 4")
	if err != nil {
		log.Fatal(err)
	}
	// fset := token.NewFileSet()
	// err = ast.Print(fset, expr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	res, err := Evaluate(expr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Eval: %d\n", res)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
