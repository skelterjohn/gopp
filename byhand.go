package gopp

var ByHandGrammar = Grammar{
	Rules: []*Rule{
		&Rule{ // Grammar => Rules=<<Rule>>+ Symbols=<<Symbol>>+
			Name: "Grammar",
			Expr: &Expr{ // Rules=<<Rule>>+ Symbols=<<Symbol>>+
				Terms: []*Term{
					&Term{ // Rules=<<Rule>>+
						Operator: "=",
						Field:    "Rules",
						Term: &Term{ // <<Rule>>+
							Operator: "+",
							Term: &Term{ // <<Rule>>
								Operator: "<<",
								Name:     "Rule",
							},
						},
					},
					&Term{ // Symbols=<<Symbol>>+
						Operator: "=",
						Field:    "Symbols",
						Term: &Term{ // <<Symbol>>+
							Operator: "+",
							Term: &Term{ // <<Symbol>>
								Operator: "<<",
								Name:     "Symbol",
							},
						},
					},
				},
			},
		},
		&Rule{ // Rule => Name=<identifier> '=>' Expr=<<Expr>> '\n'+
			Name: "Rule",
			Expr: &Expr{ // Name=<identifier> '=>' Expr=<<Expr>> '\n'+
				Terms: []*Term{
					&Term{ // Name=<identifier>
						Operator: "=",
						Field:    "Name",
						Term: &Term{ // <identifier>
							Operator: "<",
							Name:     "identifier",
						},
					},
					&Term{ // '=>'
						Literal: "=>",
					},
					&Term{ // Expr=<<Expr>>
						Operator: "=",
						Field:    "Expr",
						Term: &Term{ // <<Expr>>
							Operator: "<<",
							Name:     "Expr",
						},
					},
					&Term{ // '\n'+
						Operator: "+",
						Term: &Term{ // '\n'
							Literal: "\n",
						},
					},
				},
			},
		},
		&Rule{ // Symbol => Name=<identifier> '=' Pattern=<regexp> '\n'+
			Name: "Symbol",
			Expr: &Expr{ // Name=<identifier> '=' Pattern=<regexp> '\n'+
				Terms: []*Term{
					&Term{ // Name=<identifier>
						Operator: "=",
						Field:    "Name",
						Term: &Term{ // <identifier>
							Operator: "<",
							Name:     "identifier",
						},
					},
					&Term{ // '='
						Literal: "=",
					},
					&Term{ // Pattern=<regexp>
						Operator: "=",
						Field:    "Pattern",
						Term: &Term{ // <regexp>
							Operator: "<",
							Name:     "regexp",
						},
					},
					&Term{ // '\n'+
						Operator: "+",
						Term: &Term{ // '\n'
							Literal: "\n",
						},
					},
				},
			},
		},
		&Rule{ // Expr => Terms=<<Term>>+
			Name: "Expr",
			Expr: &Expr{
				Terms: []*Term{
					&Term{ // Terms=<<Term>>+
						Operator: "=",
						Field:    "Terms",
						Term: &Term{ // <<Term>>+
							Operator: "+",
							Term: &Term{ // <<Term>>
								Operator: "<<",
								Name:     "Term",
							},
						},
					},
				},
			},
		},
		&Rule{ // Term => Term=<<Term>> Operator='*'
			Name: "Term",
			Expr: &Expr{ // Term=<<Term>> Operator='*'
				Terms: []*Term{
					&Term{ // Term=<<Term>>
						Operator: "=",
						Field:    "Term",
						Term: &Term{ // <<Term>>
							Operator: "<<",
							Name:     "Term",
						},
					},
					&Term{ // Operator='*'
						Operator: "=",
						Field:    "Operator",
						Term: &Term{ // '*'
							Literal: "*",
						},
					},
				},
			},
		},
		&Rule{ // Term => Term=<<Term>> Operator='+'
			Name: "Term",
			Expr: &Expr{ // Term=<<Term>> Operator='+'
				Terms: []*Term{
					&Term{ // Term=<<Term>>
						Operator: "=",
						Field:    "Term",
						Term: &Term{ // <<Term>>
							Operator: "<<",
							Name:     "Term",
						},
					},
					&Term{ // Operator='+'
						Operator: "=",
						Field:    "Operator",
						Term: &Term{ // '+'
							Literal: "+",
						},
					},
				},
			},
		},
		&Rule{ // Term => Operator='[' Expr=<<Expr>> ']'
			Name: "Term",
			Expr: &Expr{ // Operator='[' Expr=<<Expr>> ']'
				Terms: []*Term{
					&Term{ // Operator='['
						Operator: "=",
						Field:    "Operator",
						Term: &Term{ // '['
							Literal: "[",
						},
					},
					&Term{ // Expr=<<Expr>>
						Operator: "=",
						Field:    "Expr",
						Term: &Term{ // <<Expr>>
							Operator: "<<",
							Name:     "Expr",
						},
					},
					&Term{ // ']'
						Literal: "]",
					},
				},
			},
		},
		&Rule{ // Term => Operator='(' Expr=<<Expr>> ')'
			Name: "Term",
			Expr: &Expr{ // Operator='(' Expr=<<Expr>> ')'
				Terms: []*Term{
					&Term{ // Operator='('
						Operator: "=",
						Field:    "Operator",
						Term: &Term{ // '('
							Literal: "(",
						},
					},
					&Term{ // Expr=<<Expr>>
						Operator: "=",
						Field:    "Expr",
						Term: &Term{ // <<Expr>>
							Operator: "<<",
							Name:     "Expr",
						},
					},
					&Term{ // ')'
						Literal: ")",
					},
				},
			},
		},
		&Rule{ // Term => Operator='<<' Name=<identifier> '>>'
			Name: "Term",
			Expr: &Expr{ // Operator='<<' Name=<identifier> '>>'
				Terms: []*Term{
					&Term{ // Operator='<<'
						Operator: "=",
						Field:    "Operator",
						Term: &Term{ // '<<'
							Literal: "<<",
						},
					},
					&Term{ // Name=<identifier>
						Operator: "=",
						Field:    "Name",
						Term: &Term{ // <identifier>
							Operator: "<",
							Name:     "identifier",
						},
					},
					&Term{ // '>>'
						Literal: ">>",
					},
				},
			},
		},
		&Rule{ // Term => Operator='<' Name=<identifier> '>'
			Name: "Term",
			Expr: &Expr{ // Operator='<' Name=<identifier> '>'
				Terms: []*Term{
					&Term{ // Operator='<'
						Operator: "=",
						Field:    "Operator",
						Term: &Term{ // '<'
							Literal: "<",
						},
					},
					&Term{ // Name=<identifier>
						Operator: "=",
						Field:    "Name",
						Term: &Term{ // <identifier>
							Operator: "<",
							Name:     "identifier",
						},
					},
					&Term{ // '>'
						Literal: ">",
					},
				},
			},
		},
		&Rule{ // Term => Field=<indentifier> Operator='=' Term=<<Term>>
			Name: "Term",
			Expr: &Expr{ // Field=<indentifier> Operator='=' Term=<<Term>>
				Terms: []*Term{
					&Term{ // Field=<indentifier>
						Operator: "=",
						Field:    "Field",
						Term: &Term{ // <indentifier>
							Operator: "<",
							Name:     "identifier",
						},
					},
					&Term{ // Operator='='
						Operator: "=",
						Field:    "Operator",
						Term: &Term{ // '='
							Literal: "=",
						},
					},
					&Term{ // Term=<<Term>>
						Operator: "=",
						Field:    "Term",
						Term: &Term{ // <<Term>>
							Operator: "<<",
							Name:     "Term",
						},
					},
				},
			},
		},
		&Rule{ // Term => Literal=<literal>
			Name: "Term",
			Expr: &Expr{ // Literal=<literal>
				Terms: []*Term{
					&Term{ // Literal=<literal>
						Operator: "=",
						Field:    "Literal",
						Term: &Term{ // <literal>
							Operator: "<",
							Name:     "literal",
						},
					},
				},
			},
		},
	},
	Symbols: []*Symbol{
		&Symbol{
			Name:    "identifier",
			Pattern: `([a-zA-Z][a-zA-Z0-9_]*)`,
		},
		&Symbol{
			Name:    "literal",
			Pattern: `'((?:[\\']|[^'])+?)'`,
		},
		&Symbol{
			Name:    "regexp",
			Pattern: `\/((?:\\/|[^\n])+?)\/`,
		},
	},
}
