//go:generate go run golang.org/x/tools/cmd/goyacc@v0.18.0 -o parser.go jabl.y
package jabl

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"text/scanner"
	"unicode"
)

type program struct {
	body stmt
}

type stmt interface{}

type blockStmt struct {
	stmt stmt
}

type seqStmt struct {
	first stmt
	rest  stmt
}

type fnStmt struct {
	fn    int
	expr  expr
	expr2 expr
	block stmt
}

type ifStmt struct {
	cond  expr
	block stmt
	other stmt
}

type expr interface{}

type parenExpr struct {
	expr expr
}

type notExpr struct {
	expr expr
}

type cmpExpr struct {
	op    int
	t     int
	left  expr
	right expr
}

type mathExpr struct {
	op    int
	left  expr
	right expr
}

type rollExpr struct {
	num   float64
	sides float64
}

type lexer struct {
	s   scanner.Scanner
	err error
	ast *program
}

func newLexer(src io.Reader) *lexer {
	var s scanner.Scanner
	s.Init(src)
	// Accept tokens with "-"
	s.IsIdentRune = func(ch rune, i int) bool {
		return unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0 || ch == '-' && i > 0
	}
	return &lexer{
		s: s,
	}
}

// The parser expects the lexer to return 0 on EOF.  Give it a name for clarity.
const EOF = 0

func (l *lexer) Lex(lval *yySymType) int {
	if token := l.s.Scan(); token == scanner.EOF {
		return EOF
	}
	text := l.s.TokenText()
	switch text {
	case "{":
		return START
	case "}":
		return END
	case "=":
		if l.s.Peek() == '=' {
			l.s.Scan()
			return CMP_EQ
		}
	case "<":
		if l.s.Peek() == '=' {
			l.s.Scan()
			return CMP_LTE
		}
		return CMP_LT
	case ">":
		if l.s.Peek() == '=' {
			l.s.Scan()
			return CMP_GTE
		}
		return CMP_GT
	case "!":
		if l.s.Peek() == '=' {
			l.s.Scan()
			return CMP_NEQ
		}
		return CMP_NOT
	case "&":
		if l.s.Peek() == '&' {
			l.s.Scan()
			return CMP_AND
		}
	case "|":
		if l.s.Peek() == '|' {
			l.s.Scan()
			return CMP_OR
		}
	case "-":
		// check if this is a negative number
		if unicode.IsDigit(rune(l.s.Peek())) {
			l.s.Scan()
			text = "-" + l.s.TokenText()

			numVal, err := strconv.ParseFloat(text, 64)
			if err != nil {
				l.Error(fmt.Sprintf("invalid number: %s", text))
				return EOF
			}
			lval.Number = numVal
			return NUMBER
		} else {
			// or the subtraction operator
			lval.String = text
			return int('-')
		}
	case "+", "*", "/", "(", ")", ",", ".":
		lval.String = text
		// Use the value of the operator itself as the identifier
		return int(text[0])
	case "if":
		return IF
	case "else":
		return ELSE
	case "get":
		return GET
	case "set":
		return SET
	case "print":
		return PRINT
	case "choice":
		return CHOICE
	case "goto":
		return GOTO
	default:
		if text[0] == '"' && text[len(text)-1] == '"' {
			// trim the start and end quotes and add a newline
			lval.String = text[1 : len(text)-1]
			return STRING
		} else if text == "true" {
			lval.Boolean = true
			return BOOLEAN
		} else if text == "false" {
			lval.Boolean = false
			return BOOLEAN
		} else if rune(text[0]) == 'd' {
			// check if this is a dice roll
			rollVal, err := strconv.ParseFloat(text[1:], 64)
			if err != nil {
				l.Error(fmt.Sprintf("invalid dice roll: %s", text))
				return EOF
			}
			lval.Number = rollVal
			return DICE
		} else if unicode.IsDigit(rune(text[0])) {
			numVal, err := strconv.ParseFloat(text, 64)
			if err != nil {
				l.Error(fmt.Sprintf("invalid number: %s", text))
				return EOF
			}
			lval.Number = numVal
			return NUMBER
		}
	}

	l.Error(fmt.Sprintf("unexpected token: %s", text))
	return EOF
}

func (l *lexer) Error(msg string) {
	l.err = errors.Join(l.err, fmt.Errorf("lex error: %s", msg))
}
