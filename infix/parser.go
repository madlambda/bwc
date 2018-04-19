package infix

import (
	"fmt"
	"strconv"
)

type parser struct {
	tokens    chan Tokval
	lookahead []Tokval
}

var TokEOF = Tokval{
	Type:  EOF,
	Value: "<eof>",
	Pos: -1,
}

func eoferr(expect string) error {
	return fmt.Errorf("premature eof, expects %s",
		expect)
}

func parserErr(expected string, tok Tokval) error {
	return fmt.Errorf("expected %s but got %s at position %d", 
			expected, tok, tok.Pos)
}

func Parse(code string) (Node, error) {
	p := &parser{
		tokens: Lex(code),
	}

	return p.parse()
}

// scry foretell the future using a crystal ball. Amount is how much
// of the future you want to foresee.
func (p *parser) scry(amount int) []Tokval {
	if amount > 2 {
		panic("lookahead > 2")
	}

	sz := len(p.lookahead)
	for i := 0; i < amount-sz; i++ {
		val, ok := <-p.tokens
		if !ok {
			val = TokEOF
		}

		p.lookahead = append(p.lookahead, val)
	}

	return p.lookahead
}

// forget what you had foresee
func (p *parser) forget(amount int) {
	for i := 0; i < amount; i++ {
		p.lookahead = p.lookahead[1:]
	}
}

// next returns the next token and consume it.
func (p *parser) next() Tokval {
	if len(p.lookahead) > 0 {
		tok := p.lookahead[0]
		p.forget(1)
		return tok
	}

	tok, ok := <-p.tokens
	if !ok {
		return TokEOF
	}
	return tok
}

func (p *parser) parse() (Node, error) {
	toks := p.scry(2)
	if toks[0].Type == EOF {
		return nil, eoferr("assign || expr")
	}

	// <ident>
	// <ident> = <expr>
	// <ident> <op> <expr>
	// requires one lookahead

	ident := toks[0]
	next := toks[1]
	if ident.Type == Ident && next.Type == Equal {
		return p.parseAssign()
	}

	return p.parseExpr()
}

func (p *parser) parseAssign() (Node, error) {
	p.scry(2)

	id := p.lookahead[0]
	eq := p.lookahead[1]

	if id.Type != Ident {
		return nil, parserErr("IDENT", id)
	}
	if eq.Type != Equal {
		return nil, parserErr("EQUAL", eq)
	}

	p.forget(2)

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return Assign{
		Varname: id.Value,
		Expr:    expr,
	}, nil
}

func (p *parser) parseOperand() (n Node, eof bool, err error) {
	var hasparens bool

	p.scry(1)

	tok := p.lookahead[0]
	if tok.Type == EOF {
		return nil, true, eoferr("expr || number || ident || unary")
	}

	if tok.Type == LParen {
		hasparens = true
		p.forget(1)
		n, err = p.parseExpr()
	} else if tok.Type == NOT {
		n, err = p.parseUnary()
	} else if tok.Type == Ident {
		n = Var(tok.Value)
		p.forget(1)
	} else {
		n, eof, err = p.parseNum()
	}

	if hasparens {
		tok = p.next()
		if tok.Type != RParen {
			return nil, tok.Type == EOF, parserErr("RPAREN", tok)
		}
	}

	if err != nil {
		return nil, false, err
	}
	if eof {
		return n, true, nil
	}

	return n, false, nil
}

func (p *parser) parseExpr() (Node, error) {
	// left hand side of expr
	lhs, eof, err := p.parseOperand()
	if err != nil {
		return nil, err
	}

	if eof {
		return lhs, nil
	}

	p.scry(1)
	tok := p.lookahead[0]
	if tok.Type == RParen {
		return lhs, nil
	}

	op, eof, err := p.parseBinOP()
	if err != nil {
		return nil, err
	}

	if eof {
		return lhs, nil
	}

	// right hand side of expr
	rhs, eof, err := p.parseOperand()
	if err != nil {
		return nil, err
	}

	expr := BinExpr{
		Op:  op,
		Lhs: lhs,
		Rhs: rhs,
	}

	if eof {
		return expr, nil
	}

	// additional, non grouped, expressions
	for !eof {
		op, eof, err = p.parseBinOP()
		if eof || err != nil {
			return expr, nil
		}

		rhs, err = p.parseExpr()
		if err != nil {
			return nil, err
		}

		// 0|1&2|3 == (((0|1)&2)|3)
		// operation order is left to right
		expr = BinExpr{
			Op:  op,
			Lhs: expr,
			Rhs: rhs,
		}
	}
	return expr, nil
}

func (p *parser) parseNum() (a Int, eof bool, err error) {
	tok := p.next()
	if tok.Type == EOF {
		return 0, true, eoferr("expected number")
	}

	if tok.Type != Number {
		return 0, false, parserErr("NUMBER", tok)
	}

	intstr := tok.Value
	if len(intstr) > 2 {
		if intstr[1] == 'b' {
			val, err := strconv.ParseInt(intstr[2:], 2, 64)
			return Int(val), false, err
		}

		if intstr[1] == 'x' {
			val, err := strconv.ParseInt(intstr[2:], 16, 64)
			return Int(val), false, err
		}
	}

	val, err := strconv.ParseInt(intstr, 10, 64)
	return Int(val), false, err
}

func (p *parser) parseUnary() (n Node, err error) {
	tok := p.next()

	var val UnaryExpr
	if tok.Type == NOT {
		val.Op = OpNOT
	} else {
		return nil, fmt.Errorf("invalid unary: %q", tok.Value)
	}

	num, eof, err := p.parseNum()
	if err != nil {
		return nil, err
	}
	if eof {
		return nil, fmt.Errorf("expected number")
	}

	val.Value = num
	return val, nil
}

func (p *parser) parseBinOP() (a Optype, eof bool, err error) {
	p.scry(1)

	optok := p.lookahead[0]
	if optok.Type == EOF {
		return 0, true, nil
	}

	op, ok := validBinOP(optok.Type)
	if !ok {
		return 0, false, parserErr("OPERATION", optok)
	}

	p.forget(1)
	return op, false, nil
}

func validBinOP(tok Token) (Optype, bool) {
	switch tok {
	case AND:
		return OpAND, true
	case OR:
		return OpOR, true
	case XOR:
		return OpXOR, true
	case SHL:
		return OpSHL, true
	case SHR:
		return OpSHR, true
	}

	return -1, false
}