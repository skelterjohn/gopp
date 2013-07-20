package gopp

type Grammar struct {
	Rules []*Rule
	Symbols []*Symbol
}

type Rule struct {
	Name string
	Expr *Expr
}

type Symbol struct {
	Name string
	Pattern string
}

type Expr struct {
	Terms []*Term
}

type Term struct {
	Operator string
	Term *Term
	Expr *Expr
	Field string
	Name string
	Literal string
}

var ByHandGrammar = Grammar{
	Rules: []*Rule{
		&Rule{ // Grammar => Rules=<<Rule>>+ Symbols=<<Symbol>>+
			Name: "Grammar",
			Expr: &Expr{ // Rules=<<Rule>>+ Symbols=<<Symbol>>+
				Terms: []*Term{
					&Term{ // Rules=<<Rule>>+
						Operator: "=",
						Field: "Rules",
						Term: &Term{ // <<Rule>>+
							Operator: "+",
							Term: &Term{ // <<Rule>>
								Operator: "<<",
								Name: "Rule",
							},
						},
					},
					&Term{ // Symbols=<<Symbol>>+
						Operator: "=",
						Field: "Symbols",
						Term: &Term{ // <<Symbol>>+
							Operator: "+",
							Term: &Term{ // <<Symbol>>
								Operator: "<<",
								Name: "Symbol",
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
						Field: "Name",
						Term: &Term { // <identifier>
							Operator: "<",
							Name: "identifier",
						},
					},
					&Term{ // '=>'
						Literal: "=>",
					},
					&Term{ // Expr=<<Expr>>
						Operator: "=",
						Field: "Expr",
						Term: &Term { // <<Expr>>
							Operator: "<<",
							Name: "Expr",
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
						Field: "Name",
						Term: &Term { // <identifier>
							Operator: "<",
							Name: "identifier",
						},
					},
					&Term{ // '='
						Literal: "=",
					},
					&Term{ // Pattern=<regexp>
						Operator: "=",
						Field: "Pattern",
						Term: &Term { // <regexp>
							Operator: "<",
							Name: "regexp",
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
						Field: "Terms",
						Term: &Term{ // <<Term>>+
							Operator: "+",
							Term: &Term{ // <<Term>>
								Operator: "<<",
								Name: "Term",
							},
						},
					},
				},
			},
		},
	},
}