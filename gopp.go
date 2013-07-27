package gopp

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const REGEXP_PREFIX = `^(?: )*`

type Grammar struct {
	Rules   []Rule
	Symbols []Symbol
}

func (g Grammar) RulesForName(name string) (rs []Rule) {
	for _, rule := range g.Rules {
		if rule.Name == name {
			rs = append(rs, rule)
		}
	}
	return
}

func (g Grammar) Symbol(name string) (s Symbol, ok bool) {
	for _, symb := range g.Symbols {
		if symb.Name == name {
			s = symb
			ok = true
			return
		}
	}
	return
}

func (g Grammar) CollectLiterals(literals map[string]bool) {
	for _, rule := range g.Rules {
		rule.CollectLiterals(literals)
	}
	return
}

type TypedRegexp struct {
	Type string
	*regexp.Regexp
}

func (g Grammar) TokenREs() (res []TypedRegexp, err error) {
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
		res = append(res, TypedRegexp{"RAW", re})
	}
	for _, symbol := range g.Symbols {
		var re *regexp.Regexp
		re, err = regexp.Compile(REGEXP_PREFIX + symbol.Pattern)
		if err != nil {
			return
		}
		res = append(res, TypedRegexp{symbol.Name, re})
	}
	return
}

type Rule struct {
	Name string
	Expr
}

type Symbol struct {
	Name    string
	Pattern string
}

type Expr []Term

func (e Expr) CollectLiterals(literals map[string]bool) {
	for _, term := range e {
		term.CollectLiterals(literals)
	}
	return
}

type Term interface {
	CollectLiterals(literals map[string]bool)
	Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error)
}

type RepeatZeroTerm struct {
	Term
}

type RepeatOneTerm struct {
	Term
}

type OptionalTerm struct {
	Expr
}

type GroupTerm struct {
	Expr
}

type noLiterals struct{}

func (n noLiterals) CollectLiterals(literals map[string]bool) {
	return
}

type RuleTerm struct {
	Name string
	noLiterals
}

type InlineRuleTerm struct {
	Name string
	noLiterals
}

type TagTerm struct {
	Tag string
	noLiterals
}

type LiteralTerm struct {
	Literal string
}

func (l LiteralTerm) CollectLiterals(literals map[string]bool) {
	literals[l.Literal] = true
	return
}

type AST []Node

type Node interface{}
type Tag string

func (t Tag) String() string {
	return fmt.Sprintf("Tag(%s)", string(t))
}

type Literal string

func (l Literal) String() string {
	return fmt.Sprintf("Literal(%s)", string(strings.Replace(string(l), "\n", `\n`, -1)))
}

type Identifier string

func (i Identifier) String() string {
	return fmt.Sprintf("Identifier(%s)", string(i))
}

type Regexp string

func (r Regexp) String() string {
	return fmt.Sprintf("Regexp(%s)", string(r))
}

type SymbolText struct {
	Type string
	Text string
}

func (s SymbolText) String() string {
	return fmt.Sprintf("<%s:%q>", s.Type, s.Text)
}
