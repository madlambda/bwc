package infix

import (
	"fmt"
	"strings"
	"unicode/utf8"
	"unicode"
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

func (t Tokval) String() string {
	return fmt.Sprintf("Token(%s, %s)", t.T, t.V)
}

func Lex(input string) chan Tokval {
	l := &Lexer{
		input:  input,
		tokens: make(chan Tokval),
	}

	go l.run()
	return l.tokens
}

func (l *Lexer) run() {
	for state := lexStart; state != nil; {
		state = state(l)
	}
	close(l.tokens)
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

func (l *Lexer) current() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
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

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	
	}

	// no rewind in case of end of input
	if l.current() != eof {
		l.backup()
	}
}

func (l *Lexer) acceptRunfn(fn func (r rune) bool) {
	for fn(l.next()) {

	}

	// no rewind in case of end of input
	if l.current() != eof {
		l.backup()
	}
}

func lexStart(l *Lexer) stateFn {
	r := l.next()
	switch {
	case unicode.IsSpace(r):
		l.acceptRunfn(unicode.IsSpace)
		l.ignore()
		return lexStart
	case r == eof:
		return nil
	case r >= '0' || r <= '9':
		l.acceptRun("0123456789")
		l.emit(Number)
		return lexStart
	case r == '|':
		l.emit(OR)
		return lexStart
	case r == '&':
		l.emit(AND)
		return lexStart
	case r == '(':
		l.emit(LParen)
		return lexStart
	case r == ')':
		l.emit(RParen)
		return lexStart
	default:
		l.errorf("Unexpected %q at %d", r, l.pos)
	}
	return nil
}