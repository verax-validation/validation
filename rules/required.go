// Package rules provides ready-to-use general validation rules.
//
// Empty-value semantics: default strict mode, rules always run;
// optional fields are declared explicitly with Optional, zero values short-circuit and pass.
package rules

import (
	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
	"github.com/verax-validation/validation/internal/zero"
)

// ErrRequired is the base error for a failed Required validation.
var ErrRequired = verax.NewError(codes.CodeRequired, "cannot be blank")

// Required requires the value not to be the zero value of its type.
// As a bare generic function, type parameters are inferred automatically when inlined into Validate/Field;
// on failure, a new error instance is rendered in the currently active language.
// Note: nil pointer fields bound to *FieldBuilder are skipped directly without triggering this rule;
// for required semantics bind a non-pointer field or dereference to a value type first.
func Required[T any](value T) error {
	if zero.IsZero(value) {
		return verax.NewMessage(codes.CodeRequired, nil)
	}
	return nil
}

// Optional wraps a set of rules as optional rules:
// zero values pass directly, otherwise the given rules run in order.
func Optional[T any](rules ...verax.Rule[T]) verax.Rule[T] {
	return func(value T) error {
		if zero.IsZero(value) {
			return nil
		}
		return verax.Validate(value, rules...)
	}
}

// RequiredIf requires the value not to be the zero value of its type when condition is true.
// Used for conditional required scenarios, e.g. a card number is required once a payment method is selected.
func RequiredIf[T any](condition func() bool) verax.Rule[T] {
	return func(value T) error {
		if condition() && zero.IsZero(value) {
			return verax.NewMessage(codes.CodeRequired, nil)
		}
		return nil
	}
}
