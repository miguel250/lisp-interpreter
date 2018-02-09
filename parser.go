package main

import (
	"fmt"
)

// parser hold a instance of the scanner as well as
// the current token name and value.
type parser struct {
	sc         *scanner
	tokenName  token
	tokenValue *value
}

// parse takes a input and passes to the scanner then
// it parses all the tokens return by the scanner.
func parse(src interface{}) (ss []sexpr, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	p := parser{sc: newScanner(src)}
	p.nextToken()

	for p.tokenName != EOF {
		ss = append(ss, p.parseNext())
	}

	return ss, nil
}

// nextToken gets the next token from the scanner and
// update the token name and value in parser.
func (p *parser) nextToken() {
	p.tokenValue, p.tokenName = p.sc.nextToken()
}

// parseNext parses the next token return by the scanner
// and return a s-expression for a given token.
func (p *parser) parseNext() (expr sexpr) {

	switch p.tokenName {
	case SYMBOL:
		expr = p.parseSymbol()
	case LPAREN:
		p.nextToken()
		expr = p.parseCons()
	case FLOAT:
		expr = p.parseAtom()
	case INT:
		expr = p.parseAtom()
	case STRING:
		expr = p.parseAtom()
	}

	// make sure we are closing the last parenthese
	if p.tokenName == EOF && p.sc.depth >= 1 {
		panic("Parsing error: parenthese missing")
	}

	p.consume(NEWLINE)

	return
}

// consume the next token if its name
// match the current token.
func (p *parser) consume(tok token) {
	if p.tokenName == tok {
		p.nextToken()
	}
}

// parseSymbol parses a symbol by wrapping it in a
// symbolExpr struct then returns a s-expression.
func (p *parser) parseSymbol() sexpr {
	tok := p.tokenName
	name := p.tokenValue.raw
	p.nextToken()

	p.consume(WHITESPACE)

	return &symbolExpr{
		token: tok,
		name:  name,
	}
}

// parseCons parses everything inside of
// the parentheses then returning a consExpr as
// s-expression interface
func (p *parser) parseCons() sexpr {
	tok := p.tokenName
	if tok == RPAREN {
		p.nextToken()
		return &nilExpr{}
	}

	car := p.parseNext()
	cdr := p.parseCons()

	p.consume(WHITESPACE)

	return &consExpr{car: car, cdr: cdr}
}

// parseAtom parses all string, integers and float points
// wrapped in an atomExpr.
func (p *parser) parseAtom() sexpr {
	var value interface{}
	tok := p.tokenName
	raw := p.tokenValue.raw

	switch tok {
	case STRING:
		value = p.tokenValue.string
	case INT:
		value = p.tokenValue.int
	case FLOAT:
		value = p.tokenValue.float
	}
	p.nextToken()

	p.consume(WHITESPACE)

	return &atomExpr{
		token: tok,
		raw:   raw,
		value: value,
	}
}
