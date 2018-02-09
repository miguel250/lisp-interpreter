package main

import (
	"bytes"
	"fmt"
)

// sexpr is a S-expression
type sexpr interface {
	expr()
	String() string
}

// A consExpr is a set of values inside of a
// parentheses.
type consExpr struct {
	car sexpr
	cdr sexpr
}

func (*consExpr) expr() {}
func (c *consExpr) String() string {
	return fmt.Sprintf("(cons %s %s)", c.car, c.cdr)
}

// a nilExpr represent a "nil".
type nilExpr struct{}

func (*nilExpr) expr()          {}
func (*nilExpr) String() string { return "nil" }

// A symbolExpr represent the name of a symbol to be able
// to access variables.
type symbolExpr struct {
	token token
	name  string
}

func (*symbolExpr) expr() {}
func (s *symbolExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.name)
	return buf.String()
}

// An atomExpr represent all variables types.
type atomExpr struct {
	token token
	raw   string
	value interface{}
}

func (*atomExpr) expr() {}
func (a *atomExpr) String() string {
	var buf bytes.Buffer

	switch a.token {
	case STRING:
		fmt.Fprintf(&buf, "%q", a.value)
	case INT:
		fmt.Fprintf(&buf, "%d", a.value)
	case FLOAT:
		fmt.Fprintf(&buf, "%g", a.value)
	}
	return buf.String()
}

// function is a function to be added to the scope and make it accessible to be called.
type function func(*scope, []sexpr) (sexpr, error)

// A funcExpr represent a callable function s-expression.
type funcExpr struct {
	name string
	fn   function
}

func (*funcExpr) expr() {}
func (f funcExpr) String() string {
	return fmt.Sprintf("fn: %s", f.name)
}
