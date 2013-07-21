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
				TagTerm{Tag: "field=."},
				RepeatOneTerm{
					RuleTerm{Name: "Term"},
				},
			},
		},
		Rule{ // Term => {type=RepeatZeroTerm} {field=Term} <<Term>> '*'
			Name: "Term",
			Expr: Expr{
				TagTerm{Tag: "type=RepeatZeroTerm"},
				TagTerm{Tag: "field=Term"},
				RuleTerm{Name: "Term"},
				LiteralTerm{Literal: "*"},
			},
		},
		Rule{ // Term => {type=RepeatOneTerm} {field=Term} <<Term>> '+'
			Name: "Term",
			Expr: Expr{
				TagTerm{Tag: "type=RepeatOneTerm"},
				TagTerm{Tag: "field=Term"},
				RuleTerm{Name: "Term"},
				LiteralTerm{Literal: "+"},
			},
		},
		Rule{ // Term => {type=OptionalTerm} '[' {field=Expr} <<Expr>> ']'
			Name: "Term",
			Expr: Expr{
				TagTerm{Tag: "type=OptionalTerm"},
				LiteralTerm{Literal: "["},
				TagTerm{Tag: "field=Expr"},
				RuleTerm{Name: "Expr"},
				LiteralTerm{Literal: "]"},
			},
		},
		Rule{ // Term => {type=GroupTerm} '(' {field=Expr} <<Expr>> ')'
			Name: "Term",
			Expr: Expr{
				TagTerm{Tag: "type=GroupTerm"},
				LiteralTerm{Literal: "("},
				TagTerm{Tag: "field=Expr"},
				RuleTerm{Name: "Expr"},
				LiteralTerm{Literal: ")"},
			},
		},
		Rule{ // Term => {type=RuleTerm} '<<' {field=Name} <identifier> '>>'
			Name: "Term",
			Expr: Expr{
				TagTerm{Tag: "type=RuleTerm"},
				LiteralTerm{Literal: "<<"},
				TagTerm{Tag: "field=Name"},
				InlineRuleTerm{Name: "identifier"},
				LiteralTerm{Literal: ">>"},
			},
		},
		Rule{ // Term => {type=InlineRuleTerm} '<' {field=Name} <identifier> '>'
			Name: "Term",
			Expr: Expr{
				TagTerm{Tag: "type=InlineRuleTerm"},
				LiteralTerm{Literal: "<"},
				TagTerm{Tag: "field=Name"},
				InlineRuleTerm{Name: "identifier"},
				LiteralTerm{Literal: ">"},
			},
		},
		Rule{ // Term => {type=TagTerm} {field=Tag} <tag>
			Name: "Term",
			Expr: Expr{
				TagTerm{Tag: "type=TagTerm"},
				TagTerm{Tag: "field=Tag"},
				InlineRuleTerm{Name: "tag"},
			},
		},
		Rule{ // Term => {type=LiteralTerm} {field=Literal} <literal> 
			Name: "Term",
			Expr: Expr{
				TagTerm{Tag: "type=LiteralTerm"},
				TagTerm{Tag: "field=Literal"},
				InlineRuleTerm{Name: "literal"},
			},
		},
	},
	Symbols: []Symbol{
		Symbol{
			Name: "identifier",
			Pattern: `([a-zA-Z][a-zA-Z0-9_]*)`,
		},
		Symbol{
			Name: "literal",
			Pattern: `'((?:[\\']|[^'])+?)'`,
		},
		Symbol{
			Name: "tag",
			Pattern: `\{((?:[\\']|[^'])+?)\}`,
		},
		Symbol{
			Name: "regexp",
			Pattern: `\/((?:\\/|[^\n])+?)\/`,
		},
	},
}
