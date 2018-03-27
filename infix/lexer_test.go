package infix_test

import (
	"github.com/madlambda/bwc/infix"
	"testing"
)

func consume(tokens chan infix.Tokval) []infix.Tokval {
	var toks []infix.Tokval
	for tok := range tokens {
		toks = append(toks, tok)
	}
	return toks
}

func TestLexer(t *testing.T) {
	for _, tc := range []struct {
		in  string
		out []infix.Tokval
	}{
		{
			in: "0",
			out: []infix.Tokval{
				{
					T: infix.Number,
					V: "0",
				},
			},
		},
		{
			in: "0123456789",
			out: []infix.Tokval{
				{
					T: infix.Number,
					V: "0123456789",
				},
			},
		},
		{
			in: "0&0",
			out: []infix.Tokval{
				{
					T: infix.Number,
					V: "0",
				},
				{
					T: infix.AND,
					V: "&",
				},
				{
					T: infix.Number,
					V: "0",
				},
			},
		},

		{
			in: "0&0|1",
			out: []infix.Tokval{
				{
					T: infix.Number,
					V: "0",
				},
				{
					T: infix.AND,
					V: "&",
				},
				{
					T: infix.Number,
					V: "0",
				},
				{
					T: infix.OR,
					V: "|",
				},
				{
					T: infix.Number,
					V: "1",
				},
			},
		},
	} {
		tc := tc
		got := consume(infix.Lex(tc.in))
		if len(got) != len(tc.out) {
			t.Logf("test data: %v", tc.in)
			t.Logf("got: %v", got)
			t.Fatalf("expect %d elems but got %d",
				len(tc.out), len(got))
		}

		for i := 0; i < len(tc.out); i++ {
			e := tc.out[i]
			g := got[i]
			if e.T != g.T {
				t.Fatalf("tok differs: %v != %v", e, g)
			}
			if e.V != g.V {
				t.Fatalf("tok differs: %s != %s", e.V, g.V)
			}
		}
	}
}