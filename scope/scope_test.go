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
			syntax.SymbolExpr{Token: syntax.SYMBOL, Name: "x"},
			syntax.AtomExpr{Token: syntax.INT, Raw: "2", Value: 2},
			nil,
			"data: map[{4 x}:atom(2)]",
		},
		{
			syntax.SymbolExpr{Token: syntax.SYMBOL, Name: "z"},
			syntax.AtomExpr{Token: syntax.STRING, Raw: "hello", Value: "hello"},
			&testScope{
				syntax.SymbolExpr{Token: syntax.SYMBOL, Name: "x"},
				syntax.AtomExpr{Token: syntax.INT, Raw: "2", Value: 2},
				nil,
				"",
			},
			"parent(data: map[{4 x}:atom(2)]) data: map[{4 z}:atom(\"hello\")]",
		},
	} {
		s := NewScope(nil)

		if test.parent != nil {
			parentScope := NewScope(nil)
			parentScope.Set(test.parent.symbol, &test.parent.atom)
			s = NewScope(parentScope)
		}

		s.Set(test.symbol, &test.atom)

		expr, err := s.Get(test.symbol)

		if err != nil {
			t.Error(err)
		}

		exprStr := expr.String()
		atomStr := test.atom.String()

		if atomStr != exprStr {
			t.Errorf("go: %s want %s", exprStr, atomStr)
		}

		if test.parent != nil {
			expr, err := s.Get(test.parent.symbol)

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
