package gopp

/*
Grammar => {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>+

Rule => {field=Name} <identifier> '=>' {field=Expr} <<Expr>> '\n'+

Symbol => {field=Name} <identifier> '=' {field=Pattern} <regexp> '\n'+

Expr => {field=Terms} <<Term>>+

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

regexp = /\/((?:\\/|[^\n])+?)\//
*/

var ByHandGrammar = Grammar{
	Rules: []Rule{
		Rule{ // Grammar => {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*
			Name: "Grammar",
			Expr: Expr{ // {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*
				TagTerm{Tag: "field=Rules"},
				RepeatOneTerm{
					RuleTerm{Name: "Rule"},
				},
				TagTerm{Tag: "field=Symbols"},
				RepeatZeroTerm{
					RuleTerm{Name: "Symbol"},
				},
			},
		},
		Rule{ // Rule => {field=Name} <identifier> '=>' {field=Expr} <<Expr>> '\n'+
			Name: "Rule",
			Expr: Expr{
				TagTerm{Tag: "field=Name"},
				InlineRuleTerm{Name: "identifier"},
				LiteralTerm{Literal: "=>"},
				TagTerm{Tag: "field=Expr"},
				RuleTerm{Name: "Expr"},
				RepeatOneTerm{
					LiteralTerm{Literal: "\n"},
				},
			},
		},
		Rule{ // Symbol => {field=Name} <identifier> '=' {field=Pattern} <regexp> '\n'+
			Name: "Symbol",
			Expr: Expr{
				TagTerm{Tag: "field=Name"},
				InlineRuleTerm{Name: "identifier"},
				LiteralTerm{Literal: "="},
				TagTerm{Tag: "field=Pattern"},
				InlineRuleTerm{Name: "regexp"},
				RepeatOneTerm{
					LiteralTerm{Literal: "\n"},
				},
			},
		},
		Rule{ // Expr => <<Term>>+
			Name: "Expr",
			Expr: Expr{
				RepeatOneTerm{
					InlineRuleTerm{Name: "Term"},
				},
			},
		},
		Rule{ // Term => Term1
			Name: "Term",
			Expr: Expr{
				InlineRuleTerm{Name: "Term1"},
			},
		},
		Rule{ // Term => Term2
			Name: "Term",
			Expr: Expr{
				InlineRuleTerm{Name: "Term2"},
			},
		},
		Rule{ // Term => {type=OptionalTerm} '[' {field=Expr} <<Expr>> ']'
			Name: "Term2",
			Expr: Expr{
				TagTerm{Tag: "type=OptionalTerm"},
				LiteralTerm{Literal: "["},
				TagTerm{Tag: "field=Expr"},
				RuleTerm{Name: "Expr"},
				LiteralTerm{Literal: "]"},
			},
		},
		Rule{ // Term => {type=GroupTerm} '(' {field=Expr} <<Expr>> ')'
			Name: "Term2",
			Expr: Expr{
				TagTerm{Tag: "type=GroupTerm"},
				LiteralTerm{Literal: "("},
				TagTerm{Tag: "field=Expr"},
				RuleTerm{Name: "Expr"},
				LiteralTerm{Literal: ")"},
			},
		},
		Rule{ // Term => {type=RuleTerm} '<<' {field=Name} <identifier> '>>'
			Name: "Term2",
			Expr: Expr{
				TagTerm{Tag: "type=RuleTerm"},
				LiteralTerm{Literal: "<<"},
				TagTerm{Tag: "field=Name"},
				InlineRuleTerm{Name: "identifier"},
				LiteralTerm{Literal: ">>"},
			},
		},
		Rule{ // Term => {type=InlineRuleTerm} '<' {field=Name} <identifier> '>'
			Name: "Term2",
			Expr: Expr{
				TagTerm{Tag: "type=InlineRuleTerm"},
				LiteralTerm{Literal: "<"},
				TagTerm{Tag: "field=Name"},
				InlineRuleTerm{Name: "identifier"},
				LiteralTerm{Literal: ">"},
			},
		},
		Rule{ // Term => {type=TagTerm} {field=Tag} <tag>
			Name: "Term2",
			Expr: Expr{
				TagTerm{Tag: "type=TagTerm"},
				TagTerm{Tag: "field=Tag"},
				InlineRuleTerm{Name: "tag"},
			},
		},
		Rule{ // Term => {type=LiteralTerm} {field=Literal} <literal>
			Name: "Term2",
			Expr: Expr{
				TagTerm{Tag: "type=LiteralTerm"},
				TagTerm{Tag: "field=Literal"},
				InlineRuleTerm{Name: "literal"},
			},
		},
		Rule{ // Term => {type=RepeatZeroTerm} {field=Term} <<Term>> '*'
			Name: "Term1",
			Expr: Expr{
				TagTerm{Tag: "type=RepeatZeroTerm"},
				TagTerm{Tag: "field=Term"},
				RuleTerm{Name: "Term2"},
				LiteralTerm{Literal: "*"},
			},
		},
		Rule{ // Term => {type=RepeatOneTerm} {field=Term} <<Term>> '+'
			Name: "Term1",
			Expr: Expr{
				TagTerm{Tag: "type=RepeatOneTerm"},
				TagTerm{Tag: "field=Term"},
				RuleTerm{Name: "Term2"},
				LiteralTerm{Literal: "+"},
			},
		},
	},
	Symbols: []Symbol{
		Symbol{
			Name:    "identifier",
			Pattern: `([a-zA-Z][a-zA-Z0-9_]*)`,
		},
		Symbol{
			Name:    "literal",
			Pattern: `'((?:[\\']|[^'])+?)'`,
		},
		Symbol{
			Name:    "tag",
			Pattern: `\{((?:[\\']|[^'])+?)\}`,
		},
		Symbol{
			Name:    "regexp",
			Pattern: `\/((?:\\/|[^\n])+?)\/`,
		},
	},
}

