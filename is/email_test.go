package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestEmailFormat(t *testing.T) {
	// same format check as Email but without DNS lookup, suitable for offline and unit test environments
	checkRules(t, "EmailFormat", is.EmailFormat,
		[]string{"alice@example.com", "a.b+tag@sub.example.co.uk"},
		[]string{"", "plain", "@example.com", "a @example.com"})
}

func TestEmailWithMXLookup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping DNS dependent test in short mode")
	}

	checkRules(t, "Email", is.Email,
		[]string{"alice@gmail.com"},
		[]string{"plain@example.invalid", "no-such-domain-xyz123.example"})
}
