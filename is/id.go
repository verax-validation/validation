package is

import (
	"slices"

	"github.com/verax-validation/validation"
)

// UUIDv3 requires the string to be a version 3 UUID.
var UUIDv3 verax.Rule[string] = uuidVersionRule('3', ErrUUIDv3)

// UUIDv4 requires the string to be a version 4 UUID.
var UUIDv4 verax.Rule[string] = uuidVersionRule('4', ErrUUIDv4)

// UUIDv5 requires the string to be a version 5 UUID.
var UUIDv5 verax.Rule[string] = uuidVersionRule('5', ErrUUIDv5)

// UUIDv7 requires the string to be a version 7 UUID, and checks the variant bits conform to RFC 9562.
var UUIDv7 verax.Rule[string] = func(value string) error {
	if !hasUUIDVersion(value, '7') {
		return ErrUUIDv7
	}
	switch value[19] {
	case '8', '9', 'a', 'b', 'A', 'B':
		return nil
	default:
		return ErrUUIDv7
	}
}

// ULID requires the string to be a 26-character ULID (Crockford Base32, excluding I/L/O/U).
var ULID verax.Rule[string] = func(value string) error {
	if ulidPattern.MatchString(value) {
		return nil
	}
	return ErrULID
}

// MongoID requires the string to be a 24-hex-digit MongoDB ObjectID.
var MongoID verax.Rule[string] = func(value string) error {
	if mongoIDPattern.MatchString(value) {
		return nil
	}
	return ErrMongoID
}

// hasUUIDVersion reports whether the string is a valid UUID whose version matches one of the given set.
func hasUUIDVersion(value string, versions ...byte) bool {
	if len(value) != 36 || !uuidPattern.MatchString(value) {
		return false
	}
	version := value[14]
	return slices.Contains(versions, version)
}

func uuidVersionRule(version byte, err *verax.Error) verax.Rule[string] {
	return func(value string) error {
		if hasUUIDVersion(value, version) {
			return nil
		}
		return err
	}
}
