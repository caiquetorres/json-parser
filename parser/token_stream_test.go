package jsonparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Add more tests

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
