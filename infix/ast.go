package infix

import (
	"fmt"
	"strconv"
)

type (
	// Int represents integer numbers
	Int int64

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
	NodeInt

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

	return "<invalid op>"
}

func (_ Int) Type() Nodetype { return NodeInt }
func (i Int) String() string { return strconv.Itoa(int(i)) }

func (_ BinExpr) Type() Nodetype { return NodeBinExpr }
func (a BinExpr) String() string {
	lhs := a.Lhs.String()
	rhs := a.Rhs.String()
	op := a.Op.String()

	return fmt.Sprintf("%s %s %s", lhs, rhs, op)
}

func (_ UnaryExpr) Type() Nodetype { return NodeUnaryExpr }
func (a UnaryExpr) String() string {
	return fmt.Sprintf("%s%s",
		a.Op.String(), a.Value.String())
}