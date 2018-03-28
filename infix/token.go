package infix

import (
	"strconv"
)

type (
	Token int
)

const (
	Illegal Token = iota
	Number
	LParen
	RParen
	OR
	AND
	SHL
	SHR
	EOF
)

func (t Token) String() string {
	switch t {
	case LParen: return "("
	case RParen: return ")"
	case OR: return "|"
	case AND: return "&"
	case Illegal: return "<ileggal>"
	case EOF: return "EOF"
	}
	return strconv.Itoa(int(t))
}