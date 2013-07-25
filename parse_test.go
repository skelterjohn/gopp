package gopp

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var ByHandGrammarREs []TypedRegexp

func init() {
	var err error
	ByHandGrammarREs, err = ByHandGrammar.TokenREs()
	if err != nil {
		panic(err)
	}
}

func printNode(node Node, indentCount int) {
	indent := func(tag string) {
		for i := 0; i < indentCount; i++ {
			fmt.Print(" ")
		}
		fmt.Println(tag)
	}
	switch node := node.(type) {
	case []Node:
		indent("[")
		for _, n := range node {
			printNode(n, indentCount+1)
		}
		indent("]")
	case AST:
		indent("[")
		for _, n := range node {
			printNode(n, indentCount+1)
		}
		indent("]")
	default:
		indent(fmt.Sprint(node))
	}
}

func xTestParseFullGrammar(t *testing.T) {
	tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader(goppgopp))
	if err != nil {
		t.Error(err)
		return
	}
	start := ByHandGrammar.RulesForName("Grammar")[0]
	items, remaining, err := start.Parse(ByHandGrammar, tokens)
	if err != nil {
		t.Error(err)
		return
	}
	if len(remaining) != 0 {
		t.Error("Leftover tokens: %v.", remaining)
	}
	fmt.Println(items)
}

func xTestParseEasyGrammar(t *testing.T) {
	byHandAST := mkGrammar(
		[]Node{
			mkRule("X",
				mkLiteralTerm("y"),
			),
		},
		[]Node{
			mkSymbol("w", "z"),
		},
	)

	tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader("X => 'y'\nw = /z/\n"))
	if err != nil {
		t.Error(err)
		return
	}
	start := ByHandGrammar.RulesForName("Grammar")[0]
	items, remaining, err := start.Parse(ByHandGrammar, tokens)
	if err != nil {
		t.Error(err)
		return
	}
	if len(remaining) != 0 {
		t.Errorf("Leftover tokens: %v.", remaining)
	}

	if !reflect.DeepEqual(byHandAST, AST(items)) {
		t.Error("Generated AST doesn't match by-hand AST.")
	}
}

func TestParseMultiRule(t *testing.T) {
	byHandAST := mkGrammar(
		[]Node{
			mkRule("X",
				mkOptionalTerm(
					mkLiteralTerm("y"),
				),
			),
			mkRule("Z",
				mkRepeatOneTerm(
					mkRuleTerm("X"),
				),
			),
		},
		[]Node{
			mkSymbol("w", "z"),
		},
	)

	tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader(`X => ['y']
Z => <<X>>+
w = /z/
`))
	if err != nil {
		t.Error(err)
		return
	}
	start := ByHandGrammar.RulesForName("Grammar")[0]
	items, remaining, err := start.Parse(ByHandGrammar, tokens)
	if err != nil {
		t.Error(err)
		return
	}
	if len(remaining) != 0 {
		t.Errorf("Leftover tokens: %v.", remaining)
	}

	if !reflect.DeepEqual(byHandAST, AST(items)) {
		t.Error("Generated AST doesn't match by-hand AST.")
	}

	if true {
		dig := func(top AST) interface{} {
			return top
		}
		fmt.Println("byhand")
		printNode(dig(byHandAST), 0)
		fmt.Println("generated")
		printNode(dig(AST(items)), 0)
	}
}

func TestParseSymbol(t *testing.T) {
	tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader("'junkinthetrunk' stuff"))
	if err != nil {
		t.Error(err)
		return
	}
	term := InlineRuleTerm{Name: "literal"}
	items, _, err := term.Parse(ByHandGrammar, tokens)
	if err != nil {
		t.Error(err)
		return
	}
	st, ok := items[0].(SymbolText)
	if !ok {
		t.Errorf("Got a %T, expected a SymbolText.", items[0])
		return
	}
	if st.Type != "literal" {
		t.Errorf("Got a %q, expected a %q.", st.Type, "literal")
		return
	}
	if st.Text != "junkinthetrunk" {
		t.Errorf("Got %q, expected %q.", st.Text, "junkinthetrunk")
		return
	}
}

func TestParseTag(t *testing.T) {
	tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader("=> stuff"))
	if err != nil {
		t.Error(err)
		return
	}
	term := TagTerm{Tag: "hello"}
	items, remaining, err := term.Parse(ByHandGrammar, tokens)
	if err != nil {
		t.Error(err)
		return
	}
	if tag, ok := items[0].(Tag); ok {
		if tag != "hello" {
			t.Errorf("Expected %q, got %q.", "hello", tag)
			return
		}
	} else {
		t.Errorf("Expected Tag, got %T.", items[0])
		return
	}
	if !reflect.DeepEqual(remaining, tokens) {
		t.Errorf("Got wrong tokens remaining.")
		return
	}
}

func TestParseLiteral(t *testing.T) {
	tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader("=> stuff"))
	if err != nil {
		t.Error(err)
		return
	}
	term := LiteralTerm{Literal: "=>"}
	items, remaining, err := term.Parse(ByHandGrammar, tokens)
	if err != nil {
		t.Error(err)
		return
	}
	if lit, ok := items[0].(Literal); ok {
		if lit != "=>" {
			t.Errorf("Expected %q, got %q.", "=>", lit)
			return
		}
	} else {
		t.Errorf("Expected Literal, got %T.", items[0])
		return
	}
	if !reflect.DeepEqual(remaining, tokens[1:]) {
		t.Errorf("Got wrong tokens remaining.")
		return
	}
}
