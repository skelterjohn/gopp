package gopp

import (
	"fmt"
	"strings"
	"testing"
)

// tests where we create a grammar and parse a document

const mathgopp = `
Eqn => {field=Left} <<Expr>> '=' {field=Right} <<Expr>> '\n'
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
	Expr
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
	
	_ = mathgopp
	df, err := NewDecoderFactory(mathgopp, "Eqn")
	if err != nil {
		t.Error(err)
		return
	}
	df.RegisterType(MathExprFactor{})
	df.RegisterType(MathNumberFactor{})
	df.RegisterType(MathSum{})
	df.RegisterType(MathProduct{})
	dec := df.NewDecoder(strings.NewReader("5+1=6"))
	var eqn MathEqn
	err = dec.Decode(&eqn)
	if err != nil {
		t.Error(err)
		return
	}
	if eqn.String() != "5+1=6" {
		t.Errorf("Expected %q, got %q.", "5+1=6", eqn)
	}
}
