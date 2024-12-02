package jsonparser

import (
	"bufio"
	"encoding/hex"
	"errors"
	"io"
	"unicode"
)

var errBad = errors.New("bad token")

type tokenStream struct {
	r   *bufio.Reader
	s   uint
	ptr uint
	tok token
	err error
}

func newTokenStream(r io.Reader) *tokenStream {
	t := &tokenStream{
		r:   bufio.NewReader(r),
		s:   0,
		ptr: 0,
	}
	t.tok, t.err = t.get()
	return t
}

func (t *tokenStream) peekByte() (byte, error) {
	data, err := t.r.Peek(1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (t *tokenStream) nextByte() (byte, error) {
	t.ptr++
	return t.r.ReadByte()
}

func (t *tokenStream) peek() (token, error) {
	return t.tok, t.err
}

func (t *tokenStream) next() (token, error) {
	tok, err := t.tok, t.err
	t.tok, t.err = t.get()
	return tok, err
}

func (t *tokenStream) get() (token, error) {
	// REVIEW: improve this function name
	t.skipWhitespace()
	t.s = t.ptr
	ch, err := t.nextByte()
	if err != nil {
		return token{}, err
	}
	switch ch {
	case '{':
		return t.newToken(LeftBracket), nil
	case '}':
		return t.newToken(RightBracket), nil
	case '[':
		return t.newToken(LeftBrace), nil
	case ']':
		return t.newToken(RightBrace), nil
	case ',':
		return t.newToken(Comma), nil
	case ':':
		return t.newToken(Colon), nil
	case '"':
		return t.tokString()
	}
	if unicode.IsLetter(rune(ch)) {
		return t.tokKeyword(ch)
	} else if ch == '-' || unicode.IsNumber(rune(ch)) {
		return t.tokNumber(ch)
	}
	return token{}, errBad
}

func (t *tokenStream) newToken(kind tokenKind) token {
	span := span{s: uint32(t.s), l: uint16(t.ptr - t.s)}
	return token{s: span, k: kind}
}

func (t *tokenStream) skipWhitespace() {
	for {
		ch, err := t.peekByte()
		if err != nil {
			break
		}
		if !unicode.IsSpace(rune(ch)) {
			break
		}
		t.nextByte()
	}
}

func (t *tokenStream) tokKeyword(firstCh byte) (token, error) {
	identifier := string(firstCh)
	for {
		ch, err := t.peekByte()
		if err != nil || !unicode.IsLetter(rune(ch)) {
			break
		}
		identifier += string(ch)
		t.nextByte()
	}
	switch identifier {
	case "true":
		return t.newToken(Bool), nil
	case "false":
		return t.newToken(Bool), nil
	case "null":
		return t.newToken(Null), nil
	}
	return token{}, errBad
}

func (t *tokenStream) tokString() (token, error) {
	for {
		ch, err := t.peekByte()
		if err != nil {
			return token{}, errBad
		}
		if ch == '"' {
			break
		}
		if ch == '\n' || ch == '\t' {
			return token{}, errBad
		}
		if ch == '\\' {
			t.nextByte() // '\'
			ch, err := t.peekByte()
			if err != nil {
				return token{}, errBad
			}
			if ch == 'u' {
				t.nextByte() // 'u'
				var h []byte
				for range 4 {
					ch, err := t.nextByte()
					if err != nil {
						return token{}, errBad
					}
					h = append(h, ch)
				}
				_, err := hex.DecodeString(string(h))
				if err != nil {
					return token{}, errBad
				}
			}
			if !isEscapingChar(ch) {
				return token{}, errBad
			}
		}
		t.nextByte()
	}
	_, err := t.nextByte()
	// Reached this point means that the current character is a double quote (")
	if err != nil {
		return token{}, errBad
	}
	return t.newToken(String), nil
}

func (t *tokenStream) tokNumber(firstCh byte) (token, error) {
	ch := firstCh
	if ch == '-' {
		var err error
		ch, err = t.nextByte() // '-'
		if err != nil {
			return token{}, nil
		}
	}
	if ch == '0' {
		ch, err := t.peekByte()
		if err != nil {
			return token{}, errBad
		}
		if ch != '.' {
			return token{}, errBad
		}
	}
	for {
		ch, err := t.peekByte()
		if err != nil || !unicode.IsDigit(rune(ch)) {
			break
		}
		t.nextByte() // digit
	}
	ch, err := t.peekByte()
	if err != nil {
		return t.newToken(Number), nil
	}
	if ch == '.' {
		t.nextByte() // '.'
		ch, err := t.peekByte()
		if err != nil || !unicode.IsNumber(rune(ch)) {
			return token{}, errBad
		}
		for {
			ch, err := t.peekByte()
			if err != nil || !unicode.IsDigit(rune(ch)) {
				break
			}
			t.nextByte() // digit
		}
	}
	ch, err = t.peekByte()
	if err != nil {
		return t.newToken(Number), nil
	}
	if ch == 'e' || ch == 'E' {
		t.nextByte() // 'e' or 'E'
		ch, err := t.peekByte()
		if err != nil {
			return token{}, errBad
		}
		if ch == '+' || ch == '-' {
			t.nextByte() // '+' or '-'
		}
		ch, err = t.peekByte()
		if err != nil {
			return token{}, errBad
		}
		if !unicode.IsNumber(rune(ch)) {
			return token{}, errBad
		}
		for {
			ch, err := t.peekByte()
			if err != nil || !unicode.IsDigit(rune(ch)) {
				break
			}
			t.nextByte() // digit
		}
	}
	return t.newToken(Number), nil
}

func isEscapingChar(ch byte) bool {
	escapingChars := map[byte]struct{}{
		'b':  {},
		'f':  {},
		'n':  {},
		'u':  {},
		'r':  {},
		't':  {},
		'"':  {},
		'\\': {},
		'/':  {},
	}
	_, ok := escapingChars[ch]
	return ok
}
