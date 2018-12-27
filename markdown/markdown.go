package markdown

type style rune

const (
	bold    style = '*'
	cursive style = '_'
	quote   style = '`'
)

// Bold returns the input string escaped, with bold anotations
func Bold(input string) string {
	return string(bold) + Escape(input) + string(bold)
}

// Cursive returns the input string escaped, with cursive anotations
func Cursive(input string) string {
	return string(cursive) + Escape(input) + string(cursive)
}

// Quote returns the input string escaped, with quote anotations
func Quote(input string) string {
	return string(quote) + Escape(input) + string(quote)
}

// Escape returns a string with escaped markdown characters
func Escape(input string) string {
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
