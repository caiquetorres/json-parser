package jsonparser

import (
	"errors"
	"io"
)

var errUnexpectedTok = errors.New("unexpected token")

type parseStream struct {
	ts *tokenStream
}

func newParseStream(r io.Reader) *parseStream {
	return &parseStream{
		ts: newTokenStream(r),
	}
}

func (p *parseStream) peek() (token, error) {
	return p.ts.peek()
}

func (p *parseStream) next() (token, error) {
	return p.ts.next()
}

func (p *parseStream) expect(kind tokenKind) (token, error) {
	tok, err := p.ts.next()
	if err != nil {
		return token{}, err
	}
	if tok.k != kind {
		return token{}, errUnexpectedTok
	}
	return tok, nil
}
