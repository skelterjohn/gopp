package gopp

import (
	"regexp"
	"sort"
)

const REGEXP_PREFIX = `^(?:\s)*`

type Grammar struct {
	Rules   []*Rule
	Symbols []*Symbol
}

func (g *Grammar) CollectLiterals(literals map[string]bool) {
	for _, rule := range g.Rules {
		rule.CollectLiterals(literals)
	}
	return
}

func (g *Grammar) TokenREs() (res []*regexp.Regexp, err error) {
	// first get all the literals, and sort them longest first (so smaller ones don't eat larger ones).
	literals := map[string]bool{}
	g.CollectLiterals(literals)
	sortedLiterals := literalSorter{}
	for literal := range literals {
		sortedLiterals = append(sortedLiterals, literal)
	}
	sort.Sort(sortedLiterals)
	for _, literal := range sortedLiterals {
		re, err := regexp.Compile(REGEXP_PREFIX + "(" + regexp.QuoteMeta(literal) + ")")
		if err != nil {
			panic("regexp.QuoteMeta returned something that didn't compile")
		}
		res = append(res, re)
	}
	for _, symbol := range g.Symbols {
		var re *regexp.Regexp
		re, err = regexp.Compile(REGEXP_PREFIX + symbol.Pattern)
		if err != nil {
			return
		}
		res = append(res, re)
	}
	return
}

type Rule struct {
	Name string
	Expr *Expr
}

func (r *Rule) CollectLiterals(literals map[string]bool) {
	r.Expr.CollectLiterals(literals)
	return
}

type Symbol struct {
	Name    string
	Pattern string
}

type Expr struct {
	Terms []*Term
}

func (e *Expr) CollectLiterals(literals map[string]bool) {
	for _, term := range e.Terms {
		term.CollectLiterals(literals)
	}
	return
}

type Term struct {
	Operator string
	Term     *Term
	Expr     *Expr
	Field    string
	Name     string
	Literal  string
}

func (t *Term) CollectLiterals(literals map[string]bool) {
	if t.Literal != "" {
		literals[t.Literal] = true
	}
	if t.Expr != nil {
		t.Expr.CollectLiterals(literals)
	}
	if t.Term != nil {
		t.Term.CollectLiterals(literals)
	}
	return
}
