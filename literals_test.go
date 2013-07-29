package gopp

import (
	"fmt"
	"strings"
	"testing"
)

var _ = fmt.Println

func TestCollectLiterals(t *testing.T) {
	correctLiterals := []string{
		"=>",
		"=",
		"[",
		"]",
		"(",
		")",
		"<",
		">",
		"<<",
		">>",
		"*",
		"+",
		"\n",
	}

	literals := map[string]bool{}
	ByHandGrammar.CollectLiterals(literals)
	if len(literals) != len(correctLiterals) {
		t.Error("Wrong number of literals")
	}
	for _, literal := range correctLiterals {
		if !literals[literal] {
			t.Errorf("Could not find %q", literal)
		}
	}
}

var symbolTests = map[string][][]string{
	"identifier": [][]string{
		{"stuff", "stuff"},
		{"xyz123", "xyz123"},
		{"x_b", "x_b"},
	},
}
var symbolFailTests = map[string][]string{
	"identifier": []string{
		"123",
		".sdf-",
		"!",
	},
}

func TestSymbolTokenize(t *testing.T) {
	for typ, examples := range symbolTests {
		for _, example := range examples {
			tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader(example[0]))
			if err != nil {
				t.Error(err)
				continue
			}
			if len(tokens) == 0 {
				t.Error("No tokens for %q.", example[0])
				continue
			}
			if typ != tokens[0].Type {
				t.Errorf("Expected type %q, got %q.", typ, tokens[0].Type)
				continue
			}
			if example[1] != tokens[0].Text {
				t.Errorf("Expected %q, got %q.", example[1], tokens[0].Text)
				continue
			}
		}
	}
}

func TestSymbolFailTokenize(t *testing.T) {
	for typ, examples := range symbolFailTests {
		for _, example := range examples {
			tokens, err := Tokenize(ByHandGrammarREs, strings.NewReader(example))
			if err != nil {
				continue
			}
			if len(tokens) == 0 {
				continue
			}
			if typ == tokens[0].Type {
				t.Errorf("Mistakenly parsed %q as %q.", example, typ)
			}
		}
	}
}

var goppgopp = `
Grammar => '\n'* {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*
Rule => {field=Name} <identifier> '=>' {field=Expr} <Expr> '\n'+
Symbol => {field=Name} <identifier> '=' {field=Pattern} <regexp> '\n'+
Expr => <<Term>>+
Term => <Term1>
Term => <Term2>
Term1 => {type=RepeatZeroTerm} {field=Term} <<Term2>> '*'
Term1 => {type=RepeatOneTerm} {field=Term} <<Term2>> '+'
Term2 => {type=OptionalTerm} '[' {field=Expr} <Expr> ']'
Term2 => {type=GroupTerm} '(' {field=Expr} <Expr> ')'
Term2 => {type=RuleTerm} '<<' {field=Name} <identifier> '>>'
Term2 => {type=InlineRuleTerm} '<' {field=Name} <identifier> '>'
Term2 => {type=TagTerm} {field=Tag} <tag>
Term2 => {type=LiteralTerm} {field=Literal} <literal>
identifier = /([a-zA-Z][a-zA-Z0-9_]*)/
literal = /'((?:[\\']|[^'])+?)'/
tag = /\{((?:[\\']|[^'])+?)\}/
regexp = /\/((?:\\/|[^\n])+?)\//
`

func xTestTokenREs(t *testing.T) {
	res, err := ByHandGrammar.TokenREs()
	if err != nil {
		t.Error(err)
	}

	counter := 0
	r := strings.NewReader(goppgopp)
	tokens, err := Tokenize(res, r)
	if err != nil {
		t.Error(err)
	}
	for _, token := range tokens {
		if token != goppTokens[counter] {
			t.Errorf("Expected %v, got %v.", goppTokens[counter], token)
		}
		counter++
	}
	if counter != len(goppTokens) {
		t.Errorf("Expected %d tokens, got %d.", len(goppTokens), counter)
	}
}

