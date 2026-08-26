package is

import (
	"strings"

	"github.com/verax-validation/validation"
)

// ISBN10 requires the string to be a valid ISBN-10, allowing hyphens and spaces as separators.
var ISBN10 verax.Rule[string] = func(value string) error {
	if isValidISBN(compactDigits(value), 10) {
		return nil
	}
	return ErrISBN10
}

// ISBN13 requires the string to be a valid ISBN-13, allowing hyphens and spaces as separators.
var ISBN13 verax.Rule[string] = func(value string) error {
	if isValidISBN(compactDigits(value), 13) {
		return nil
	}
	return ErrISBN13
}

// ISBN requires the string to be a valid ISBN-10 or ISBN-13.
var ISBN verax.Rule[string] = func(value string) error {
	digits := compactDigits(value)
	if isValidISBN(digits, 10) || isValidISBN(digits, 13) {
		return nil
	}
	return ErrISBN
}

// compactDigits removes hyphens and spaces from the string.
func compactDigits(value string) string {
	replacer := strings.NewReplacer("-", "", " ", "")
	return replacer.Replace(value)
}

// isValidISBN validates a digit-only ISBN against the given version:
// ISBN-10 weights 10..2 summed and mod 11, check digit may be X;
// ISBN-13 alternates 1/3 weights summed and mod 10.
func isValidISBN(digits string, version int) bool {
	if len(digits) != version {
		return false
	}

	sum := 0
	switch version {
	case 10:
		for i := range 9 {
			if digits[i] < '0' || digits[i] > '9' {
				return false
			}
			sum += int(digits[i]-'0') * (10 - i)
		}
		last := digits[9]
		switch {
		case last == 'X' || last == 'x':
			sum += 10
		case last >= '0' && last <= '9':
			sum += int(last - '0')
		default:
			return false
		}
		return sum%11 == 0

	case 13:
		for i := range 13 {
			if digits[i] < '0' || digits[i] > '9' {
				return false
			}
			n := int(digits[i] - '0')
			if i%2 == 1 {
				n *= 3
			}
			sum += n
		}
		return sum%10 == 0

	default:
		return false
	}
}
