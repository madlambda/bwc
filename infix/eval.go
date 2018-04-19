package infix	

import (
	"fmt"
)

type interp struct {
	environ map[string]Int
}

func NewInterp() *interp {
	return &interp{
		environ: make(map[string]Int),
	}
}

func (e *interp) Exec(code string) (Int, error) {
	n, err := Parse(code)
	if err != nil {
		return 0, err
	}
	return e.Eval(n)
}

func (e *interp) Eval(n Node) (Int, error) {
	switch n.Type() {
	case NodeInt:
		return n.(Int), nil
	case NodeVar:
		return e.evalVar(n.(Var))
	case NodeUnaryExpr:
		return e.evalUnaryExpr(n.(UnaryExpr))
	case NodeBinExpr:
		return e.evalBinExpr(n.(BinExpr))
	case NodeAssign:
		return e.evalAssign(n.(Assign))
	}


	return 0, fmt.Errorf("unexpected %s", n) 
}

func (e *interp) evalVar(v Var) (Int, error) {
	if val, ok := e.environ[string(v)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("undefined variable %s", v)
}

func (e *interp) evalUnaryExpr(expr UnaryExpr) (Int, error) {
	num, err := e.Eval(expr.Value)
	if err != nil {
		return 0, err
	}

	switch expr.Op {
	case OpNOT:
		return Int(^int64(num)), nil
	}

	return 0, fmt.Errorf("invalid unary expr: %s", expr.Op)
}

func (e *interp) evalOperand(n Node) (Int, error) {
	if n.Type() != NodeInt {
		return e.Eval(n)
	}

	return n.(Int), nil
}
 
func (e *interp) evalBinExpr(expr BinExpr) (Int, error) {
	lhs, err := e.evalOperand(expr.Lhs)
	if err != nil {
		return 0, err
	}
	
	rhs, err := e.evalOperand(expr.Rhs)
	if err != nil {
		return 0, err
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
		return 0	, fmt.Errorf("invalid op (%v)", expr.Op)
	}
	return ret, nil
}

func (e *interp) evalAssign(assign Assign) (Int, error) {
	ret, err := e.Eval(assign.Expr)
	if err != nil {
		return 0, err
	}
	e.environ[assign.Varname] = ret
	return ret, nil
}