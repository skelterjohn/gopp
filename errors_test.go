// Copyright 2013 The gopp AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopp_test

import (
	"github.com/skelterjohn/gopp"
	"strings"
	"testing"
)

type ErrorCase struct {
	Document      string
	ExpectedError string
}

type ErrorSubject struct {
	Name    string
	Gopp    string
	Grammar gopp.Grammar
	Cases   []ErrorCase
}

var ErrorSubjects = []ErrorSubject{
	ErrorSubject{
		Name: "LiteralConjunction",
		Gopp: `
Start => 'x' 'y' 'z'
`,
		Cases: []ErrorCase{
			ErrorCase{`xyz`, ``},
			ErrorCase{`xzy`, `Start -> 'x' 'y' 'z' -> 'y': Expected "y" at 0:1.`},
			ErrorCase{`x`, `Start -> 'x' 'y' 'z' -> 'y': Expected "y" at EOF.`},
		},
	},
}

func TestErrors(t *testing.T) {
subject:
	for _, s := range ErrorSubjects {
		df, err := gopp.NewDecoderFactory(s.Gopp, "Start")
		if err != nil {
			t.Error(err)
			continue subject
		}
	scase:
		for _, c := range s.Cases {
			dec := df.NewDecoder(strings.NewReader(c.Document))
			err = dec.Decode(&XYZ{})
			if !(err == nil && c.ExpectedError == "") && (err != nil && err.Error() != c.ExpectedError) {
				t.Error(err)
				continue scase
			}
		}
	}
}
