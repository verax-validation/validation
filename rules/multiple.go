package rules

import (
	"fmt"
	"math"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

// numberish constrains numeric types, excluding strings.
type numberish interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// MultipleOf requires the value to be an integer multiple of the given base, for integer and floating-point types.
// Computed via float64 internally; very large integers may lose precision.
func MultipleOf[T numberish](step T) verax.Rule[T] {
	return func(value T) error {
		denom := float64(step)
		if denom == 0 || math.Mod(float64(value), denom) != 0 {
			return verax.NewMessage(codes.CodeMultipleOf, map[string]string{
				"step": fmt.Sprintf("%v", step),
			})
		}
		return nil
	}
}
