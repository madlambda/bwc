package bwc

import (
	"fmt"
	"strconv"
)

type (
	// Int represents integer numbers
	Int int64

	// Var is a variable
	Var string

	// UnaryExpr holds unary operations like ~
	UnaryExpr struct {
		Op    Optype
		Value Node
	}

	// BinExpr holds binary operations like &,|,etc
	BinExpr struct {
		Op  Optype
		Lhs Node
		Rhs Node
	}

	Assign struct {
		Varname string // Varname is the lhs of the assignment
		Expr    Node   // Expr is the rhs of the assignment
	}

	Node interface {
		Type() Nodetype
		String() string
	}

	Nodetype int
	Optype   int
)

const (
	NodeUnaryExpr Nodetype = iota + 1
	NodeBinExpr
	NodeAssign
	NodeInt
	NodeVar

	binaryOPbegin Optype = iota + 1
	OpAND
	OpOR
	OpXOR
	OpSHL
	OpSHR
	binaryOPend

	unaryOPbegin
	OpNOT
	unaryOPend
)

func (o Optype) String() string {
	switch o {
	case OpAND:
		return "&"
	case OpOR:
		return "|"
	case OpXOR:
		return "^"
	case OpNOT:
		return "~"
	case OpSHL:
		return "<<"
	case OpSHR:
		return ">>"
	}

	panic(fmt.Sprintf("invalid operation: %d", o))
}

func (nt Nodetype) String() string {
	if nt == NodeUnaryExpr {
		return "NodeUnaryExpr"
	} else if nt == NodeBinExpr {
		return "NodeBinExpr"
	} else if nt == NodeInt {
		return "NodeInt"
	} else if nt == NodeAssign {
		return "NodeAssign"
	} else if nt == NodeVar {
		return "NodeVar"
	}
	panic(fmt.Sprintf("invalid node: %d", nt))
}

func (_ Int) Type() Nodetype { return NodeInt }
func (i Int) String() string { return strconv.Itoa(int(i)) }

func (_ Var) Type() Nodetype { return NodeVar }
func (a Var) String() string { return string(a) }

func (_ BinExpr) Type() Nodetype { return NodeBinExpr }
func (a BinExpr) String() string {
	return fmt.Sprintf("%s%s%s", a.Lhs, a.Op, a.Rhs)
}

func (_ UnaryExpr) Type() Nodetype { return NodeUnaryExpr }
func (a UnaryExpr) String() string {
	return fmt.Sprintf("%s%s",
		a.Op, a.Value)
}

func (_ Assign) Type() Nodetype { return NodeAssign }
func (a Assign) String() string {
	return fmt.Sprintf("%s = %s", a.Varname, a.Expr)
}
