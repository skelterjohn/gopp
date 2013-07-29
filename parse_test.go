package gopp

import (
	"testing"
)

// tests where we create a grammar and parse a document

func TestMath(t *testing.T) {
	var mathgopp = `
Eqn => <<Expr>> '=' <<Expr>>
Expr => <<Term>> '+' <<Term>>
Expr => <<Term>>
Term => <<Factor>> '*' <<Factor>>
Term => <<Factor>>
Factor => '(' <<Expr>> ')'
Factor => <number>
number = /(\d*)/
`
	_ = mathgopp
}
