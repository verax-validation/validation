package is

import "github.com/verax-validation/validation"

// Multibyte requires the string to contain at least one multibyte (non-ASCII) character.
var Multibyte verax.Rule[string] = func(value string) error {
	if multibyteRe.MatchString(value) {
		return nil
	}
	return ErrMultibyte
}

// FullWidth requires the string to contain at least one full-width or wide character.
// Inclusive semantics: passes when present, not requiring the whole string to be full-width;
// highly overlapping with Multibyte but with a slightly different charset.
var FullWidth verax.Rule[string] = func(value string) error {
	if fullWidthSearch.MatchString(value) {
		return nil
	}
	return ErrFullWidth
}

// HalfWidth requires the string to contain at least one half-width character (ASCII and half-width katakana etc.).
var HalfWidth verax.Rule[string] = func(value string) error {
	if halfWidthSearch.MatchString(value) {
		return nil
	}
	return ErrHalfWidth
}

// VariableWidth requires the string to contain both full-width and half-width characters.
var VariableWidth verax.Rule[string] = func(value string) error {
	if fullWidthSearch.MatchString(value) && halfWidthSearch.MatchString(value) {
		return nil
	}
	return ErrVariableWidth
}
