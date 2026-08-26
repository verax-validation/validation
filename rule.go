package verax

// Rule validates a value of type T; returning nil means pass, returning non-nil means fail.
// Built-in rules are generic functions matching this signature, instantiated by the compiler at the Validate call site,
// e.g. in Validate("abc", rules.Required), Required is auto-instantiated as Rule[string].
type Rule[T any] func(value T) error

// When runs rules only when condition is true and returns their result, otherwise treated as pass.
// condition is evaluated when the rule is constructed, suitable for conditions that are already determined before assembling the validation logic.
func When[T any](condition bool, rules ...Rule[T]) Rule[T] {
	return func(value T) error {
		if !condition {
			return nil
		}
		return Validate(value, rules...)
	}
}