func mki(text string) SymbolText {
	return SymbolText{
		Type: "identifier",
		Text: text,
	}
}

func mkr(text string) SymbolText {
	return SymbolText{
		Type: "regexp",
		Text: text,
	}
}

func mkt(text string) SymbolText {
	return SymbolText{
		Type: "tag",
		Text: text,
	}
}

func mkl(text string) SymbolText {
	return SymbolText{
		Type: "literal",
		Text: text,
	}
}

func mkGrammar(rules, symbols []Node) []Node {
	return []Node{
		Tag("field=Rules"),
		rules,
		Tag("field=Symbols"),
		symbols,
	}
}

func mkRule(name string, nodes ...Node) []Node {
	return []Node{
		Tag("field=Name"),
		mki(name),
		Literal("=>"),
		Tag("field=Expr"),
		mkExpr(nodes...),
		[]Node{
			Literal("\n"),
		},
	}
	return nodes
}

func mkSymbol(name, pattern string) []Node {
	return []Node{
		Tag("field=Name"),
		mki(name),
		Literal("="),
		Tag("field=Pattern"),
		mkr(pattern),
		[]Node{
			Literal("\n"),
		},
	}
}

func mkExpr(nodes ...Node) []Node {
	return nodes
}
func mkRepeatZeroTerm(node Node) []Node {
	return []Node{
		Tag("type=RepeatZeroTerm"),
		Tag("field=Term"),
		node,
		Literal("*"),
	}
}

func mkRepeatOneTerm(node Node) []Node {
	return []Node{
		Tag("type=RepeatOneTerm"),
		Tag("field=Term"),
		node,
		Literal("+"),
	}
}

func mkOptionalTerm(node Node) []Node {
	return []Node{
		Tag("type=OptionalTerm"),
		Literal("["),
		Tag("field=Expr"),
		node,
		Literal("]"),
	}
}

func mkRuleTerm(text string) []Node {
	return []Node{
		Tag("type=RuleTerm"),
		Literal("<<"),
		Tag("field=Name"),
		mki(text),
		Literal(">>"),
	}
}

