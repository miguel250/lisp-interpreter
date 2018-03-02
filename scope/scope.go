package scope

import (
	"bytes"
	"fmt"

	"github.com/miguel250/lisp-interpreter/syntax"
)

// A Scope holds the current set of symbols available
// for a s-expression.
type Scope struct {
	data   map[syntax.SymbolExpr]syntax.Sexpr
	parent *Scope
}

// NewScope returns a new instance of a scope.
// It can take a parent scope for nesting scope.
func NewScope(parent *Scope) *Scope {
	return &Scope{
		data:   make(map[syntax.SymbolExpr]syntax.Sexpr),
		parent: parent,
	}
}

func (s *Scope) String() string {
	var buf bytes.Buffer
	if s.parent != nil {
		fmt.Fprintf(&buf, "parent(%s) ", s.parent)
	}
	fmt.Fprintf(&buf, "data: %v", s.data)
	return buf.String()
}

// Set adds a symbol and s-expression into the scope.
func (s *Scope) Set(symbol syntax.SymbolExpr, expr syntax.Sexpr) {
	s.data[symbol] = expr
}

// Get returns a s-expression from a symbol.
func (s *Scope) Get(symbol syntax.SymbolExpr) (syntax.Sexpr, error) {
	v, ok := s.data[symbol]

	if ok {
		return v, nil
	}

	if s.parent != nil {
		v, err := s.parent.Get(symbol)

		if err == nil {
			return v, nil
		}
	}

	return nil, fmt.Errorf("Symbol not found in scope: {%s}", symbol.Name)
}

// Function is a function to be added to the scope and make it accessible to be called.
type Function func(*Scope, []syntax.Sexpr) (syntax.Sexpr, error)

// A FuncExpr represent a callable function s-expression.
type FuncExpr struct {
	Name string
	Fn   Function
}

func (*FuncExpr) expr() {}
func (f FuncExpr) String() string {
	return fmt.Sprintf("fn: %s", f.Name)
}
