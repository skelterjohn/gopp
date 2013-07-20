package gopp

type Grammar struct {
	Rules []Rule
	Symbols []Symbol
}

type Rule struct {
	Name string
	Expr Expr
}

type Symbol struct {
	Name string
	Pattern string
}

type Expr struct {
	Terms []Term
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
	Rules: []Rule{
		Rule{
			Name: "Grammar",
			Expr: Expr{
				Terms: []Term{
					Term{
						Operator: "=",
						Field: "Rules",
						Term: &Term{
							Operator: "+",
							Term: &Term{
								Operator: "<<",
								Name: "Rule",
							},
						},
					},
					Term{
						Operator: "=",
						Field: "Symbols",
						Term: &Term{
							Operator: "+",
							Term: &Term{
								Operator: "<<",
								Name: "Symbol",
							},
						},
					},
				},
			},
		},
	},
}