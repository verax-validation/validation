package collections

import (
	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

// Unique requires the elements of the slice to be mutually distinct, for constraints like "tags must not repeat".
func Unique[T comparable]() verax.Rule[[]T] {
	return func(items []T) error {
		seen := make(map[T]struct{}, len(items))
		for _, item := range items {
			if _, dup := seen[item]; dup {
				return verax.NewMessage(codes.CodeCollectionUnique, nil)
			}
			seen[item] = struct{}{}
		}
		return nil
	}
}
