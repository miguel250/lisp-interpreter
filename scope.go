package main

import (
	"bytes"
	"fmt"
)

// A scope holds the current set of symbols available
// for a s-expression.
type scope struct {
	data   map[symbolExpr]sexpr
	parent *scope
}

// newScope returns a new instance of a scope.
// It can take a parent scope for nesting scope.
func newScope(parent *scope) *scope {
	return &scope{
		data:   make(map[symbolExpr]sexpr),
		parent: parent,
	}
}

func (s *scope) String() string {
	var buf bytes.Buffer
	if s.parent != nil {
		fmt.Fprintf(&buf, "parent(%s) ", s.parent)
	}
	fmt.Fprintf(&buf, "data: %v", s.data)
	return buf.String()
}

// set adds a symbol and s-expression into the scope.
func (s *scope) set(symbol symbolExpr, expr sexpr) {
	s.data[symbol] = expr
}

// get returns a s-expression from a symbol.
func (s *scope) get(symbol symbolExpr) (sexpr, error) {
	v, ok := s.data[symbol]

	if ok {
		return v, nil
	}

	if s.parent != nil {
		v, err := s.parent.get(symbol)

		if err == nil {
			return v, nil
		}
	}

	return nil, fmt.Errorf("Symbol not found in scope: {%s}", symbol.name)
}
