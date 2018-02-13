package main

import (
	"testing"
)

func TestEval(t *testing.T) {

	for _, test := range []struct {
		input, want string
	}{
		//{`3`, "3"},
		//{`(setq a 5)`, "5"},
		//{`(setq a "6")`, "\"6\""},
		//{`(list 4 5 6)`, "(cons 4 (cons 5 (cons 6 nil)))"},
		//{`(first (list 4 5 6))`, "4"},
		{`(+ 1 2)`, "3"},
		{`(+ 1.4 5.0)`, "6.4"},
	} {

		expr, err := parse(test.input)
		if err != nil {
			t.Fatalf("%s", err)
		}

		s := newScope(nil)

		b := newBuiltins()
		for k, v := range b.fn {
			s.set(k, v)
		}

		e, err := eval(expr[0], s)

		if err != nil {
			t.Fatalf("%s", err)
		}

		if got := e.String(); got != test.want {
			t.Errorf("eval `%s` = %s, want %s", test.input, got, test.want)

		}
	}
}
