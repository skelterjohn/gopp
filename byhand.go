package gopp

var ByHandGrammar = Grammar{
	LexSteps: []LexStep{
		LexStep{
			Name:    "ignore",
			Pattern: `^#.*\n`,
		},
		LexStep{
			Name:    "ignore",
			Pattern: `^(?:[ \t])+`,
		},
	},
	Rules: []Rule{
		Rule{ // Grammar => {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*
			Name: "Grammar",
			Expr: Expr{ // '\n'* {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*
				TagTerm{Tag: "type=Grammar"},
				RepeatZeroTerm{
					LiteralTerm{Literal: "\n"},
				},
				TagTerm{Tag: "field=LexSteps"},
				RepeatZeroTerm{
					RuleTerm{Name: "LexStep"},
				},
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
		Rule{ // Symbol => {field=Name} <identifier> '=' {field=Pattern} <regexp> '\n'+
			Name: "LexStep",
			Expr: Expr{
				TagTerm{Tag: "field=Name"},
				InlineRuleTerm{Name: "identifier"},
				LiteralTerm{Literal: ":"},
				TagTerm{Tag: "field=Pattern"},
				InlineRuleTerm{Name: "regexp"},
				RepeatOneTerm{
					LiteralTerm{Literal: "\n"},
				},
			},
		},
		Rule{ // Rule => {field=Name} <identifier> '=>' {field=Expr} <Expr> '\n'+
			Name: "Rule",
			Expr: Expr{
				TagTerm{Tag: "field=Name"},
				InlineRuleTerm{Name: "identifier"},
				LiteralTerm{Literal: "=>"},
				TagTerm{Tag: "field=Expr"},
				InlineRuleTerm{Name: "Expr"},
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
					RuleTerm{Name: "Term"},
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
		Rule{ // Term => {type=OptionalTerm} '[' {field=Expr} <<Expr>> ']'
			Name: "Term2",
			Expr: Expr{
				TagTerm{Tag: "type=OptionalTerm"},
				LiteralTerm{Literal: "["},
				TagTerm{Tag: "field=Expr"},
				InlineRuleTerm{Name: "Expr"},
				LiteralTerm{Literal: "]"},
			},
		},
		Rule{ // Term => {type=GroupTerm} '(' {field=Expr} <<Expr>> ')'
			Name: "Term2",
			Expr: Expr{
				TagTerm{Tag: "type=GroupTerm"},
				LiteralTerm{Literal: "("},
				TagTerm{Tag: "field=Expr"},
				InlineRuleTerm{Name: "Expr"},
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
		Text: escapeString(text),
	}
}

func mkGrammar(lexsteps, rules, symbols []Node) AST {
	return []Node{
		Tag("type=Grammar"),
		[]Node{
			Literal("\n"),
		},
		Tag("field=LexSteps"),
		lexsteps,
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

func mkLexStep(name, pattern string) []Node {
	return []Node{
		Tag("field=Name"),
		mki(name),
		Literal(":"),
		Tag("field=Pattern"),
		mkr(pattern),
		[]Node{
			Literal("\n"),
		},
	}
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
		mkExpr(node),
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
		Tag("field=Tag"),
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
		mkLexStep("ignore", `^#.*\n`),
		mkLexStep("ignore", `^(?:[ \t])+`),
	},
	[]Node{
		mkRule("Grammar",
			mkTagTerm("type=Grammar"),
			mkRepeatZeroTerm(mkLiteralTerm("\n")),
			mkTagTerm("field=LexSteps"),
			mkRepeatZeroTerm(
				mkRuleTerm("LexStep"),
			),
			mkTagTerm("field=Rules"),
			mkRepeatOneTerm(
				mkRuleTerm("Rule"),
			),
			mkTagTerm("field=Symbols"),
			mkRepeatZeroTerm(
				mkRuleTerm("Symbol"),
			),
		),
		mkRule("LexStep",
			mkTagTerm("field=Name"),
			mkInlineRuleTerm("identifier"),
			mkLiteralTerm(":"),
			mkTagTerm("field=Pattern"),
			mkInlineRuleTerm("regexp"),
			mkRepeatOneTerm(mkLiteralTerm("\n")),
		),
		mkRule("Rule",
			mkTagTerm("field=Name"),
			mkInlineRuleTerm("identifier"),
			mkLiteralTerm("=>"),
			mkTagTerm("field=Expr"),
			mkInlineRuleTerm("Expr"),
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
			mkInlineRuleTerm("Expr"),
			mkLiteralTerm("]"),
		),
		mkRule("Term2",
			mkTagTerm("type=GroupTerm"),
			mkLiteralTerm("("),
			mkTagTerm("field=Expr"),
			mkInlineRuleTerm("Expr"),
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
			mkTagTerm("field=Tag"),
			mkInlineRuleTerm("tag"),
		),
		mkRule("Term2",
			mkTagTerm("type=LiteralTerm"),
			mkTagTerm("field=Literal"),
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
