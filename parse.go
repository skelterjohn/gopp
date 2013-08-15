// Copyright 2013 The gopp AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopp

import (
	"errors"
	"fmt"
	"github.com/skelterjohn/debugtags"
	"strconv"
)

func Parse(g Grammar, startRule string, document []byte) (ast AST, err error) {
	tokenREs, err := g.TokenREs()
	if err != nil {
		return
	}
	ignoreREs, err := g.IgnoreREs()
	if err != nil {
		return
	}
	ti := TokenizeInfo{
		TokenREs:  tokenREs,
		IgnoreREs: ignoreREs,
	}
	tokens, err := Tokenize(ti, document)
	if err != nil {
		return
	}
	rules := g.RulesForName(startRule)
	if len(rules) != 1 {
		err = fmt.Errorf("Rule %q had %d definitions.", startRule, len(rules))
		return
	}
	start := rules[0]
	pd := NewParseData()
	items, remaining, err := start.Parse(g, tokens, pd, []string{})

	if err != nil {
		// TODO: use pd to return informative error messages.
		err = pd.FarthestErrors[0]
		return
	}
	if len(remaining) != 0 {
		err = errors.New("Did not parse entire file.")
	}

	ast = items

	return
}

const debug = false

func SetTr(e bool) {
	tr.Enabled = e
}

var tr = debugtags.Tracer{Enabled: false}

type ParseData struct {
	accepted             bool
	LastUnacceptedTokens []Token
	errored              bool
	FarthestErrors       []error
	TokensForError       []Token
}

func NewParseData() (pd *ParseData) {
	pd = &ParseData{}
	return
}

func (pd *ParseData) AcceptUpTo(remaining []Token) {
	if !pd.accepted || len(remaining) < len(pd.LastUnacceptedTokens) {
		pd.LastUnacceptedTokens = remaining
	}
	pd.accepted = true
}

func (pd *ParseData) ErrorWith(err error, remaining []Token) {
	if !pd.errored || len(remaining) < len(pd.TokensForError) {
		pd.FarthestErrors = append(pd.FarthestErrors, err)
		pd.TokensForError = remaining
	}
	pd.errored = true
}

func (r Rule) Parse(g Grammar, tokens []Token, pd *ParseData, parentRuleNames []string) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("Rule(%q)", r.Name)
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	for _, n := range parentRuleNames {
		if n == r.Name {
			err = fmt.Errorf("Rule cycle with %q.", r.Name)
			return
		}
	}

	items, remainingTokens, err = r.Expr.Parse(g, tokens, pd, append(parentRuleNames, r.Name))
	return
}

func (e Expr) Parse(g Grammar, tokens []Token, pd *ParseData, parentRuleNames []string) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("Expr")
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	startTokens := tokens

	for _, term := range e {
		var newItems []Node
		var prns []string
		if len(startTokens) == len(tokens) {
			prns = parentRuleNames
		}
		newItems, tokens, err = term.Parse(g, tokens, pd, prns)
		if err != nil {
			return
		}
		items = append(items, newItems...)
	}
	remainingTokens = tokens
	return
}

func (t RepeatZeroTerm) Parse(g Grammar, tokens []Token, pd *ParseData, parentRuleNames []string) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("RepeatZeroTerm")
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	remainingTokens = tokens
	var myitems []Node
	first := true
	for {
		var prns []string
		if first {
			prns = parentRuleNames
			first = false
		}
		subitems, subtokens, suberr := t.Term.Parse(g, remainingTokens, pd, prns)
		if suberr != nil {
			break
		}
		myitems = append(myitems, subitems...)
		remainingTokens = subtokens
	}
	items = []Node{myitems}
	return
}

func (t RepeatOneTerm) Parse(g Grammar, tokens []Token, pd *ParseData, parentRuleNames []string) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("RepeatOneTerm")
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	remainingTokens = tokens
	var myitems []Node
	first := true
	var suberr error
	for {
		var prns []string
		if first {
			prns = parentRuleNames
			first = false
		}
		var subitems []Node
		var subtokens []Token
		subitems, subtokens, suberr = t.Term.Parse(g, remainingTokens, pd, prns)
		if suberr != nil {
			break
		}
		myitems = append(myitems, subitems...)
		remainingTokens = subtokens
	}
	items = []Node{myitems}
	if len(items) == 0 {
		err = suberr
		pd.ErrorWith(err, tokens)
	}
	return
}

