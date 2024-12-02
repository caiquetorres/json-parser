package jsonparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenStream_BasicTokens(t *testing.T) {
	input := `{"key":"value"}`
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

func TestTokenStream_TokString_ValidStrings(t *testing.T) {
	tests := []string{
		`"hello"`,
		`"hello world"`,
		`"with \\ backslash"`,
		`"with \" escaped quotes"`,
		`"new\nline"`,
		`"unicode \u1234 sequence"`,
		`""`,
	}
	for _, input := range tests {
		r := strings.NewReader(input)
		ts := newTokenStream(r)
		tok, err := ts.next()
		assert.NoError(t, err)
		assert.Equal(t, String, tok.k)
	}
}

func TestTokenStream_TokString_InvalidStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{`"unterminated`, errBad},
		{`"invalid escape sequence \q"`, errBad},
		{`"incomplete unicode \u12"`, errBad},
		{`"incomplete escape \`, errBad},
	}
	for _, test := range tests {
		r := strings.NewReader(test.input)
		ts := newTokenStream(r)
		_, err := ts.next()
		assert.Error(t, err)
		assert.Equal(t, test.expected, err)
	}
}

func TestTokenStream_TokString_EdgeCases(t *testing.T) {
	tests := []string{
		`" "`,
		`"special chars !@#$%^&*"`,
		`"123"`,
		`"true"`,
	}
	for _, input := range tests {
		r := strings.NewReader(input)
		ts := newTokenStream(r)
		tok, err := ts.next()
		assert.NoError(t, err)
		assert.Equal(t, String, tok.k)
	}
}

func TestTokenStream_TokNumber_ValidNumbers(t *testing.T) {
	tests := []string{
		"0",
		"123",
		"-123",
		"123.456",
		"-123.456",
		"1e3",
		"1.23e3",
		"1.234e-5",
		"1.234E+56",
	}
	for _, test := range tests {
		r := strings.NewReader(test)
		ts := newTokenStream(r)
		tok, err := ts.next()
		assert.NoError(t, err)
		assert.Equal(t, Number, tok.k)
	}
}

func TestTokenStream_TokNumber_InvalidNumbers(t *testing.T) {
	tests := []string{
		"0123",
		"123.",
		".123",
		"1e",
		"1e+abc",
	}
	for _, test := range tests {
		r := strings.NewReader(test)
		ts := newTokenStream(r)
		_, err := ts.next()
		assert.Error(t, err)
	}
}
