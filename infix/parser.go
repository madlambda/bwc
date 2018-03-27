package infix

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type (
	Lexer struct {
		input  string // expression being lex'ed
		start  int    // start position of token
		pos    int    // pos in the input
		width  int    // width of last rune
		tokens chan Tokval
	}

	stateFn func(*Lexer) stateFn

	// Tokval is the token type + value
	Tokval struct {
		T Token
		V string
	}
)

const eof = -1

func Lex(input string) chan Tokval {
	l := &Lexer{
		input:  input,
		tokens: make(chan Tokval),
	}

	go l.run()
	return l.tokens
}

func (l *Lexer) run() {
	for state := lexExpr; state != nil; {
		state = state(l)
	}
}

func (l *Lexer) emit(tok Token) {
	l.tokens <- Tokval{
		T: tok,
		V: l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *Lexer) errorf(msg string, args ...interface{}) {
	l.tokens <- Tokval{
		T: Illegal,
		V: fmt.Sprintf(msg, args...),
	}
	l.start = len(l.input)
	l.pos = l.start
}

func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func lexExpr(l *Lexer) stateFn {
	r := l.next()
	switch {
	case unicode.IsDigit(r):
		l.acceptRun("xb0123456789")
		l.emit(Number)
		return nil
	default:
		l.errorf("Unexpected %q at %d", r, l.pos)
	}
	return nil
}

func Parse(code string) (Expr, error) {
	for tok := range Lex(code) {
		fmt.Printf("tok: %s\n", tok)
	}
	return Expr{
		Op:  1,
		Lhs: Int(1),
		Rhs: Int(1),
	}, nil
}