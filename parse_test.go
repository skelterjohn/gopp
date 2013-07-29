package gopp

import (
	"fmt"
	"strings"
	"testing"
)

// tests where we create a grammar and parse a document

const mathgopp = `
Eqn => {field=Left} <<Expr>> '=' {field=Right} <<Expr>> '\n'
Expr => {field=First} <<Term>> '+' {field=Second} <<Term>>
Expr => {field=First} <<Term>>
Term => {field=First} <<Factor>> '*' {field=Second} <<Factor>>
Term => {field=First} <<Factor>>
Factor => {type=MathExprFactor} '(' {field=Expr} <<Expr>> ')'
Factor => {type=MathNumberFactor} {field=Number} <number>
number = /(\d+)/
`

type MathEqn struct {
	Left, Right MathExpr
}

type MathExpr struct {
	First, Second MathTerm
}

type MathTerm struct {
	First, Second interface{}
}

type MathExprFactor struct {
	Expr
}

type MathNumberFactor struct {
	Number string
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
	dec := df.NewDecoder(strings.NewReader("5+1=6"))
	var eqn MathEqn
	err = dec.Decode(&eqn)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", eqn)
}
