package is

import (
	"encoding/json"

	"github.com/verax-validation/validation"
)

// UUID requires the string to be the standard 36-character hyphenated format, without version restriction, case-insensitive.
var UUID verax.Rule[string] = func(value string) error {
	if uuidPattern.MatchString(value) {
		return nil
	}
	return ErrUUID
}

// Base64 requires the string to be a valid standard base64 encoding (RFC 4648, with padding).
// An empty string is considered invalid, consistent with the strict mode of other rules in the package.
var Base64 verax.Rule[string] = func(value string) error {
	if len(value) == 0 || len(value)%4 != 0 {
		return ErrBase64
	}
	if _, err := base64Std.DecodeString(value); err != nil {
		return ErrBase64
	}
	return nil
}

// JSON requires the string to be a syntactically complete JSON document, supporting objects, arrays, scalars, and all other valid forms.
var JSON verax.Rule[string] = func(value string) error {
	if json.Valid([]byte(value)) {
		return nil
	}
	return ErrJSON
}
