package parser

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/dcaiafa/go-expr/expr/internal/ast"
)

var keywords = map[string]int{
	"and":   kAND,
	"false": kFALSE,
	"in":    kIN,
	"not":   kNOT,
	"or":    kOR,
	"true":  kTRUE,
}

type lex struct {
	Program *ast.Program

	input *strings.Reader
	buf   bytes.Buffer
	pos   int
	err   error
}

func newLex(input string) *lex {
	return &lex{
		input: strings.NewReader(input),
	}
}

func (l *lex) Lex(lval *yySymType) int {
	return l.scan(lval)
}

func (l *lex) scan(lval *yySymType) int {
	for {
		r := l.read()
		if r == 0 {
			return 0
		}
		if isSpace(r) {
			continue
		}
		switch r {
		case '&':
			r = l.read()
			if r != '&' {
				return LEXERR
			}
			return AND
		case '|':
			r = l.read()
			if r != '|' {
				return LEXERR
			}
			return OR
		case '=':
			r = l.read()
			if r != '=' {
				return LEXERR
			}
			return EQ
		case '<':
			r = l.read()
			if r != '=' {
				l.unread()
				return '<'
			}
			return LE
		case '>':
			r = l.read()
			if r != '=' {
				l.unread()
				return '>'
			}
			return GE
		case '!':
			r = l.read()
			if r != '=' {
				l.unread()
				return '!'
			}
			return NE
		case '"':
			l.unread()
			return l.scanQuotedString(lval)
		case '+', '-', '*', '/', ';', '(', ')', ',', '[', ']':
			return int(r)
		default:
			if isNumber(r) {
				l.unread()
				return l.scanNumber(lval)
			} else if isLetter(r) || r == '_' {
				l.unread()
				return l.scanIdentifier(lval)
			} else {
				return LEXERR
			}
		}
	}
}

func (l *lex) scanIdentifier(lval *yySymType) int {
	l.buf.Reset()

	r := l.read()
	if !isLetter(r) && r != '_' {
		return LEXERR
	}
	l.buf.WriteRune(r)

	for {
		r := l.read()
		if !isLetter(r) && !isNumber(r) && r != '_' {
			l.unread()
			break
		}
		l.buf.WriteRune(r)
	}

	lval.str = l.buf.String()

	keyword, ok := keywords[lval.str]
	if ok {
		return keyword
	}

	return ID
}

func (l *lex) scanQuotedString(lval *yySymType) int {
	l.buf.Reset()

	if l.read() != '"' {
		return LEXERR
	}

	for {
		r := l.read()
		if r == '"' {
			break
		} else if r == '\\' {
			r = l.read()
			switch r {
			case '\\', '"':
				l.buf.WriteRune(r)
			default:
				return LEXERR
			}
		} else if r == '\n' || r == '\r' {
			return LEXERR
		} else {
			l.buf.WriteRune(r)
		}
	}

	lval.str = l.buf.String()
	return STRING
}

func (l *lex) scanNumber(lval *yySymType) int {
	l.buf.Reset()
	l.buf.WriteRune(l.read())

	var r rune
	for {
		r = l.read()
		if !isNumber(r) {
			break
		}
		l.buf.WriteRune(r)
	}

	if r == '.' {
		l.buf.WriteRune(r)
		r = l.read()
		if !isNumber(r) {
			return LEXERR
		}
		l.buf.WriteRune(r)
		for {
			r = l.read()
			if !isNumber(r) {
				l.unread()
				break
			}
			l.buf.WriteRune(r)
		}
	}

	l.unread()

	var err error
	lval.num, err = strconv.ParseFloat(l.buf.String(), 64)
	if err != nil {
		return LEXERR
	}

	return NUMBER

}

func (l *lex) read() rune {
	r, _, err := l.input.ReadRune()
	if err != nil {
		return 0
	}
	return r
}

func (l *lex) unread() {
	l.input.UnreadRune()
}

func (l *lex) Error(s string) {
	l.err = errors.New(s)
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
