package rules

import (
	"regexp"
	"time"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

// ErrMatch is the base error for a failed Match validation.
var ErrMatch = verax.NewError(codes.CodeMatch, "must be in a valid format")

// Match requires the string value to match the given regular expression.
// pattern is compiled at construction; an invalid expression panics immediately,
// consistent with regexp.MustCompile, fitting package-level variable construction scenarios.
func Match(pattern string) verax.Rule[string] {
	re := regexp.MustCompile(pattern)
	return func(value string) error {
		if re.MatchString(value) {
			return nil
		}
		return verax.NewMessage(codes.CodeMatch, nil)
	}
}

// Date requires the string value to match the given time layout (see time.Parse layout syntax, e.g. "2006-01-02").
func Date(layout string) verax.Rule[string] {
	return func(value string) error {
		if _, perr := time.Parse(layout, value); perr != nil {
			return verax.NewMessage(codes.CodeDate, map[string]string{
				"layout": layout,
			})
		}
		return nil
	}
}
