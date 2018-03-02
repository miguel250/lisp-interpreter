package syntax

// token holds all tokens to be parse.
type Token int8

const (
	EOF Token = iota
	INVALID

	NEWLINE
	WHITESPACE

	SYMBOL
	INT
	FLOAT
	STRING

	LPAREN // (
	RPAREN // )
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
