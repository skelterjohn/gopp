// Copyright 2013 The gopp AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopp_test

import (
	"github.com/skelterjohn/gopp"
	"strings"
	"testing"
)

// test literal syntax

const literalgopp = `
Start => <Rule>

Rule => {field=R} 'Success!'
Rule => {field=S} '1'
Rule => {field=T} <symbol>
Rule => {field=X} {Success!} 'X'
Rule => {field=Y} {1} 'Y'
Rule => {field=Z} {1} 'Z'

symbol = /(2)/
`

type LiteralTester struct {
	R string
	S int
	T int
	X string
	Y int
	Z uint
}

var LiteralTestTable = [...]struct {
	expected LiteralTester
	src      string
}{
	{LiteralTester{X: "Success!"}, "X"},
	{LiteralTester{Y: 1}, "Y"},
	{LiteralTester{Z: 1}, "Z"},
	{LiteralTester{R: "Success!"}, "Success!"},
	{LiteralTester{S: 1}, "1"},
	{LiteralTester{T: 2}, "2"},
}

func TestLiteral(t *testing.T) {
	df, err := gopp.NewDecoderFactory(literalgopp, "Start")
	if err != nil {
		t.Error(err)
		return
	}
	for _, test := range LiteralTestTable {
		dec := df.NewDecoder(strings.NewReader(test.src))
		var lit LiteralTester
		err = dec.Decode(&lit)
		if err != nil {
			t.Error(err)
			return
		}
		if lit != test.expected {
			t.Errorf("Expected %+v, got %+v.", test.expected, lit)
		}
	}
}
