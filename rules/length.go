package rules

import (
	"strconv"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

// stringish constrains length-measurable types: strings and byte slices and their underlying types.
type stringish interface {
	~string | ~[]byte
}

// Length requires the byte length of the value to fall within the closed interval [min, max].
// Only strings and byte slices are constrained; collection size validation is handled by the collections package.
// Parameters are written into the template at construction; rendered per the current language on failure.
func Length[S stringish](min, max int) verax.Rule[S] {
	return func(value S) error {
		if n := len(value); n < min || n > max {
			return verax.NewMessage(codes.CodeLength, map[string]string{
				"min": strconv.Itoa(min),
				"max": strconv.Itoa(max),
			})
		}
		return nil
	}
}
