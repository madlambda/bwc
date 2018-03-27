package infix

type (
	Token int
)

const (
	Number Token = iota+1
	LParen
	RParen
	OR
	AND
	Illegal
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
	return "<invalid>"
}