package verax_test

import (
	"errors"
	"testing"

	verax "github.com/verax-validation/validation"
	"github.com/verax-validation/validation/rules"
)

const localeTestCode = "validation.test_locale"

// newLocaleErr renders a test error with a min parameter in the currently active language
func newLocaleErr() *verax.Error {
	return verax.NewMessage(localeTestCode, map[string]string{"min": "3"})
}

func TestRegisterLocaleReplacesData(t *testing.T) {
	// registration takes effect immediately: a later-registered language replaces the current table
	verax.RegisterLocale(verax.LocaleZhCN, verax.MessageMap{localeTestCode: "最小值是 {{.min}}"})

	if got := newLocaleErr().Error(); got != "最小值是 3" {
		t.Errorf("zh render = %q, want %q", got, "最小值是 3")
	}

	// registering Japanese then overwrites it
	verax.RegisterLocale(verax.LocaleJa, verax.MessageMap{localeTestCode: "最小値は {{.min}}"})
	if got := newLocaleErr().Error(); got != "最小値は 3" {
		t.Errorf("ja render = %q, want %q", got, "最小値は 3")
	}

	verax.RegisterEn()
}

func TestRegisterLocaleFullReplace(t *testing.T) {
	verax.RegisterLocale(verax.LocaleZhCN, verax.MessageMap{localeTestCode: "第一版 {{.min}}"})
	verax.RegisterLocale(verax.LocaleZhCN, verax.MessageMap{localeTestCode: "第二版 {{.min}}"})

	// direct replacement, no leftover old value
	if got := newLocaleErr().Error(); got != "第二版 3" {
		t.Errorf("replaced render = %q, want %q", got, "第二版 3")
	}
	verax.RegisterEn()
}

func TestErrorsRenderWithFieldNames(t *testing.T) {
	verax.RegisterLocale(verax.LocaleZhCN, verax.MessageMap{localeTestCode: "最小值是 {{.min}}"})
	defer verax.RegisterEn()

	errs := verax.Errors{
		verax.WithField(newLocaleErr(), "age"),
		errors.New("plain error stays as is"),
	}

	want := "age: 最小值是 3; 1: plain error stays as is"
	if got := errs.Error(); got != want {
		t.Errorf("Errors() =\n%q\nwant\n%q", got, want)
	}
}

func TestUnknownCodeFallsBackToEnglish(t *testing.T) {
	// falls back to the built-in English table when the current table misses the code
	verax.RegisterLocale(verax.LocaleZhCN, verax.MessageMap{localeTestCode: "最小值是 {{.min}}"})
	defer verax.RegisterEn()

	other := verax.NewMessage("validation.required", nil)
	if got := other.Error(); got != "cannot be blank" {
		t.Errorf("fallback = %q, want %q", got, "cannot be blank")
	}
}

func TestBuiltinRuleRendersChinese(t *testing.T) {
	verax.RegisterZhCN()
	defer verax.RegisterEn()

	err := rules.Min(18)(17)
	if got := err.Error(); got != "不能小于 18" {
		t.Errorf("rules.Min(18)(17) = %q, want %q", got, "不能小于 18")
	}
}
