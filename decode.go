// Copyright 2013 The gopp AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopp

import (
	"errors"
	"fmt"
	"github.com/skelterjohn/debugtags"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

type DecoderFactory struct {
	g     Grammar
	start string
	types map[string]reflect.Type
}

func NewDecoderFactory(gopp string, start string) (df *DecoderFactory, err error) {
	df = &DecoderFactory{
		start: start,
		types: map[string]reflect.Type{},
	}
	ast, err := Parse(ByHandGrammar, "Grammar", []byte(gopp))
	if err != nil {
		return
	}
	sa := NewStructuredAST(ast)
	sa.RegisterType(RepeatZeroTerm{})
	sa.RegisterType(RepeatOneTerm{})
	sa.RegisterType(OptionalTerm{})
	sa.RegisterType(GroupTerm{})
	sa.RegisterType(RuleTerm{})
	sa.RegisterType(InlineRuleTerm{})
	sa.RegisterType(TagTerm{})
	sa.RegisterType(LiteralTerm{})
	err = sa.Decode(&df.g)
	if err != nil {
		return
	}
	return
}

func (df *DecoderFactory) RegisterType(x interface{}) {
	typ := reflect.TypeOf(x)
	df.types[typ.Name()] = typ
}

func (df *DecoderFactory) NewDecoder(r io.Reader) (d Decoder) {
	d = Decoder{
		DecoderFactory: df,
		Reader:         r,
	}
	return
}

type Decoder struct {
	*DecoderFactory
	io.Reader
}

func (d *Decoder) Decode(obj interface{}) (err error) {
	document, err := ioutil.ReadAll(d.Reader)
	if err != nil {
		return
	}
	ast, err := Parse(d.g, d.start, document)
	if err != nil {
		return
	}
	sa := NewStructuredAST(ast)
	sa.types = d.types
	err = sa.Decode(obj)
	if err != nil {
		return
	}
	return
}

func getTagValue(typ string, t Tag) (value string, ok bool) {
	prefix := typ + "="
	if strings.HasPrefix(string(t), prefix) {
		value = string(t[len(prefix):])
		ok = true
	}
	return
}

var _ = fmt.Println

type StructuredAST struct {
	ast   AST
	types map[string]reflect.Type
}

func NewStructuredAST(ast AST) (sa StructuredAST) {
	sa = StructuredAST{
		ast:   ast,
		types: map[string]reflect.Type{},
	}
	return
}

func (sa StructuredAST) RegisterType(x interface{}) {
	t := reflect.TypeOf(x)
	sa.types[t.Name()] = t
}

func (sa StructuredAST) Decode(obj interface{}) (err error) {
	return sa.decode([]Node(sa.ast), reflect.ValueOf(obj))
}

var dtr = debugtags.Tracer{Enabled: false}

func (sa StructuredAST) decode(node Node, v reflect.Value) (err error) {
	name := fmt.Sprintf("%T", v.Interface())
	dtr.In(name, node)
	defer func() {
		dtr.Out(name, v.Interface())
	}()

	typ := v.Type()

	// deref a pointer
	if typ.Kind() == reflect.Ptr {
		// but first check if it's nil and, if so, allocate
		if v.IsNil() {
			v.Elem().Set(reflect.New(typ.Elem()))
		}
		v = v.Elem()
		typ = typ.Elem()
	}

	// populate struct fields
	if typ.Kind() == reflect.Struct {
		// we've got a struct pointer - iterate through node looking for field= tags
		nodes, ok := node.([]Node)
		if !ok {
			err = errors.New("Need to populate struct via []Node with tags.")
			return
		}
		for i := range nodes {
			if tag, ok := nodes[i].(Tag); ok {
				if typName, isType := getTagValue("type", tag); isType {
					if typName != typ.Name() {
						err = fmt.Errorf("AST wants type %q, being decoded to type %q.", typName, typ.Name())
					}
				}

				if name, isField := getTagValue("field", tag); isField {
					// if we have a field tag, that indicates that the next node should be decided into the field with the given name.
					var fv reflect.Value
					fv, err = getField(v, name)
					if err != nil {
						return
					}

					if fv.Type().Kind() == reflect.Interface {
						//dtr.Println("field of interface")
						var pv reflect.Value
						pv, err = sa.makePointerWithType(nodes[i+1])
						if err != nil {
							return
						}
						err = sa.decode(nodes[i+1], pv.Elem())
						fv.Set(pv.Elem())
						if err != nil {
							return
						}
					} else {
						sa.decode(nodes[i+1], fv)
					}
				}
			}
		}
		return
	}

	// map things into slices
	if typ.Kind() == reflect.Slice {
		//fmt.Printf("Going into %s is\n", typ.Elem().Name())
		//printNode(node, 0)
		isInterfaceSlice := typ.Elem().Kind() == reflect.Interface
		// if isInterfaceSlice {
		// 	dtr.Println("slice of interface")
		// }
		nodes, ok := node.([]Node)
		if !ok {
			err = errors.New("Need to populate slice via []Node.")
			return
		}
		for _, n := range nodes {
			// create an addressable value to put in the slice
			ev := reflect.New(typ.Elem()).Elem()
			if isInterfaceSlice {
				var pv reflect.Value
				pv, err = sa.makePointerWithType(n)
				if err != nil {
					return
				}
				err = sa.decode(n, pv.Elem())
				ev.Set(pv.Elem())
				if err != nil {
					return
				}
			} else {
				err = sa.decode(n, ev)
				if err != nil {
					return
				}
			}
			// this is how append looks w/ reflect
			v.Set(reflect.Append(v, ev))
		}
		return
	}

	// symbols and tags go into strings
	if typ.Kind() == reflect.String {
		switch nn := node.(type) {
		case SymbolText:
			ds, derr := descapeString(nn.Text)
			if derr == nil {
				v.SetString(ds)
			} else {
				v.SetString(nn.Text)
			}
		case Tag:
			v.SetString(string(nn))
		default:
			err = errors.New("Trying to store invalid type into string type.")
			return
		}
		return
	}

	err = fmt.Errorf("Unanticipated type: %s.", typ.Name())
	return
}

func (sa StructuredAST) makePointerWithType(node Node) (pointer reflect.Value, err error) {
	var ntag Tag
	nodes, ok := node.([]Node)
	if ok && len(nodes) != 0 {
		ntag, ok = nodes[0].(Tag)
		if ok {
			ok = strings.HasPrefix(string(ntag), "type=")
		}
	}
	if !ok {
		err = errors.New("Can only infer type from []Node with a type= tag.")
		return
	}
	typeName := ntag[len("type="):]
	typ := sa.types[string(typeName)]
	pointer = reflect.New(typ)
	return
}

func getField(v reflect.Value, field string) (fv reflect.Value, err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("Type %s has no field named %q.", v.Type().Name(), field)
		}
	}()
	if field == "." {
		// . means to store the next level deeper in the same value
		fv = v
	} else {
		fv = v.FieldByName(field)
	}
	return
}
