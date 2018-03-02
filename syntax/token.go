package syntax

// Token holds all tokens to be parse.
type Token int8

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

func (t Token) String() string {
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
