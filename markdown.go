package main

type style rune

const (
	mdbold    style = '*'
	mdcursive style = '_'
	mdquote   style = '`'
)

// Markdown ensures safe markdown parsing on unsafe input
func Markdown(input string, style style) string {
	return string(style) + MDEscape(input) + string(style)
}

func MDEscape(input string) string {
	result := ""

	for _, char := range input {
		switch char {
		case '_', '*', '\\', '`':
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
