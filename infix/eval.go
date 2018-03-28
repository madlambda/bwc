package infix	

import (
	"fmt"
)

func Eval(n Node) (Int, error) {
	var (
		lhs, rhs Int
		err error
	)

	if n.Type() == NodeInt {
		return n.(Int), nil
	}

	if n.Type() != NodeExpr {
		return 0, fmt.Errorf("unexpected %s", n)
	}

	expr := n.(*Expr)
	if expr.Lhs.Type() != NodeInt {
		lhs, err = Eval(expr.Lhs.(*Expr))
		if err != nil {
			return 0, err
		}
	} else {
		lhs = expr.Lhs.(Int)
	}

	if expr.Rhs.Type() != NodeInt {
		rhs, err = Eval(expr.Rhs.(*Expr))
		if err != nil {
			return 0, err
		}
	} else {
		rhs = expr.Rhs.(Int)
	}

	var ret Int
	switch expr.Op {
	case OpAND:
		ret = lhs & rhs
	case OpOR:
		ret = lhs | rhs
	case OpSHL:
		ret = lhs << uint(rhs)
	case OpSHR:
		ret = lhs >> uint(rhs)
	}
	return ret, nil
}
