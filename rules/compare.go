package rules

import (
	"cmp"
	"fmt"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

// Eq requires the value to equal target, for all comparable types.
func Eq[T comparable](target T) verax.Rule[T] {
	return func(value T) error {
		if value != target {
			return verax.NewMessage(codes.CodeEq, map[string]string{
				"value": fmt.Sprintf("%v", target),
			})
		}
		return nil
	}
}

// Ne requires the value to differ from target, constructed like Eq.
func Ne[T comparable](target T) verax.Rule[T] {
	return func(value T) error {
		if value == target {
			return verax.NewMessage(codes.CodeNe, map[string]string{
				"value": fmt.Sprintf("%v", target),
			})
		}
		return nil
	}
}

// Gt requires the value to be strictly greater than min, complementing the closed-interval Min.
func Gt[T cmp.Ordered](min T) verax.Rule[T] {
	return func(value T) error {
		if value <= min {
			return verax.NewMessage(codes.CodeGt, map[string]string{
				"value": fmt.Sprintf("%v", min),
			})
		}
		return nil
	}
}

// Lt requires the value to be strictly less than max, complementing the closed-interval Max.
func Lt[T cmp.Ordered](max T) verax.Rule[T] {
	return func(value T) error {
		if value >= max {
			return verax.NewMessage(codes.CodeLt, map[string]string{
				"value": fmt.Sprintf("%v", max),
			})
		}
		return nil
	}
}
