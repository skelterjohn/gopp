package gopp

import (
	"fmt"
)

type Token struct {
	Type string
	Raw  string
	Text string
}

func (t Token) String() string {
	return fmt.Sprintf("(%s: %q)", t.Type, t.Text)
}

func Tokenize(res []TypedRegexp, document []byte) (tokens []Token, err error) {
	for len(document) != 0 {
		var newdocument []byte
		for _, re := range res {
			matches := re.FindSubmatch(document)
			if len(matches) == 0 {
				continue
			}

			token := Token{
				Type: re.Type,
				Raw:  string(matches[0]),
			}
			if len(matches) > 1 {
				token.Text = string(matches[1])
				if err != nil {
					return
				}
			}
			newdocument = document[len(matches[0]):]
			tokens = append(tokens, token)
			break
		}
		if newdocument == nil {
			snippet := document
			if len(snippet) > 80 {
				snippet = snippet[:80]
			}
			err = fmt.Errorf("Could not match starting from %q.", snippet)
			return
		}
		document = newdocument
	}
	return
}
