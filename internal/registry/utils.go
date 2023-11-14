package registry

import (
	"strings"
	"unicode"
)

func toSnakeCase(s string) string {
	// Function to determine if a character is a separator
	isSeparator := func(r rune) bool {
		return r == '-' || unicode.IsSpace(r)
	}

	// Builder to construct the final string
	var builder strings.Builder
	for i, r := range s {
		if isSeparator(r) {
			if i > 0 {
				builder.WriteRune('_')
			}
		} else {
			builder.WriteRune(unicode.ToLower(r))
		}
	}

	return strings.Trim(builder.String(), "_")
}
