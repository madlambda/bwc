package infix

import (
	"testing"
	"reflect"
)

func TestParserBinExpr(t *testing.T) {
	for _, tc := range []struct{
		code string
		ast BinExpr
		err error
	} {
		{
			code: "0|1",
			ast: BinExpr{
				Op: OpOR,
				Lhs: Int(0),
				Rhs: Int(1),
			},
		},
		{
			code: "0|1|2",
			ast: BinExpr{
				Op: OpOR,
				Lhs: BinExpr{
					Op: OpOR,
					Lhs: Int(0),
					Rhs: Int(1),
				},
				Rhs: Int(2),
			},
		},
	} {
		got, err := Parse(tc.code)
		if err != tc.err {
			t.Fatalf("expected(%s) but got(%s)", tc.err, err)
		}

		if got.Type() != tc.ast.Type() {
			t.Fatalf("expected type(%s) but got(%s) -> %s", tc.ast.Type(),
				got.Type(), got)
		}

		if !reflect.DeepEqual(got, tc.ast) {
			t.Fatalf("node differs: (%s) != (%s)", got, tc.ast)
		}		
	}
}