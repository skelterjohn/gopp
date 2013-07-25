package gopp

import (
	"errors"
	"fmt"
)

const debug = false

var indentCount int = 0

func indentIn() {
	for i := 0; i < indentCount; i++ {
		fmt.Print(" ")
	}
	indentCount++
}
func indentOut() {
	indentCount--
	for i := 0; i < indentCount; i++ {
		fmt.Print(" ")
	}
}

func (r Rule) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
	if debug {
		indentIn()
		fmt.Printf("Rule(%q): %v\n", r.Name, tokens)
		defer func() {
			indentOut()
			fmt.Printf("-Rule(%q): ", r.Name)
			if err == nil {
				fmt.Println(items)
			} else {
				fmt.Println(err)
			}
		}()
	}
	return r.Expr.Parse(g, tokens)
}

func (e Expr) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
	if debug {
		indentIn()
		fmt.Printf("Expr: %v\n", tokens)
		defer func() {
			indentOut()
			fmt.Printf("-Expr: ")
			if err == nil {
				fmt.Println(items)
			} else {
				fmt.Println(err)
			}
		}()
	}
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
	if debug {
		indentIn()
		fmt.Printf("RepeatZeroTerm: %v\n", tokens)
		defer func() {
			indentOut()
			fmt.Printf("-RepeatZeroTerm: ")
			if err == nil {
				fmt.Println(items)
			} else {
				fmt.Println(err)
			}
		}()
	}
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

func (t RepeatOneTerm) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
	if debug {
		indentIn()
		fmt.Printf("RepeatOneTerm: %v\n", tokens)
		defer func() {
			indentOut()
			fmt.Printf("-RepeatOneTerm: ")
			if err == nil {
				fmt.Println(items)
			} else {
				fmt.Println(err)
			}
		}()
	}
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

func (t OptionalTerm) Parse(g Grammar, tokens []Token) (items []Node, remainingTokens []Token, err error) {
	if debug {
		indentIn()
		fmt.Printf("OptionalTerm: %v\n", tokens)
		defer func() {
			indentOut()
			fmt.Printf("-OptionalTerm: ")
			if err == nil {
				fmt.Println(items)
			} else {
				fmt.Println(err)
			}
		}()
	}
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
	if debug {
		indentIn()
		fmt.Printf("RuleTerm(%q): %v\n", t.Name, tokens)
		defer func() {
			indentOut()
			fmt.Printf("-RuleTerm(%q): ", t.Name)
			if err == nil {
				fmt.Println(items)
			} else {
				fmt.Println(err)
			}
		}()
	}
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
	if debug {
		indentIn()
		fmt.Printf("InlineRuleTerm(%q): %v\n", t.Name, tokens)
		defer func() {
			indentOut()
			fmt.Printf("-InlineRuleTerm(%q): ", t.Name)
			if err == nil {
				fmt.Println(items)
			} else {
				fmt.Println(err)
			}
		}()
	}
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
	if debug {
		indentIn()
		fmt.Printf("LiteralTerm(%q): %v\n", t.Literal, tokens)
		defer func() {
			indentOut()
			fmt.Printf("-LiteralTerm(%q): ", t.Literal)
			if err == nil {
				fmt.Println(items)
			} else {
				fmt.Println(err)
			}
		}()
	}
	if len(tokens) == 0 {
		err = errors.New("Not enough tokens.")
		return
	}
	if tokens[0].Type != "RAW" {
		err = errors.New("Incorrect literal.")
		return
	}
	if tokens[0].Text != t.Literal {
		err = errors.New("Incorrect literal.")
		return
	}
	items = []Node{Literal(tokens[0].Text)}
	remainingTokens = tokens[1:]
	return
}
