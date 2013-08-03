gopp
====

A GO Parser Parser.

Pronounced 'gahp', rather than 'go pee pee'.

gopp is a library that takes a grammar, specified in .gopp format, a document, and an object, parses the document using the grammar, and decodes the resulting tree into the provided object.

.gopp is a BNF-like format for describing context-free grammars.

This README does not attempt to describe the use and purpose of context-free grammars - see google for more information about grammars and recursive descent parsing.

Example
-------

The following grammar can be used to parse simple arithmetic equations.

```
Eqn => {type=MathEqn} {field=Left} <<Expr>> '=' {field=Right} <<Expr>> '\n'
Expr => {type=MathSum} {field=First} <<Term>> '+' {field=Second} <<Term>>
Expr => <Term>
Term => {type=MathProduct} {field=First} <<Factor>> '*' {field=Second} <<Factor>>
Term => <Factor>
Factor => {type=MathExprFactor} '(' {field=Expr} <<Expr>> ')'
Factor => {type=MathNumberFactor} {field=Number} <number>
number = /(\d+)/
```

A grammar is made up of rules ```<<Name>>```, inline rules ```<Name>```, literals ```'string'```, and tags ```{tag}```.

When parsing a document, a rule creates a new subtree as a child of the current tree, and an inline rule creates a new tree and adds its children to the current tree (the difference between ```[1,2,3,[a,b,c]]``` and ```[1,2,3,a,b,c]```).

Literals are strings that must appear exactly in the document text. To have other kinds of text matched, a .gopp also defines a set of symbols using regular expressions, and they are brought into the main tree by using inline rules.

Tags are elements that are put into the AST if their rule can be parsed. They do not match anything in the actual document text, but they can be used to provide information about the tree structure. For things to be decoded into objects, the "type=" and "field=" tags are used. A "type=" tag tells the decoder what type to allocate in the case that the field or slice element being decoded into is an interface without concrete type. A "field=" tag tells the decoder that, if the current object is a struct, the subtree in the next element is decoded into the field with the given name. Tags can be anything, and can be seen if the AST is accessed directly before decoding.

The grammar above can be used to decode documents into objects of type MathEqn, with the following types defined.

```
type MathEqn struct {
	Left, Right interface{}
}

type MathSum struct {
	First, Second interface{}
}

type MathProduct struct {
	First, Second interface{}
}

type MathExprFactor struct {
	Expr
}

type MathNumberFactor struct {
	Number string
}
```

So, the document "5+1=6" would get the AST

```
AST{
	Tag("type=MathEqn"),
	Tag("field=Left"),
	[]Node{
		Tag("type=MathSum"),
		Tag("field=First"),
		[]Node{
			Tag("type=MathNumberFactor"),
			Tag("field=Number"),
			return SymbolText{
				Type: "number",
				Text: "5",
			}
		},
		Literal("+"),
		Tag("field=Second"),
		[]Node{
			Tag("type=MathNumberFactor"),
			Tag("field=Number"),
			return SymbolText{
				Type: "number",
				Text: "1",
			}
		},
	},
	Literal("=")
	Tag("field=Right"),
	[]Node{
		Tag("type=MathNumberFactor"),
		Tag("field=Number"),
		return SymbolText{
			Type: "number",
			Text: "6",
		}
	},
	Literal("\n"),
}
```

and the object

```
MathEqn{
	Left:MathSum{
		First:MathNumberFactor{"5"},
		Second:MathNumberFactor{"5"},
	},
	Right: MathNumberFactor{"6"},
}
```

Grammar
-------

The following .gopp grammar describes .gopp grammars, and how to decode them into gopp.Grammar objects.

```
Grammar => {type=Grammar} '\n'* {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*
Rule => {field=Name} <identifier> '=>' {field=Expr} <Expr> '\n'+
Symbol => {field=Name} <identifier> '=' {field=Pattern} <regexp> '\n'+
Expr => <<Term>>+
Term => <Term1>
Term => <Term2>
Term1 => {type=RepeatZeroTerm} {field=Term} <<Term2>> '*'
Term1 => {type=RepeatOneTerm} {field=Term} <<Term2>> '+'
Term2 => {type=OptionalTerm} '[' {field=Expr} <Expr> ']'
Term2 => {type=GroupTerm} '(' {field=Expr} <Expr> ')'
Term2 => {type=RuleTerm} '<<' {field=Name} <identifier> '>>'
Term2 => {type=InlineRuleTerm} '<' {field=Name} <identifier> '>'
Term2 => {type=TagTerm} {field=Tag} <tag>
Term2 => {type=LiteralTerm} {field=Literal} <literal>
identifier = /([a-zA-Z][a-zA-Z0-9_]*)/
literal = /'((?:[\\']|[^'])+?)'/
tag = /\{((?:[\\']|[^'])+?)\}/
regexp = /\/((?:\\/|[^\n])+?)\//
```

