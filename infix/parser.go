package infix

import (
	"fmt"
	"strconv"
)

type parser struct {
	tokens chan Tokval
}

func Parse(code string) (*Expr, error) {
	p := &parser{
		tokens: Lex(code),
	}
	return p.parseExpr()
}

func parserErr(tok Tokval) error {
	return fmt.Errorf("unexpected %s", tok)
}

func (p *parser) parseExpr() (*Expr, error) {
	ltok := <-p.tokens
	if ltok.T != Number {
		return nil, parserErr(ltok)
	}

	optok := <-p.tokens
	op, ok := validOperation(optok.T)
	if !ok {
		return nil, parserErr(optok)
	}

	rtok := <-p.tokens
	if rtok.T != Number {
		return nil, parserErr(rtok)
	}

	lhs, err := parseNum(ltok.V)
	if err != nil {
		return nil, err
	}

	rhs, err := parseNum(rtok.V)
	if err != nil {
		return nil, err
	}

	return &Expr{
		Op:  op,
		Lhs: lhs,
		Rhs: rhs,
	}, nil
}

func parseNum(num string) (Int, error) {
	i, err := strconv.Atoi(num)
	return Int(i), err
}

func validOperation(tok Token) (Optype, bool) {
	switch tok {
	case AND:
		return OpAND, true
	case OR:
		return OpOR, true
	}

	return -1, false
}