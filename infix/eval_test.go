package infix

import (
	"fmt"
	"testing"
)

var format = fmt.Sprintf

func desc(n Node, res Int) string {
	if n.Type() == NodeBinExpr {
		a := n.(BinExpr)
		return format("%s %s %s = %s",
			a.Lhs, a.Rhs, a.Op, res)
	} else if n.Type() == NodeUnaryExpr {
		a := n.(UnaryExpr)
		return format("%s%s", a.Op, a.Value)
	}

	a := n.(Int)
	return format("%d", a)
}

func TestEvalGrammar(t *testing.T) {
	interp := NewInterp()

	for _, tc := range []struct {
		code string
		res  Int
	}{
		{
			code: "0",
			res:  0,
		},
		{
			code: "0|1",
			res:  1,
		},
		{
			code: "0|1|2",
			res:  3,
		},
		{
			code: "0|1|2|3",
			// 00000000
			// 00000001
			// 00000010
			// 00000011
			res: 3,
		},
		{
			code: "0|1|2|3|4",
			res:  7,
		},
		{
			code: "0&1",
			res:  0,
		},
		{
			code: "7&3",
			res:  3,
		},
		{
			code: "1|2&1",
			res:  1,
		},
		{
			code: "a = 0",
			res: 0,
		},
		{
			code: "b = a",
			res: 0,
		},
		{
			code: "a | 1",
			res: 1,
		},
	} {
		
		got, err := interp.Exec(tc.code)
		if err != nil {
			t.Fatal(err)
		}

		if got != tc.res {
			t.Fatalf("got(%s) != expected(%s)", got, tc.res)
		}
	}
}

func TestEval(t *testing.T) {
	for _, tc := range []struct {
		expr Node
		res  Int
	}{
		{
			expr: BinExpr{
				Op:  OpAND,
				Lhs: Int(0),
				Rhs: Int(0),
			},
			res: 0,
		},
		{
			expr: BinExpr{
				Op:  OpAND,
				Lhs: Int(0),
				Rhs: Int(1),
			},
			res: 0,
		},
		{
			expr: BinExpr{
				Op:  OpAND,
				Lhs: Int(1),
				Rhs: Int(1),
			},
			res: 1,
		},
		{
			expr: BinExpr{
				Op:  OpAND,
				Lhs: Int(0xffff0000),
				Rhs: Int(0x0000ffff),
			},
			res: 0,
		},
		{
			expr: BinExpr{
				Op:  OpOR,
				Lhs: Int(0xffff0000),
				Rhs: Int(0x0000ffff),
			},
			res: 0xffffffff,
		},

		// groups
		{
			expr: BinExpr{
				Op: OpAND,
				Lhs: BinExpr{
					Op:  OpOR,
					Lhs: Int(0xffff0000),
					Rhs: Int(0x0000ffff),
				},
				Rhs: Int(0x000000ff),
			},
			res: 0x000000ff,
		},
		{
			expr: BinExpr{
				Op: OpOR,
				Lhs: BinExpr{
					Op:  OpOR,
					Lhs: Int(0x000000ff),
					Rhs: Int(0xff000000),
				},
				Rhs: BinExpr{
					Op:  OpOR,
					Lhs: Int(0x00ff0000),
					Rhs: Int(0x0000ff00),
				},
			},
			res: 0xffffffff,
		},
		{
			expr: UnaryExpr{
				Op:    OpNOT,
				Value: Int(7),
			},
			res: -8,
		},
	} {
		tc := tc
		t.Run(desc(tc.expr, tc.res), func(t *testing.T) {
			interp := NewInterp()
			got, err := interp.Eval(tc.expr)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.res {
				t.Fatalf("Fail: %d != %d", got, tc.res)
			}
		})
	}
}