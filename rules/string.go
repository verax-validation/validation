package rules

import (
	"strconv"
	"strings"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

// Contains requires the string to contain the given substring.
func Contains(substr string) verax.Rule[string] {
	return func(value string) error {
		if !strings.Contains(value, substr) {
			return verax.NewMessage(codes.CodeContains, map[string]string{
				"substr": substr,
			})
		}
		return nil
	}
}

// StartWith requires the string to start with the given prefix.
func StartWith(prefix string) verax.Rule[string] {
	return func(value string) error {
		if !strings.HasPrefix(value, prefix) {
			return verax.NewMessage(codes.CodeStartWith, map[string]string{
				"prefix": prefix,
			})
		}
		return nil
	}
}

// EndWith requires the string to end with the given suffix.
func EndWith(suffix string) verax.Rule[string] {
	return func(value string) error {
		if !strings.HasSuffix(value, suffix) {
			return verax.NewMessage(codes.CodeEndWith, map[string]string{
				"suffix": suffix,
			})
		}
		return nil
	}
}

// Excludes requires the string not to contain the given substring.
func Excludes(substr string) verax.Rule[string] {
	return func(value string) error {
		if strings.Contains(value, substr) {
			return verax.NewMessage(codes.CodeExcludes, map[string]string{
				"substr": substr,
			})
		}
		return nil
	}
}

// ContainsAny requires the string to contain at least one of the given characters,
// used for constraints like "password must contain special characters".
func ContainsAny(chars string) verax.Rule[string] {
	return func(value string) error {
		if !strings.ContainsAny(value, chars) {
			return verax.NewMessage(codes.CodeContainsAny, map[string]string{
				"chars": chars,
			})
		}
		return nil
	}
}

// TrimSpace returns a rule that trims leading/trailing whitespace from the value before applying inner,
// used for "user input with spaces should still pass" scenarios; validation is based on the cleaned value.
func TrimSpace(inner verax.Rule[string]) verax.Rule[string] {
	return func(value string) error {
		return inner(strings.TrimSpace(value))
	}
}

// Len requires the byte length of the string or byte slice to be exactly n; use Length for interval semantics.
func Len[S stringish](n int) verax.Rule[S] {
	return func(value S) error {
		if len(value) != n {
			return verax.NewMessage(codes.CodeExactLen, map[string]string{
				"n": strconv.Itoa(n),
			})
		}
		return nil
	}
}
