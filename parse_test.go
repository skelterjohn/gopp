package gopp

import (
	"testing"
	"strings"
	"reflect"
)

var res []TypedRegexp
func init() {
	var err error
	res, err = ByHandGrammar.TokenREs()
	if err != nil {
		panic(err)
	}
}

func TestParseSymbol(t *testing.T) {
	tokens := Tokenize(res, strings.NewReader("stuff <="))
	term := InlineRuleTerm{Name: "identifier"}
	items, remaining, err := term.Parse(ByHandGrammar, tokens)
	if err != nil {
		t.Error(err)
	}
	if st, ok := items[0].(SymbolText); ok {
		if st.Text != "stuff" {
			t.Errorf("Expected %q, got %q.", "stuff", st.Text)
		}
	} else {
		t.Errorf("Expected SymbolText, got %T.", items[0])
	}
	if !reflect.DeepEqual(remaining, tokens[1:]) {
		t.Errorf("Got wrong tokens remaining.")
	}
}

func TestParseTag(t *testing.T) {
	tokens := Tokenize(res, strings.NewReader("=> stuff"))
	term := TagTerm{Tag: "hello"}
	items, remaining, err := term.Parse(ByHandGrammar, tokens)
	if err != nil {
		t.Error(err)
	}
	if tag, ok := items[0].(Tag); ok {
		if tag != "hello" {
			t.Errorf("Expected %q, got %q.", "hello", tag)
		}
	} else {
		t.Errorf("Expected Tag, got %T.", items[0])
	}
	if !reflect.DeepEqual(remaining, tokens) {
		t.Errorf("Got wrong tokens remaining.")
	}
}

func TestParseLiteral(t *testing.T) {
	tokens := Tokenize(res, strings.NewReader("=> stuff"))
	term := LiteralTerm{Literal: "=>"}
	items, remaining, err := term.Parse(ByHandGrammar, tokens)
	if err != nil {
		t.Error(err)
	}
	if lit, ok := items[0].(Literal); ok {
		if lit != "=>" {
			t.Errorf("Expected %q, got %q.", "=>", lit)
		}
	} else {
		t.Errorf("Expected Literal, got %T.", items[0])
	}
	if !reflect.DeepEqual(remaining, tokens[1:]) {
		t.Errorf("Got wrong tokens remaining.")
	}
}
