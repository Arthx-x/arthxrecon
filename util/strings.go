package util

import "strings"

// SanitizeString trims spaces from the input string.
func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}
