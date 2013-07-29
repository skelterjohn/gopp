package gopp

import (
	"reflect"
	"fmt"
	"errors"
	"strings"
)

func getTagValue(typ string, t Tag) (value string, ok bool) {
	prefix := typ + "="
	if strings.HasPrefix(string(t), prefix) {
		value = string(t[len(prefix):])
		ok = true
	}
	return
}

var _ = fmt.Println

func Decode(ast AST, obj interface{}) (err error) {
	return decode([]Node(ast), reflect.ValueOf(obj))
}

func decode(node Node, v reflect.Value) (err error) {
	// if we've got a []Node with one element that's also a []Node, just dig deeper
	if nodes, isNodeSlice := node.([]Node); isNodeSlice {
		if len(nodes) == 1 {
			if nodesElem, isNodeSliceSlice := nodes[0].([]Node); isNodeSliceSlice {
				fmt.Println("dropping deeper")
				return decode(nodesElem, v)
			}
		}
	}

	fmt.Printf("Decoding into a %T\n", v.Interface())
	fmt.Println(node)

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
				name, isField := getTagValue("field", tag)
				if isField {
					// if we have a field tag, that indicates that the next node should be decided into the field with the given name.
					var fv reflect.Value
					fv, err = getField(v, name)
					if err != nil {
						return
					}
					decode(nodes[i+1], fv)
				}
			}
		}
	}

	// map things into slices
	if typ.Kind() == reflect.Slice {
		//fmt.Printf("Going into %s is\n", typ.Elem().Name())
		//printNode(node, 0)
		nodes, ok := node.([]Node)
		if !ok {
			err = errors.New("Need to populate slice via []Node.")
			return
		}
		for _, n := range nodes {
			// create an addressable value to put in the slice
			ev := reflect.New(typ.Elem()).Elem()
			err = decode(n, ev)
			if err != nil {
				return
			}
			// this is how append looks w/ reflect
			v.Set(reflect.Append(v, ev))
		}
	}

	// symbols go into strings
	if typ.Kind() == reflect.String {
		st, isSymbol := node.(SymbolText)
		if !isSymbol {
			err = errors.New("Trying to store non-symbol into string type.")
			return
		}
		v.SetString(st.Text)
		fmt.Println(v.Interface())
	}
	return
}

func getField(v reflect.Value, field string) (fv reflect.Value, err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("Type %s has no field named %q.", v.Type().Name(), field)
		}
	}()
	fv = v.FieldByName(field)
	return
}
