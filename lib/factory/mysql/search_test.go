package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanFulltextOperator(t *testing.T) {
	result := CleanFulltextOperator(`a'b+c-d*e"f\g`)
	expected := `abcdefg`
	assert.Equal(t, expected, result)
}

func TestFindInJSON(t *testing.T) {
	result := FindInJSON(`categories`, 10, `*`)
	expected := "? MEMBER OF(`categories`->'$[*]')"
	assert.Equal(t, expected, result.String())

	result = FindInJSON(`categories`, 10)
	expected = "? MEMBER OF(`categories`)"
	assert.Equal(t, expected, result.String())

	result = FindInJSON(`tags`, `文章`, `name`)
	expected = "? MEMBER OF(`tags`->'$.name')"
	assert.Equal(t, expected, result.String())

	result = FindInJSON(`categories`, 20, `[1]`)
	expected = "? MEMBER OF(`categories`->'$[1]')"
	assert.Equal(t, expected, result.String())
}
