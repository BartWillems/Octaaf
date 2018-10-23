package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MDEscape(t *testing.T) {
	assert.Equal(t, MDEscape(`some_string`), `some\_string`)
}

func Test_MDStyle(t *testing.T) {
	assert.Equal(t, Markdown(`some_cursive_string`, mdcursive), `_some\_cursive\_string_`)
	assert.Equal(t, Markdown(`some bold* string`, mdbold), `*some bold\* string*`)
	assert.Equal(t, Markdown("some` quoted string", mdquote), "`some quoted string`")
}
