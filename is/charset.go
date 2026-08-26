package is

import (
	"strings"
	"unicode"

	"github.com/verax-validation/validation"
)

// Alpha requires the string to contain English letters only.
var Alpha verax.Rule[string] = func(value string) error {
	if alphaPattern.MatchString(value) {
		return nil
	}
	return ErrAlpha
}

// Alphanumeric requires the string to contain English letters and digits only.
var Alphanumeric verax.Rule[string] = func(value string) error {
	if alphanumericPattern.MatchString(value) {
		return nil
	}
	return ErrAlphanumeric
}

// Digit requires the string to contain ASCII digits only.
var Digit verax.Rule[string] = func(value string) error {
	if digitPattern.MatchString(value) {
		return nil
	}
	return ErrDigit
}

// UTFLetter requires the string to contain only unicode letters of any language, and be non-empty.
var UTFLetter verax.Rule[string] = func(value string) error {
	return matchRunes(value, ErrUTFLetter, unicode.IsLetter)
}

// UTFDigit requires the string to contain only unicode decimal digits, and be non-empty.
var UTFDigit verax.Rule[string] = func(value string) error {
	return matchRunes(value, ErrUTFDigit, unicode.IsDigit)
}

// UTFNumeric requires the string to contain only unicode number characters (including fractions, roman numerals, etc., the Number class), and be non-empty.
var UTFNumeric verax.Rule[string] = func(value string) error {
	return matchRunes(value, ErrUTFNumeric, unicode.IsNumber)
}

// UTFLetterNumeric requires the string to contain only unicode letters or number characters, and be non-empty.
var UTFLetterNumeric verax.Rule[string] = func(value string) error {
	return matchRunes(value, ErrUTFLetterNumeric, func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsNumber(r)
	})
}

// LowerCase requires the string to be non-empty and all lowercase.
var LowerCase verax.Rule[string] = func(value string) error {
	if len(value) > 0 && value == strings.ToLower(value) {
		return nil
	}
	return ErrLowerCase
}

// UpperCase requires the string to be non-empty and all uppercase.
var UpperCase verax.Rule[string] = func(value string) error {
	if len(value) > 0 && value == strings.ToUpper(value) {
		return nil
	}
	return ErrUpperCase
}

// ASCII requires the string to be non-empty and contain only ASCII characters.
var ASCII verax.Rule[string] = func(value string) error {
	if hasASCIIOnly(value) && len(value) > 0 {
		return nil
	}
	return ErrASCII
}

// PrintableASCII requires the string to be non-empty and contain only printable ASCII characters (no control characters).
var PrintableASCII verax.Rule[string] = func(value string) error {
	for i := 0; i < len(value); i++ {
		if value[i] < 0x20 || value[i] > 0x7E {
			return ErrPrintableASCII
		}
	}
	if len(value) == 0 {
		return ErrPrintableASCII
	}
	return nil
}

// Hexadecimal requires the string to be a hexadecimal digit sequence.
var Hexadecimal verax.Rule[string] = func(value string) error {
	if hexadecimalPattern.MatchString(value) {
		return nil
	}
	return ErrHexadecimal
}

// matchRunes runs a predicate character by character; empty strings are rejected directly.
// The return type must be error rather than *verax.Error,
// otherwise a nil pointer boxed into the error interface becomes a non-nil error (classic Go pitfall).
func matchRunes(value string, err *verax.Error, predicate func(rune) bool) error {
	if len(value) == 0 {
		return err
	}
	for _, r := range value {
		if !predicate(r) {
			return err
		}
	}
	return nil
}

func hasASCIIOnly(value string) bool {
	for i := 0; i < len(value); i++ {
		if value[i] >= 0x80 {
			return false
		}
	}
	return true
}
