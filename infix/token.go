package infix

import "fmt"

type (
	Token int
)

const (
	Illegal Token = iota
	Ident
	Number
	LParen
	RParen
	Equal
	OR
	AND
	XOR
	NOT
	SHL
	SHR
	EOF
)

func (t Token) String() string {
	switch t {
	case Ident:
		return "IDENT"
	case Number:
		return "NUMBER"
	case LParen:
		return "("
	case RParen:
		return ")"
	case Equal:
		return "="
	case OR:
		return "|"
	case AND:
		return "&"
	case XOR:
		return "^"
	case NOT:
		return "~"
	case SHL:
		return "<<"
	case SHR:
		return ">>"
	case Illegal:
		return "<ileggal>"
	case EOF:
		return "EOF"
	}
	panic(fmt.Sprintf("invalid token: %d", t))
}