package syntax

import (
	"bytes"
	"fmt"
)

// Sexpr is a S-expression
type Sexpr interface {
	Expr()
	String() string
}

// A ConsExpr is a set of values inside of a
// parentheses.
type ConsExpr struct {
	Car Sexpr
	Cdr Sexpr
}

// Expr is use to satified Sexpr interface
func (*ConsExpr) Expr() {}
func (c *ConsExpr) String() string {
	return fmt.Sprintf("(cons %s %s)", c.Car, c.Cdr)
}

// A NilExpr represent a "nil".
type NilExpr struct{}

// Expr is use to satified Sexpr interface
func (*NilExpr) Expr()          {}
func (*NilExpr) String() string { return "nil" }

// A SymbolExpr represent the name of a symbol to be able
// to access variables.
type SymbolExpr struct {
	Token Token
	Name  string
}

// Expr is use to satified Sexpr interface
func (*SymbolExpr) Expr() {}
func (s *SymbolExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.Name)
	return buf.String()
}

// An AtomExpr represent all variables types.
type AtomExpr struct {
	Token Token
	Raw   string
	Value interface{}
}

// Expr is use to satified Sexpr interface
func (*AtomExpr) Expr() {}
func (a *AtomExpr) String() string {
	var buf bytes.Buffer

	switch a.Token {
	case STRING:
		fmt.Fprintf(&buf, "%q", a.Value)
	case INT:
		fmt.Fprintf(&buf, "%d", a.Value)
	case FLOAT:
		fmt.Fprintf(&buf, "%g", a.Value)
	}
	return buf.String()
}
