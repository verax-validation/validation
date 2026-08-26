package rules

import (
	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

// NotNil requires the pointer value not to be nil.
// Note: nil pointer fields are skipped directly by FieldBuilder, so this rule is for validating pointers
// directly with Validate, or validating pointer elements one by one in collections.
func NotNil[T any](value *T) error {
	if value == nil {
		return verax.NewMessage(codes.CodeNotNil, nil)
	}
	return nil
}
