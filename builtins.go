package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/miguel250/lisp-interpreter/scope"
	"github.com/miguel250/lisp-interpreter/syntax"
)

// builtins holds all go fuction to make it available at run time.
type builtins struct {
	fn map[syntax.SymbolExpr]syntax.Sexpr
}

// add new function to builtins internal map.
func (b *builtins) add(name string, fn scope.Function) {
	s := syntax.SymbolExpr{Token: syntax.SYMBOL, Name: name}
	f := scope.FuncExpr{Name: name, Fn: fn}
	b.fn[s] = &f
}

// newBuiltins returns an instance of builtins with all built-in functions
// added.
func newBuiltins() *builtins {
	b := &builtins{fn: make(map[syntax.SymbolExpr]syntax.Sexpr)}

	b.add("setq", builtinSetq)
	b.add("print", builtinPrint)
	b.add("list", builtinList)
	b.add("first", builtinFirst)
	b.add("+", builtinAdd)
	return b
}

// builtinSetq adds a s-expression into scope. It will failed
// if not enough arguments are pass to it.
func builtinSetq(s *scope.Scope, ss []syntax.Sexpr) (syntax.Sexpr, error) {
	if len(ss) < 2 {
		return nil, fmt.Errorf("setq needs two arguments")
	}

	symbol := ss[0].(*syntax.SymbolExpr)
	expr, err := eval(ss[1], s)
	if err != nil {
		return nil, err
	}

	s.Set(*symbol, expr)
	return expr, nil
}

// builtinPrint prints a s-expression into the stdout.
func builtinPrint(s *scope.Scope, ss []syntax.Sexpr) (syntax.Sexpr, error) {
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
	case *syntax.ConsExpr:
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
func formatCons(src io.Writer, cons *syntax.ConsExpr, parenthese bool) {
	if parenthese {
		fmt.Fprint(src, "( ")
	}

	// car can be a consExpr or atomExpr
	car, ok := cons.Car.(*syntax.ConsExpr)

	if ok {
		formatCons(src, car, true)
	} else {
		fmt.Fprintf(src, "%s ", cons.Car)
	}

	cdr, ok := cons.Cdr.(*syntax.ConsExpr)
	if ok {
		formatCons(src, cdr, false)
	}

	if parenthese {
		fmt.Fprint(src, ")")
	}

}

// builtinList creates a list by linking a set of const together.
// (cons 4 (cons 5 (cons 6 nil)))
func builtinList(s *scope.Scope, ss []syntax.Sexpr) (syntax.Sexpr, error) {
	if len(ss) < 1 {
		return nil, fmt.Errorf("list needs an argument")
	}

	ss, err := evalArgs(s, ss)

	if err != nil {
		return nil, err
	}

	var expr syntax.Sexpr
	expr = &syntax.NilExpr{}

	for i := len(ss) - 1; i >= 0; i-- {
		expr = &syntax.ConsExpr{Car: ss[i], Cdr: expr}
	}

	return expr, nil
}

// builtinFirst returns the first value of a list (const.car).
func builtinFirst(s *scope.Scope, ss []syntax.Sexpr) (syntax.Sexpr, error) {
	if len(ss) < 1 {
		return nil, fmt.Errorf("first needs an argument")
	}

	args, err := evalArgs(s, ss)

	if err != nil {
		return nil, err
	}

	cons, ok := args[0].(*syntax.ConsExpr)

	if !ok {
		return nil, fmt.Errorf("Unable to convert expression to cons: {%v}", args[0])
	}
	return cons.Car, nil
}

// builtinAdd adds two number of the same type together
func builtinAdd(s *scope.Scope, ss []syntax.Sexpr) (syntax.Sexpr, error) {
	if len(ss) < 2 && len(ss) > 2 {
		return nil, fmt.Errorf("Addition only takes 2 args")
	}

	args, err := evalArgs(s, ss)

	if err != nil {
		return nil, err
	}

	firstArg, ok := args[0].(*syntax.AtomExpr)

	if !ok {
		return nil, fmt.Errorf("Failed to add %s", ss)
	}

	secondArg, ok := args[1].(*syntax.AtomExpr)

	if !ok {
		return nil, fmt.Errorf("Failed to add %s", ss)
	}

	if firstArg.Token != secondArg.Token {
		return nil, fmt.Errorf("Arguments have to be the same type got: %s %s", firstArg.Token, secondArg.Token)
	}

	var total int64
	var totalFloat float64

	switch firstArg.Token {
	case syntax.INT:
		total = firstArg.Value.(int64) + secondArg.Value.(int64)
	case syntax.FLOAT:
		totalFloat = firstArg.Value.(float64) + secondArg.Value.(float64)
	default:
		return nil, fmt.Errorf("Unsupported type")
	}

	totalAtom := &syntax.AtomExpr{}
	if firstArg.Token == syntax.INT {
		totalAtom.Value = total
		totalAtom.Token = syntax.INT
	} else {
		totalAtom.Value = totalFloat
		totalAtom.Token = syntax.FLOAT
	}
	return totalAtom, nil
}

// evalArgs evaluate all arguments pass to a function when necessary and
// returns a slice with the s-expressions.
func evalArgs(s *scope.Scope, ss []syntax.Sexpr) ([]syntax.Sexpr, error) {
	args := make([]syntax.Sexpr, 0, 0)
	for _, e := range ss {
		e, err := eval(e, s)
		if err != nil {
			return nil, err
		}
		args = append(args, e)
	}
	return args, nil
}
