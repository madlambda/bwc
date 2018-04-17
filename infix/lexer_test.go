package infix_test

import (
	"github.com/madlambda/bwc/infix"
	"testing"
)

type testcase struct {
	in  string
	out []infix.Tokval
}

func consume(tokens chan infix.Tokval) []infix.Tokval {
	var toks []infix.Tokval
	for tok := range tokens {
		toks = append(toks, tok)
	}
	return toks
}

func test(t *testing.T, tc testcase) {
	t.Helper()
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

func TestLexer(t *testing.T) {
	for _, tc := range []testcase{
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
			in: "0^0",
			out: []infix.Tokval{
				{
					T: infix.Number,
					V: "0",
				},
				{
					T: infix.XOR,
					V: "^",
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
		{
			in: "(0&0)",
			out: []infix.Tokval{
				{
					T: infix.LParen,
					V: "(",
				},
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
					T: infix.RParen,
					V: ")",
				},
			},
		},
		{
			in: "0xf",
			out: []infix.Tokval{
				{
					T: infix.Number,
					V: "0xf",
				},
			},
		},
		{
			in: "(1>>2)",
			out: []infix.Tokval{
				{
					T: infix.LParen,
					V: "(",
				},
				{
					T: infix.Number,
					V: "1",
				},
				{
					T: infix.SHR,
					V: ">>",
				},
				{
					T: infix.Number,
					V: "2",
				},
				{
					T: infix.RParen,
					V: ")",
				},
			},
		},
		{
			in: "a = 0",
			out: []infix.Tokval{
				{
					T: infix.Ident,
					V: "a",
				},
				{
					T: infix.Equal,
					V: "=",
				},
				{
					T: infix.Number,
					V: "0",
				},
			},
		},
		{
			in: "aa = aa",
			out: []infix.Tokval{
				{
					T: infix.Ident,
					V: "aa",
				},
				{
					T: infix.Equal,
					V: "=",
				},
				{
					T: infix.Ident,
					V: "aa",
				},
			},
		},
		{
			in: "_a = 0b10000",
			out: []infix.Tokval{
				{
					T: infix.Ident,
					V: "_a",
				},
				{
					T: infix.Equal,
					V: "=",
				},
				{
					T: infix.Number,
					V: "0b10000",
				},
			},
		},
		{
			in: "1invalid = 0b10000",
			out: []infix.Tokval{
				{
					T: infix.Illegal,
					V: "malformed number",
				},
			},
		},
	} {
		tc := tc
		test(t, tc)
	}
}

func TestLexerNumbers(t *testing.T) {
	for _, tc := range []testcase{
		{
			in: "01",
			out: []infix.Tokval{
				{
					T: infix.Number,
					V: "01",
				},
			},
		},
		{
			in: "97497239472938",
			out: []infix.Tokval{
				{
					T: infix.Number,
					V: "97497239472938",
				},
			},
		},
		{
			in: "0b",
			out: []infix.Tokval{
				{
					T: infix.Illegal,
					V: "malformed binary number",
				},
			},
		},
		{
			in: "0bf",
			out: []infix.Tokval{
				{
					T: infix.Illegal,
					V: "malformed binary number",
				},
			},
		},
		{
			in: "0b111112",
			out: []infix.Tokval{
				{
					T: infix.Illegal,
					V: "malformed number",
				},
			},
		},
		{
			in: "0xffg",
			out: []infix.Tokval{
				{
					T: infix.Illegal,
					V: "malformed number",
				},
			},
		},
		{
			in: "0b11111111",
			out: []infix.Tokval{
				{
					T: infix.Number,
					V: "0b11111111",
				},
			},
		},
	} {
		tc := tc
		test(t, tc)
	}
}