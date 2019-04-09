package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Escape(t *testing.T) {
	assert.Equal(t, Escape(`some\_string`), `some\\\_string`)
	assert.Equal(t, Escape(`some_string`), `some\_string`)
	assert.Equal(t, Escape(`some string\`), `some string\\`)
	assert.Equal(t, Escape(`_some * string`), `\_some \* string`)
	assert.Equal(t, Escape(`pee_is_stored_in_the_brain`), `pee\_is\_stored\_in\_the\_brain`)
	assert.Equal(t, Escape(`¯\_(ツ)_/¯`), `¯\\\_(ツ)\_/¯`)
	assert.Equal(t, Escape("Quotes ` `"), "Quotes \\` \\`")
}

func Test_Prune(t *testing.T) {
	assert.Equal(t, Prune(`*some string*`), `some string`)
	assert.Equal(t, Prune(`*some string`), `some string`)
	assert.Equal(t, Prune(`_some_ string`), `some string`)
	assert.Equal(t, Prune(`some_ string*`), `some string`)
	assert.Equal(t, Prune(`## some string`), ` some string`)
	assert.Equal(t, Prune(`##__**`), ``)
	assert.Equal(t, Prune("Quotes ` `"), "Quotes  ")
}

func Test_MDStyle(t *testing.T) {
	assert.Equal(t, Cursive(`some_cursive_string`), `_some\_cursive\_string_`)
	assert.Equal(t, Bold(`some bold* string`), `*some bold\* string*`)
	assert.Equal(t, Quote("some` quoted string"), "`some\\` quoted string`")
}