Parsing
-------

In english, from top to bottom,
```
Grammar => '\n'* {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*
```
A "Grammar" is made up of one or more "Rule"s followed by some "Symbol"s.
```
Rule => {field=Name} <identifier> '=>' {field=Expr} <Expr> '\n'+
```
A "Rule" is an identifier, followed by a literal '=>', followed by an Expr, and ends with a newline.
```
Symbol => {field=Name} <identifier> '=' {field=Pattern} <regexp> '\n'+
```
A "Symbol" is an identifier, followed by a literal '=', followed by a regexp, and ends with a newline.
```
Expr => <<Term>>+
```
An "Expr" is made up of one ore more "Term"s.
```
Term => <Term1>
```
A "Term" can be either a "Term1",
```
Term => <Term2>
```
or a "Term2".
```
Term1 => {type=RepeatZeroTerm} {field=Term} <<Term2>> '*'
```
A "Term1" can be a "Term2" followed by a literal '*',
```
Term1 => {type=RepeatOneTerm} {field=Term} <<Term2>> '+'
```
or a "Term2" followed by a literal '+'.
```
Term2 => {type=OptionalTerm} '[' {field=Expr} <Expr> ']'
```
A "Term2" can be a literal '[', then an "Expr", then a literal ']',
```
Term2 => {type=GroupTerm} '(' {field=Expr} <Expr> ')'
```
or '(', ')' instead of '[', ']',
```
Term2 => {type=RuleTerm} '<<' {field=Name} <identifier> '>>'
```
or a literal '<<', followed by an identifier, followed by a literal '>>',
```
Term2 => {type=InlineRuleTerm} '<' {field=Name} <identifier> '>'
```
or '<', '>', instead of '<<', '>>',
```
Term2 => {type=TagTerm} {field=Tag} <tag>
```
or a tag,
```
Term2 => {type=LiteralTerm} {field=Literal} <literal>
```
or a literal.
```
identifier = /([a-zA-Z][a-zA-Z0-9_]*)/
```
An identifier is a letter followed by some number or letters, digits or underscores.
```
literal = /'((?:[\\']|[^'])+?)'/
```
A literal is any string inside single quotes, provided that any single quotes in the string itself are properly escaped.
```
tag = /\{((?:[\\']|[^'])+?)\}/
```
A tag is a string surrounded by curly braces.
```
regexp = /\/((?:\\/|[^\n])+?)\//
```
A regexp is a string surrounded by forward slashes, provided that any forward slashes inside the string itself are properly escaped.

The ```<<X>>``` and ```<Y>``` indicate recursively evaluated rules and inline rules. A rule will create an AST subtree in its parent. An inline rule will expand its children into its parent, rather than creating a new subtree. In otherword, if the child evaluates to [1,2,3], if that child were from a rule, the parent that already had [a,b,c] would become [a,b,c,[1,2,3]] when adding that child. For an inline rule, that same parent becomes [a,b,c,1,2,3]. This inlining is useful for keeping trees compact and easy to work with.

Anything within a '[' and ']' is optional: if it cannot be parsed, the parent rule may still successfully parse without the optional component.

Anything within a '(' and ')' is grouped, and external operators (like '*' or '+', or a forthcoming '|') apply to the group as a whole.

The '*' and '+' operators indicate that the rule should be applied as many times as possible, with the '+' requiring at least one successful application for the '+' to succeed.

A tag inserts a gopp.Tag into the tree when evaluated, and is always evaluated successfully when reached. This element is useful for inserting information into the tree that can be looked at by a post-processor. gopp itself makes use of several tags to help it decode into objects, described in the decoding section.

Decoding
--------

A parsed tree is decoded into an object.

If that object is a slice and the tree is also a slice, each element of the tree-slice is decoded into a new element for the object slice.

If that object is a struct, a tag of the form "{field=X}" indicates that the subsequent tree element should be decoded into the object's .X field. As a special case, "{field=.}" will apply the subsequent tree element to the current object.

If a field or a slice element is an interface type, the tree needs to have a tag of the form "{type=T}", indicating that the type T should be used to allocate the element for decoding. T must have been registered before-hand.
