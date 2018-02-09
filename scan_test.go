package main

import (
	"bytes"
	"fmt"
	"testing"
)

func scan(src interface{}) (string, error) {
	sc := newScanner(src)

	var buf bytes.Buffer

	for {
		val, token := sc.nextToken()

		if buf.Len() > 0 {
			buf.WriteString(" ")
		}

		switch token {
		case SYMBOL:
			buf.WriteString(val.raw)
		case STRING:
			fmt.Fprintf(&buf, "%q", val.string)
		case INT:
			fmt.Fprintf(&buf, "%d", val.int)
		case FLOAT:
			fmt.Fprintf(&buf, "%e", val.float)
		case EOF:
			buf.WriteString("EOF")
		default:
			buf.WriteString(token.String())
		}

		if token == EOF {
			break
		}
	}

	return buf.String(), nil
}

func TestScanner(t *testing.T) {
	for _, test := range []struct {
		input, want string
	}{
		{``, "EOF"},
		{`(1e-1 1e1)`, "( 1.000000e-01 whitespace 1.000000e+01 ) EOF"},
		{`(y 3.14159265 .1e+1)`, "( y whitespace 3.141593e+00 whitespace 1.000000e+00 ) EOF"},
		{`(x 2 "3")`, "( x whitespace 2 whitespace \"3\" ) EOF"},
		{`b^2-4*a*c`, "b^2-4*a*c EOF"},
		{`+1`, "+1 EOF"},
		{`+$`, "+$ EOF"},
		{`(first (list 1 (+ 2 3) 9))`, "( first whitespace ( list whitespace 1 whitespace ( + whitespace 2 whitespace 3 ) whitespace 9 ) ) EOF"},
		{"(1e-1 1e1)\n(x 2 \"3\")", "( 1.000000e-01 whitespace 1.000000e+01 ) newline ( x whitespace 2 whitespace \"3\" ) EOF"},
	} {

		got, err := scan(test.input)
		if err != nil {
			got = err.(error).Error()
		}

		if test.want != got {
			t.Errorf("scan `%s` = [%s], want [%s]", test.input, got, test.want)
		}
	}
}
