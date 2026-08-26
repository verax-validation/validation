package verax

// Validatable indicates a type that carries its own validation logic.
// Nested objects participate in outer validation by implementing this interface,
// e.g. passed into verax.Valid as a field rule.
type Validatable interface {
	Validate() error
}

// Valid requires the value to pass its own Validate method, used for validating nested objects.
func Valid[T Validatable](value T) error {
	return value.Validate()
}
