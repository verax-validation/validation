package collections

import (
	"fmt"
	"strconv"

	"github.com/verax-validation/validation"
)

// Map requires every value in the map to pass rules, with failures aggregated into verax.Errors by key name.
// Key names are converted to strings via their default formatting (%v).
func Map[K comparable, V any](rules ...verax.Rule[V]) verax.Rule[map[K]V] {
	return func(items map[K]V) error {
		return collectErrors(func(visit func(string, V)) {
			for key, item := range items {
				visit(fmt.Sprintf("%v", key), item)
			}
		}, rules...)
	}
}

// MapLen requires the number of key-value pairs in the map to fall within the closed interval [min, max].
func MapLen[K comparable, V any](min, max int) verax.Rule[map[K]V] {
	return func(items map[K]V) error {
		if n := len(items); n < min || n > max {
			return verax.NewMessage(verax.CodeCollectionLen, map[string]string{
				"min": strconv.Itoa(min),
				"max": strconv.Itoa(max),
			})
		}
		return nil
	}
}
