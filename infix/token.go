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