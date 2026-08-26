package is

import (
	"time"

	"github.com/verax-validation/validation"
)

// MobileCN requires the string to be a valid mainland China mobile number (11 digits, starts with 1, second digit 3-9).
var MobileCN verax.Rule[string] = func(value string) error {
	if mobileCNPattern.MatchString(value) {
		return nil
	}
	return ErrMobileCN
}

// IDCardCN requires the string to be a valid 18-digit mainland China resident ID card number,
// checking the birth date and the trailing check digit.
var IDCardCN verax.Rule[string] = func(value string) error {
	if !idCardCNPattern.MatchString(value) {
		return ErrIDCardCN
	}
	if _, err := time.Parse("20060102", value[6:14]); err != nil {
		return ErrIDCardCN
	}
	check := value[17]
	if check == 'x' {
		check = 'X'
	}
	if check != idCardCNCheck(value[:17]) {
		return ErrIDCardCN
	}
	return nil
}

// PostalCodeCN requires the string to be a 6-digit mainland China postal code.
var PostalCodeCN verax.Rule[string] = func(value string) error {
	if postalCodeCNPattern.MatchString(value) {
		return nil
	}
	return ErrPostalCodeCN
}

// idCardCNWeights are the weighting factors for ID card check-digit calculation (GB 11643-1999).
var idCardCNWeights = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// idCardCNCodes are the check-digit characters corresponding to each mod-11 remainder.
var idCardCNCodes = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

// idCardCNCheck computes the check-digit character for the leading 17 digits.
func idCardCNCheck(digits string) byte {
	sum := 0
	for i := 0; i < len(digits); i++ {
		sum += int(digits[i]-'0') * idCardCNWeights[i]
	}
	return idCardCNCodes[sum%11]
}
