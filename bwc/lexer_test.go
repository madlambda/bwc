package bwc_test

import (
	"testing"

	"github.com/madlambda/bwc/bwc"
)

type testcase struct {
	in  string
	out []bwc.Tokval
}

func consume(tokens <-chan bwc.Tokval) []bwc.Tokval {
	var toks []bwc.Tokval
	for tok := range tokens {
		toks = append(toks, tok)
	}
	return toks
}

func test(t *testing.T, tc testcase) {
	t.Helper()
	got := consume(bwc.Lex(tc.in))
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
			out: []bwc.Tokval{
				{
					Type:  bwc.Number,
					Value: "0",
				},
			},
		},
		{
			in: "0123456789",
			out: []bwc.Tokval{
				{
					Type:  bwc.Number,
					Value: "0123456789",
				},
			},
		},
		{
			in: "0&0",
			out: []bwc.Tokval{
				{
					Type:  bwc.Number,
					Value: "0",
				},
				{
					Type:  bwc.AND,
					Value: "&",
				},
				{
					Type:  bwc.Number,
					Value: "0",
				},
			},
		},
		{
			in: "0^0",
			out: []bwc.Tokval{
				{
					Type:  bwc.Number,
					Value: "0",
				},
				{
					Type:  bwc.XOR,
					Value: "^",
				},
				{
					Type:  bwc.Number,
					Value: "0",
				},
			},
		},
		{
			in: "0&0|1",
			out: []bwc.Tokval{
				{
					Type:  bwc.Number,
					Value: "0",
				},
				{
					Type:  bwc.AND,
					Value: "&",
				},
				{
					Type:  bwc.Number,
					Value: "0",
				},
				{
					Type:  bwc.OR,
					Value: "|",
				},
				{
					Type:  bwc.Number,
					Value: "1",
				},
			},
		},
		{
			in: "(0&0)",
			out: []bwc.Tokval{
				{
					Type:  bwc.LParen,
					Value: "(",
				},
				{
					Type:  bwc.Number,
					Value: "0",
				},
				{
					Type:  bwc.AND,
					Value: "&",
				},
				{
					Type:  bwc.Number,
					Value: "0",
				},
				{
					Type:  bwc.RParen,
					Value: ")",
				},
			},
		},
		{
			in: "0xf",
			out: []bwc.Tokval{
				{
					Type:  bwc.Number,
					Value: "0xf",
				},
			},
		},
		{
			in: "(1>>2)",
			out: []bwc.Tokval{
				{
					Type:  bwc.LParen,
					Value: "(",
				},
				{
					Type:  bwc.Number,
					Value: "1",
				},
				{
					Type:  bwc.SHR,
					Value: ">>",
				},
				{
					Type:  bwc.Number,
					Value: "2",
				},
				{
					Type:  bwc.RParen,
					Value: ")",
				},
			},
		},
		{
			in: "a = 0",
			out: []bwc.Tokval{
				{
					Type:  bwc.Ident,
					Value: "a",
				},
				{
					Type:  bwc.Equal,
					Value: "=",
				},
				{
					Type:  bwc.Number,
					Value: "0",
				},
			},
		},
		{
			in: "aa = aa",
			out: []bwc.Tokval{
				{
					Type:  bwc.Ident,
					Value: "aa",
				},
				{
					Type:  bwc.Equal,
					Value: "=",
				},
				{
					Type:  bwc.Ident,
					Value: "aa",
				},
			},
		},
		{
			in: "_a = 0b10000",
			out: []bwc.Tokval{
				{
					Type:  bwc.Ident,
					Value: "_a",
				},
				{
					Type:  bwc.Equal,
					Value: "=",
				},
				{
					Type:  bwc.Number,
					Value: "0b10000",
				},
			},
		},
		{
			in: "1invalid = 0b10000",
			out: []bwc.Tokval{
				{
					Type:  bwc.Illegal,
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
			out: []bwc.Tokval{
				{
					Type:  bwc.Number,
					Value: "01",
				},
			},
		},
		{
			in: "97497239472938",
			out: []bwc.Tokval{
				{
					Type:  bwc.Number,
					Value: "97497239472938",
				},
			},
		},
		{
			in: "0b",
			out: []bwc.Tokval{
				{
					Type:  bwc.Illegal,
					Value: "malformed binary number",
				},
			},
		},
		{
			in: "0bf",
			out: []bwc.Tokval{
				{
					Type:  bwc.Illegal,
					Value: "malformed binary number",
				},
			},
		},
		{
			in: "0b111112",
			out: []bwc.Tokval{
				{
					Type:  bwc.Illegal,
					Value: "malformed number",
				},
			},
		},
		{
			in: "0xffg",
			out: []bwc.Tokval{
				{
					Type:  bwc.Illegal,
					Value: "malformed number",
				},
			},
		},
		{
			in: "0b11111111",
			out: []bwc.Tokval{
				{
					Type:  bwc.Number,
					Value: "0b11111111",
				},
			},
		},
	} {
		tc := tc
		test(t, tc)
	}
}
