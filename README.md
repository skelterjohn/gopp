gopp
====

A go parser parser.

gopp is a library that takes a grammar, specified in .gopp format, a document, and an object, parses the document using the grammar, and decodes the resulting tree into the provided object.

.gopp is a BNF-like format for describing context-free grammars.

The following .gopp grammar describes .gopp grammars, and how to decode them into gopp.Grammar objects.

```
Grammar => '\n'* {field=Rules} <<Rule>>+ {field=Symbols} <<Symbol>>*
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
