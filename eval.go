package main

import (
	"github.com/miguel250/lisp-interpreter/scope"
	"github.com/miguel250/lisp-interpreter/syntax"
)

// eval evaluate s-expressions by looking in the current scope
// or by running a function.
func eval(e syntax.Sexpr, s *scope.Scope) (syntax.Sexpr, error) {
	switch e := e.(type) {
	case *syntax.ConsExpr:
		car, err := eval(e.Car, s)

		if err != nil {
			return nil, err
		}

		args := make([]syntax.Sexpr, 0, 0)

		cdr, ok := e.Cdr.(*syntax.ConsExpr)

		// make sure we have a list of arguments ready to
		// pass to function.
		for ok {
			args = append(args, cdr.Car)
			cdr, ok = cdr.Cdr.(*syntax.ConsExpr)
		}

		f := car.(*scope.FuncExpr)
		// call function with arguments
		return f.Fn(s, args)
	case *syntax.SymbolExpr:
		return s.Get(*e)
	}
	return e, nil
}
