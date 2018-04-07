package infix	

import (
	"fmt"
)

func Exec(code string) (Int, error) {
	n, err := Parse(code)
	if err != nil {
		return 0, err
	}
	return Eval(n)
}

func Eval(n Node) (Int, error) {
	switch n.Type() {
	case NodeInt:
		return n.(Int), nil
	case NodeUnaryExpr:
		return evalUnaryExpr(n.(UnaryExpr))
	case NodeBinExpr:
		return evalBinExpr(n.(BinExpr))
	}


	return 0, fmt.Errorf("unexpected %s", n) 
}

func evalUnaryExpr(expr UnaryExpr) (Int, error) {
	num, err := Eval(expr.Value)
	if err != nil {
		return 0, err
	}

	switch expr.Op {
	case OpNOT:
		return Int(^int64(num)), nil
	}

	return 0, fmt.Errorf("invalid unary expr: %s", expr.Op)
}
 
func evalBinExpr(expr BinExpr) (Int, error) {
	var lhs, rhs Int
	var err error
	if expr.Lhs.Type() != NodeInt {
		lhs, err = Eval(expr.Lhs)
		if err != nil {
			return 0, err
		}
	} else {
		lhs = expr.Lhs.(Int)
	}

	if expr.Rhs.Type() != NodeInt {
		rhs, err = Eval(expr.Rhs)
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
	default:
		return 0, fmt.Errorf("invalid op (%v)", expr.Op)
	}
	return ret, nil
}
