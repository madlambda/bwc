package infix

import (
	"reflect"
	"testing"
)

type testcase struct {
	code string
	ast  Node
	err  error
}

func test(t *testing.T, tc testcase) {
	t.Helper()

	got, err := Parse(tc.code)
	if err != tc.err {
		t.Fatalf("expected(%v) but got(%s)", tc.err, err)
	}

	if got.Type() != tc.ast.Type() {
		t.Fatalf("expected type(%s) but got(%s) -> %s", tc.ast.Type(),
			got.Type(), got)
	}

	if !reflect.DeepEqual(got, tc.ast) {
		t.Fatalf("node differs: (%s) != (%s)", got, tc.ast)
	}
}

func TestParseAssign(t *testing.T) {
	for _, tc := range []testcase{
		{
			code: "a = 1",
			ast: Assign{
				Varname: "a",
				Expr:    Int(1),
			},
		},
		{
			code: "a = a",
			ast: Assign{
				Varname: "a",
				Expr:    Var("a"),
			},
		},
		{
			code: "a = a | b",
			ast: Assign{
				Varname: "a",
				Expr: BinExpr{
					Lhs: Var("a"),
					Rhs: Var("b"),
					Op:  OpOR,
				},
			},
		},
	} {

		test(t, tc)
	}
}

func TestParserBinExpr(t *testing.T) {
	for _, tc := range []testcase{
		{
			code: "0|1",
			ast: BinExpr{
				Op:  OpOR,
				Lhs: Int(0),
				Rhs: Int(1),
			},
		},
		{
			code: "0|1|2",
			ast: BinExpr{
				Op: OpOR,
				Lhs: BinExpr{
					Op:  OpOR,
					Lhs: Int(0),
					Rhs: Int(1),
				},
				Rhs: Int(2),
			},
		},
		{
			code: "a|1",
			ast: BinExpr{
				Op:  OpOR,
				Lhs: Var("a"),
				Rhs: Int(1),
			},
		},
		{
			code: "(a|1)",
			ast: BinExpr{
				Op:  OpOR,
				Lhs: Var("a"),
				Rhs: Int(1),
			},
		},
		{
			code: "((a|1))",
			ast: BinExpr{
				Op:  OpOR,
				Lhs: Var("a"),
				Rhs: Int(1),
			},
		},
		{
			code: "((a|1)|2)",
			ast: BinExpr{
				Op: OpOR,
				Lhs: BinExpr{
					Op:  OpOR,
					Lhs: Var("a"),
					Rhs: Int(1),
				},
				Rhs: Int(2),
			},
		},
	} {
		test(t, tc)
	}
}