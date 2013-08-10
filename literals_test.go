package gopp

import (
	"fmt"
	"testing"
)

var _ = fmt.Println

func TestCollectLiterals(t *testing.T) {
	correctLiterals := []string{
		"=>",
		"=",
		":",
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
			tokens, err := Tokenize(TokenizeInfo{TokenREs: ByHandGrammarREs}, []byte(example[0]))
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
			tokens, err := Tokenize(TokenizeInfo{TokenREs: ByHandGrammarREs}, []byte(example))
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
ignore: /^#.*\n/
ignore: /^(?:[ \t])+/
Grammar => {type=Grammar} '\n'* {field=LexSteps} <<LexStep>>* {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*
LexStep => {field=Name} <identifier> ':' {field=Pattern} <regexp> '\n'+
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
