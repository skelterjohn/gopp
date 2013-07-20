package gopp

type Document struct {
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
}