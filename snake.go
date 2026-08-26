package verax

import "strings"

// Snake converts camelCase names to snake_case.
// A run of uppercase letters is treated as one acronym, e.g. OrderID -> order_id, HTTPServer -> http_server.
func Snake(s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 4)

	runes := []rune(s)
	for i, r := range runes {
		if isUpperRune(r) {
			prevLower := i > 0 && !isUpperRune(runes[i-1])
			nextLower := i+1 < len(runes) && !isUpperRune(runes[i+1])
			if i > 0 && (prevLower || nextLower) {
				sb.WriteByte('_')
			}
			sb.WriteRune(toLowerRune(r))
			continue
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

func isUpperRune(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func toLowerRune(r rune) rune {
	if isUpperRune(r) {
		return r + ('a' - 'A')
	}
	return r
}
