package gopp

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestFixedPointDecoder(t *testing.T) {
	df, err := NewDecoderFactory(goppgopp, "Grammar")
	if err != nil {
		t.Error(err)
		return
	}
	df.RegisterType(RepeatZeroTerm{})
	df.RegisterType(RepeatOneTerm{})
	df.RegisterType(OptionalTerm{})
	df.RegisterType(GroupTerm{})
	df.RegisterType(RuleTerm{})
	df.RegisterType(InlineRuleTerm{})
	df.RegisterType(TagTerm{})
	df.RegisterType(LiteralTerm{})
	var g Grammar
	dec := df.NewDecoder(strings.NewReader(goppgopp))
	err = dec.Decode(&g)
	if err != nil {
		t.Error(err)
		return
	}
	err = compareGrammars(g, ByHandGrammar)
	if err != nil {
		t.Error(err)
		return
	}

	// now see if the just-populated grammar can generate itself
	df.g = g
	dec = df.NewDecoder(strings.NewReader(goppgopp))
	var g2 Grammar
	err = dec.Decode(&g2)
	if err != nil {
		t.Error(err)
		return
	}
	err = compareGrammars(g, g2)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDecodeGrammar(t *testing.T) {
	var g Grammar
	ast, err := Parse(ByHandGrammar, "Grammar", strings.NewReader(goppgopp))
	if err != nil {
		t.Error(err)
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
	err = sa.Decode(&g)
	if err != nil {
		t.Error(err)
	}
	err = compareGrammars(g, ByHandGrammar)
	if err != nil {
		t.Error(err)
	}
}

func compareGrammars(g1, g2 Grammar) (err error) {
	for i := range g1.Rules {
		r1 := g1.Rules[i]
		if len(g2.Rules) <= i {
			err = fmt.Errorf("Not enough rules in g2: %d.", len(g1.Rules))
			return
		}
		r2 := g2.Rules[i]
		err = compareRules(r1, r2)
		if err != nil {
			return
		}
	}
	return
}

func compareRules(r1, r2 Rule) (err error) {
	if r1.Name != r2.Name {
		err = fmt.Errorf("Rule names %q and %q don't match.", r1.Name, r2.Name)
	}
	err = compareExprs(r1.Expr, r2.Expr)
	return
}

func compareExprs(e1, e2 Expr) (err error) {
	if len(e1) != len(e2) {
		err = fmt.Errorf("Expr lengths %d and %d don't match.", len(e1), len(e2))
		return
	}
	for i := range e1 {
		err = compareTerms(e1[i], e2[i])
		if err != nil {
			return
		}
	}
	return
}

func compareTerms(t1, t2 Term) (err error) {
	if reflect.TypeOf(t1) != reflect.TypeOf(t2) {
		err = fmt.Errorf("Term types %T and %T don't match.", t1, t2)
		return
	}
	switch t1 := t1.(type) {
	case RepeatZeroTerm:
		err = compareTerms(t1.Term, t2.(RepeatZeroTerm).Term)
	case RepeatOneTerm:
		err = compareTerms(t1.Term, t2.(RepeatOneTerm).Term)
	case LiteralTerm:
		if t1.Literal != t2.(LiteralTerm).Literal {
			err = fmt.Errorf("Literals %q and %q don't match.", t1.Literal, t2.(LiteralTerm).Literal)
			return
		}
	case TagTerm:
		if t1.Tag != t2.(TagTerm).Tag {
			err = fmt.Errorf("Tags %q and %q don't match.", t1.Tag, t2.(TagTerm).Tag)
			return
		}
	case RuleTerm:
		if t1.Name != t2.(RuleTerm).Name {
			err = fmt.Errorf("Names %q and %q don't match.", t1.Name, t2.(RuleTerm).Name)
			return
		}
	case InlineRuleTerm:
		if t1.Name != t2.(InlineRuleTerm).Name {
			err = fmt.Errorf("Names %q and %q don't match.", t1.Name, t2.(InlineRuleTerm).Name)
			return
		}
	default:
		fmt.Printf("%T\n", t1)
	}
	return
}

var ByHandGrammarREs []TypedRegexp

func init() {
	var err error
	ByHandGrammarREs, err = ByHandGrammar.TokenREs()
	if err != nil {
		panic(err)
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
				fmt.Printf("one has length %d, other %d\n", len(nl1), len(nl2))
				fmt.Println(nl1)
				fmt.Println(nl2)
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
		fmt.Printf("one is []Node, other is %T\n", n2)
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
		`Grammar => {type=Grammar} '\n'* {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*`,
		ByHandGoppAST[3].([]Node)[0],
	},
	{
		"Rule",
		`Rule => {field=Name} <identifier> '=>' {field=Expr} <Expr> '\n'+`,
		ByHandGoppAST[3].([]Node)[1],
	},
	{
		"Symbol",
		`Symbol => {field=Name} <identifier> '=' {field=Pattern} <regexp> '\n'+`,
		ByHandGoppAST[3].([]Node)[2],
	},
	{
		"Expr",
		`Expr => <<Term>>+`,
		ByHandGoppAST[3].([]Node)[3],
	},
	{
		"Term.1",
		`Term => <Term1>`,
		ByHandGoppAST[3].([]Node)[4],
	},
	{
		"Term.2",
		`Term => <Term2>`,
		ByHandGoppAST[3].([]Node)[5],
	},
	{
		"Term1.1",
		`Term1 => {type=RepeatZeroTerm} {field=Term} <<Term2>> '*'`,
		ByHandGoppAST[3].([]Node)[6],
	},
	{
		"Term1.2",
		`Term1 => {type=RepeatOneTerm} {field=Term} <<Term2>> '+'`,
		ByHandGoppAST[3].([]Node)[7],
	},
	{
		"Term2.1",
		`Term2 => {type=OptionalTerm} '[' {field=Expr} <Expr> ']'`,
		ByHandGoppAST[3].([]Node)[8],
	},
	{
		"Term2.2",
		`Term2 => {type=GroupTerm} '(' {field=Expr} <Expr> ')'`,
		ByHandGoppAST[3].([]Node)[9],
	},
	{
		"Term2.3",
		`Term2 => {type=RuleTerm} '<<' {field=Name} <identifier> '>>'`,
		ByHandGoppAST[3].([]Node)[10],
	},
	{
		"Term2.4",
		`Term2 => {type=InlineRuleTerm} '<' {field=Name} <identifier> '>'`,
		ByHandGoppAST[3].([]Node)[11],
	},
	{
		"Term2.5",
		`Term2 => {type=TagTerm} {field=Tag} <tag>`,
		ByHandGoppAST[3].([]Node)[12],
	},
	{
		"Term2.6",
		`Term2 => {type=LiteralTerm} {field=Literal} <literal>`,
		ByHandGoppAST[3].([]Node)[13],
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

		if false && th.Name == "Expr" {
			dig := func(top AST) interface{} {
				return top[2].([]Node)[0].([]Node)[4].([]Node)[0]
			}
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

		ok, indices := compareNodes(byHandAST, AST(items))
		if !ok {
			t.Errorf("%s: Generated AST doesn't match by-hand AST at %v.", th.Name, indices)
		}
	}
}

func TestParseFullGrammar(t *testing.T) {
	tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader(goppgopp))
	if err != nil {
		t.Error(err)
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
		t.Errorf("%s", err)
		return
	}
	if len(remaining) != 0 {
		t.Errorf("leftover tokens: %v.", remaining)
	}

	if false {
		dig := func(top AST) interface{} {
			return top[2].([]Node)[1].([]Node)[4].([]Node)[4]
		}
		byhand := dig(ByHandGoppAST)
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

	ok, indices := compareNodes(ByHandGoppAST, AST(items))
	if !ok {
		t.Errorf("Generated AST doesn't match by-hand AST at %v.", indices)
	}
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

	ok, indices := compareNodes(byHandAST, AST(items))
	if !ok {
		t.Errorf("Generated AST doesn't match by-hand AST at %v.", indices)
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

	ok, indices := compareNodes(byHandAST, AST(items))
	if !ok {
		t.Errorf("Generated AST doesn't match by-hand AST at %v.", indices)
	}

	if false {
		dig := func(top AST) interface{} {
			return top[2].([]Node)[0].([]Node)[4].([]Node)[0] //.([]Node)[3]
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
