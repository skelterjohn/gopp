package gopp_test

import (
	"fmt"
	"github.com/skelterjohn/gopp"
	"strings"
	"testing"
)

// tests where we create a grammar and parse a document

const mathgopp = `
Eqn => {type=MathEqn} {field=Left} <<Expr>> '=' {field=Right} <<Expr>> '\n'
Expr => {type=MathSum} {field=First} <<Term>> '+' {field=Second} <<Term>>
Expr => <Term>
Term => {type=MathProduct} {field=First} <<Factor>> '*' {field=Second} <<Factor>>
Term => <Factor>
Factor => {type=MathExprFactor} '(' {field=Expr} <<Expr>> ')'
Factor => {type=MathNumberFactor} {field=Number} <number>
number = /(\d+)/
`

type MathEqn struct {
	Left, Right interface{}
}

func (e MathEqn) String() string {
	return fmt.Sprintf("%s=%s", e.Left, e.Right)
}

type MathSum struct {
	First, Second interface{}
}

func (s MathSum) String() string {
	return fmt.Sprintf("%s+%s", s.First, s.Second)
}

type MathProduct struct {
	First, Second interface{}
}

func (p MathProduct) String() string {
	return fmt.Sprintf("%s*%s", p.First, p.Second)
}

type MathExprFactor struct {
	Expr interface{}
}

func (ef MathExprFactor) String() string {
	return fmt.Sprintf("(%s)", ef.Expr)
}

type MathNumberFactor struct {
	Number string
}

func (nf MathNumberFactor) String() string {
	return nf.Number
}

func TestMath(t *testing.T) {
	df, err := gopp.NewDecoderFactory(mathgopp, "Eqn")
	if err != nil {
		t.Error(err)
		return
	}
	df.RegisterType(MathExprFactor{})
	df.RegisterType(MathNumberFactor{})
	df.RegisterType(MathSum{})
	df.RegisterType(MathProduct{})
	dec := df.NewDecoder(strings.NewReader("5+1=6\n"))
	var eqn MathEqn
	err = dec.Decode(&eqn)
	if err != nil {
		t.Error(err)
		return
	}

	expectedEqn := MathEqn{
		Left: MathSum{
			First:  MathNumberFactor{"5"},
			Second: MathNumberFactor{"1"},
		},
		Right: MathNumberFactor{"6"},
	}

	if eqn != expectedEqn {
		t.Errorf("Expected %q, got %q.", expectedEqn, eqn)
	}
}
