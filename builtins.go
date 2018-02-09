package main

import (
	"bytes"
	"fmt"
	"io"
)

// builtins holds all go fuction to make it available at run time.
type builtins struct {
	fn map[symbolExpr]sexpr
}

// add new function to builtins internal map.
func (b *builtins) add(name string, fn function) {
	s := symbolExpr{SYMBOL, name}
	f := funcExpr{name, fn}
	b.fn[s] = &f
}

// newBuiltins returns an instance of builtins with all built-in functions
// added.
func newBuiltins() *builtins {
	b := &builtins{fn: make(map[symbolExpr]sexpr)}

	b.add("setq", builtinSetq)
	b.add("print", builtinPrint)
	b.add("list", builtinList)
	b.add("first", builtinFirst)

	return b
}

// builtinSetq adds a s-expression into scope. It will failed
// if not enough arguments are pass to it.
func builtinSetq(s *scope, ss []sexpr) (sexpr, error) {
	if len(ss) < 2 {
		return nil, fmt.Errorf("setq needs two arguments")
	}

	symbol := ss[0].(*symbolExpr)
	expr, err := eval(ss[1], s)
	if err != nil {
		return nil, err
	}

	s.set(*symbol, expr)
	return expr, nil
}

// builtinPrint prints a s-expression into the stdout.
func builtinPrint(s *scope, ss []sexpr) (sexpr, error) {
	if len(ss) < 1 {
		return nil, fmt.Errorf("print needs an argument")
	}

	args, err := evalArgs(s, ss)

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	expr := args[0]

	switch e := expr.(type) {
	case *consExpr:
		formatCons(&buf, e, true)
	default:
		fmt.Fprintf(&buf, "%s", expr)
	}

	fmt.Println(buf.String())
	return nil, nil
}

// formatCons formats consExpr in a more readable way only adding parentheses
// when necessary.
// ( 4 1 ( 4 1 2 ))
func formatCons(src io.Writer, cons *consExpr, parenthese bool) {
	if parenthese {
		fmt.Fprint(src, "( ")
	}

	// car can be a consExpr or atomExpr
	car, ok := cons.car.(*consExpr)

	if ok {
		formatCons(src, car, true)
	} else {
		fmt.Fprintf(src, "%s ", cons.car)
	}

	cdr, ok := cons.cdr.(*consExpr)
	if ok {
		formatCons(src, cdr, false)
	}

	if parenthese {
		fmt.Fprint(src, ")")
	}

}

// builtinList creates a list by linking a set of const together.
// (cons 4 (cons 5 (cons 6 nil)))
func builtinList(s *scope, ss []sexpr) (sexpr, error) {
	if len(ss) < 1 {
		return nil, fmt.Errorf("list needs an argument")
	}

	ss, err := evalArgs(s, ss)

	if err != nil {
		return nil, err
	}

	var expr sexpr
	expr = &nilExpr{}

	for i := len(ss) - 1; i >= 0; i-- {
		expr = &consExpr{ss[i], expr}
	}

	return expr, nil
}

// builtinFirst returns the first value of a list (const.car).
func builtinFirst(s *scope, ss []sexpr) (sexpr, error) {
	if len(ss) < 1 {
		return nil, fmt.Errorf("first needs an argument")
	}

	args, err := evalArgs(s, ss)

	if err != nil {
		return nil, err
	}

	cons, ok := args[0].(*consExpr)

	if !ok {
		return nil, fmt.Errorf("Unable to convert expression to cons: {%v}", args[0])
	}
	return cons.car, nil
}

// evalArgs evaluate all arguments pass to a function when necessary and
// returns a slice with the s-expressions.
func evalArgs(s *scope, ss []sexpr) ([]sexpr, error) {
	args := make([]sexpr, 0, 0)
	for _, e := range ss {
		e, err := eval(e, s)
		if err != nil {
			return nil, err
		}
		args = append(args, e)
	}
	return args, nil
}
