package is

import (
	"strings"
	"time"

	"github.com/verax-validation/validation"
)

// Int requires the string to be a valid integer literal, allowing a sign.
var Int verax.Rule[string] = func(value string) error {
	if intPattern.MatchString(value) {
		return nil
	}
	return ErrInt
}

// Float requires the string to be a valid decimal floating-point literal, allowing a sign and scientific notation.
var Float verax.Rule[string] = func(value string) error {
	if floatPattern.MatchString(value) {
		return nil
	}
	return ErrFloat
}

// HexColor requires the string to be a three- or six-digit hex color code starting with #.
var HexColor verax.Rule[string] = func(value string) error {
	if hexColorPattern.MatchString(value) {
		return nil
	}
	return ErrHexColor
}

// RGBColor requires the string to be an "rgb(r, g, b)" color code, with components in 0-255.
var RGBColor verax.Rule[string] = func(value string) error {
	if rgbColorPattern.MatchString(value) {
		return nil
	}
	return ErrRGBColor
}

// Boolean requires the string to be a recognizable boolean value, accepting true/false/t/f/1/0, case-insensitive.
var Boolean verax.Rule[string] = func(value string) error {
	switch strings.ToLower(value) {
	case "true", "false", "1", "0", "t", "f":
		return nil
	}
	return ErrBoolean
}

// TimeZone requires the string to be a system-recognizable timezone name, e.g. "Asia/Shanghai" or "UTC".
var TimeZone verax.Rule[string] = func(value string) error {
	if len(value) == 0 {
		return ErrTimeZone
	}
	if _, err := time.LoadLocation(value); err != nil {
		return ErrTimeZone
	}
	return nil
}

// RGBA requires the string to be an "rgba(r, g, b, a)" color code,
// r/g/b in 0-255, a is a 0-1 decimal or a 0-100% percentage.
var RGBA verax.Rule[string] = func(value string) error {
	if rgbaPattern.MatchString(value) {
		return nil
	}
	return ErrRGBA
}

// HSL requires the string to be an "hsl(h, s%, l%)" color code, h in 0-360, s/l in 0-100.
var HSL verax.Rule[string] = func(value string) error {
	if hslPattern.MatchString(value) {
		return nil
	}
	return ErrHSL
}

// HSLA requires the string to be an "hsla(h, s%, l%, a)" color code,
// h in 0-360, s/l in 0-100, a is a 0-1 decimal or a 0-100% percentage.
var HSLA verax.Rule[string] = func(value string) error {
	if hslaPattern.MatchString(value) {
		return nil
	}
	return ErrHSLA
}
