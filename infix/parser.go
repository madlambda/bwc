package infix

import (
	"fmt"
)

func Parse(code string) (Expr, error) {
	tokens := Lex(code)
	for tok := range tokens {
		if tok.T == EOF {
			break
		}

		fmt.Printf("tok: %v\n", tok)
	}
	return Expr{
		Op:  1,
		Lhs: Int(1),
		Rhs: Int(1),
	}, nil
}