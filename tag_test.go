// Copyright 2013 The gopp AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopp_test

import (
	"github.com/skelterjohn/gopp"
	"strings"
	"testing"
)

// test {field=X} {literal value} syntax

const literalgopp = `
Start => <Rule>

Rule => {field=R} 'Success!'
Rule => {field=X} {Success!} 'X'
Rule => {field=Y} {1} 'Y'
Rule => {field=Z} {1} 'Z'
`

type LiteralTester struct {
	R string
	X string
	Y int
	Z uint
}

func TestTagStringLiteral(t *testing.T) {
	df, err := gopp.NewDecoderFactory(literalgopp, "Start")
	if err != nil {
		t.Error(err)
		return
	}
	dec := df.NewDecoder(strings.NewReader("X"))
	var lit LiteralTester
	err = dec.Decode(&lit)
	if err != nil {
		t.Error(err)
		return
	}
	expected := LiteralTester{X: "Success!"}
	if lit != expected {
		t.Errorf("Expected %#v, got %#v.", expected, lit)
	}
}

func TestTagIntLiteral(t *testing.T) {
	df, err := gopp.NewDecoderFactory(literalgopp, "Start")
	if err != nil {
		t.Error(err)
		return
	}
	dec := df.NewDecoder(strings.NewReader("Y"))
	var lit LiteralTester
	err = dec.Decode(&lit)
	if err != nil {
		t.Error(err)
		return
	}
	expected := LiteralTester{Y: 1}
	if lit != expected {
		t.Errorf("Expected %#v, got %#v.", expected, lit)
	}
}

func TestTagUintLiteral(t *testing.T) {
	df, err := gopp.NewDecoderFactory(literalgopp, "Start")
	if err != nil {
		t.Error(err)
		return
	}
	dec := df.NewDecoder(strings.NewReader("Z"))
	var lit LiteralTester
	err = dec.Decode(&lit)
	if err != nil {
		t.Error(err)
		return
	}
	expected := LiteralTester{Z: 1}
	if lit != expected {
		t.Errorf("Expected %#v, got %#v.", expected, lit)
	}
}

func TestLiteral(t *testing.T) {
	df, err := gopp.NewDecoderFactory(literalgopp, "Start")
	if err != nil {
		t.Error(err)
		return
	}
	dec := df.NewDecoder(strings.NewReader("Success!"))
	var lit LiteralTester
	err = dec.Decode(&lit)
	if err != nil {
		t.Error(err)
		return
	}
	expected := LiteralTester{R: "Success!"}
	if lit != expected {
		t.Errorf("Expected %#v, got %#v.", expected, lit)
	}
}
