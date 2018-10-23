package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MDEscape(t *testing.T) {
	assert.Equal(t, MDEscape(`some\_string`), `some\_string`)
	assert.Equal(t, MDEscape(`some\\_string`), `some\\\_string`)
	assert.Equal(t, MDEscape(`some string\`), `some string\\`)
	assert.Equal(t, MDEscape(`_some * string`), `\_some \* string`)
}

func Test_MDStyle(t *testing.T) {
	assert.Equal(t, Markdown(`some_cursive_string`, mdcursive), `_some\_cursive\_string_`)
	assert.Equal(t, Markdown(`some bold* string`, mdbold), `*some bold\* string*`)
	assert.Equal(t, Markdown("some` quoted string", mdquote), "`some quoted string`")
}
