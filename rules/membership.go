package rules

import (
	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

var (
	// ErrIn is the base error for a failed In validation.
	ErrIn = verax.NewError(codes.CodeIn, "must be a valid value")
	// ErrNotIn is the base error for a failed NotIn validation.
	ErrNotIn = verax.NewError(codes.CodeNotIn, "must not be in list")
)

// In requires the value to appear in the given allowed list.
// The list is hashed into a set at construction, giving O(1) lookups at validation time;
// the comparable constraint ensures the value type and the list element type agree at compile time.
func In[T comparable](values ...T) verax.Rule[T] {
	set := make(map[T]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return func(value T) error {
		if _, ok := set[value]; ok {
			return nil
		}
		return verax.NewMessage(codes.CodeIn, nil)
	}
}

// NotIn requires the value not to appear in the given forbidden list, constructed like In.
func NotIn[T comparable](values ...T) verax.Rule[T] {
	banned := make(map[T]struct{}, len(values))
	for _, v := range values {
		banned[v] = struct{}{}
	}
	return func(value T) error {
		if _, ok := banned[value]; ok {
			return verax.NewMessage(codes.CodeNotIn, nil)
		}
		return nil
	}
}
