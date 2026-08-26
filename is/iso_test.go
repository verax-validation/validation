package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestCountryCode(t *testing.T) {
	checkRules(t, "CountryCode", is.CountryCode,
		[]string{"CN", "US", "JP", "DE", "br"},
		[]string{"", "XX", "USA", "CN2"})
}

func TestCurrencyCode(t *testing.T) {
	checkRules(t, "CurrencyCode", is.CurrencyCode,
		[]string{"USD", "CNY", "EUR", "jpy"},
		[]string{"", "ABC", "usd1", "US"})
}

func TestLanguageCode(t *testing.T) {
	checkRules(t, "LanguageCode", is.LanguageCode,
		[]string{"zh", "en", "ja", "de", "EN"},
		[]string{"", "xx", "eng", "zh-CN"})
}
