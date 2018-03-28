package infix

import (
	"fmt"
	"strconv"
)

type parser struct {
	tokens chan Tokval
	pocket *Tokval // lookahead
}

var TokEOF = Tokval{
	T: EOF,
	V: "<eof>",
}

func parserErr(tok Tokval) error {
	return fmt.Errorf("unexpected %s", tok)
}

func Parse(code string) (Node, error) {
	p := &parser{
		tokens: Lex(code),
	}
	return p.parseExpr()
}

func (p *parser) peek() Tokval {
	if p.pocket != nil {
		val := *p.pocket
		p.pocket = nil
		return val
	}

	val, ok := <-p.tokens
	if !ok {
		return TokEOF
	}
	p.pocket = &val
	return val
}

func (p *parser) next() Tokval {
	if p.pocket != nil {
		val := *p.pocket
		p.pocket = nil
		return val
	}
	tok, ok := <-p.tokens
	if !ok {
		return TokEOF
	}
	return tok
}

func (p *parser) parseExpr() (Node, error) {
	var eoferr = func(expect string) error {
		return fmt.Errorf("premature eof, expects %s",
			expect)
	}

	tok := p.peek()
	if tok.T == EOF {
		return nil, eoferr("expression")
	}

	// left hand side of expr

	var lhs Node
	var err error
	var eof bool

	var hasparens bool
	if tok.T == LParen {
		hasparens = true
		p.next()
		lhs, err = p.parseExpr()
	} else if tok.T == NOT {
		lhs, eof, err = p.parseUnary()
	} else {
		lhs, eof, err = p.parseNum()
	}

	if err != nil {
		return nil, err
	}
	if eof {
		return nil, eoferr("expr || number")
	}

	if hasparens {
		tok = p.next()
		if tok.T != RParen {
			return nil, parserErr(tok)
		}
	}

	op, eof, err := p.parseBinOP()
	if err != nil {
		return nil, err
	}
	if eof {
		return lhs, nil
	}

	// right hand side of expr
	var rhs Node
	hasparens = false

	tok = p.peek()
	if tok.T == EOF {
		return nil, eoferr("rhs of expr")
	}
 
	if tok.T == LParen {
		hasparens = true
		p.next()
		rhs, err = p.parseExpr()
	} else if tok.T == NOT {
		lhs, eof, err = p.parseUnary()
	} else {
		rhs, eof, err = p.parseNum()
	}

	if err != nil {
		return nil, err
	}
	if eof {
		return nil, eoferr("number")
	}

	if hasparens {
		tok := p.next()
		if tok.T != RParen {
			return nil, parserErr(tok)
		}
	}

	return BinExpr{
		Op:  op,
		Lhs: lhs,
		Rhs: rhs,
	}, nil
}

func (p *parser) parseNum() (a Int, eof bool, err error) {
	tok := p.next()
	if tok.T == EOF {
		return 0, true, nil
	}

	if tok.T != Number {
		return 0, false, parserErr(tok)
	}

	intstr := tok.V
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

func (p *parser) parseUnary() (n Node, eof bool, err error) {
	tok := p.next()
	
	var val UnaryExpr
	if tok.T == NOT {
		val.Op = OpNOT
	} else {
		return nil, false, fmt.Errorf("invalid unary: %q", tok.V)
	}

	num, eof, err := p.parseNum()
	if err != nil {
		return nil, false, err
	}
	if eof {
		return nil, true, nil
	}

	val.Value = num
	return val, false, nil
}

func (p *parser) parseBinOP() (a Optype, eof bool, err error) {
	optok := p.next()
	if optok.T == EOF {
		return 0, true, nil
	}

	op, ok := validBinOP(optok.T)
	if !ok {
		return 0, false, parserErr(optok)
	}
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