var goppTokens = []Token{
	Token{"identifier", "Grammar", "Grammar"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {field=Rules}", "field=Rules"},
	Token{"RAW", " <<", "<<"},
	Token{"identifier", "Rule", "Rule"},
	Token{"RAW", ">>", ">>"},
	Token{"RAW", "+", "+"},
	Token{"tag", " {field=Symbols}", "field=Symbols"},
	Token{"RAW", " <<", "<<"},
	Token{"identifier", "Symbol", "Symbol"},
	Token{"RAW", ">>", ">>"},
	Token{"RAW", "+", "+"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Rule", "Rule"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {field=Name}", "field=Name"},
	Token{"RAW", " <", "<"},
	Token{"identifier", "identifier", "identifier"},
	Token{"RAW", ">", ">"},
	Token{"literal", " '=>'", "=>"},
	Token{"tag", " {field=Expr}", "field=Expr"},
	Token{"RAW", " <<", "<<"},
	Token{"identifier", "Expr", "Expr"},
	Token{"RAW", ">>", ">>"},
	Token{"literal", " '\\n'", "\\n"},
	Token{"RAW", "+", "+"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Symbol", "Symbol"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {field=Name}", "field=Name"},
	Token{"RAW", " <", "<"},
	Token{"identifier", "identifier", "identifier"},
	Token{"RAW", ">", ">"},
	Token{"literal", " '='", "="},
	Token{"tag", " {field=Pattern}", "field=Pattern"},
	Token{"RAW", " <", "<"},
	Token{"identifier", "regexp", "regexp"},
	Token{"RAW", ">", ">"},
	Token{"literal", " '\\n'", "\\n"},
	Token{"RAW", "+", "+"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Expr", "Expr"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {field=.}", "field=."},
	Token{"RAW", " <<", "<<"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", ">>", ">>"},
	Token{"RAW", "+", "+"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {type=RepeatZeroTerm}", "type=RepeatZeroTerm"},
	Token{"tag", " {field=Term}", "field=Term"},
	Token{"RAW", " <<", "<<"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", ">>", ">>"},
	Token{"literal", " '*'", "*"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {type=RepeatOneTerm}", "type=RepeatOneTerm"},
	Token{"tag", " {field=Term}", "field=Term"},
	Token{"RAW", " <<", "<<"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", ">>", ">>"},
	Token{"literal", " '+'", "+"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {type=OptionalTerm}", "type=OptionalTerm"},
	Token{"literal", " '['", "["},
	Token{"tag", " {field=Expr}", "field=Expr"},
	Token{"RAW", " <<", "<<"},
	Token{"identifier", "Expr", "Expr"},
	Token{"RAW", ">>", ">>"},
	Token{"literal", " ']'", "]"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {type=GroupTerm}", "type=GroupTerm"},
	Token{"literal", " '('", "("},
	Token{"tag", " {field=Expr}", "field=Expr"},
	Token{"RAW", " <<", "<<"},
	Token{"identifier", "Expr", "Expr"},
	Token{"RAW", ">>", ">>"},
	Token{"literal", " ')'", ")"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {type=RuleTerm}", "type=RuleTerm"},
	Token{"literal", " '<<'", "<<"},
	Token{"tag", " {field=Name}", "field=Name"},
	Token{"RAW", " <", "<"},
	Token{"identifier", "identifier", "identifier"},
	Token{"RAW", ">", ">"},
	Token{"literal", " '>>'", ">>"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {type=InlineRuleTerm}", "type=InlineRuleTerm"},
	Token{"literal", " '<'", "<"},
	Token{"tag", " {field=Name}", "field=Name"},
	Token{"RAW", " <", "<"},
	Token{"identifier", "identifier", "identifier"},
	Token{"RAW", ">", ">"},
	Token{"literal", " '>'", ">"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {type=TagTerm}", "type=TagTerm"},
	Token{"tag", " {field=.}", "field=."},
	Token{"RAW", " <", "<"},
	Token{"identifier", "tag", "tag"},
	Token{"RAW", ">", ">"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "Term", "Term"},
	Token{"RAW", " =>", "=>"},
	Token{"tag", " {type=LiteralTerm}", "type=LiteralTerm"},
	Token{"tag", " {field=Literal}", "field=Literal"},
	Token{"RAW", " <", "<"},
	Token{"identifier", "literal", "literal"},
	Token{"RAW", ">", ">"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "identifier", "identifier"},
	Token{"RAW", " =", "="},
	Token{"regexp", " /([a-zA-Z][a-zA-Z0-9_]*)/", "([a-zA-Z][a-zA-Z0-9_]*)"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "literal", "literal"},
	Token{"RAW", " =", "="},
	Token{"regexp", " /'((?:[\\\\']|[^'])+?)'/", "'((?:[\\\\']|[^'])+?)'"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "tag", "tag"},
	Token{"RAW", " =", "="},
	Token{"regexp", " /\\{((?:[\\\\']|[^'])+?)\\}/", "\\{((?:[\\\\']|[^'])+?)\\}"},
	Token{"RAW", "\n", "\n"},
	Token{"identifier", "regexp", "regexp"},
	Token{"RAW", " =", "="},
	Token{"regexp", " /\\/((?:\\\\/|[^\\n])+?)\\//", "\\/((?:\\\\/|[^\\n])+?)\\/"},
	Token{"RAW", "\n", "\n"},
}
