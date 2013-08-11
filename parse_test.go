// Copyright 2013 The gopp AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopp_test

import (
	"fmt"
	"github.com/skelterjohn/gopp"
	"strings"
	"testing"
)

// tests where we create a grammar and parse a document

const mathgopp = `
# The root is an equation, with a left-hand and right-hand side.
Eqn => {type=MathEqn} {field=Left} <<Expr>> '=' {field=Right} <<Expr>> '\n'

# An Expr is either the sum of two terms,
Expr => {type=MathSum} {field=First} <<Term>> '+' {field=Second} <<Term>>
# or just another term.
Expr => <Term>

# A Term is either the product of two factors,
Term => {type=MathProduct} {field=First} <<Factor>> '*' {field=Second} <<Factor>>
# or just another factor.
Term => <Factor>

# A factor is either a parenthesized expression,
Factor => {type=MathExprFactor} '(' {field=Expr} <<Expr>> ')'
# or just a number.
Factor => {type=MathNumberFactor} {field=Number} <number>

# A number is a string of consecutive digits.
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

func TestMathPrecedence(t *testing.T) {
	df, err := gopp.NewDecoderFactory(mathgopp, "Eqn")
	if err != nil {
		t.Error(err)
		return
	}
	df.RegisterType(MathExprFactor{})
	df.RegisterType(MathNumberFactor{})
	df.RegisterType(MathSum{})
	df.RegisterType(MathProduct{})
	dec := df.NewDecoder(strings.NewReader("5+5*2=6*2+3\n"))
	var eqn MathEqn
	err = dec.Decode(&eqn)
	if err != nil {
		t.Error(err)
		return
	}

	expectedEqn := MathEqn{
		Left: MathSum{
			MathNumberFactor{"5"},
			MathProduct{
				MathNumberFactor{"5"},
				MathNumberFactor{"2"},
			},
		},
		Right: MathSum{
			MathProduct{
				MathNumberFactor{"6"},
				MathNumberFactor{"2"},
			},
			MathNumberFactor{"3"},
		},
	}

	if eqn != expectedEqn {
		t.Errorf("Expected %q, got %q.", expectedEqn, eqn)
	}
}
