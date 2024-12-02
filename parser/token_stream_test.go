package jsonparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenStream_BasicTokens(t *testing.T) {
	input := `{ "key": "value" }`
	r := strings.NewReader(input)
	ts := newTokenStream(r)

	tok, err := ts.next()
	assert.NoError(t, err)
	assert.Equal(t, LeftBracket, tok.k)

	tok, err = ts.next()
	assert.NoError(t, err)
	assert.Equal(t, String, tok.k)

	tok, err = ts.next()
	assert.NoError(t, err)
	assert.Equal(t, Colon, tok.k)

	tok, err = ts.next()
	assert.NoError(t, err)
	assert.Equal(t, String, tok.k)

	tok, err = ts.next()
	assert.NoError(t, err)
	assert.Equal(t, RightBracket, tok.k)
}

func TestTokenStream_WhitespaceHandling(t *testing.T) {
	input := `  [`
	r := strings.NewReader(input)
	ts := newTokenStream(r)

	tok, err := ts.next()
	assert.NoError(t, err)
	assert.Equal(t, LeftBrace, tok.k)
}

func TestTokenStream_InvalidCharacter(t *testing.T) {
	input := `#`
	r := strings.NewReader(input)
	ts := newTokenStream(r)

	_, err := ts.next()
	assert.Error(t, err)
	assert.Equal(t, errBad, err)
}

func TestTokenStream_Numbers(t *testing.T) {
	input := "123"
	r := strings.NewReader(input)
	ts := newTokenStream(r)
	tok, err := ts.next()

	assert.NoError(t, err)
	assert.Equal(t, Number, tok.k)

	input = "-123"
	r = strings.NewReader(input)
	ts = newTokenStream(r)
	tok, err = ts.next()
	assert.NoError(t, err)
	assert.Equal(t, Number, tok.k)

	input = "123.456"
	r = strings.NewReader(input)
	ts = newTokenStream(r)
	tok, err = ts.next()
	assert.NoError(t, err)
	assert.Equal(t, Number, tok.k)

	input = "-123.456"
	r = strings.NewReader(input)
	ts = newTokenStream(r)
	tok, err = ts.next()
	assert.NoError(t, err)
	assert.Equal(t, Number, tok.k)

	input = "1e3"
	r = strings.NewReader(input)
	ts = newTokenStream(r)
	tok, err = ts.next()
	assert.NoError(t, err)
	assert.Equal(t, Number, tok.k)

	input = "1.23e3"
	r = strings.NewReader(input)
	ts = newTokenStream(r)
	tok, err = ts.next()
	assert.NoError(t, err)
	assert.Equal(t, Number, tok.k)

	input = "1.234e-5"
	r = strings.NewReader(input)
	ts = newTokenStream(r)
	tok, err = ts.next()
	assert.NoError(t, err)
	assert.Equal(t, Number, tok.k)

	input = "123."
	r = strings.NewReader(input)
	ts = newTokenStream(r)
	tok, err = ts.next()
	assert.Error(t, err)
	assert.Equal(t, errBad, err)

	input = ".123"
	r = strings.NewReader(input)
	ts = newTokenStream(r)
	tok, err = ts.next()
	assert.Error(t, err)
	assert.Equal(t, errBad, err)

	input = "1e"
	r = strings.NewReader(input)
	ts = newTokenStream(r)
	tok, err = ts.next()
	assert.Error(t, err)
	assert.Equal(t, errBad, err)

	input = "1e+abc"
	r = strings.NewReader(input)
	ts = newTokenStream(r)
	tok, err = ts.next()
	assert.Error(t, err)
	assert.Equal(t, errBad, err)
}
