package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// token holds all tokens to be parse.
type token int8

const (
	// EOF presents the end of a file
	EOF token = iota

	// INVALID token
	INVALID

	// NEWLINE \n
	NEWLINE

	// WHITESPACE mark with \t or \s
	WHITESPACE

	// SYMBOL (list 1 2 3)
	SYMBOL

	// INT atom
	INT

	// FLOAT atom
	FLOAT

	// STRING atom
	STRING

	// LPAREN (
	LPAREN

	// RPAREN )
	RPAREN
)

func (t token) String() string {
	return tokenNames[t]
}

// tokenNames holds all token with their string names.
var tokenNames = [...]string{
	EOF:        "end of file",
	INVALID:    "invalid token",
	NEWLINE:    "newline",
	WHITESPACE: "whitespace",
	SYMBOL:     "symbol",
	INT:        "int literal",
	FLOAT:      "float literal",
	STRING:     "string literal",
	LPAREN:     "(",
	RPAREN:     ")",
}

// A position what we are reading.
type position struct {
	Line int32
	Col  int32
}

// A scanner represent the input to be parser.
type scanner struct {
	input []byte   // entire input
	rest  []byte   // rest of input
	token []byte   // token being scanned
	pos   position // current input position
	depth int      // nesting of ( )
}

// newScanner creates a new instace of scanner
// with its input. It converts the src instance into
// a byte slice.
func newScanner(src interface{}) *scanner {
	input := []byte(src.(string))

	return &scanner{
		input: input,
		rest:  input,
		pos:   position{Line: 1, Col: 1},
	}
}

// a value represent the value of a token and its position.
type value struct {
	raw    string   // raw text of token
	int    int64    // decoded int
	float  float64  // decoded float
	string string   // decoded string
	pos    position // start position of token
}

// nextToken scan next token and determines the token value.
func (sc *scanner) nextToken() (*value, token) {
	var (
		c   rune
		val = new(value)
	)

	c = sc.peek()

	// brackets
	switch c {
	case '(':
		sc.depth++
		sc.next()

		return val, LPAREN
	case ')':
		if sc.depth == 0 {
			panic("Parsing error: parenthese missing")
		}

		sc.depth--
		sc.next()

		return val, RPAREN
	}

	// Newline. We only need to worry about \n
	// since peek() is coverting \r to \n.
	if c == '\n' {
		sc.next()
		sc.startToken(val)
		sc.endToken(val)
		return val, NEWLINE
	}

	// Spaces and tabs
	if c == ' ' || c == '\t' {
		sc.next()
		sc.startToken(val)
		sc.endToken(val)
		return val, WHITESPACE
	}

	// end of file
	if c == 0 {
		sc.startToken(val)
		sc.endToken(val)
		return val, EOF
	}

	sc.startToken(val)

	// strings atom
	if c == '"' {
		return sc.scanString(val, c)
	}

	// integers and floats atom
	if isdigit(c) || c == '.' {
		return sc.scanNumber(val, c)
	}

	// symbols
	if isSymbolStart(c) {
		for isSymbol(c) {
			sc.next()
			c = sc.peek()
		}
		sc.endToken(val)
		return val, SYMBOL
	}

	sc.next()
	return val, INVALID
}

// Peek return next rune without consuming it.
func (sc *scanner) peek() rune {
	if len(sc.rest) == 0 {
		return 0
	}

	r, _ := utf8.DecodeRune(sc.rest)

	// convert windows newline to \n
	if r == '\r' {
		r = '\n'
	}
	return r
}

// next consumes the next rune and update the current
// position.
func (sc *scanner) next() rune {
	if len(sc.rest) == 0 {
		panic("next at EOF")
	}

	r, size := utf8.DecodeRune(sc.rest)
	sc.rest = sc.rest[size:]

	if r == '\r' {
		r = '\n'
		sc.pos.Line++
		sc.pos.Col = 1
	} else {
		sc.pos.Col++
	}

	return r
}

// startToken collecting processing the token value.
func (sc *scanner) startToken(val *value) {
	sc.token = sc.rest
	val.raw = ""
	val.pos = sc.pos
}

// endToken ends collection of token value.
func (sc *scanner) endToken(val *value) {
	if val.raw == "" {
		val.raw = string(sc.token[:len(sc.token)-len(sc.rest)])
	}
}

// scanString collects all runes for a string.
func (sc *scanner) scanString(val *value, quote rune) (*value, token) {
	sc.next() // handle first quote

	for {
		c := sc.next()

		if c == quote {
			break
		}

		// make sure we collect string without escape
		if c == '\\' {
			sc.next()
		}
	}

	sc.endToken(val)
	r, _ := strconv.Unquote(string(val.raw))
	val.string = string(r)

	return val, STRING
}

// scanNumber collects value for a string token.
func (sc *scanner) scanNumber(val *value, c rune) (*value, token) {
	fraction, exponent := false, false

	// check if number starts with a dot for decimal.
	if c == '.' {
		sc.next()
		c = sc.peek()
		fraction = true
	} else {
		for isdigit(c) {
			sc.next()
			c = sc.peek()
		}

		// check if we have a decimal point with a whole number
		// or we have an exponent.
		if c == '.' {
			fraction = true
		} else if c == 'e' || c == 'E' {
			exponent = true
		}
	}

	if fraction {
		sc.next() // consume dot
		c = sc.next()
		for isdigit(c) {
			// Make sure we don't cosume a closing parenthese, spaces or tabs.
			if next := sc.peek(); next == ')' || next == ' ' || next == '\t' {
				break
			}

			sc.next()
			c = sc.peek()
		}

		if c == 'e' || c == 'E' {
			exponent = true
		}
	}

	if exponent {
		sc.next() // consume e
		c = sc.peek()
		if c == '+' || c == '-' {
			sc.next()
			c = sc.peek()
			// make sure we don't have any invalid runes.
			if !isdigit(c) {
				panic("invalid float literal")
			}
		}
		for isdigit(c) {
			sc.next()
			c = sc.peek()
		}
	}

	sc.endToken(val)

	// convert value into a float point
	if fraction || exponent {
		var err error
		val.float, err = strconv.ParseFloat(strings.TrimSpace(val.raw), 64)
		if err != nil {
			panic(fmt.Sprintf("invalid float literal: %s", err))
		}
		return val, FLOAT
	}

	// covert value into a integer
	var err error
	s := val.raw
	val.int, err = strconv.ParseInt(s, 0, 64)

	if err != nil {
		panic(err)
	}

	return val, INT
}

// isSymbolStart return true if rune is in list of
// valid runes a symbol can start with.
func isSymbolStart(c rune) bool {
	return 'a' <= c && c <= 'z' ||
		'A' <= c && c <= 'Z' ||
		c == '+' ||
		c == '-' ||
		c == '*' ||
		c == '/' ||
		c == '@' ||
		c == '$' ||
		c == '%' ||
		c == '^' ||
		c == '&' ||
		c == '_' ||
		c == '=' ||
		c == '<' ||
		c == '>' ||
		c == '~' ||
		c == '.' ||
		unicode.IsLetter(c)
}

// isSymbol check if rune can be use to make a valid
// symbol name.
func isSymbol(c rune) bool {
	return isdigit(c) || isSymbolStart(c)
}

// isdigit checks if rune is a valid digit.
func isdigit(c rune) bool {
	return '0' <= c && c <= '9'
}
