package scope

import (
	"testing"

	"github.com/miguel250/lisp-interpreter/syntax"
)

type testScope struct {
	symbol      syntax.SymbolExpr
	atom        syntax.AtomExpr
	parent      *testScope
	scopeString string
}

func TestScope(t *testing.T) {
	for _, test := range []testScope{
		{
			syntax.SymbolExpr{syntax.SYMBOL, "x"},
			syntax.AtomExpr{syntax.INT, "2", 2},
			nil,
			"data: map[{4 x}:atom(2)]",
		},
		{
			syntax.SymbolExpr{syntax.SYMBOL, "z"},
			syntax.AtomExpr{syntax.STRING, "hello", "hello"},
			&testScope{
				syntax.SymbolExpr{syntax.SYMBOL, "x"},
				syntax.AtomExpr{syntax.INT, "2", 2},
				nil,
				"",
			},
			"parent(data: map[{4 x}:atom(2)]) data: map[{4 z}:atom(\"hello\")]",
		},
	} {
		s := newScope(nil)

		if test.parent != nil {
			parentScope := newScope(nil)
			parentScope.set(test.parent.symbol, &test.parent.atom)
			s = newScope(parentScope)
		}

		s.set(test.symbol, &test.atom)

		expr, err := s.get(test.symbol)

		if err != nil {
			t.Error(err)
		}

		exprStr := expr.String()
		atomStr := test.atom.String()

		if atomStr != exprStr {
			t.Errorf("go: %s want %s", exprStr, atomStr)
		}

		if test.parent != nil {
			expr, err := s.get(test.parent.symbol)

			if err != nil {
				t.Fatalf("%s", err)
			}

			exprStr := expr.String()
			atomStr := test.parent.atom.String()

			if atomStr != exprStr {
				t.Errorf("go: %s want %s", exprStr, atomStr)
			}
		}
	}
}