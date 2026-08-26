package is

import (
	"strings"

	"github.com/verax-validation/validation"
)

// CountryCode requires the string to be an ISO 3166-1 alpha-2 country code, case-insensitive.
var CountryCode verax.Rule[string] = func(value string) error {
	if _, ok := countryCodes[strings.ToUpper(value)]; ok {
		return nil
	}
	return ErrCountryCode
}

// CurrencyCode requires the string to be an ISO 4217 currency code, case-insensitive.
var CurrencyCode verax.Rule[string] = func(value string) error {
	if _, ok := currencyCodes[strings.ToUpper(value)]; ok {
		return nil
	}
	return ErrCurrencyCode
}

// LanguageCode requires the string to be an ISO 639-1 language code, case-insensitive.
var LanguageCode verax.Rule[string] = func(value string) error {
	if _, ok := languageCodes[strings.ToLower(value)]; ok {
		return nil
	}
	return ErrLanguageCode
}
