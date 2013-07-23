package gopp

import (
	"fmt"
	"regexp"
	"errors"
)

func (r Rule) Parse(g Grammar, tokens []Token) (items []interface{}, remainingTokens []Token, err error) {
	return
}

func (e Expr) Parse(g Grammar, tokens []Token) (items []interface{}, remainingTokens []Token, err error) {
	for _, term := range e {
		var newItem interface{}
		newItem, tokens, err = term.Parse(g, tokens)
		if err != nil {
			return
		}
		items = append(items, newItem)
	}
	remainingTokens = tokens
	return
}

func (t RepeatZeroTerm) Parse(g Grammar, tokens []Token) (items []interface{}, remainingTokens []Token, err error) {
	remainingTokens = tokens
	for {
		subitem, subtokens, suberr := t.Term.Parse(g, remainingTokens)
		if suberr != nil {
			break
		}
		items = append(items, subitem)
		remainingTokens = subtokens
	}
	return
}

func (t RepeatOneTerm) Parse(g Grammar, tokens []Token) (items []interface{}, remainingTokens []Token, err error) {
	remainingTokens = tokens
	for {
		subitem, subtokens, suberr := t.Term.Parse(g, remainingTokens)
		if suberr != nil {
			break
		}
		items = append(items, subitem)
		remainingTokens = subtokens
	}
	if len(items) == 0 {
		err = errors.New("RepeatOneTerm found zero.")
	}
	return
}

func (t OptionalTerm) Parse(g Grammar, tokens []Token) (items []interface{}, remainingTokens []Token, err error) {
	subitem, subtokens, suberr := t.Expr.Parse(g, remainingTokens)
	if suberr != nil {
		remainingTokens = tokens
		return
	}
	items = append(items, subitem)
	remainingTokens = subtokens
	return
}

func (t RuleTerm) Parse(g Grammar, tokens []Token) (items []interface{}, remainingTokens []Token, err error) {
	r, ok := g.Rule(t.Name)
	if !ok {
		err  = fmt.Errorf("Unknown rule name: %q.", t.Name)
		return
	}

	var subitems []interface{}
	subitems, remainingTokens, err = r.Parse(g, tokens)

	if err != nil {
		items = []interface{}{subitems}
	}
	return
}

func (t InlineRuleTerm) Parse(g Grammar, tokens []Token) (items []interface{}, remainingTokens []Token, err error) {
	if r, ok := g.Rule(t.Name); ok {
		items, remainingTokens, err = r.Parse(g, tokens)
		return
	}

	if s, ok := g.Symbol(t.Name); ok {
		re := regexp.MustCompile(s.Pattern)
		matches := re.FindStringSubmatch(tokens[0].Text)
		if len(matches) == 0 {
			err = errors.New("Could not match symbol.")
			return
		}
		text := matches[0]
		if len(matches) > 0 {
			text = matches[1]
		}
		st := SymbolText{
			Type: t.Name,
			Text: text,
		}
		items = []interface{}{st}
		remainingTokens = tokens[1:]
		return
	}
		
	err  = fmt.Errorf("Unknown rule name: %q.", t.Name)

	return
}

func (t TagTerm) Parse(g Grammar, tokens []Token) (items []interface{}, remainingTokens []Token, err error) {
	items = []interface{}{Tag(t.Tag)}
	remainingTokens = tokens
	return
}

func (t LiteralTerm) Parse(g Grammar, tokens []Token) (items []interface{}, remainingTokens []Token, err error) {
	if tokens[0].Type != "RAW" || tokens[0].Text != t.Literal {
		err = errors.New("Incorrect literal.")
	}
	items = []interface{}{Literal(tokens[0].Text)}
	remainingTokens = tokens[1:]
	return
}