func mkInlineRuleTerm(text string) []Node {
	return []Node{
		Tag("type=InlineRuleTerm"),
		Literal("<"),
		Tag("field=Name"),
		mki(text),
		Literal(">"),
	}
}

func mkTagTerm(text string) []Node {
	return []Node{
		Tag("type=TagTerm"),
		Tag("field=."),
		mkt(text),
	}
}

func mkLiteralTerm(text string) []Node {
	return []Node{
		Tag("type=LiteralTerm"),
		Tag("field=Literal"),
		mkl(text),
	}
}

var ByHandGoppAST = mkGrammar(
	[]Node{
		mkRule("Grammar",
			Tag("field=Rules"),
			mkRepeatOneTerm(
				mkRuleTerm("Rule"),
			),
			Tag("field=Symbols"),
			mkRepeatOneTerm(
				mkRuleTerm("Symbol"),
			),
		),
		mkRule("Rule",
			mkTagTerm("field=Name"),
			mkInlineRuleTerm("identifier"),
			mkLiteralTerm("=>"),
			mkTagTerm("field=Expr"),
			mkRuleTerm("Expr"),
			mkRepeatOneTerm(mkLiteralTerm("\n")),
		),
		mkRule("Symbol",
			mkTagTerm("field=Name"),
			mkInlineRuleTerm("identifier"),
			mkLiteralTerm("="),
			mkTagTerm("field=Pattern"),
			mkInlineRuleTerm("regexp"),
			mkRepeatOneTerm(mkLiteralTerm("\n")),
		),
		mkRule("Expr",
			mkTagTerm("field=."),
			mkRepeatOneTerm(mkRuleTerm("Term")),
		),
		mkRule("Term",
			mkInlineRuleTerm("Term1"),
		),
		mkRule("Term",
			mkInlineRuleTerm("Term2"),
		),
		mkRule("Term1",
			mkTagTerm("type=RepeatZeroTerm"),
			mkTagTerm("field=Term"),
			mkRuleTerm("Term2"),
			mkLiteralTerm("*"),
		),
		mkRule("Term1",
			mkTagTerm("type=RepeatOneTerm"),
			mkTagTerm("field=Term"),
			mkRuleTerm("Term2"),
			mkLiteralTerm("+"),
		),
		mkRule("Term2",
			mkTagTerm("type=OptionalTerm"),
			mkLiteralTerm("["),
			mkTagTerm("field=Expr"),
			mkRuleTerm("Expr"),
			mkLiteralTerm("]"),
		),
		mkRule("Term2",
			mkTagTerm("type=GroupTerm"),
			mkLiteralTerm("("),
			mkTagTerm("field=Expr"),
			mkRuleTerm("Expr"),
			mkLiteralTerm(")"),
		),
		mkRule("Term2",
			mkTagTerm("type=RuleTerm"),
			mkLiteralTerm("<<"),
			mkTagTerm("field=Name"),
			mkInlineRuleTerm("identifier"),
			mkLiteralTerm(">>"),
		),
		mkRule("Term2",
			mkTagTerm("type=InlineRuleTerm"),
			mkLiteralTerm("<"),
			mkTagTerm("field=Name"),
			mkInlineRuleTerm("identifier"),
			mkLiteralTerm(">"),
		),
		mkRule("Term2",
			mkTagTerm("type=TagTerm"),
			mkTagTerm("field=."),
			mkInlineRuleTerm("tag"),
		),
		mkRule("Term2",
			mkTagTerm("type=LiteralTerm"),
			mkTagTerm("field=."),
			mkInlineRuleTerm("literal"),
		),
	},
	[]Node{
		mkSymbol("identifier", `([a-zA-Z][a-zA-Z0-9_]*)`),
		mkSymbol("literal", `'((?:[\\']|[^'])+?)'`),
		mkSymbol("tag", `\{((?:[\\']|[^'])+?)\}`),
		mkSymbol("regexp", `\/((?:\\/|[^\n])+?)\/`),
	},
)
