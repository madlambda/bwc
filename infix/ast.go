package infix

import (
	"fmt"
	"strconv"
)

type (
	Int  int64
	Expr struct {
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
	NodeExpr Nodetype = iota + 1
	NodeInt

	OpAND Optype = iota + 1
	OpOR
	OpSHL
	OpSHR
)

func (o Optype) String() string {
	if o == OpAND {
		return "&"
	}
	if o == OpOR {
		return "|"
	}
	return "<invalid op>"
}

func (_ Int) Type() Nodetype { return NodeInt }
func (i Int) String() string { return strconv.Itoa(int(i)) }

func (_ Expr) Type() Nodetype { return NodeExpr }
func (a Expr) String() string {
	lhs := a.Lhs.String()
	rhs := a.Rhs.String()
	op := a.Op.String()

	return fmt.Sprintf("%s %s %s", lhs, rhs, op)
}