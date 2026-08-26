package collections

import "github.com/verax-validation/validation"

// collectErrors runs rules on every element and aggregates failures by the element's locating identifier,
// returning nil when all pass.
// locate returns the element's locating identifier: index for slices, key name for maps.
func collectErrors[T any](each func(func(locate string, item T)), rules ...verax.Rule[T]) error {
	var errs verax.Errors
	each(func(locate string, item T) {
		if err := verax.Validate(item, rules...); err != nil {
			errs = append(errs, verax.WithField(err, locate))
		}
	})
	if len(errs) == 0 {
		return nil
	}
	return errs
}