func (t OptionalTerm) Parse(g Grammar, tokens []Token, pd *ParseData, parentRuleNames []string) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("OptionalTerm")
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	subitems, subtokens, suberr := t.Expr.Parse(g, tokens, pd, parentRuleNames)
	if suberr != nil {
		remainingTokens = tokens
		return
	}
	items = subitems
	remainingTokens = subtokens
	return
}

func (t RuleTerm) Parse(g Grammar, tokens []Token, pd *ParseData, parentRuleNames []string) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("RuleTerm(%q)", t.Name)
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	rules := g.RulesForName(t.Name)
	if len(rules) == 0 {
		err = fmt.Errorf("Unknown rule name: %q.", t.Name)
		pd.ErrorWith(err, tokens)
		return
	}

	var subitems []Node
	//fmt.Printf("%d rules for %q.\n", len(rules), t.Name)
	for _, rule := range rules {
		// if tt, ok := rule.Expr[0].(TagTerm); ok {
		// 	fmt.Printf("Trying %q.\n", tt.Tag)
		// }
		subitems, remainingTokens, err = rule.Parse(g, tokens, pd, parentRuleNames)

		if err == nil {
			items = []Node{subitems}
			return
		}
	}

	return
}

func (t InlineRuleTerm) Parse(g Grammar, tokens []Token, pd *ParseData, parentRuleNames []string) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("InlineRuleTerm(%q)", t.Name)
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	rules := g.RulesForName(t.Name)
	for _, rule := range rules {
		items, remainingTokens, err = rule.Parse(g, tokens, pd, parentRuleNames)

		if err == nil {
			return
		}
	}
	err = nil
	if _, ok := g.Symbol(t.Name); ok {
		if len(tokens) < 1 {
			err = errors.New("Need at least one token to make a symbol.")
			pd.ErrorWith(err, tokens)
			return
		}
		if t.Name == tokens[0].Type {
			st := SymbolText{
				Type: t.Name,
				Text: tokens[0].Text,
			}
			items = []Node{st}
			remainingTokens = tokens[1:]
			pd.AcceptUpTo(remainingTokens)
			return
		}
		err = fmt.Errorf("Expected %s at %d:%d.", t.Name, tokens[0].Row, tokens[0].Col)
		pd.ErrorWith(err, tokens)
		return
	}

	err = fmt.Errorf("Unknown rule name: %q.", t.Name)
	pd.ErrorWith(err, tokens)

	return
}

func (t TagTerm) Parse(g Grammar, tokens []Token, pd *ParseData, parentRuleNames []string) (items []Node, remainingTokens []Token, err error) {
	tr.Println(Tag(t.Tag))
	items = []Node{Tag(t.Tag)}
	remainingTokens = tokens
	return
}

func (t LiteralTerm) Parse(g Grammar, tokens []Token, pd *ParseData, parentRuleNames []string) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("LiteralTerm(%q)", t.Literal)
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	if len(tokens) == 0 {
		err = fmt.Errorf("Expected %q at EOF.", t.Literal)
		pd.ErrorWith(err, tokens)
		return
	}
	if tokens[0].Type != "RAW" {
		err = fmt.Errorf("Expected %q at %d:%d.", t.Literal, tokens[0].Row, tokens[0].Col)
		pd.ErrorWith(err, tokens)
		return
	}

	literalText := t.Literal
	// quoted := fmt.Sprintf("\"%s\"", t.Literal)
	// _ = quoted
	// unquoted, qerr := strconv.Unquote(quoted)
	unquoted, qerr := descapeString(t.Literal)

	if qerr == nil && false {
		literalText = unquoted
	}

	if tokens[0].Text != literalText {
		err = fmt.Errorf("Expected %q at %d:%d.", t.Literal, tokens[0].Row, tokens[0].Col)
		pd.ErrorWith(err, tokens)
		return
	}
	items = []Node{Literal(literalText)}
	remainingTokens = tokens[1:]
	pd.AcceptUpTo(remainingTokens)
	return

}

var _ = strconv.Unquote
