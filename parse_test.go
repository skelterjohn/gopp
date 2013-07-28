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

func compareNodes(n1, n2 Node) (ok bool, indicesToError []int) {
	if a, ok := n1.(AST); ok {
		n1 = []Node(a)
	}
	if a, ok := n2.(AST); ok {
		n2 = []Node(a)
	}
	if nl1, isList1 := n1.([]Node); isList1 {
		if nl2, isList2 := n2.([]Node); isList2 {
			if len(nl1) != len(nl2) {
				ok = false
				return
			}
			ok = true
			for i := range nl1 {
				var subindices []int
				ok, subindices = compareNodes(nl1[i], nl2[i])
				if !ok {
					indicesToError = append([]int{i}, subindices...)
					return
				}
			}
			return
		}
		ok = false
		return
	}
	ok = n1 == n2
	return
}

type textByHand struct {
	Name   string
	Text   string
	ByHand Node
}

var rulesTextAndByHand = []textByHand{
	{
		"Grammar",
		`Grammar => '\n'* {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>+`,
		ByHandGoppAST[2].([]Node)[0],
	},
	{
		"Rule",
		`Rule => {field=Name} <identifier> '=>' {field=Expr} <<Expr>> '\n'+`,
		ByHandGoppAST[2].([]Node)[1],
	},
}

func TestParseRulesIndividual(t *testing.T) {
	for _, th := range rulesTextAndByHand {
		rule := th.ByHand
		byHandAST := mkGrammar(
			[]Node{rule},
			[]Node{},
		)

		txt := fmt.Sprintf("\n%s\n", th.Text)
		tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader(txt))
		if err != nil {
			t.Errorf("%s: %s", th.Name, err)
			return
		}
		start := ByHandGrammar.RulesForName("Grammar")[0]
		// tr.Enabled = true
		pd := &ParseData{}
		items, remaining, err := start.Parse(ByHandGrammar, tokens, pd)
		// tr.Enabled = false
		if err != nil {
			fmt.Printf("Remaining: %v\n", pd.TokensForError)
			for _, err := range pd.FarthestErrors {
				fmt.Printf(" - %s\n", err)
			}
			t.Errorf("%s: %s", th.Name, err)
			return
		}
		if len(remaining) != 0 {
			t.Errorf("%s: leftover tokens: %v.", th.Name, remaining)
		}

		dig := func(top AST) interface{} {
			return top
		}

		if false && th.Name == "Grammar" {
			byhand := dig(byHandAST)
			gen := dig(AST(items))
			ok, indices := compareNodes(byhand, gen)
			if !ok {
				fmt.Println("byhand")
				printNode(byhand, 0)
				fmt.Println("generated")
				printNode(gen, 0)
				fmt.Println(ok, indices)
			}
		}

		ok, indices := compareNodes(dig(byHandAST), dig(AST(items)))
		if !ok {
			t.Errorf("%s: Generated AST doesn't match by-hand AST at %v.", th.Name, indices)
		}
	}
}

func xTestParseFullGrammar(t *testing.T) {
	tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader(goppgopp))
	if err != nil {
		t.Error(err)
		return
	}
	start := ByHandGrammar.RulesForName("Grammar")[0]
	pd := &ParseData{}
	_, remaining, err := start.Parse(ByHandGrammar, tokens, pd)
	if err != nil {
		t.Error(err)
		return
	}
	if len(remaining) != 0 {
		t.Error("Leftover tokens: %v.", remaining)
	}
	//printNode(items, 0)
}

func TestParseEasyGrammar(t *testing.T) {
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

	tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader(`
X => 'y'
w = /z/
`))
	if err != nil {
		t.Error(err)
		return
	}
	start := ByHandGrammar.RulesForName("Grammar")[0]
	pd := &ParseData{}
	items, remaining, err := start.Parse(ByHandGrammar, tokens, pd)
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

	if false {
		dig := func(top AST) interface{} {
			return top
		}
		fmt.Println("byhand")
		printNode(dig(byHandAST), 0)
		fmt.Println("generated")
		printNode(dig(AST(items)), 0)
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

	tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader(`
X => ['y']
Z => <<X>>+
w = /z/
`))
	if err != nil {
		t.Error(err)
		return
	}
	start := ByHandGrammar.RulesForName("Grammar")[0]
	pd := &ParseData{}
	items, remaining, err := start.Parse(ByHandGrammar, tokens, pd)
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

	if false {
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
	pd := &ParseData{}
	items, _, err := term.Parse(ByHandGrammar, tokens, pd)
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
	pd := &ParseData{}
	items, remaining, err := term.Parse(ByHandGrammar, tokens, pd)
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
	pd := &ParseData{}
	items, remaining, err := term.Parse(ByHandGrammar, tokens, pd)
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
