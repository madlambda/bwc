package infix 

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
	XOR
	NOT
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
	case XOR: return "^"
	case NOT: return "~"
	case SHL: return "<<"
	case SHR: return ">>"
	case Illegal: return "<ileggal>"
	case EOF: return "EOF"
	}
	panic("invalid token")
}