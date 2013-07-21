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

var goppgopp = `Grammar => {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>+

Rule => {field=Name} <identifier> '=>' {field=Expr} <<Expr>> '\n'+

Symbol => {field=Name} <identifier> '=' {field=Pattern} <regexp> '\n'+

Expr => {field=.} <<Term>>+

Term => {type=RepeatZeroTerm} {field=Term} <<Term>> '*'
Term => {type=RepeatOneTerm} {field=Term} <<Term>> '+'
Term => {type=OptionalTerm} '[' {field=Expr} <<Expr>> ']'
Term => {type=GroupTerm} '(' {field=Expr} <<Expr>> ')'
Term => {type=RuleTerm} '<<' {field=Name} <identifier> '>>'
Term => {type=InlineRuleTerm} '<' {field=Name} <identifier> '>'
Term => {type=TagTerm} '{' {field=Tag} <identifier> '}'
Term => {type=LiteralTerm} {field=Literal} <literal>

identifier = /([a-zA-Z][a-zA-Z0-9_]*)/

literal = /'((?:[\\']|[^'])+?)'/

tag = /\{((?:[\\']|[^'])+?)\}/

regexp = /\/((?:\\/|[^\n])+?)\//
`

var gopptokens = []string{"Grammar","=>","field=Rules","<<","Rule",">>","+","field=Symbols","<<","Symbol",">>","+","Rule","=>","field=Name","<","identifier",">","=>","field=Expr","<<","Expr",">>","\\n","+","Symbol","=>","field=Name","<","identifier",">","=","field=Pattern","<","regexp",">","\\n","+","Expr","=>","field=.","<<","Term",">>","+","Term","=>","type=RepeatZeroTerm","field=Term","<<","Term",">>","*","Term","=>","type=RepeatOneTerm","field=Term","<<","Term",">>","+","Term","=>","type=OptionalTerm","[","field=Expr","<<","Expr",">>","]","Term","=>","type=GroupTerm","(","field=Expr","<<","Expr",">>",")","Term","=>","type=RuleTerm","<<","field=Name","<","identifier",">",">>","Term","=>","type=InlineRuleTerm","<","field=Name","<","identifier",">",">","Term","=>","type=TagTerm","{","field=Tag","<","identifier",">","}","Term","=>","type=LiteralTerm","field=Literal","<","literal",">","identifier","=","([a-zA-Z][a-zA-Z0-9_]*)","literal","=","'((?:[\\\\']|[^'])+?)'","tag","=","\\{((?:[\\\\']|[^'])+?)\\}","regexp","=","\\/((?:\\\\/|[^\\n])+?)\\/"}

func TestTokenREs(t *testing.T) {
	res, err := ByHandGrammar.TokenREs()
	if err != nil {
		t.Error(err)
	}

	counter := 0
	r := strings.NewReader(goppgopp)
	for token := range Tokenize(res, r) {
		if false && token != gopptokens[counter] {
			t.Errorf("Token %d: expected %q, got %q.", counter, gopptokens[counter], token)
		}
		counter++
	}
	if counter != len(gopptokens) {
		t.Errorf("Expected %d tokens, got %d.", len(gopptokens), counter)
	}
}
