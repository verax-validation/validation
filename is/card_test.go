package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestCreditCard(t *testing.T) {
	// common test card numbers with valid check digits
	valid := []string{
		"4111111111111111", // visa test number
		"5500005555555559", // mastercard test number
		"79927398713",      // Luhn algorithm Wikipedia example
	}
	for _, v := range valid {
		if err := is.CreditCard(v); err != nil {
			t.Errorf("CreditCard(%q) = %v, want nil", v, err)
		}
	}

	invalid := []string{
		"",
		"4111111111111112",    // wrong check digit
		"1234567890",          // does not pass Luhn
		"4111-1111-1111-1111", // contains separators
	}
	for _, v := range invalid {
		if err := is.CreditCard(v); err == nil {
			t.Errorf("CreditCard(%q) should fail", v)
		}
	}
}
