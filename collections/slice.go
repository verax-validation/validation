package collections

import (
	"iter"
	"strconv"

	"github.com/verax-validation/validation"
)

// Slice requires every element of the slice to pass rules in order; a failing element does not affect the validation of the others,
// and failures are aggregated into verax.Errors by index, e.g. {"1": error of element 1}.
func Slice[T any](rules ...verax.Rule[T]) verax.Rule[[]T] {
	return func(items []T) error {
		return collectErrors(func(visit func(string, T)) {
			for i, item := range items {
				visit(strconv.Itoa(i), item)
			}
		}, rules...)
	}
}

// Each requires every element of any iterable sequence (iter.Seq) to pass rules,
// with the same semantics as Slice, failures aggregated by the traversal-order index.
// Slices can be converted to a sequence with slices.Values.
func Each[T any](rules ...verax.Rule[T]) verax.Rule[iter.Seq[T]] {
	return func(seq iter.Seq[T]) error {
		if seq == nil {
			return nil
		}
		return collectErrors(func(visit func(string, T)) {
			i := 0
			for item := range seq {
				visit(strconv.Itoa(i), item)
				i++
			}
		}, rules...)
	}
}

// SliceLen requires the slice length to fall within the closed interval [min, max],
// used for collection size constraints like "select at least N".
func SliceLen[T any](min, max int) verax.Rule[[]T] {
	return func(items []T) error {
		if n := len(items); n < min || n > max {
			return verax.NewMessage(verax.CodeCollectionLen, map[string]string{
				"min": strconv.Itoa(min),
				"max": strconv.Itoa(max),
			})
		}
		return nil
	}
}
