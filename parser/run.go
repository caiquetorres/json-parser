package jsonparser

import (
	"fmt"
	"io"
)

func Run(r io.ReadSeeker) error {
	ts := newTokenStream(r)
	for {
		tok, err := ts.next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		fmt.Println(tok.string())
	}
	return nil
}
