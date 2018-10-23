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

func mdUnquote(input string) string {
	return strings.Replace(input, "`", "", -1)
}

func MDEscape(input string) string {
	// Remove backslashes as they can't be escaped
	input = mdUnquote(input)

	result := ""

	// Boolean that tests if previous char is a backslash
	shouldSkip := false

	for pos, char := range input {
		// Previous character was \ so this character shouldn't be escaped
		if shouldSkip {
			shouldSkip = false
			result += string(char)
			continue
		}

		switch char {
		case '\\':
			// If the last character is a backslash, it should always be escaped
			if pos+1 == len(input) {
				result += escape(char)
			} else {
				result += string(char)
				shouldSkip = true
			}
		case '_', '*':
			result += escape(char)
		default:
			result += string(char)
		}
	}
	return result
}

func escape(input rune) string {
	return `\` + string(input)
}
