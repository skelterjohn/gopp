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
		"[",
		"]",
		")",
		"<",
		">",
		"\n",
		"=",
		"*",
		"+",
		"(",
		"<<",
		">",
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

var goppgopp = `Grammar => Rules=<<Rule>>+ Symbols=<<Symbol>>+

Rule => Name=<identifier> '=>' Expr=<<Expr>> '\n'+

Symbol => Name=<identifier> '=' Pattern=<regexp> '\n'+

Expr => Terms=<<Term>>+

Term => Term=<<Term>> Operator='*'
Term => Term=<<Term>> Operator='+'
Term => Operator='[' Expr=<<Expr>> ']'
Term => Operator='(' Expr=<<Expr>> ')'
Term => Operator='<<' Name=<identifier> '>>'
Term => Operator='<' Name=<identifier> '>'
Term => Field=<indentifier> Operator='=' Term=<<Term>>
Term => Literal=<literal>

identifier = /([a-zA-Z][a-zA-Z0-9_]*)/

literal = /'((?:[\\']|[^'])+?)'/

regexp = /\/((?:\\/|[^\n])+?)\//
`

var gopptokens = []string{"Grammar","=>","Rules","=","<<","Rule",">>","+","Symbols","=","<<","Symbol",">>","+","Rule","=>","Name","=","<","identifier",">","=>","Expr","=","<<","Expr",">>","\\n","+","Symbol","=>","Name","=","<","identifier",">","=","Pattern","=","<","regexp",">","\\n","+","Expr","=>","Terms","=","<<","Term",">>","+","Term","=>","Term","=","<<","Term",">>","Operator","=","*","Term","=>","Term","=","<<","Term",">>","Operator","=","+","Term","=>","Operator","=","[","Expr","=","<<","Expr",">>","]","Term","=>","Operator","=","(","Expr","=","<<","Expr",">>",")","Term","=>","Operator","=","<<","Name","=","<","identifier",">",">>","Term","=>","Operator","=","<","Name","=","<","identifier",">",">","Term","=>","Field","=","<","indentifier",">","Operator","=","=","Term","=","<<","Term",">>","Term","=>","Literal","=","<","literal",">","identifier","=","([a-zA-Z][a-zA-Z0-9_]*)","literal","=","'((?:[\\\\']|[^'])+?)'","regexp","=","\\/((?:\\\\/|[^\\n])+?)\\/"}

func TestTokenREs(t *testing.T) {
	res, err := ByHandGrammar.TokenREs()
	if err != nil {
		t.Error(err)
	}

	counter := 0
	r := strings.NewReader(goppgopp)
	for token := range Tokenize(res, r) {
		if token != gopptokens[counter] {
			t.Errorf("Token %d: expected %q, got %q.", counter, gopptokens[counter], token)
		}
		counter++
	}
}
