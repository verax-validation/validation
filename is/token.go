package is

import (
	"encoding/base64"
	"strings"

	"github.com/verax-validation/validation"
)

// Base64URL requires the string to be a valid URL-safe base64 encoding (RFC 4648 section 5, without padding),
// with the charset A-Za-z0-9_-.
var Base64URL verax.Rule[string] = func(value string) error {
	if len(value) > 0 && len(value)%4 != 1 && base64URLPattern.MatchString(value) {
		return nil
	}
	return ErrBase64URL
}

// JWT requires the string to be a valid JWT, i.e. three base64url segments joined by ".",
// each segment must decode as valid base64url.
var JWT verax.Rule[string] = func(value string) error {
	if isJWT(value) {
		return nil
	}
	return ErrJWT
}

// isJWT validates the three-segment JWT structure and that each segment decodes as base64url.
func isJWT(value string) bool {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return false
	}
	for _, part := range parts {
		if len(part) == 0 {
			return false
		}
		if _, err := base64.RawURLEncoding.DecodeString(part); err != nil {
			return false
		}
	}
	return true
}
