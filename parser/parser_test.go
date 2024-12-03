package jsonparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_ValidJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"Empty Object", `{}`},
		{"Empty Array", `[]`},
		{"Simple Object", `{"key": "value"}`},
		{"Nested Object", `{"key": {"nestedKey": "nestedValue"}}`},
		{"Array of Objects", `[{"key": "value"}, {"key2": "value2"}]`},
		{"Mixed Array", `[{"key": "value"}, "string", 123, null, true, false]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.json)
			err := Parse(reader)
			assert.NoError(t, err, "Parse() failed for valid JSON: %s", tt.json)
		})
	}
}

func TestParse_ComplexJSON(t *testing.T) {
	json := `
	{
		"object": {
			"array": [1, 2, 3],
			"nestedObject": {"key": "value"}
		},
		"array": [{"key": "value"}, 123, null, true, false],
		"emptyObject": {},
		"emptyArray": []
	}`
	reader := strings.NewReader(json)
	err := Parse(reader)
	assert.NoError(t, err, "Parse() failed for complex JSON")
}

func TestParse_InvalidJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"Empty Input", ``},
		{"Whitespace Only", "   \n\t  "},
		{"Single Expression", `{}{}`},
		{"Missing Closing Bracket", `{"key": "value"`},
		{"Missing Comma", `{"key1": "value1" "key2": "value2"}`},
		{"Unexpected Token in Object", `{"key": "value", ]`},
		{"Unexpected Token in Array", `[1, 2, }`},
		{"Extra Comma in Object", `{"key1": "value1",}`},
		{"Extra Comma in Array", `[1, 2,]`},
		{"Unterminated String", `{"key": "unterminated}`},
		{"Missing Colon", `{"key" "value"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.json)
			err := Parse(reader)
			assert.Error(t, err, "Parse() did not fail for invalid JSON: %s", tt.json)
		})
	}
}
