// Copyright 2013 The gopp AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopp_test

import (
	"github.com/skelterjohn/gopp"
	"reflect"
	"strings"
	"testing"
)

type Case struct {
	Document         string
	Object, Expected interface{}
}

type Subject struct {
	Name    string
	Gopp    string
	Grammar gopp.Grammar
	Cases   []Case
}

type XYZ struct {
	X, Y, Z string
}

var Subjects = []Subject{
	Subject{
		Name: "OptionalTest",
		Gopp: `
Start => {field=Y} <Y> [{field=Z} <Z>]
Y = /(y)/
Z = /(z)/
`,
		Cases: []Case{
			Case{`yz`, &XYZ{}, &XYZ{Y: "y", Z: "z"}},
			Case{`y`, &XYZ{}, &XYZ{Y: "y"}},
		},
	},
}

func TestSubjects(t *testing.T) {
subject:
	for _, s := range Subjects {
		df, err := gopp.NewDecoderFactory(s.Gopp, "Start")
		if err != nil {
			t.Error(err)
			continue subject
		}
	scase:
		for _, c := range s.Cases {
			dec := df.NewDecoder(strings.NewReader(c.Document))
			err = dec.Decode(c.Object)
			if err != nil {
				t.Error(err)
				continue scase
			}
			if !reflect.DeepEqual(c.Object, c.Expected) {
				t.Errorf("(%s) With %q, got %+v, expected %+v.", s.Name, c.Document, c.Object, c.Expected)
			}
		}
	}
}
