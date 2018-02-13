package main

import (
	"testing"
)

func TestParser(t *testing.T) {
	for _, test := range []struct {
		input, want string
	}{
		{"x\n", "x"},
		{`"x"`, "\"x\""},
		{`()`, "nil"},
		{`(9)`, "(cons 9 nil)"},
		{`(list 1 9 1)`, "(cons list (cons 1 (cons 9 (cons 1 nil))))"},
		{`(first (list 1 7))`, "(cons first (cons (cons list (cons 1 (cons 7 nil))) nil))"},
		{`(1 (2 3) ())`, "(cons 1 (cons (cons 2 (cons 3 nil)) (cons nil nil)))"},
		{`(setq c (list 1.4 "1" 3))`, "(cons setq (cons c (cons (cons list (cons 1.4 (cons \"1\" (cons 3 nil)))) nil)))"},
		{`(())`, "(cons nil nil)"},
		{`(+ 1.4 5.0)`, "(cons + (cons 1.4 (cons 5 nil)))"},
	} {
		expr, err := parse(test.input)

		if err != nil {
			t.Errorf("Parser failed %s", err)
		}

		for _, e := range expr {
			got := e.String()

			if got != test.want {
				t.Errorf("parse `%s` = %s, want %s", test.input, got, test.want)
			}
		}
	}
}
