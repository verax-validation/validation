package is

import "github.com/verax-validation/validation"

// CreditCard requires the string to pass the Luhn algorithm, for mainstream credit and debit card numbers.
// Only the check digit and pure digits are validated, without judging card issuer ranges.
var CreditCard verax.Rule[string] = func(value string) error {
	if passesLuhn(value) {
		return nil
	}
	return ErrCreditCard
}

// passesLuhn accumulates from right to left per the Luhn algorithm:
// odd positions are added as-is, even positions are doubled and their digits are added separately;
// input containing non-digit characters or empty is rejected.
func passesLuhn(digits string) bool {
	if len(digits) == 0 {
		return false
	}

	sum := 0
	double := false
	for i := len(digits) - 1; i >= 0; i-- {
		c := digits[i]
		if c < '0' || c > '9' {
			return false
		}
		d := int(c - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}
