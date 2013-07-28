package gopp

import (
	"errors"
	"fmt"
	"github.com/skelterjohn/debugtags"
	"strconv"
)

const debug = false

var tr = debugtags.Tracer{Enabled: false}

func (r Rule) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("Rule(%q)", r.Name)
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	return r.Expr.Parse(g, tokens)
}

func (e Expr) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("Expr")
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	for _, term := range e {
		var newItems []Node
		newItems, tokens, err = term.Parse(g, tokens)
		if err != nil {
			return
		}
		items = append(items, newItems...)
	}
	remainingTokens = tokens
	return
}

func (t RepeatZeroTerm) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
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
	for {
		subitems, subtokens, suberr := t.Term.Parse(g, remainingTokens)
		if suberr != nil {
			break
		}
		myitems = append(myitems, subitems...)
		remainingTokens = subtokens
	}
	items = []Node{myitems}
	return
}

func (t RepeatOneTerm) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
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
	for {
		subitems, subtokens, suberr := t.Term.Parse(g, remainingTokens)
		if suberr != nil {
			break
		}
		myitems = append(myitems, subitems...)
		remainingTokens = subtokens
	}
	items = []Node{myitems}
	if len(items) == 0 {
		err = errors.New("RepeatOneTerm found zero.")
	}
	return
}

func (t OptionalTerm) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
	rName := fmt.Sprintf("OptionalTerm")
	tr.In(rName, tokens)
	defer func() {
		if err == nil {
			tr.Out(rName, items)
		} else {
			tr.Out(rName, err)
		}
	}()

	subitem, subtokens, suberr := t.Expr.Parse(g, remainingTokens)
	if suberr != nil {
		remainingTokens = tokens
		return
	}
	items = append(items, subitem)
	remainingTokens = subtokens
	return
}

func (t RuleTerm) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
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
		return
	}

	var subitems []Node
	//fmt.Printf("%d rules for %q.\n", len(rules), t.Name)
	for _, rule := range rules {
		// if tt, ok := rule.Expr[0].(TagTerm); ok {
		// 	fmt.Printf("Trying %q.\n", tt.Tag)
		// }
		subitems, remainingTokens, err = rule.Parse(g, tokens)

		if err == nil {
			items = []Node{subitems}
			return
		}
	}

	return
}

func (t InlineRuleTerm) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
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
		items, remainingTokens, err = rule.Parse(g, tokens)

		if err == nil {
			return
		}
	}

	if _, ok := g.Symbol(t.Name); ok {
		if len(tokens) < 1 {
			err = errors.New("Need at least one token to make a symbol.")
			return
		}
		if t.Name == tokens[0].Type {
			st := SymbolText{
				Type: t.Name,
				Text: tokens[0].Text,
			}
			items = []Node{st}
			remainingTokens = tokens[1:]
			return
		}
	}

	err = fmt.Errorf("Unknown rule name: %q.", t.Name)

	return
}

func (t TagTerm) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
	items = []Node{Tag(t.Tag)}
	remainingTokens = tokens
	return
}

func (t LiteralTerm) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
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
		err = errors.New("Not enough tokens.")
		return
	}
	if tokens[0].Type != "RAW" {
		err = errors.New("Incorrect literal.")
		return
	}

	literalText := t.Literal
	quoted := fmt.Sprintf("\"%s\"", t.Literal)
	_ = quoted
	unquoted, qerr := strconv.Unquote(quoted)
	if qerr == nil {
		literalText = unquoted
	}

	if tokens[0].Text != literalText {
		err = errors.New("Incorrect literal.")
		return
	}
	items = []Node{Literal(literalText)}
	remainingTokens = tokens[1:]
	return
}

var _ = strconv.Unquote
