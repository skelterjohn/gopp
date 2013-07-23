package gopp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Token struct {
	Type string
	Raw  string
	Text string
}

func (t Token) String() string {
	return fmt.Sprintf("(%s: %q)", t.Type, t.Text)
}

func Tokenize(res []TypedRegexp, r io.Reader) (tokens []Token) {

	scanner := bufio.NewScanner(r)
	var buf bytes.Buffer
	// true if the last iteration through the loop did not result in a match (more data needed)
	noMatch := true
	eof := false
	for {
		if noMatch {
			if scanner.Scan() {
				// TODO: panic? error?
				// add a chunk
				buf.Write(scanner.Bytes())
				buf.WriteString("\n")
			} else {
				eof = true
			}
		}
		noMatch = true
		// try to match one of our regexps to the current buffer
		for _, re := range res {
			matches := re.FindSubmatch(buf.Bytes())
			if len(matches) == 0 {
				continue
			}

			token := Token{
				Type: re.Type,
				Raw:  string(matches[0]),
			}
			if len(matches) > 1 {
				token.Text = string(matches[1])
			}
			buf.Read(matches[0])
			tokens = append(tokens, token)
			noMatch = false
			break
		}

		if noMatch && eof {
			// TODO: panic?
			break
		}
		if eof && buf.Len() == 0 {
			break
		}
	}
	return
}
