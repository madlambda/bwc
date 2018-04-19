package infix

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type (
	lexer struct {
		input  string // expression being lex'ed
		start  int    // start position of token
		pos    int    // pos in the input
		width  int    // width of last rune
		tokens chan Tokval
	}

	stateFn func(*lexer) stateFn

	// Tokval is the token type + value
	Tokval struct {
		Type  Token
		Value string
		Pos   int
	}
)

// eof mimics the EOF but for the input string.
// The name isn't eoi to avoid confusion.
const eof = -1

func (t Tokval) String() string {
	return fmt.Sprintf("Token(%s, %s)", t.Type, t.Value)
}

// Lex creates a concurrent lexer and returns a
// channel of tokens processed from input.
func Lex(input string) chan Tokval {
	l := &lexer{
		input:  input,
		tokens: make(chan Tokval),
	}

	go l.run()
	return l.tokens
}

// run the state machine
func (l *lexer) run() {
	for state := lexStart; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// emit a token.
func (l *lexer) emit(tok Token) {
	l.tokens <- Tokval{
		Type:   tok,
		Value:   l.input[l.start:l.pos],
		Pos: l.start,
	}
	l.start = l.pos
}

// errorf emits an illegal token. This token carries
// the lexer error.
func (l *lexer) errorf(msg string, args ...interface{}) stateFn {
	l.tokens <- Tokval{
		Type: Illegal,
		Value: fmt.Sprintf(msg, args...),
		Pos: l.start,
	}
	return nil
}

// peek looks for the next rune from the input
// but do not increases the cursor.
// It returns eof when reaches the end of input.
func (l *lexer) peek() rune {
	r := l.next()
	if r == eof {
		return eof
	}
	l.backup()
	return r
}

// next returns the next rune from input or
// eof if reached end of input.
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		// zeroes the width to simplify handling at
		// end of input. Some functions like acceptRun*
		// do a explicit rewind/backup of last processed
		// character, but in case of eof we dont want
		// that.
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// backup rollback one rune.
func (l *lexer) backup() {
	l.pos -= l.width
}

// ignore bytes processed since last token.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept checks if next is in the valid string.
func (l *lexer) accept(valid string) bool {
	return strings.IndexRune(valid, l.peek()) != -1
}

// acceptRun consumes next runes if valid.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}

	l.backup()
}

// acceptRunfn is like acceptRun but uses fn as
// function comparator.
func (l *lexer) acceptRunfn(fn func(r rune) bool) {
	for fn(l.next()) {

	}

	l.backup()
}

// lexStart is the start of state machine
func lexStart(l *lexer) stateFn {
	r := l.next()
	switch {
	case isIdentBegin(r):
		l.acceptRunfn(func(r rune) bool {
			return unicode.IsLetter(r) ||
				unicode.IsDigit(r) ||
				r == '_'
		})
		l.emit(Ident)
		return lexStart
	case unicode.IsSpace(r):
		l.acceptRunfn(unicode.IsSpace)
		l.ignore()
		return lexStart
	case r == eof:
		return nil
	case r >= '0' && r <= '9':
		l.backup()
		return lexNumber
	case r == '|':
		l.emit(OR)
		return lexStart
	case r == '&':
		l.emit(AND)
		return lexStart
	case r == '^':
		l.emit(XOR)
		return lexStart
	case r == '~':
		l.emit(NOT)
		return lexStart
	case r == '<':
		next := l.next()
		if next != '<' {
			return l.errorf("unexpected %q", next)
		}
		l.emit(SHL)
		return lexStart
	case r == '>':
		next := l.next()
		if next != '>' {
			return l.errorf("unexpected %q", next)
		}
		l.emit(SHR)
		return lexStart
	case r == '(':
		l.emit(LParen)
		return lexStart
	case r == ')':
		l.emit(RParen)
		return lexStart
	case r == '=':
		l.emit(Equal)
		return lexStart
	default:
		return l.errorf("Unexpected %q at %d", r, l.pos)
	}
	return nil
}

// lexer of dec, hex and bin numbers.
func lexNumber(l *lexer) stateFn {
	r := l.next()
	next := l.peek()
	if r == '0' && next == 'b' {
		// 0bnnnnnnnn
		l.next()
		digits := "01"
		if !l.accept(digits) {
			return l.errorf("malformed binary number")
		}

		l.acceptRun(digits)
	} else if r == '0' && next == 'x' {
		// 0xnnnnnnnn
		l.next()
		digits := "0123456789abcdef"
		if !l.accept(digits) {
			return l.errorf("malformed hex number")
		}

		l.acceptRun(digits)
	} else {
		// decimal
		l.acceptRun("0123456789")
	}

	if isAlphaNumeric(l.peek()) {
		return l.errorf("malformed number")
	}

	l.emit(Number)
	return lexStart
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isIdentBegin tells if r could be the first rune of an ident.
func isIdentBegin(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}