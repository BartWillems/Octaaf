package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MDEscape(t *testing.T) {
	assert.Equal(t, MDEscape(`some\_string`), `some\\\_string`)
	assert.Equal(t, MDEscape(`some_string`), `some\_string`)
	assert.Equal(t, MDEscape(`some string\`), `some string\\`)
	assert.Equal(t, MDEscape(`_some * string`), `\_some \* string`)
	assert.Equal(t, MDEscape(`pee_is_stored_in_the_brain`), `pee\_is\_stored\_in\_the\_brain`)
	assert.Equal(t, MDEscape(`¯\_(ツ)_/¯`), `¯\\\_(ツ)\_/¯`)
	assert.Equal(t, MDEscape("Quotes ` `"), "Quotes \\` \\`")
}

func Test_MDStyle(t *testing.T) {
	assert.Equal(t, Markdown(`some_cursive_string`, mdcursive), `_some\_cursive\_string_`)
	assert.Equal(t, Markdown(`some bold* string`, mdbold), `*some bold\* string*`)
	assert.Equal(t, Markdown("some` quoted string", mdquote), "`some\\` quoted string`")
}
