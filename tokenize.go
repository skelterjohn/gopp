// Copyright 2013 The gopp AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopp

import (
	"fmt"
	"regexp"
)

type Token struct {
	Type     string
	Raw      string
	Text     string
	Row, Col int
}

func (t Token) String() string {
	return fmt.Sprintf("(%s: %q)", t.Type, t.Text)
}

type TokenizeInfo struct {
	TokenREs  []TypedRegexp
	IgnoreREs []*regexp.Regexp
}

func Tokenize(ti TokenizeInfo, document []byte) (tokens []Token, err error) {
	var row, col int
tokenloop:
	for len(document) != 0 {

		snippet := document
		if len(snippet) > 20 {
			snippet = snippet[:20]
		}

		// If something to ignore, trim it off.
		for _, re := range ti.IgnoreREs {
			matches := re.FindSubmatch(document)
			if len(matches) == 0 {
				continue
			}
			document = document[len(matches[0]):]
			continue tokenloop
		}

		var newdocument []byte
		for _, re := range ti.TokenREs {

			matches := re.FindSubmatch(document)
			if len(matches) == 0 {
				continue
			}

			matchedText := matches[0]
			capturedText := matches[1]

			token := Token{
				Type: re.Type,
				Raw:  string(matchedText),
				Row:  row,
				Col:  col,
			}
			if len(matches) > 1 {
				token.Text = string(capturedText)
				if err != nil {
					return
				}
			}
			for _, c := range matchedText {
				if c == '\n' {
					row++
					col = 0
				} else {
					col++
				}
			}
			newdocument = document[len(matchedText):]
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
