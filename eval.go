package main

// eval evaluate s-expressions by looking in the current scope
// or by running a function.
func eval(e sexpr, s *scope) (sexpr, error) {
	switch e := e.(type) {
	case *consExpr:
		car, err := eval(e.car, s)

		if err != nil {
			return nil, err
		}

		args := make([]sexpr, 0, 0)

		cdr, ok := e.cdr.(*consExpr)

		// make sure we have a list of arguments ready to
		// pass to function.
		for ok {
			args = append(args, cdr.car)
			cdr, ok = cdr.cdr.(*consExpr)
		}

		f := car.(*funcExpr)
		// call function with arguments
		return f.fn(s, args)
	case *symbolExpr:
		return s.get(*e)
	}
	return e, nil
}
