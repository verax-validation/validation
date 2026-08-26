package rules

import (
	"cmp"
	"fmt"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

// Min requires the value to be greater than or equal to min, for all ordered types (including lexicographic string comparison).
func Min[T cmp.Ordered](min T) verax.Rule[T] {
	return func(value T) error {
		if value < min {
			return verax.NewMessage(codes.CodeMin, map[string]string{
				"min": fmt.Sprintf("%v", min),
			})
		}
		return nil
	}
}

// Max requires the value to be less than or equal to max, for all ordered types (including lexicographic string comparison).
func Max[T cmp.Ordered](max T) verax.Rule[T] {
	return func(value T) error {
		if value > max {
			return verax.NewMessage(codes.CodeMax, map[string]string{
				"max": fmt.Sprintf("%v", max),
			})
		}
		return nil
	}
}

// Between requires the value to fall within the closed interval [lo, hi], for all ordered types.
func Between[T cmp.Ordered](lo, hi T) verax.Rule[T] {
	return func(value T) error {
		if value < lo || value > hi {
			return verax.NewMessage(codes.CodeBetween, map[string]string{
				"lo": fmt.Sprintf("%v", lo),
				"hi": fmt.Sprintf("%v", hi),
			})
		}
		return nil
	}
}
