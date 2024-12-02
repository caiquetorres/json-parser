package jsonparser

import (
	"fmt"
	"io"
)

func Run(r io.ReadSeeker) error {
	ts := newTokenStream(r)
	var toks []token
	for {
		tok, err := ts.next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		toks = append(toks, tok)
	}
	for _, tok := range toks {
		txt, err := tok.textContent(r)
		if err != nil {
			return err
		}
		fmt.Print(txt)
	}
	fmt.Println()
	return nil
}
