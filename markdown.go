package main

import (
	"regexp"
	"strings"
)

type style rune

const (
	mdbold    style = '*'
	mdcursive style = '_'
	mdquote   style = '`'
)

var re = regexp.MustCompile(`(\*|_)`)

// Markdown ensures safe markdown parsing on unsafe input
func Markdown(input string, style style) string {
	// When quoting, there is no need to escape other characters
	if style == mdquote {
		return string(style) + mdUnquote(input) + string(style)
	}
	return string(style) + MDEscape(input) + string(style)
}

// MDEscape escapes all telegram markdown characters & removes `
func MDEscape(input string) string {
	return mdUnquote(re.ReplaceAllString(input, `\$1`))
}

func mdUnquote(input string) string {
	return strings.Replace(input, "`", "", -1)
}
