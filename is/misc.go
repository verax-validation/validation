package is

import (
	"math"
	"strconv"
	"strings"

	"github.com/verax-validation/validation"
)

// Latitude requires the string to be a latitude within [-90, 90].
var Latitude verax.Rule[string] = coordinateRule(90, ErrLatitude)

// Longitude requires the string to be a longitude within [-180, 180].
var Longitude verax.Rule[string] = coordinateRule(180, ErrLongitude)

// SSN requires the string to be a US social security number;
// the region cannot be 000/666/9xx, the group cannot be 00, and the serial cannot be 0000.
var SSN verax.Rule[string] = func(value string) error {
	if !ssnPattern.MatchString(value) {
		return ErrSSN
	}
	area := value[:3]
	if area == "000" || area == "666" || area[0] == '9' {
		return ErrSSN
	}
	if value[4:6] == "00" || value[7:] == "0000" {
		return ErrSSN
	}
	return nil
}

// Semver requires the string to be a version conforming to semantic versioning 2.0.0.
var Semver verax.Rule[string] = func(value string) error {
	if semverPattern.MatchString(value) {
		return nil
	}
	return ErrSemver
}

// Origin requires the string to be a valid CORS origin:
// http or https scheme + host + optional port, no path allowed.
var Origin verax.Rule[string] = func(value string) error {
	if len(value) <= 255 && originPattern.MatchString(value) {
		return nil
	}
	return ErrOrigin
}

// DataURI requires the string to be an RFC 2397 data URI;
// when the base64 marker is present, the payload is decoded and validated.
var DataURI verax.Rule[string] = func(value string) error {
	m := dataURIPattern.FindStringSubmatch(value)
	if m == nil {
		return ErrDataURI
	}
	if m[3] == ";base64" {
		if _, err := base64Std.DecodeString(m[4]); err != nil {
			return ErrDataURI
		}
	}
	return nil
}

// coordinateRule builds a rule requiring the string to parse as a float within [-bound, bound].
func coordinateRule(bound float64, err *verax.Error) verax.Rule[string] {
	return func(value string) error {
		n, perr := strconv.ParseFloat(strings.TrimSpace(value), 64)
		// ParseFloat("NaN") reports no error and returns NaN, and NaN compares false against any number, so it must be intercepted separately
		if perr != nil || math.IsNaN(n) || n < -bound || n > bound {
			return err
		}
		return nil
	}
}
