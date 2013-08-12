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
Rule => {field=X} {Success!} 'X'
`

type LiteralTester struct {
	X string
}

func TestLiteral(t *testing.T) {
	df, err := gopp.NewDecoderFactory(literalgopp, "Rule")
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
	expected := LiteralTester{"Success!"}
	if lit != expected {
		t.Errorf("Expected %#v, got %#v.", expected, lit)
	}
}
