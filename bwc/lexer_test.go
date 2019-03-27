package bwc_test

import (
	"testing"

	"github.com/madlambda/bwc/infix"
)

type testcase struct {
	in  string
	out []infix.Tokval
}

func consume(tokens <-chan infix.Tokval) []infix.Tokval {
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
		if e.Type != g.Type {
			t.Fatalf("tok differs: %v != %v", e, g)
		}
		if e.Value != g.Value {
			t.Fatalf("tok differs: %s != %s", e.Value, g.Value)
		}
	}
}

func TestLexer(t *testing.T) {
	for _, tc := range []testcase{
		{
			in: "0",
			out: []infix.Tokval{
				{
					Type:  infix.Number,
					Value: "0",
				},
			},
		},
		{
			in: "0123456789",
			out: []infix.Tokval{
				{
					Type:  infix.Number,
					Value: "0123456789",
				},
			},
		},
		{
			in: "0&0",
			out: []infix.Tokval{
				{
					Type:  infix.Number,
					Value: "0",
				},
				{
					Type:  infix.AND,
					Value: "&",
				},
				{
					Type:  infix.Number,
					Value: "0",
				},
			},
		},
		{
			in: "0^0",
			out: []infix.Tokval{
				{
					Type:  infix.Number,
					Value: "0",
				},
				{
					Type:  infix.XOR,
					Value: "^",
				},
				{
					Type:  infix.Number,
					Value: "0",
				},
			},
		},
		{
			in: "0&0|1",
			out: []infix.Tokval{
				{
					Type:  infix.Number,
					Value: "0",
				},
				{
					Type:  infix.AND,
					Value: "&",
				},
				{
					Type:  infix.Number,
					Value: "0",
				},
				{
					Type:  infix.OR,
					Value: "|",
				},
				{
					Type:  infix.Number,
					Value: "1",
				},
			},
		},
		{
			in: "(0&0)",
			out: []infix.Tokval{
				{
					Type:  infix.LParen,
					Value: "(",
				},
				{
					Type:  infix.Number,
					Value: "0",
				},
				{
					Type:  infix.AND,
					Value: "&",
				},
				{
					Type:  infix.Number,
					Value: "0",
				},
				{
					Type:  infix.RParen,
					Value: ")",
				},
			},
		},
		{
			in: "0xf",
			out: []infix.Tokval{
				{
					Type:  infix.Number,
					Value: "0xf",
				},
			},
		},
		{
			in: "(1>>2)",
			out: []infix.Tokval{
				{
					Type:  infix.LParen,
					Value: "(",
				},
				{
					Type:  infix.Number,
					Value: "1",
				},
				{
					Type:  infix.SHR,
					Value: ">>",
				},
				{
					Type:  infix.Number,
					Value: "2",
				},
				{
					Type:  infix.RParen,
					Value: ")",
				},
			},
		},
		{
			in: "a = 0",
			out: []infix.Tokval{
				{
					Type:  infix.Ident,
					Value: "a",
				},
				{
					Type:  infix.Equal,
					Value: "=",
				},
				{
					Type:  infix.Number,
					Value: "0",
				},
			},
		},
		{
			in: "aa = aa",
			out: []infix.Tokval{
				{
					Type:  infix.Ident,
					Value: "aa",
				},
				{
					Type:  infix.Equal,
					Value: "=",
				},
				{
					Type:  infix.Ident,
					Value: "aa",
				},
			},
		},
		{
			in: "_a = 0b10000",
			out: []infix.Tokval{
				{
					Type:  infix.Ident,
					Value: "_a",
				},
				{
					Type:  infix.Equal,
					Value: "=",
				},
				{
					Type:  infix.Number,
					Value: "0b10000",
				},
			},
		},
		{
			in: "1invalid = 0b10000",
			out: []infix.Tokval{
				{
					Type:  infix.Illegal,
					Value: "malformed number",
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
					Type:  infix.Number,
					Value: "01",
				},
			},
		},
		{
			in: "97497239472938",
			out: []infix.Tokval{
				{
					Type:  infix.Number,
					Value: "97497239472938",
				},
			},
		},
		{
			in: "0b",
			out: []infix.Tokval{
				{
					Type:  infix.Illegal,
					Value: "malformed binary number",
				},
			},
		},
		{
			in: "0bf",
			out: []infix.Tokval{
				{
					Type:  infix.Illegal,
					Value: "malformed binary number",
				},
			},
		},
		{
			in: "0b111112",
			out: []infix.Tokval{
				{
					Type:  infix.Illegal,
					Value: "malformed number",
				},
			},
		},
		{
			in: "0xffg",
			out: []infix.Tokval{
				{
					Type:  infix.Illegal,
					Value: "malformed number",
				},
			},
		},
		{
			in: "0b11111111",
			out: []infix.Tokval{
				{
					Type:  infix.Number,
					Value: "0b11111111",
				},
			},
		},
	} {
		tc := tc
		test(t, tc)
	}
}
