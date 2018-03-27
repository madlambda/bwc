package infix	

func Eval(tree Expr) Int {
	var lhs, rhs Int

	if tree.Lhs.Type() != NodeInt {
		lhs = Eval(tree.Lhs.(Expr))
	} else {
		lhs = tree.Lhs.(Int)
	}

	if tree.Rhs.Type() != NodeInt {
		rhs = Eval(tree.Rhs.(Expr))
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
