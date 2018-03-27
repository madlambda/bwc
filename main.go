package main

import (
	"fmt"
)

type (
	Expr struct {
		Op  int
		Lhs Node
		Rhs Node
	}

	Node interface {
		Type() int
	}

	Int int64
)

const (
	Nodetype = iota
	NodeExpr
	NodeInt

	Optype = iota
	OpAND
	OpOR
)

func (_ Int) Type() int { return NodeInt }
func (_ Expr) Type() int { return NodeExpr }

func eval(tree Expr) Int {
	var lhs, rhs Int

	if tree.Lhs.Type() != NodeInt {
		lhs = eval(tree.Lhs.(Expr))
	} else {
		lhs = tree.Lhs.(Int)
	}

	if tree.Rhs.Type() != NodeInt {
		rhs = eval(tree.Rhs.(Expr))
	} else {
		rhs = tree.Rhs.(Int)
	}

	var ret Int
	switch tree.Op {
	case OpAND:
		ret = lhs & rhs
	case OpOR:
		ret = lhs | rhs
	}
	return ret
}

func main() {
	expr := Expr{
		Op:  OpAND,
		Lhs: Int(1),
		Rhs: Int(2),
	}

	fmt.Printf("%d\n", eval(expr))
}