gopp
====

A GO Parser Parser.

Pronounced 'gahp', rather than 'go pee pee'.

gopp is a library that takes a grammar, specified in .gopp format, a document, and an object, parses the document using the grammar, and decodes the resulting tree into the provided object.

.gopp is a BNF-like format for describing context-free grammars.

This README does not attempt to describe the use and purpose of context-free grammars - see google for more information about grammars and recursive descent parsing. Or try http://blog.reverberate.org/2013/07/ll-and-lr-parsing-demystified.html.

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

Clearly, the object is a more reasonable representation than the AST for actually dealing with in your code, which is why the decoding step was created.

Grammar
-------

The following .gopp grammar describes .gopp grammars, and how to decode them into gopp.Grammar objects.

```
# The first things are lex steps, which are for use by the tokenizer. 
# Currently the only recognized lex step is stuff to ignore.

# We ignore comments,
ignore: /^#.*\n/
# and whitespace that preceeds something more interesting.
ignore: /^(?:[ \t])+/

# After the lex steps are the rules.
# The fact that Grammar is first is irrelevant. The name of the starting rule
# needs to be provided in code.
# A Grammar is made up of lists of LexSteps, Rules, and Symbols, in that order,
# and there may be zero LexSteps or Symbols. There must be at least one Rule.
Grammar => {type=Grammar} '\n'* {field=LexSteps} <<LexStep>>* {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*

# The next three rules define the major types of elements in a grammar.

# A LexStep is an identifier, a literal ':', and a regexp pattern. If the name
# is 'ignore', then when the lexer goes to get the next token, it will try to
# trim the remaining document using the provided pattern. No other names are
# used, currently.
LexStep => {field=Name} <identifier> ':' {field=Pattern} <regexp> '\n'+

# A Rule is an identifier, a literal '=>', an Expr, and ends with one or more
# newlines.
Rule => {field=Name} <identifier> '=>' {field=Expr} <Expr> '\n'+
# A Symbol is an identifier, a literal '=', a regexp, and ends with one or more
# newlines.
Symbol => {field=Name} <identifier> '=' {field=Pattern} <regexp> '\n'+

# An Expr is one or more Terms.
Expr => <<Term>>+

# A Term can be a Term1,
Term => <Term1>
# or a Term2.
Term => <Term2>

# A Term1 can be a Term2 followed by a literal '*',
Term1 => {type=RepeatZeroTerm} {field=Term} <<Term2>> '*'
# or a Term2 followd by a literal '+'.
Term1 => {type=RepeatOneTerm} {field=Term} <<Term2>> '+'

# A Term2 can be an Expr surrounded by '[' and ']',
Term2 => {type=OptionalTerm} '[' {field=Expr} <Expr> ']'
# or by '(' and ')',
Term2 => {type=GroupTerm} '(' {field=Expr} <Expr> ')'
# or an identifier surrounded by '<<' and '>>',
Term2 => {type=RuleTerm} '<<' {field=Name} <identifier> '>>'
# or by '<' and '>',
Term2 => {type=InlineRuleTerm} '<' {field=Name} <identifier> '>'
# or a tag,
Term2 => {type=TagTerm} {field=Tag} <tag>
# or a literal.
Term2 => {type=LiteralTerm} {field=Literal} <literal>

# And last is the symbols, which are regular expressions that can be found in
# the document. Their order is important - it indicates the order in which the
# tokenizer attempts to match them against the rest of the document. So, if two
# symbols could be used starting at the same point in the document, the one
# that is listed first will win.
identifier = /([a-zA-Z][a-zA-Z0-9_]*)/
literal = /'((?:[\\']|[^'])+?)'/
tag = /\{((?:[\\']|[^'])+?)\}/
regexp = /\/((?:\\/|[^\n])+?)\//

```


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
