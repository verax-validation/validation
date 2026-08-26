package is

import (
	"net/url"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

var (
	// ErrRequestURL is the error returned when RequestURL validation fails.
	ErrRequestURL = verax.NewError(codes.CodeRequestURL, "must be a valid request URL")
	// ErrRequestURI is the error returned when RequestURI validation fails.
	ErrRequestURI = verax.NewError(codes.CodeRequestURI, "must be a valid request URI")
)

// RequestURL requires the string to be an absolute URL with a scheme and host.
// The current implementation is the same as URL, kept as a separate name to express a clearer semantic intent.
var RequestURL = URL

// RequestURI requires the string to be a valid request URI,
// accepting both absolute URLs and relative paths starting with /.
var RequestURI verax.Rule[string] = func(value string) error {
	if _, err := url.ParseRequestURI(value); err != nil {
		return ErrRequestURI
	}
	return nil
}
