// Copyright 2013 The gopp AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopp_test

import (
	"strings"
	"testing"

	"github.com/skelterjohn/gopp"
)

type Node struct {
	Val  string
	Kids []*Node
}

func TestDecodePtrSlice(t *testing.T) {
	grammar := `
ignore: /^\s+/

Start => {field=Kids} <<Node>>*
Node => {field=Val} <dig>

dig = /(\d+)/
`
	data := []string{"1", "4", "9", "42"}
	in := strings.Join(data, " ")

	df, err := gopp.NewDecoderFactory(grammar, "Start")
	if err != nil {
		t.Error(err)
	}
	dec := df.NewDecoder(strings.NewReader(in))
	out := &Node{}
	err = dec.Decode(out)
	if err != nil {
		t.Error(err)
	}
	if len(data) != len(out.Kids) {
		t.Fatalf("Expected %d nodes, got %d", len(data), len(out.Kids))
	}
	for i, s := range data {
		if out.Kids[i].Val != s {
			t.Errorf("Expected node %s, got %s", s, out.Kids[i])
		}
	}
}
