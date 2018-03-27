package infix

import (
	"fmt"
	"testing"
)

var format = fmt.Sprintf

func desc(a Expr, res Int) string {
	return format("%s %s %s = %s",
		a.Lhs, a.Rhs, a.Op, res)
}

func TestEval(t *testing.T) {
	for _, tc := range []struct {
		expr Expr
		res  Int
	}{
		{
			expr: Expr{
				Op:  OpAND,
				Lhs: Int(0),
				Rhs: Int(0),
			},
			res: 0,
		},
		{
			expr: Expr{
				Op:  OpAND,
				Lhs: Int(0),
				Rhs: Int(1),
			},
			res: 0,
		},
		{
			expr: Expr{
				Op:  OpAND,
				Lhs: Int(1),
				Rhs: Int(1),
			},
			res: 1,
		},
		{
			expr: Expr{
				Op:  OpAND,
				Lhs: Int(0xffff0000),
				Rhs: Int(0x0000ffff),
			},
			res: 0,
		},
		{
			expr: Expr{
				Op:  OpOR,
				Lhs: Int(0xffff0000),
				Rhs: Int(0x0000ffff),
			},
			res: 0xffffffff,
		},

		// groups
		{
			expr: Expr{
				Op: OpAND,
				Lhs: Expr{
					Op:  OpOR,
					Lhs: Int(0xffff0000),
					Rhs: Int(0x0000ffff),
				},
				Rhs: Int(0x000000ff),
			},
			res: 0x000000ff,
		},
		{
			expr: Expr{
				Op: OpOR,
				Lhs: Expr{
					Op:  OpOR,
					Lhs: Int(0x000000ff),
					Rhs: Int(0xff000000),
				},
				Rhs: Expr{
					Op:  OpOR,
					Lhs: Int(0x00ff0000),
					Rhs: Int(0x0000ff00),
				},
			},
			res: 0xffffffff,
		},
	} {
		tc := tc
		t.Run(desc(tc.expr, tc.res), func(t *testing.T) {
			got := Eval(tc.expr)
			if got != tc.res {
				t.Fatalf("Fail: %d != %d", got, tc.res)
			}
		})
	}
}