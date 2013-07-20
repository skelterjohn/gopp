package gopp

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
)

func Tokenize(res []*regexp.Regexp, r io.Reader) (tokens <-chan string) {
	ch := make(chan string)
	tokens = ch
	go func() {
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

				token := matches[0]
				if len(matches) > 1 {
					token = matches[1]
				}
				buf.Read(matches[0])
				ch <- string(token)
				noMatch = false
				break
			}

			if noMatch && eof {
				// TODO: panic?
				close(ch)
				break
			}
			if eof && buf.Len() == 0 {
				close(ch)
				break
			}
		}
	}()
	return tokens
}
