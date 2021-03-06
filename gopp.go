// Copyright 2013 The gopp AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopp

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Grammar struct {
	LexSteps []LexStep
	Rules    []Rule
	Symbols  []Symbol
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
		re, err := regexp.Compile("^(" + regexp.QuoteMeta(literal) + ")")
		if err != nil {
			panic("regexp.QuoteMeta returned something that didn't compile")
		}
		res = append(res, TypedRegexp{"RAW", re})
	}
	for _, symbol := range g.Symbols {
		var re *regexp.Regexp
		re, err = regexp.Compile("^" + symbol.Pattern)
		if err != nil {
			return
		}
		res = append(res, TypedRegexp{symbol.Name, re})
	}
	return
}

func (g Grammar) IgnoreREs() (res []*regexp.Regexp, err error) {
	for _, ls := range g.LexSteps {
		if ls.Name == "ignore" {
			var re *regexp.Regexp
			re, err = regexp.Compile(ls.Pattern)
			if err != nil {
				return
			}
			res = append(res, re)
		}
	}
	return
}

type LexStep struct {
	Name    string
	Pattern string
}

type Rule struct {
	Name string
	Expr
}

func (r Rule) String() string {
	return fmt.Sprintf("Rule(%s:%v)", r.Name, r.Expr)
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
	Parse(g Grammar, tokens []Token, pd *ParseData, parentRuleNames []string) (items []Node, remainingTokens []Token, err error)
}

type RepeatZeroTerm struct {
	Term
}

func (rzt RepeatZeroTerm) String() string {
	return fmt.Sprintf("RepeatZeroTerm(%v)", rzt.Term)
}

type RepeatOneTerm struct {
	Term
}

func (rot RepeatOneTerm) String() string {
	return fmt.Sprintf("RepeatOneTerm(%v)", rot.Term)
}

type OptionalTerm struct {
	Expr
}

func (ot OptionalTerm) String() string {
	return fmt.Sprintf("OptionalTerm(%v)", ot.Expr)
}

type GroupTerm struct {
	Expr
}

func (gt GroupTerm) String() string {
	return fmt.Sprintf("GroupTerm(%v)", gt.Expr)
}

type noLiterals struct{}

func (n noLiterals) CollectLiterals(literals map[string]bool) {
	return
}

type RuleTerm struct {
	Name string
	noLiterals
}

func (rt RuleTerm) String() string {
	return fmt.Sprintf("RuleTerm(%s)", rt.Name)
}

type InlineRuleTerm struct {
	Name string
	noLiterals
}

func (irt InlineRuleTerm) String() string {
	return fmt.Sprintf("InlineRuleTerm(%s)", irt.Name)
}

type TagTerm struct {
	Tag string
	noLiterals
}

func (tt TagTerm) String() string {
	return fmt.Sprintf("TagTerm(%q)", tt.Tag)
}

type LiteralTerm struct {
	Literal string
}

func (lt LiteralTerm) String() string {
	return fmt.Sprintf("LiteralTerm(%q)", lt.Literal)
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
