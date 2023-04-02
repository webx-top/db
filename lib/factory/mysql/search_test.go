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
