package jsonparser

import (
	"fmt"
	"io"
)

type tokenKind byte

const (
	LeftBracket tokenKind = iota
	RightBracket
	LeftBrace
	RightBrace

	Comma
	Colon

	Null
	Bool
	Number
	String
)

func (k tokenKind) String() string {
	switch k {
	case LeftBracket:
		return "left bracket"
	case RightBracket:
		return "right bracket"
	case LeftBrace:
		return "left brace"
	case RightBrace:
		return "right brace"
	case Comma:
		return "comma"
	case Colon:
		return "colon"
	case Null:
		return "null"
	case Bool:
		return "bool"
	case String:
		return "string"
	case Number:
		return "number"
	default:
		return "unknown"
	}
}

type token struct {
	s span
	k tokenKind
}

func (t *token) string() string {
	return fmt.Sprintf("token { k: %s, span: { s: %v, l: %v } }", t.k.String(), t.s.s, t.s.l)
}

func (t *token) textContent(r io.ReadSeeker) (string, error) {
	return t.s.textContent(r)
}
